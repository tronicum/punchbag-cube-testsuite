package main

import "fmt"

// generateAksTerraformBlock generates the Terraform block for an AKS cluster from properties
func generateAksTerraformBlock(props map[string]interface{}) string {
	name := safeString(props, "name", "example-aks")
	location := safeString(props, "location", "eastus")
	resourceGroup := safeString(props, "resourceGroup", "example-rg")
	nodeCount := safeInt(props, "nodeCount", 3)
	networkPlugin := safeString(props, "networkPlugin", "azure")
	networkPolicy := safeString(props, "networkPolicy", "azure")
	dnsPrefix := safeString(props, "dnsPrefix", "exampleaks")
	identity := safeString(props, "identity", "")
	tags := safeString(props, "tags", "")
	// AKS availability zones (if applicable)
	var zones string
	if z, ok := props["availabilityZones"].([]interface{}); ok && len(z) > 0 {
		zones = "  availability_zones = ["
		for i, v := range z {
			if i > 0 {
				zones += ", "
			}
			zones += fmt.Sprintf("\"%s\"", v)
		}
		zones += "]\n"
	}
	// AKS node pool labels and taints
	labels := safeString(props, "nodePoolLabels", "")
	tagsLine := ""
	if tags != "" {
		tagsLine = fmt.Sprintf("  tags = %s\n", tags)
	}
	// TODO: Map additional AKS fields as needed
	// Example: Enable RBAC, API server authorized IP ranges, etc.
	enableRBAC := safeBool(props, "enableRBAC", true)
	rbacLine := ""
	if enableRBAC {
		rbacLine = "  role_based_access_control {\n    enabled = true\n  }\n"
	}
	apiServerAuthorizedIPRanges := safeString(props, "apiServerAuthorizedIPRanges", "")
	apiServerIPLine := ""
	if apiServerAuthorizedIPRanges != "" {
		apiServerIPLine = fmt.Sprintf("  api_server_authorized_ip_ranges = [%s]\n", apiServerAuthorizedIPRanges)
	}
	return fmt.Sprintf(`resource "azurerm_kubernetes_cluster" "example" {
  name                = "%s"
  location            = "%s"
  resource_group_name = "%s"
  default_node_pool {
    name       = "default"
    node_count = %d
  }
  identity {
    type = "%s"
  }
  sku {
    name     = "Standard_DS2_v2"
    tier     = "Standard"
    capacity = %d
  }
  network_profile {
    network_plugin = "%s"
    network_policy = "%s"
  }
  dns_prefix          = "%s"
%s%s%s%s%s  // ...map more fields from JSON as needed
}`,
		name, location, resourceGroup, nodeCount, identity, nodeCount, networkPlugin, networkPolicy, dnsPrefix, tagsLine, zones, labels, rbacLine, apiServerIPLine)
}
