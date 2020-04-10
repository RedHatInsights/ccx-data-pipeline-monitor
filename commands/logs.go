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

	"github.com/RedHatInsights/ccx-data-pipeline-monitor/analyser"
)

func loadAggregatorLogs() {
	fmt.Println(colorizer.Blue("Aggregator logs"))
	entries, err := analyser.ReadAggregatorLogFiles()
	if err != nil {
		fmt.Println(colorizer.Red(err))
	}
	fmt.Println(colorizer.Green("Success:"), "read", colorizer.Blue(entries), "entries\n")
}

func loadPipelineLogs() {
	fmt.Println(colorizer.Blue("CCX data pipeline logs"))
	entries, err := analyser.ReadPipelineLogFiles()
	if err != nil {
		fmt.Println(colorizer.Red(err))
	}
	fmt.Println(colorizer.Green("Success:"), "read", colorizer.Blue(entries), "entries\n")
}

func LoadLogs() {
	fmt.Println(colorizer.Magenta("Loading logs"))
	loadAggregatorLogs()
	loadPipelineLogs()
}
