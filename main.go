package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"sigs.k8s.io/yaml"
)

func main() {
	const inputFile = "input.yaml"
	const outputFile = "output.json"

	if err := convertYAMLFileToJSONFile(inputFile, outputFile); err != nil {
		exitWithError("converting input.yaml to output.json", err)
	}

	fmt.Println("Converted input.yaml to output.json")
}

func convertYAMLFileToJSONFile(inputPath, outputPath string) error {
	yamlContent, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", inputPath, err)
	}

	jsonContent, err := yaml.YAMLToJSON(yamlContent)
	if err != nil {
		return fmt.Errorf("converting YAML to JSON: %w", err)
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, jsonContent, "", "  "); err != nil {
		return fmt.Errorf("formatting JSON: %w", err)
	}

	if err := os.WriteFile(outputPath, prettyJSON.Bytes(), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", outputPath, err)
	}

	return nil
}

func exitWithError(action string, err error) {
	fmt.Fprintf(os.Stderr, "Error %s: %v\n", action, err)
	os.Exit(1)
}
