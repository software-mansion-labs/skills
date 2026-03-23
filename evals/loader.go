package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func loadEvalSuite(path string) (*EvalSuite, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	var suite EvalSuite
	if err := json.Unmarshal(data, &suite); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}

	return &suite, nil
}
