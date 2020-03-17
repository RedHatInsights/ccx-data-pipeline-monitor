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
	Level   string `json:"level"`
	Time    string `json:"time"`
	Message string `json:"message"`
	Type    string `json:"type"`
	Error   string `json:"error"`
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

func analyse() {
	entries, err := readPipelineLogFile("pipeline2.log")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(entries)

	/*
			entries, err := readAggregatorLogFile("aggregator1.log")
			if err != nil {
				log.Fatal(err)
			}
		fmt.Println(entries)
	*/
}
