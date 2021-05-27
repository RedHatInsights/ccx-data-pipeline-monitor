// Copyright 2020 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package analyser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/logrusorgru/aurora"

	"github.com/RedHatInsights/ccx-data-pipeline-monitor/config"
)

// PipelineLogEntry represents one log entry (record) read from log file.
type PipelineLogEntry struct {
	Level    string `json:"levelname"`
	Time     string `json:"asctime"`
	Name     string `json:"name"`
	Filename string `json:"filename"`
	Message  string `json:"message"`
}

var pipelineEntries []PipelineLogEntry = nil

func readPipelineLogFile(filename string) ([]PipelineLogEntry, error) {
	entries := []PipelineLogEntry{}

	// disable "G304 (CWE-22): Potential file inclusion via variable"
	// #nosec G304
	file, err := os.Open(filename)
	if err != nil {
		return entries, err
	}
	// log file needs to be closed properly
	defer func() {
		// try to close the file
		err := file.Close()

		// in case of error all we can do is to just log the error
		if err != nil {
			log.Println(err)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entry := PipelineLogEntry{}
		err = json.Unmarshal([]byte(scanner.Text()), &entry)
		if err != nil {
			log.Println(err)
		} else {
			entries = append(entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return entries, err
	}

	return entries, nil
}

func filterPipelineMessagesByMessage(entries []PipelineLogEntry, prefix string) []PipelineLogEntry {
	filtered := []PipelineLogEntry{}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Message, prefix) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func printStatisticLinePipeline(colorizer aurora.Aurora, what string, entries []PipelineLogEntry) {
	e := strconv.Itoa(len(entries))
	fmt.Printf("%-26s %s messages\n", what, colorizer.Blue(e))
}

func printPipelineStatistic(colorizer aurora.Aurora, entries []PipelineLogEntry) {
	validated1 := filterPipelineMessagesByMessage(entries, "JSON schema validated")
	validated2 := filterPipelineMessagesByMessage(entries, "Identity schema validated")
	downloaded := filterPipelineMessagesByMessage(entries, "Downloading ")
	saved := filterPipelineMessagesByMessage(entries, "Saved ")
	sendStart := filterPipelineMessagesByMessage(entries, "Sending response to the ")
	sendSuccess := filterPipelineMessagesByMessage(entries, "Message has been sent successfully")
	contextRetrieved := filterPipelineMessagesByMessage(entries, "Message context: ")
	success := filterPipelineMessagesByMessage(entries, "Status: Success; ")

	printStatisticLinePipeline(colorizer, "JSON schema validated", validated1)
	printStatisticLinePipeline(colorizer, "Identity schema validated", validated2)
	printStatisticLinePipeline(colorizer, "Downloaded", downloaded)
	printStatisticLinePipeline(colorizer, "Saved", saved)
	printStatisticLinePipeline(colorizer, "Sending start", sendStart)
	printStatisticLinePipeline(colorizer, "Sending successful", sendSuccess)
	printStatisticLinePipeline(colorizer, "Context retrieved", contextRetrieved)
	printStatisticLinePipeline(colorizer, "Success", success)
}

// ReadPipelineLogFiles reads all log files gathered from CCX data pipeline pods.
func ReadPipelineLogFiles() (int, error) {
	var err error
	pipelineEntries, err = readPipelineLogFile(config.PipelineLogFileName)
	if err != nil {
		return 0, err
	}
	return len(pipelineEntries), nil
}

// PrintPipelineStatistic prints statistic gathered from CCX data pipeline logs.
func PrintPipelineStatistic(colorizer aurora.Aurora) {
	if pipelineEntries == nil {
		fmt.Println(colorizer.Red("logs are not loaded"))
		return
	}
	if len(pipelineEntries) == 0 {
		fmt.Println(colorizer.Red("empty log"))
		return
	}
	printPipelineStatistic(colorizer, pipelineEntries)
}
