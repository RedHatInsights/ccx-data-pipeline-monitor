// Copyright 2020, 2021 Red Hat, Inc
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

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/ccx-data-pipeline-monitor/packages/analyser/aggregator.html

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

// Filters
const (
	readFilter              = "Read"
	marshalledFilter        = "Marshalled"
	timeOkFilter            = "Time ok"
	storedFilter            = "Stored"
	consumedFilter          = "Consumed"
	organizationWhitelisted = "Organization whitelisted"
)

// Messages
const (
	emptyLog         = "Empty log"
	logsAreNotLoaded = "Logs are not loaded"
)

// Log leves for analyzed files
const (
	entryLevelError = "error"
)

// AggregatorLogEntry represents one log entry (record) read from log file.
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

	// disable "G304 (CWE-22): Potential file inclusion via variable"
	// #nosec G304
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
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

	// log file needs to be closed properly
	// try to close the file
	err = file.Close()

	// in case of error all we can do is to just log the error
	if err != nil {
		log.Println(err)
	}

	return entries, nil
}

func filterConsumedMessages(entries []AggregatorLogEntry) []AggregatorLogEntry {
	consumed := []AggregatorLogEntry{}

	for _, entry := range entries {
		if entry.Message == consumedFilter && entry.Group != "" {
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
	read := filterByMessage(entries, readFilter)
	whitelisted := filterByMessage(entries, organizationWhitelisted)
	marshalled := filterByMessage(entries, marshalledFilter)
	checked := filterByMessage(entries, timeOkFilter)
	stored := filterByMessage(entries, storedFilter)

	printStatisticLine(colorizer, consumedFilter, consumed, consumed)
	printStatisticLine(colorizer, readFilter, read, consumed)
	printStatisticLine(colorizer, "Whitelisted", whitelisted, read)
	printStatisticLine(colorizer, marshalledFilter, marshalled, whitelisted)
	printStatisticLine(colorizer, "Checked", checked, marshalled)
	printStatisticLine(colorizer, storedFilter, stored, checked)
}

func printConsumedEntry(colorizer aurora.Aurora, i int, entry AggregatorLogEntry) {
	e := strconv.Itoa(i)
	fmt.Printf("%5s  %s  %s  %s  %d\t", colorizer.Blue(e), colorizer.Gray(8, entry.Time), entry.Group, entry.Topic, colorizer.Cyan(entry.Offset))
}

func printReadEntry(colorizer aurora.Aurora, i int, entry AggregatorLogEntry) {
	e := strconv.Itoa(i)
	fmt.Printf("%5s  %s  %s  %s  %d  %d  %s\t", colorizer.Blue(e), colorizer.Gray(8, entry.Time), entry.Group, entry.Topic, colorizer.Cyan(entry.Offset), colorizer.Yellow(entry.Organization), entry.Cluster)
}

func printErrorsForMessageWithOffset(colorizer aurora.Aurora, entries []AggregatorLogEntry, offset int) {
	for _, entry := range entries {
		if entry.Offset == offset && entry.Level == entryLevelError {
			fmt.Printf("\t%s  %s\n", colorizer.Gray(8, entry.Time), colorizer.Red(entry.Error))
		}
	}
}

func printMessageForErrorsMessageWithOffset(colorizer aurora.Aurora, entries []AggregatorLogEntry, offset int) {
	for _, entry := range entries {
		if entry.Offset == offset && entry.Level == "error" {
			fmt.Printf("\t%s  %s\n", colorizer.Gray(8, entry.Time), colorizer.Red(entry.Message))
		}
	}
}

func printInfoForMessageWithOffset(colorizer aurora.Aurora, entries []AggregatorLogEntry, offset int) {
	for _, entry := range entries {
		if entry.Offset == offset && entry.Level == "info" {
			fmt.Printf("\t%s  %s\n", colorizer.Gray(8, entry.Time), colorizer.Red(entry.Error))
		}
	}
}

func printConsumedEntries(colorizer aurora.Aurora, entries []AggregatorLogEntry, notRead []AggregatorLogEntry) {
	for i, entry := range notRead {
		printConsumedEntry(colorizer, i+1, entry)
		printErrorsForMessageWithOffset(colorizer, entries, entry.Offset)
	}
	fmt.Println()
}

func printReadEntries(colorizer aurora.Aurora, entries []AggregatorLogEntry, notRead []AggregatorLogEntry) {
	for i, entry := range notRead {
		printReadEntry(colorizer, i+1, entry)
		printMessageForErrorsMessageWithOffset(colorizer, entries, entry.Offset)
	}
	fmt.Println()
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
	read := filterByMessage(entries, readFilter)
	return diffEntryListsByOffset(consumed, read)
}

func getNotWhitelistedMessages(entries []AggregatorLogEntry) []AggregatorLogEntry {
	read := filterByMessage(entries, readFilter)
	whitelisted := filterByMessage(entries, "Organization whitelisted")
	return diffEntryListsByOffset(read, whitelisted)
}

func getNotMarshalledMessages(entries []AggregatorLogEntry) []AggregatorLogEntry {
	whitelisted := filterByMessage(entries, "Organization whitelisted")
	marshalled := filterByMessage(entries, marshalledFilter)
	return diffEntryListsByOffset(whitelisted, marshalled)
}

func getNotCheckedMessages(entries []AggregatorLogEntry) []AggregatorLogEntry {
	marshalled := filterByMessage(entries, marshalledFilter)
	checked := filterByMessage(entries, timeOkFilter)
	return diffEntryListsByOffset(marshalled, checked)
}

func getNotStoredMessages(entries []AggregatorLogEntry) []AggregatorLogEntry {
	checked := filterByMessage(entries, timeOkFilter)
	stored := filterByMessage(entries, storedFilter)
	return diffEntryListsByOffset(checked, stored)
}

func printConsumedNotRead(colorizer aurora.Aurora, entries []AggregatorLogEntry) {
	notRead := getConsumedNotReadMessages(entries)
	printConsumedEntries(colorizer, entries, notRead)
}

func printNotWhitelisted(colorizer aurora.Aurora, entries []AggregatorLogEntry) {
	notWhitelisted := getNotWhitelistedMessages(entries)
	printReadEntries(colorizer, entries, notWhitelisted)
}

func printWhitelistedNotMarshalled(colorizer aurora.Aurora, entries []AggregatorLogEntry) {
	notMarshalled := getNotMarshalledMessages(entries)
	printReadEntries(colorizer, entries, notMarshalled)
}

func printMarshalledNotChecked(colorizer aurora.Aurora, entries []AggregatorLogEntry) {
	notChecked := getNotCheckedMessages(entries)
	printReadEntries(colorizer, entries, notChecked)
}

func printCheckedNotStored(colorizer aurora.Aurora, entries []AggregatorLogEntry) {
	notStored := getNotStoredMessages(entries)
	printReadEntries(colorizer, entries, notStored)
}

// ReadAggregatorLogFiles reads all log files gathered from aggregator pods.
func ReadAggregatorLogFiles() (int, error) {
	var err error
	aggregatorEntries, err = readAggregatorLogFile(config.AggregatorLogFileName)
	if err != nil {
		return 0, err
	}
	return len(aggregatorEntries), nil
}

// PrintAggregatorStatistic prints statistic gathered from aggregator logs.
func PrintAggregatorStatistic(colorizer aurora.Aurora) {
	if aggregatorEntries == nil {
		fmt.Println(colorizer.Red(logsAreNotLoaded))
		return
	}
	if len(aggregatorEntries) == 0 {
		fmt.Println(colorizer.Red(emptyLog))
		return
	}
	printAggregatorStatistic(colorizer, aggregatorEntries)
}

// PrintAggregatorConsumedNotReadMessages function prints all messages that are consumer (from input) but not read for any reason
func PrintAggregatorConsumedNotReadMessages(colorizer aurora.Aurora) {
	if aggregatorEntries == nil {
		fmt.Println(colorizer.Red(logsAreNotLoaded))
		return
	}
	if len(aggregatorEntries) == 0 {
		fmt.Println(colorizer.Red(emptyLog))
		return
	}
	printConsumedNotRead(colorizer, aggregatorEntries)
}

// PrintAggregatorConsumedNotWhitelisted function prints all consumed, but not whitelisted messages, ie. messages that have been filtered
func PrintAggregatorConsumedNotWhitelisted(colorizer aurora.Aurora) {
	if aggregatorEntries == nil {
		fmt.Println(colorizer.Red(logsAreNotLoaded))
		return
	}
	if len(aggregatorEntries) == 0 {
		fmt.Println(colorizer.Red(emptyLog))
		return
	}
	printNotWhitelisted(colorizer, aggregatorEntries)
}

// PrintAggregatorWhitelistedNotMarshalled function prints whitelisted messages (that are supposed to be processed) that can't be marshalled for any reason
func PrintAggregatorWhitelistedNotMarshalled(colorizer aurora.Aurora) {
	if aggregatorEntries == nil {
		fmt.Println(colorizer.Red(logsAreNotLoaded))
		return
	}
	if len(aggregatorEntries) == 0 {
		fmt.Println(colorizer.Red(emptyLog))
		return
	}
	printWhitelistedNotMarshalled(colorizer, aggregatorEntries)
}

// PrintAggregatorMarshalledNotChecked function prints messages that can be marshalled but are not checked for any reason (improper internal structure etc.)
func PrintAggregatorMarshalledNotChecked(colorizer aurora.Aurora) {
	if aggregatorEntries == nil {
		fmt.Println(colorizer.Red(logsAreNotLoaded))
		return
	}
	if len(aggregatorEntries) == 0 {
		fmt.Println(colorizer.Red(emptyLog))
		return
	}
	printMarshalledNotChecked(colorizer, aggregatorEntries)
}

// PrintAggregatorCheckedNotStored function prints all messages that have been checked but not stored into database for whatever reason
func PrintAggregatorCheckedNotStored(colorizer aurora.Aurora) {
	if aggregatorEntries == nil {
		fmt.Println(colorizer.Red(logsAreNotLoaded))
		return
	}
	if len(aggregatorEntries) == 0 {
		fmt.Println(colorizer.Red(emptyLog))
		return
	}
	printCheckedNotStored(colorizer, aggregatorEntries)
}
