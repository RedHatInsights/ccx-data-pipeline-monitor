/*
Copyright © 2020 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/RedHatInsights/ccx-data-pipeline-monitor/commands"
	"github.com/RedHatInsights/ccx-data-pipeline-monitor/config"
	"github.com/RedHatInsights/ccx-data-pipeline-monitor/oc"
	"github.com/RedHatInsights/ccx-data-pipeline-monitor/server"
)

var openShiftConfig config.OpenShiftConfig

var ocLogin string
var colorizer aurora.Aurora

// BuildVersion contains the major.minor version of the CLI client
var BuildVersion string = "*not set*"

// BuildTime contains timestamp when the CLI client has been built
var BuildTime string = "*not set*"

func printVersion() {
	fmt.Println(colorizer.Blue("Insights operator CLI client "), "version", colorizer.Yellow(BuildVersion), "compiled", colorizer.Yellow(BuildTime))
}

func tryToLogin(ocLogin string) {
	stdout, stderr, err := oc.Login(openShiftConfig.URL, ocLogin)
	fmt.Println(stdout)
	fmt.Println(stderr)
	fmt.Println(err == nil)
	fmt.Println(colorizer.Blue("\nDone"))
}

func login() {
	fmt.Print("login: ")
	p, err := terminal.ReadPassword(0)
	if err != nil {
		fmt.Println(colorizer.Red("not set"))
	} else {
		ocLogin = string(p)
		tryToLogin(ocLogin)
	}
}

type simpleCommand struct {
	prefix  string
	handler func()
}

var simpleCommands = []simpleCommand{
	{"bye", commands.Quit},
	{"exit", commands.Quit},
	{"quit", commands.Quit},
	{"login", login},
	{"?", commands.PrintHelp},
	{"help", commands.PrintHelp},
	{"version", printVersion},
	{"license", commands.PrintLicense},
	{"authors", commands.PrintAuthors},
}

func executeFixedCommand(t string) {
	// simple commands without parameters
	for _, command := range simpleCommands {
		if strings.HasPrefix(t, command.prefix) {
			command.handler()
			return
		}
	}
	fmt.Println("Command not found")
}

func executor(t string) {
	executeFixedCommand(t)
}

func completer(in prompt.Document) []prompt.Suggest {
	firstWord := []prompt.Suggest{
		{Text: "exit", Description: "quit the application"},
		{Text: "quit", Description: "quit the application"},
		{Text: "bye", Description: "quit the application"},

		{Text: "help", Description: "show help with all commands"},
		{Text: "version", Description: "prints the build information for CLI executable"},
		{Text: "copyright", Description: "displays copyright notice"},
		{Text: "license", Description: "displays license used by this project"},
		{Text: "authors", Description: "displays list of authors"},

		{Text: "login", Description: "provide login info"},
	}

	blocks := strings.Split(in.TextBeforeCursor(), " ")

	// don't display compation for empty command
	if in.GetWordBeforeCursor() == "" {
		return nil
	}

	// commands consisting of just one word
	return prompt.FilterHasPrefix(firstWord, blocks[0], true)
}

func loadConfiguration(defaultConfigName string, envVar string) error {
	log.Println("Reading configuration")
	configFile, specified := os.LookupEnv("CCX_DATA_PIPELINE_MONITOR")
	if specified {
		// we need to separate the directory name and filename without extension
		directory, basename := filepath.Split(configFile)
		file := strings.TrimSuffix(basename, filepath.Ext(basename))
		// parse the configuration
		viper.SetConfigName(file)
		viper.AddConfigPath(directory)
	} else {
		// parse the configuration
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	}
	defer log.Println("Done")
	return viper.ReadInConfig()
}

func startCLI() {
	// parse command line arguments and flags
	var colors = flag.Bool("colors", true, "enable or disable colors")
	var useCompleter = flag.Bool("completer", true, "enable or disable command line completer")
	// var askForConfirmation = flag.Bool("confirmation", true, "enable or disable asking for confirmation for selected actions (like delete)")
	flag.Parse()

	colorizer = aurora.NewAurora(*colors)
	commands.SetColorizer(colorizer)

	if *useCompleter {
		p := prompt.New(executor, completer)
		p.Run()
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("> ")
		for scanner.Scan() {
			line := scanner.Text()
			executor(line)
			fmt.Print("> ")
		}
	}
}

func startWebUI() {
	serverConfig := config.ReadServerConfig()
	server := server.New(serverConfig)
	server.Start()
	defer server.Stop(context.TODO())
}

func main() {
	// read configuration first
	err := loadConfiguration("config", "CCX_DATA_PIPELINE_MONITOR")
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	uiType := viper.Sub("ui").GetString("type")
	openShiftConfig = config.ReadOpenShiftConfig()

	switch uiType {
	case "cli":
		startCLI()
	case "web":
		startWebUI()
	default:
		log.Fatal("Unknown UI type", uiType)
	}
}
