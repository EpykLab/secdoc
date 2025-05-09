package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/EpykLab/secdoc/internal"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: secparser <source-dir> [output-file]")
		os.Exit(1)
	}

	sourcePath := os.Args[1]
	outputPath := "security-report.json"
	if len(os.Args) > 2 {
		outputPath = os.Args[2]
	}

	report, err := internal.ParseDirectory(sourcePath)
	if err != nil {
		fmt.Printf("Error parsing directory: %v\n", err)
		os.Exit(1)
	}

	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(outputPath, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Report written to %s\n", outputPath)
}
