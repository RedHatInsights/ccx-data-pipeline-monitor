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
	"io/ioutil"
	"strings"

	"github.com/RedHatInsights/ccx-data-pipeline-monitor/config"
	"github.com/RedHatInsights/ccx-data-pipeline-monitor/oc"
)

var aggregatorPod string = ""
var pipelinePod string = ""

// TryToLogin tries to login to OpenShift via oc command
func TryToLogin(url string, ocLogin string) bool {
	stdout, stderr, err := oc.Login(url, ocLogin)
	if err != nil {
		fmt.Println(colorizer.Red("\nUnable to login to OpenShift"))
		fmt.Println(stdout)
		fmt.Println(stderr)
		return false
	}
	fmt.Println(colorizer.Green("\nDone: you have been loged in to OpenShift"))
	return true
}

// GetPods function retrieves list of pods available for given user
func GetPods() {
	stdout, stderr, err := oc.GetPods()
	if err != nil {
		fmt.Println(colorizer.Red("\nUnable to get pods"))
		fmt.Println(stdout)
		fmt.Println(stderr)
		return
	}
	fmt.Println(colorizer.Blue("List of available pods"))
	fmt.Println(stdout)
	lines := strings.Split(stdout, "\n")

	aggregatorPod = ""
	pipelinePod = ""
	for _, line := range lines {
		if strings.HasPrefix(line, "ccx-data-pipeline") && !strings.HasPrefix(line, "ccx-data-pipeline-db") {
			pipelinePod = strings.Fields(line)[0]
		}
		if strings.HasPrefix(line, "insights-results-aggregator") {
			aggregatorPod = strings.Fields(line)[0]
		}
	}

	fmt.Print(colorizer.Blue("Aggregator pod: "))
	if aggregatorPod != "" {
		fmt.Println(aggregatorPod)
	} else {
		fmt.Println(colorizer.Red("not found"))
	}

	fmt.Print(colorizer.Blue("Pipeline pod:   "))
	if pipelinePod != "" {
		fmt.Println(pipelinePod)
	} else {
		fmt.Println(colorizer.Red("not found"))
	}
}

// GetLogs function retrieves logs from selected pod and stores logs in file.
func GetLogs(pod string, storeto string) {
	stdout, stderr, err := oc.GetLogs(pod)
	if err != nil {
		fmt.Println(colorizer.Red("\nUnable to read logs"))
		fmt.Println(stderr)
		return
	}
	fmt.Println(colorizer.Green("Logs have been read"))
	fmt.Printf("Log file size: %d bytes\n", len(stdout))

	err = ioutil.WriteFile(storeto, []byte(stdout), 0600)
	if err != nil {
		fmt.Println(colorizer.Red("\nUnable to write logs"))
		fmt.Println(err)
		return
	}
	fmt.Println(colorizer.Blue("Written into " + storeto))
}

// GetAggregatorLogs function retrieves logs from aggregator pods and stores logs in file.
func GetAggregatorLogs() {
	if aggregatorPod == "" {
		fmt.Println(colorizer.Red("Aggregator pod was not found"))
		return
	}
	GetLogs(aggregatorPod, config.AggregatorLogFileName)

}

// GetPipelineLogs function retrieves logs from ccx-data-pipeline pods and stores logs in file.
func GetPipelineLogs() {
	if pipelinePod == "" {
		fmt.Println(colorizer.Red("Pipeline pod was not found"))
		return
	}
	GetLogs(pipelinePod, config.PipelineLogFileName)
}
