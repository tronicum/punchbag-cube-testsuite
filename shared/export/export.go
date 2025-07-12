package export

import (
	"encoding/json"
	"io"
	"os"
)

// CloudState represents exportable cloud state data
type CloudState struct {
	Provider  string                 `json:"provider"`
	Resources []map[string]interface{} `json:"resources"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// ToJSON exports cloud state to a JSON writer
func ToJSON(state *CloudState, w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(state)
}

// ToJSONFile exports cloud state to a JSON file
func ToJSONFile(state *CloudState, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	
	return ToJSON(state, f)
}
		return err
	}
	defer f.Close()
	
	return ToJSON(state, f)
}

// ToYAML exports cloud state to YAML format
func ToYAML(state *CloudState, w io.Writer) error {
	encoder := yaml.NewEncoder(w)
	encoder.SetIndent(2)
	return encoder.Encode(state)
}

// ToYAMLFile exports cloud state to a YAML file
func ToYAMLFile(state *CloudState, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	
	return ToYAML(state, f)
}

// FromJSON imports cloud state from JSON format
func FromJSON(r io.Reader) (*CloudState, error) {
	var state CloudState
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&state)
	if err != nil {
		return nil, err
	}
	return &state, nil
}

// FromJSONFile imports cloud state from a JSON file
func FromJSONFile(filename string) (*CloudState, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	
	return FromJSON(f)
}
