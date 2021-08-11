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

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/ccx-data-pipeline-monitor/packages/commands/pipeline.html

import (
	"fmt"

	// "github.com/c-bata/go-prompt"

	"github.com/RedHatInsights/ccx-data-pipeline-monitor/analyser"
)

// DisplayPipelineStatistic function displays statistic gathered from ccx-data-pipeline logs
func DisplayPipelineStatistic() {
	fmt.Println(colorizer.Magenta("Popeline statistic"))
	analyser.PrintPipelineStatistic(colorizer)
}

// DisplayPipelineLogs function displays selected types of logs gathered from ccx-data-pipeline logs
func DisplayPipelineLogs() {
}
