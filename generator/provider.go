package main

// Provider is an interface for cloud provider resource generation and validation
// This allows for plugin/config-driven codegen and easy extensibility

type Provider interface {
	GenerateTerraform(props map[string]interface{}, inputPath string) (string, string, error) // returns (tfHeader, tfBlock, error)
	ValidateResource(resourceType string, props map[string]interface{}) error
	DetectResourceType(props map[string]interface{}, inputPath string) (string, error)
}
