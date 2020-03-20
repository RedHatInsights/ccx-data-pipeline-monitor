package analyser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type PipelineLogEntry struct {
	Level    string `json:"levelname"`
	Time     string `json:"asctime"`
	Name     string `json:"name"`
	Filename string `json:"filename"`
	Message  string `json:"message"`
}

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

	file, err := os.Open(filename)
	if err != nil {
		return entries, err
	}
	defer file.Close()

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

func readAggregatorLogFile(filename string) ([]AggregatorLogEntry, error) {
	entries := []AggregatorLogEntry{}

	file, err := os.Open(filename)
	if err != nil {
		return entries, err
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

func printStatisticLine(what string, entries []AggregatorLogEntry) {
	fmt.Printf("%-12s %d messages\n", what, len(entries))
}

func printAggregatorStatistic(entries []AggregatorLogEntry) {
	consumed := filterConsumedMessages(entries)
	read := filterByMessage(entries, "Read")
	whitelisted := filterByMessage(entries, "Organization whitelisted")
	marshalled := filterByMessage(entries, "Marshalled")
	checked := filterByMessage(entries, "Time ok")
	stored := filterByMessage(entries, "Stored")

	printStatisticLine("Consumed", consumed)
	printStatisticLine("Read", read)
	printStatisticLine("Whitelisted", whitelisted)
	printStatisticLine("Marshalled", marshalled)
	printStatisticLine("Checked", checked)
	printStatisticLine("Stored", stored)
}

func printConsumedEntry(entry AggregatorLogEntry) {
	fmt.Printf("%s %s %s %d\n", entry.Time, entry.Group, entry.Topic, entry.Offset)
}

func printErrorsForMessageWithOffset(entries []AggregatorLogEntry, offset int) {
	for _, entry := range entries {
		if entry.Offset == offset && entry.Level == "error" {
			fmt.Printf("\t%s %s\n", entry.Time, entry.Error)

		}
	}
}

func printConsumedEntries(entries []AggregatorLogEntry) {
	for _, entry := range entries {
		printConsumedEntry(entry)
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

func getConsumedNotReadMessages(entries []AggregatorLogEntry) []AggregatorLogEntry {
	consumed := filterConsumedMessages(entries)
	read := filterByMessage(entries, "Read")
	notRead := []AggregatorLogEntry{}

	for _, consumed := range consumed {
		if !messageWithOffsetIn(read, consumed.Offset) {
			notRead = append(notRead, consumed)
		}
	}
	return notRead
}

func printConsumedNotRead(entries []AggregatorLogEntry) {
	notRead := getConsumedNotReadMessages(entries)

	printConsumedEntries(notRead)
}

func analyse() {
	/*
		entries, err := readPipelineLogFile("pipeline2.log")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(entries)
	*/

	entries2, err := readAggregatorLogFile("aggregator3.log")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Read", len(entries2), "entries")
	printAggregatorStatistic(entries2)
	printConsumedNotRead(entries2)
}
