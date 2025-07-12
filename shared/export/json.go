package export

import (
	"encoding/json"
	"io"
	"os"

	"punchbag-cube-testsuite/shared/providers"
)

// CloudState represents the exportable cloud state
type CloudState struct {
	Provider string                   `json:"provider"`
	Clusters []*providers.ClusterInfo `json:"clusters,omitempty"`
	Monitors []*providers.MonitorInfo `json:"monitors,omitempty"`
	Budgets  []*providers.BudgetInfo  `json:"budgets,omitempty"`
}

// ToJSON exports cloud state to JSON
func ToJSON(state *CloudState, writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(state)
}

// ToJSONFile exports cloud state to JSON file
func ToJSONFile(state *CloudState, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	return ToJSON(state, file)
}

// FromJSON imports cloud state from JSON
func FromJSON(reader io.Reader) (*CloudState, error) {
	var state CloudState
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&state)
	return &state, err
}

// FromJSONFile imports cloud state from JSON file
func FromJSONFile(filename string) (*CloudState, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	return FromJSON(file)
}

// ExportToJSON exports data to JSON writer
func ExportToJSON(data interface{}, w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// ExportToJSONFile exports data to a JSON file
func ExportToJSONFile(data interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	return ExportToJSON(data, file)
}
