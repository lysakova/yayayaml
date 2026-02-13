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

	yamlContent, err := os.ReadFile(inputFile)
	if err != nil {
		exitWithError("reading input.yaml", err)
	}

	jsonContent, err := yaml.YAMLToJSON(yamlContent)
	if err != nil {
		exitWithError("converting YAML to JSON", err)
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, jsonContent, "", "  "); err != nil {
		exitWithError("formatting JSON", err)
	}

	if err := os.WriteFile(outputFile, prettyJSON.Bytes(), 0o644); err != nil {
		exitWithError("writing output.json", err)
	}

	fmt.Println("Converted input.yaml to output.json")
}

func exitWithError(action string, err error) {
	fmt.Fprintf(os.Stderr, "Error %s: %v\n", action, err)
	os.Exit(1)
}
