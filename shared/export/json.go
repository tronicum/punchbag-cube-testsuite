package export

import (
	"encoding/json"
	"io"
	"os"

	"github.com/tronicum/punchbag-cube-testsuite/shared/types"
)

// CloudState represents the exportable cloud state
type CloudState struct {
	Provider string              `json:"provider"`
	Clusters []*types.ClusterInfo `json:"clusters,omitempty"`
	Monitors []*types.MonitorInfo `json:"monitors,omitempty"`
	Budgets  []*types.BudgetInfo  `json:"budgets,omitempty"`
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

// ExportClusters exports cluster information to JSON
func ExportClusters(clusters []*types.ClusterInfo) ([]byte, error) {
	return json.MarshalIndent(clusters, "", "  ")
}

// ExportMonitors exports monitor information to JSON
func ExportMonitors(monitors []*types.MonitorInfo) ([]byte, error) {
	return json.MarshalIndent(monitors, "", "  ")
}

// ExportBudgets exports budget information to JSON
func ExportBudgets(budgets []*types.BudgetInfo) ([]byte, error) {
	return json.MarshalIndent(budgets, "", "  ")
}

// ExportAll exports all resource types to JSON
func ExportAll(clusters []*types.ClusterInfo, monitors []*types.MonitorInfo, budgets []*types.BudgetInfo) ([]byte, error) {
	data := map[string]interface{}{
		"clusters": clusters,
		"monitors": monitors,
		"budgets":  budgets,
	}
	return json.MarshalIndent(data, "", "  ")
}
