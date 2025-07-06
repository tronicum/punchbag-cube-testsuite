package sim

import (
	"encoding/json"
	"log"
	"net/http"
	sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"
)

// ValidationResult is a generic response for validation
type ValidationResult struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
}

// HandleValidation validates resource payloads for Azure, AWS, GCP
func HandleValidation(w http.ResponseWriter, r *http.Request) {
	resource := r.URL.Query().Get("resource")
	provider := r.URL.Query().Get("provider")
	var msg string
	var valid bool = true

	dec := json.NewDecoder(r.Body)
	switch {
	case provider == "azure" && resource == "aks":
		var c sharedmodels.Cluster
		if err := dec.Decode(&c); err != nil {
			log.Printf("[ERROR] %s %s: failed to decode AKS payload for validation: %v", r.Method, r.URL.Path, err)
			valid = false
			msg = "Invalid AKS payload: " + err.Error()
		} else if c.Name == "" || c.ResourceGroup == "" || c.Location == "" {
			valid = false
			msg = "Missing required AKS fields: name, resourceGroup, location"
		} else {
			msg = "AKS payload valid"
		}
	case provider == "azure" && resource == "loganalytics":
		var ws sharedmodels.LogAnalyticsWorkspace
		if err := dec.Decode(&ws); err != nil {
			log.Printf("[ERROR] %s %s: failed to decode Log Analytics payload for validation: %v", r.Method, r.URL.Path, err)
			valid = false
			msg = "Invalid Log Analytics payload: " + err.Error()
		} else if ws.Name == "" || ws.ResourceGroup == "" || ws.Location == "" {
			valid = false
			msg = "Missing required Log Analytics fields: name, resourceGroup, location"
		} else {
			msg = "Log Analytics payload valid"
		}
	case provider == "azure" && resource == "budget":
		var b sharedmodels.AzureBudget
		if err := dec.Decode(&b); err != nil {
			log.Printf("[ERROR] %s %s: failed to decode Budget payload for validation: %v", r.Method, r.URL.Path, err)
			valid = false
			msg = "Invalid Budget payload: " + err.Error()
		} else if b.Name == "" || b.ResourceGroup == "" || b.Amount <= 0 {
			valid = false
			msg = "Missing required Budget fields: name, resourceGroup, amount (>0)"
		} else {
			msg = "Budget payload valid"
		}
	case provider == "azure" && resource == "appinsights":
		var a sharedmodels.AppInsightsResource
		if err := dec.Decode(&a); err != nil {
			log.Printf("[ERROR] %s %s: failed to decode App Insights payload for validation: %v", r.Method, r.URL.Path, err)
			valid = false
			msg = "Invalid App Insights payload: " + err.Error()
		} else if a.Name == "" || a.ResourceGroup == "" || a.Location == "" {
			valid = false
			msg = "Missing required App Insights fields: name, resourceGroup, location"
		} else {
			msg = "App Insights payload valid"
		}
	case provider == "aws" && resource == "s3":
		var b sharedmodels.ObjectStorageBucket
		if err := dec.Decode(&b); err != nil {
			valid = false
			msg = "Invalid S3 payload: " + err.Error()
		} else if b.Name == "" || b.Region == "" {
			valid = false
			msg = "Missing required S3 fields: name, region"
		} else {
			msg = "S3 payload valid"
		}
	case provider == "azure" && resource == "blob":
		var b sharedmodels.ObjectStorageBucket
		if err := dec.Decode(&b); err != nil {
			valid = false
			msg = "Invalid Blob payload: " + err.Error()
		} else if b.Name == "" || b.Region == "" {
			valid = false
			msg = "Missing required Blob fields: name, region"
		} else {
			msg = "Blob payload valid"
		}
	case provider == "gcp" && resource == "gcs":
		var b sharedmodels.ObjectStorageBucket
		if err := dec.Decode(&b); err != nil {
			valid = false
			msg = "Invalid GCS payload: " + err.Error()
		} else if b.Name == "" || b.Region == "" {
			valid = false
			msg = "Missing required GCS fields: name, region"
		} else {
			msg = "GCS payload valid"
		}
	// Add AWS/GCP resource validation as needed
	default:
		valid = false
		msg = "Unknown provider/resource for validation"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ValidationResult{Valid: valid, Message: msg})
}
