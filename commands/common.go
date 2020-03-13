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

package commands

import (
	"fmt"
	"os"

	"github.com/c-bata/go-prompt"
	"github.com/logrusorgru/aurora"
)

var colorizer aurora.Aurora

// SetColorizer set the terminal colorizer
func SetColorizer(c aurora.Aurora) {
	colorizer = c
}

// NoOpCompleter implements a no-op completer needed to input random data
func NoOpCompleter(in prompt.Document) []prompt.Suggest {
	return nil
}

// ProceedQuestion ask user about y/n answer.
func ProceedQuestion(question string) bool {
	fmt.Println(colorizer.Red(question))
	proceed := prompt.Input("proceed? [y/n] ", NoOpCompleter)
	if proceed != "y" {
		fmt.Println(colorizer.Blue("cancelled"))
		return false
	}
	return true
}

// Quit will exit from the CLI client
func Quit() {
	fmt.Println(colorizer.Magenta("Quitting"))
	os.Exit(0)
}
