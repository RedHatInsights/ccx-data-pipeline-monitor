package main

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/ccx-data-pipeline-monitor/packages/temp/analyser.html

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Filters
const (
	readFilter       = "Read"
	storedFilter     = "Stored"
	marshalledFilter = "Marshalled"
)

// PipelineLogEntry represents one log entry (record) read from log file.
type PipelineLogEntry struct {
	Level    string `json:"levelname"`
	Time     string `json:"asctime"`
	Name     string `json:"name"`
	Filename string `json:"filename"`
	Message  string `json:"message"`
}

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

func readPipelineLogFile(filename string) ([]PipelineLogEntry, error) {
	entries := []PipelineLogEntry{}

	// disable "G304 (CWE-22): Potential file inclusion via variable"
	// #nosec G304
	file, err := os.Open(filename)
	if err != nil {
		return entries, err
	}

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

	// log file needs to be closed properly

	// try to close the file
	err = file.Close()

	// in case of error all we can do is to just log the error
	if err != nil {
		log.Println(err)
	}

	return entries, nil
}

func readAggregatorLogFile(filename string) ([]AggregatorLogEntry, error) {
	entries := []AggregatorLogEntry{}

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

func printStatisticLine(what string, entries []AggregatorLogEntry) {
	fmt.Printf("%-12s %d messages\n", what, len(entries))
}

func printAggregatorStatistic(entries []AggregatorLogEntry) {
	consumed := filterConsumedMessages(entries)
	read := filterByMessage(entries, readFilter)
	whitelisted := filterByMessage(entries, "Organization whitelisted")
	marshalled := filterByMessage(entries, marshalledFilter)
	checked := filterByMessage(entries, "Time ok")
	stored := filterByMessage(entries, storedFilter)

	printStatisticLine("Consumed", consumed)
	printStatisticLine(readFilter, read)
	printStatisticLine("Whitelisted", whitelisted)
	printStatisticLine("Marshalled messages", marshalled)
	printStatisticLine("Checked", checked)
	printStatisticLine(storedFilter, stored)
}

func printConsumedEntry(entry AggregatorLogEntry) {
	fmt.Printf("%s %s %s %d\n", entry.Time, entry.Group, entry.Topic, entry.Offset)
}

func printReadEntry(entry AggregatorLogEntry) {
	fmt.Printf("%s %s %s %d %d %s\n", entry.Time, entry.Group, entry.Topic, entry.Offset, entry.Organization, entry.Cluster)
}

func printErrorsForMessageWithOffset(entries []AggregatorLogEntry, offset int) {
	for _, entry := range entries {
		if entry.Offset == offset && entry.Level == "error" {
			fmt.Printf("\t%s %s\n", entry.Time, entry.Error)

		}
	}
}

func printConsumedEntries(entries []AggregatorLogEntry, notRead []AggregatorLogEntry) {
	for _, entry := range notRead {
		printConsumedEntry(entry)
		printErrorsForMessageWithOffset(entries, entry.Offset)
	}
}

func printReadEntries(entries []AggregatorLogEntry, notRead []AggregatorLogEntry) {
	for _, entry := range notRead {
		printReadEntry(entry)
		printErrorsForMessageWithOffset(entries, entry.Offset)
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
	read := filterByMessage(entries, readFilter)
	return diffEntryListsByOffset(consumed, read)
}

func getNotWhitelistedMessages(entries []AggregatorLogEntry) []AggregatorLogEntry {
	read := filterByMessage(entries, readFilter)
	whitelisted := filterByMessage(entries, "Organization whitelisted")
	return diffEntryListsByOffset(read, whitelisted)
}

func printConsumedNotRead(entries []AggregatorLogEntry) {
	notRead := getConsumedNotReadMessages(entries)
	printConsumedEntries(entries, notRead)
}

func printAggregatorNotWhitelisted(entries []AggregatorLogEntry) {
	notWhitelisted := getNotWhitelistedMessages(entries)
	printReadEntries(entries, notWhitelisted)
}

func analyse() {
	/*
		entries, err := readPipelineLogFile("pipeline3.log")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(entries)*/

	entries2, err := readAggregatorLogFile("aggregator3.log")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Read:", len(entries2), "log entries read")
	printAggregatorStatistic(entries2)
	// printConsumedNotRead(entries2)
	printAggregatorNotWhitelisted(entries2)
}

func main() {
	analyse()
}
