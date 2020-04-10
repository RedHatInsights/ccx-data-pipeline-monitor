/*
Copyright © 2019, 2020 Red Hat, Inc.

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

package commands

import (
	"fmt"
)

// PrintHelp can be used to display help on (color) terminal.
func PrintHelp() {
	fmt.Println(colorizer.Magenta("HELP:"))
	fmt.Println()
	fmt.Println(colorizer.Blue("OC related commands:"))
	fmt.Println(colorizer.Yellow("login                    "), "login into OC")
	fmt.Println(colorizer.Yellow("get pods                 "), "get list of all pods + identify important ones")
	fmt.Println(colorizer.Yellow("get aggregator           "), "retrieve logs from aggregator pods")
	fmt.Println(colorizer.Yellow("get pipeline             "), "retrieve logs from ccx-data-pipeline pods")
	fmt.Println()
	fmt.Println(colorizer.Blue("Status commands:"))
	fmt.Println(colorizer.Yellow("status                   "), "print current status")
	fmt.Println()
	fmt.Println(colorizer.Blue("Analysis commands:"))
	fmt.Println(colorizer.Yellow("aggregator logs          "), "display aggregator logs")
	fmt.Println(colorizer.Yellow("aggregator statistic     "), "display aggregator statistic")
	fmt.Println(colorizer.Yellow("pipeline logs            "), "display pipeline logs")
	fmt.Println(colorizer.Yellow("pipeline statistic       "), "display pipeline statistic")
	fmt.Println()
	fmt.Println(colorizer.Blue("Other commands:"))
	fmt.Println(colorizer.Yellow("version                  "), "print version information")
	fmt.Println(colorizer.Yellow("authors                  "), "displays list of authors")
	fmt.Println(colorizer.Yellow("license                  "), "displays license used by this project")
	fmt.Println(colorizer.Yellow("version                  "), "print version information")
	fmt.Println(colorizer.Yellow("quit                     "), "quit the application")
	fmt.Println(colorizer.Yellow("exit                     "), "dtto")
	fmt.Println(colorizer.Yellow("bye                      "), "dtto")
	fmt.Println(colorizer.Yellow("help                     "), "this help")
	fmt.Println(colorizer.Yellow("?                        "), "this help")
	fmt.Println()
}
