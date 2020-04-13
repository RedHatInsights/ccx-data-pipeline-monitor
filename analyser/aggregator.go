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

	"github.com/logrusorgru/aurora"

	"github.com/RedHatInsights/ccx-data-pipeline-monitor/config"
)

type AggregatorLogEntry struct {
	Level        string `json:"level"`
	Time         string `json:"time"`
	Message      string `json:"message"`
	Type         string `json:"type"`
	Error        string `json:"error"`
	Topic        string `json:"topic"`
	Offset       int    `json:"offset"`
	Group        string `json:"group"`
	Organization int    `json:"organization"`
	Cluster      string `json:"cluster"`
}

var aggregatorEntries []AggregatorLogEntry = nil

func readAggregatorLogFile(filename string) ([]AggregatorLogEntry, error) {
	entries := []AggregatorLogEntry{}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entry := AggregatorLogEntry{}
		text := scanner.Text()
		err = json.Unmarshal([]byte(text), &entry)
		if err != nil {
			log.Println(err)
			log.Println(text)
		} else {
			entries = append(entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return entries, err
	}

	return entries, nil
}

func filterConsumedMessages(entries []AggregatorLogEntry) []AggregatorLogEntry {
	consumed := []AggregatorLogEntry{}

	for _, entry := range entries {
		if entry.Message == "Consumed" && entry.Group != "" {
			consumed = append(consumed, entry)
		}
	}
	return consumed
}

func filterByMessage(entries []AggregatorLogEntry, message string) []AggregatorLogEntry {
	filtered := []AggregatorLogEntry{}

	for _, entry := range entries {
		if entry.Message == message && entry.Topic != "" && entry.Organization != 0 && entry.Cluster != "" && entry.Group == "" {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func printStatisticLine(colorizer aurora.Aurora, what string, entries []AggregatorLogEntry, nextEntries []AggregatorLogEntry) {
	e := strconv.Itoa(len(entries))
	x := strconv.Itoa(len(nextEntries) - len(entries))
	fmt.Printf("%-12s %s messages (%s excluded)\n", what, colorizer.Blue(e), colorizer.Red(x))
}

func printAggregatorStatistic(colorizer aurora.Aurora, entries []AggregatorLogEntry) {
	consumed := filterConsumedMessages(entries)
	read := filterByMessage(entries, "Read")
	whitelisted := filterByMessage(entries, "Organization whitelisted")
	marshalled := filterByMessage(entries, "Marshalled")
	checked := filterByMessage(entries, "Time ok")
	stored := filterByMessage(entries, "Stored")

	printStatisticLine(colorizer, "Consumed", consumed, consumed)
	printStatisticLine(colorizer, "Read", read, consumed)
	printStatisticLine(colorizer, "Whitelisted", whitelisted, read)
	printStatisticLine(colorizer, "Marshalled", marshalled, whitelisted)
	printStatisticLine(colorizer, "Checked", checked, marshalled)
	printStatisticLine(colorizer, "Stored", stored, checked)
}

func printConsumedEntry(colorizer aurora.Aurora, i int, entry AggregatorLogEntry) {
	e := strconv.Itoa(i)
	fmt.Printf("%s %s %s %s %d\n", colorizer.Blue(e), entry.Time, entry.Group, entry.Topic, entry.Offset)
}

func printReadEntry(entry AggregatorLogEntry) {
	fmt.Printf("%s %s %s %d %d %s\n", entry.Time, entry.Group, entry.Topic, entry.Offset, entry.Organization, entry.Cluster)
}

func printErrorsForMessageWithOffset(colorizer aurora.Aurora, entries []AggregatorLogEntry, offset int) {
	for _, entry := range entries {
		if entry.Offset == offset && entry.Level == "error" {
			fmt.Printf("\t%s %s\n", entry.Time, entry.Error)

		}
	}
}

func printConsumedEntries(colorizer aurora.Aurora, entries []AggregatorLogEntry, notRead []AggregatorLogEntry) {
	for i, entry := range notRead {
		printConsumedEntry(colorizer, i+1, entry)
		printErrorsForMessageWithOffset(colorizer, entries, entry.Offset)
	}
}

func printReadEntries(colorizer aurora.Aurora, entries []AggregatorLogEntry, notRead []AggregatorLogEntry) {
	for _, entry := range notRead {
		printReadEntry(entry)
		printErrorsForMessageWithOffset(colorizer, entries, entry.Offset)
	}
}

func messageWithOffsetIn(entries []AggregatorLogEntry, offset int) bool {
	for _, entry := range entries {
		if entry.Offset == offset {
			return true
		}
	}
	return false
}

func diffEntryListsByOffset(list1 []AggregatorLogEntry, list2 []AggregatorLogEntry) []AggregatorLogEntry {
	diff := []AggregatorLogEntry{}
	for _, element := range list1 {
		if !messageWithOffsetIn(list2, element.Offset) {
			diff = append(diff, element)
		}
	}
	return diff
}

func getConsumedNotReadMessages(entries []AggregatorLogEntry) []AggregatorLogEntry {
	consumed := filterConsumedMessages(entries)
	read := filterByMessage(entries, "Read")
	return diffEntryListsByOffset(consumed, read)
}

func getNotWhitelistedMessages(entries []AggregatorLogEntry) []AggregatorLogEntry {
	read := filterByMessage(entries, "Read")
	whitelisted := filterByMessage(entries, "Organization whitelisted")
	return diffEntryListsByOffset(read, whitelisted)
}

func printConsumedNotRead(colorizer aurora.Aurora, entries []AggregatorLogEntry) {
	notRead := getConsumedNotReadMessages(entries)
	printConsumedEntries(colorizer, entries, notRead)
}

func printAggregatorNotWhitelisted(colorizer aurora.Aurora, entries []AggregatorLogEntry) {
	notWhitelisted := getNotWhitelistedMessages(entries)
	printReadEntries(colorizer, entries, notWhitelisted)
}

func ReadAggregatorLogFiles() (int, error) {
	var err error
	aggregatorEntries, err = readAggregatorLogFile(config.AggregatorLogFileName)
	if err != nil {
		return 0, err
	}
	return len(aggregatorEntries), nil
}

func PrintAggregatorStatistic(colorizer aurora.Aurora) {
	if aggregatorEntries == nil {
		fmt.Println(colorizer.Red("logs are not loaded"))
		return
	}
	if len(aggregatorEntries) == 0 {
		fmt.Println(colorizer.Red("empty log"))
		return
	}
	printAggregatorStatistic(colorizer, aggregatorEntries)
}

func PrintAggregatorConsumedNotReadMessages(colorizer aurora.Aurora) {
	if aggregatorEntries == nil {
		fmt.Println(colorizer.Red("logs are not loaded"))
		return
	}
	if len(aggregatorEntries) == 0 {
		fmt.Println(colorizer.Red("empty log"))
		return
	}
	printConsumedNotRead(colorizer, aggregatorEntries)
}
