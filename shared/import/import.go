package importpkg

import (
	"encoding/json"
	"fmt"
	"os"
)

// ImportCloudState loads cloud state from a JSON or YAML file
func ImportCloudState(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	var state map[string]interface{}
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return state, nil
}
