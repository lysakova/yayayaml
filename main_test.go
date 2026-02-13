package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConvertYAMLFileToJSONFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		inputYAML   string
		expectErr   bool
		expectedOut string
	}{
		{
			name:        "empty yaml file",
			inputYAML:   "",
			expectedOut: "null",
		},
		{
			name:      "incorrectly formatted yaml",
			inputYAML: "name: test\nitems: [1, 2\n",
			expectErr: true,
		},
		{
			name:      "valid simple object yaml",
			inputYAML: "name: Alice\nage: 30\nactive: true\n",
			expectedOut: `{
  "active": true,
  "age": 30,
  "name": "Alice"
}`,
		},
		{
			name: "valid nested object and list yaml",
			inputYAML: "service:\n" +
				"  name: api\n" +
				"  ports:\n" +
				"    - 8080\n" +
				"    - 8443\n" +
				"  labels:\n" +
				"    env: prod\n" +
				"    team: platform\n",
			expectedOut: `{
  "service": {
    "labels": {
      "env": "prod",
      "team": "platform"
    },
    "name": "api",
    "ports": [
      8080,
      8443
    ]
  }
}`,
		},
		{
			name: "valid top-level array yaml",
			inputYAML: "- id: 1\n" +
				"  role: admin\n" +
				"- id: 2\n" +
				"  role: viewer\n",
			expectedOut: `[
  {
    "id": 1,
    "role": "admin"
  },
  {
    "id": 2,
    "role": "viewer"
  }
]`,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tmpDir := t.TempDir()
			inputPath := filepath.Join(tmpDir, "input.yaml")
			outputPath := filepath.Join(tmpDir, "output.json")

			if err := os.WriteFile(inputPath, []byte(tc.inputYAML), 0o644); err != nil {
				t.Fatalf("write input file: %v", err)
			}

			err := convertYAMLFileToJSONFile(inputPath, outputPath)
			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				if !strings.Contains(err.Error(), "converting YAML to JSON") {
					t.Fatalf("expected conversion error, got: %v", err)
				}
				if _, statErr := os.Stat(outputPath); !os.IsNotExist(statErr) {
					t.Fatalf("output file should not exist on error, stat err: %v", statErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			gotBytes, err := os.ReadFile(outputPath)
			if err != nil {
				t.Fatalf("read output file: %v", err)
			}

			got := string(gotBytes)
			if got != tc.expectedOut {
				t.Fatalf("output mismatch\nexpected:\n%s\ngot:\n%s", tc.expectedOut, got)
			}
		})
	}
}
