package main

import (
	"flag"
	"fmt"
	"os"
	"text/template"

	sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"
)

type S3Bucket struct {
	Name     string
	Region   string
	Provider sharedmodels.CloudProvider
}

type AzureBlob struct {
	Name     string
	Location string
	Provider sharedmodels.CloudProvider
}

type GCSBucket struct {
	Name     string
	Location string
	Provider sharedmodels.CloudProvider
}

type StackITObjectStorage struct {
	Name     string
	Region   string
	Provider sharedmodels.CloudProvider
}

type HCloudObjectStorage struct {
	Name     string
	Location string
	Provider sharedmodels.CloudProvider
}

type IonosS3Bucket struct {
	Name     string
	Region   string
	Provider sharedmodels.CloudProvider
}

const s3BucketTemplate = `
resource "aws_s3_bucket" "example" {
  bucket = "{{ .Name }}"
  region = "{{ .Region }}"
  # provider = "{{ .Provider }}"
}
`

const azureBlobTemplate = `
resource "azurerm_storage_account" "example" {
  name                     = "{{ .Name }}"
  location                 = "{{ .Location }}"
  resource_group_name      = "example-rg"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  # provider = "{{ .Provider }}"
}
`

const gcsBucketTemplate = `
resource "google_storage_bucket" "example" {
  name     = "{{ .Name }}"
  location = "{{ .Location }}"
  # provider = "{{ .Provider }}"
}
`

const stackitObjectStorageTemplate = `
resource "stackit_object_storage_bucket" "example" {
  name   = "{{ .Name }}"
  region = "{{ .Region }}"
  # provider = "{{ .Provider }}"
}
`

const hcloudObjectStorageTemplate = `
resource "hcloud_object_storage" "example" {
  name     = "{{ .Name }}"
  location = "{{ .Location }}"
  # provider = "{{ .Provider }}"
}
`

const ionosS3BucketTemplate = `
resource "ionoscloud_s3_bucket" "example" {
  name   = "{{ .Name }}"
  region = "{{ .Region }}"
  # provider = "{{ .Provider }}"
}
`

func main() {
	cloudTarget := flag.String("cloud-provider-target", "aws", "Target cloud provider (stackit|aws|azure|gcp|hcloud|ionos)")
	service := flag.String("cloud-provider-service", "s3", "Cloud provider service to generate (s3|azureblob|gcs|objectstorage)")
	flag.Parse()

	var provider sharedmodels.CloudProvider
	switch *cloudTarget {
	case "aws":
		provider = sharedmodels.AWS
	case "azure":
		provider = sharedmodels.Azure
	case "gcp":
		provider = sharedmodels.GCP
	case "hcloud":
		provider = sharedmodels.Hetzner
	case "ionos":
		provider = sharedmodels.IONOS
	case "stackit":
		provider = sharedmodels.StackIT
	default:
		fmt.Println("Unknown cloud provider. Use --cloud-provider-target=stackit|aws|azure|gcp|hcloud|ionos")
		os.Exit(1)
	}

	switch *service {
	case "s3":
		bucket := S3Bucket{
			Name:     "my-bucket",
			Region:   "us-west-2",
			Provider: provider,
		}
		tmpl, err := template.New("s3").Parse(s3BucketTemplate)
		if err != nil {
			panic(err)
		}
		f, err := os.Create("generated_s3_bucket.tf")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		if err := tmpl.Execute(f, bucket); err != nil {
			panic(err)
		}
		fmt.Println("Terraform file generated: generated_s3_bucket.tf")
	case "azureblob":
		blob := AzureBlob{
			Name:     "mystorageacct",
			Location: "westeurope",
			Provider: provider,
		}
		tmpl, err := template.New("azureblob").Parse(azureBlobTemplate)
		if err != nil {
			panic(err)
		}
		f, err := os.Create("generated_azure_blob.tf")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		if err := tmpl.Execute(f, blob); err != nil {
			panic(err)
		}
		fmt.Println("Terraform file generated: generated_azure_blob.tf")
	case "gcs":
		gcs := GCSBucket{
			Name:     "my-gcs-bucket",
			Location: "EU",
			Provider: provider,
		}
		tmpl, err := template.New("gcs").Parse(gcsBucketTemplate)
		if err != nil {
			panic(err)
		}
		f, err := os.Create("generated_gcs_bucket.tf")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		if err := tmpl.Execute(f, gcs); err != nil {
			panic(err)
		}
		fmt.Println("Terraform file generated: generated_gcs_bucket.tf")
	case "objectstorage":
		switch *cloudTarget {
		case "stackit":
			stackit := StackITObjectStorage{
				Name:     "my-stackit-bucket",
				Region:   "eu01",
				Provider: provider,
			}
			tmpl, err := template.New("stackit").Parse(stackitObjectStorageTemplate)
			if err != nil {
				panic(err)
			}
			f, err := os.Create("generated_stackit_object_storage.tf")
			if err != nil {
				panic(err)
			}
			defer f.Close()
			if err := tmpl.Execute(f, stackit); err != nil {
				panic(err)
			}
			fmt.Println("Terraform file generated: generated_stackit_object_storage.tf")
		case "hcloud":
			hcloud := HCloudObjectStorage{
				Name:     "my-hcloud-bucket",
				Location: "fsn1",
				Provider: provider,
			}
			tmpl, err := template.New("hcloud").Parse(hcloudObjectStorageTemplate)
			if err != nil {
				panic(err)
			}
			f, err := os.Create("generated_hcloud_object_storage.tf")
			if err != nil {
				panic(err)
			}
			defer f.Close()
			if err := tmpl.Execute(f, hcloud); err != nil {
				panic(err)
			}
			fmt.Println("Terraform file generated: generated_hcloud_object_storage.tf")
		case "ionos":
			ionos := IonosS3Bucket{
				Name:     "my-ionos-bucket",
				Region:   "de",
				Provider: provider,
			}
			tmpl, err := template.New("ionos").Parse(ionosS3BucketTemplate)
			if err != nil {
				panic(err)
			}
			f, err := os.Create("generated_ionos_s3_bucket.tf")
			if err != nil {
				panic(err)
			}
			defer f.Close()
			if err := tmpl.Execute(f, ionos); err != nil {
				panic(err)
			}
			fmt.Println("Terraform file generated: generated_ionos_s3_bucket.tf")
		default:
			fmt.Println("Unknown object storage provider for --cloud-provider-target. Use stackit, hcloud, or ionos.")
		}
	default:
		fmt.Println("Unknown service. Use --cloud-provider-service=s3|azureblob|gcs|objectstorage")
	}
}
