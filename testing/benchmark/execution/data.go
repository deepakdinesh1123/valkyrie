package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type EnvironmentVariable struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type ExecutionEnvironmentSpec struct {
	EnvironmentVariables []EnvironmentVariable `json:"environment_variables,omitempty"`
	LanguageDependencies []string              `json:"languageDependencies,omitempty"`
	SystemDependencies   []string              `json:"systemDependencies,omitempty"`
	Args                 string                `json:"args,omitempty"`
}

type ExecutionRequest struct {
	Environment ExecutionEnvironmentSpec `json:"environment,omitempty"`
	Code        string                   `json:"code,omitempty"`
	Language    string                   `json:"language,omitempty"`
	MaxRetries  int                      `json:"max_retries,omitempty"`
	Timeout     int32                    `json:"timeout,omitempty"`
}

func main() {
	var data []ExecutionRequest
	var files []string

	root := "testing/benchmark/execution/fixtures"

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("error walking the path %v: %v\n", root, err)
	}
	for _, file := range files {
		var execRequest ExecutionRequest
		fileContent, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("error reading file %v: %v\n", file, err)
		}
		var execRequests []ExecutionRequest
		err = json.Unmarshal(fileContent, &execRequests)
		if err != nil {
			log.Fatalf("error unmarshalling file %v: %v\n", file, err)
		}
		data = append(execRequests, execRequest)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("error marshalling data: %v\n", err)
	}
	err = os.WriteFile("testing/benchmark/execution/load-test-data.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("error writing file: %v\n", err)
	}
	fmt.Println("Data generated successfully")
}
