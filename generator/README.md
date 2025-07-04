# Generator

The `generator` component is responsible for generating code and templates based on predefined configurations. It currently supports generating Go code from Terraform templates for Azure services.

## Usage

1. Navigate to the `generator` directory:
   ```bash
   cd generator
   ```

2. Initialize Terraform (if not already done):
   ```bash
   terraform init
   ```

3. Run the code generator:
   ```bash
   go run main.go
   ```

The generated Go code will be available in the `generated_resources.go` file.

## Features

- Reads Terraform templates.
- Generates Go code for Azure resources.
- Provides a starting point for extending to other cloud providers.
