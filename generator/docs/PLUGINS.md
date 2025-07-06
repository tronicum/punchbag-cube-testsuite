# Plugin System Guide

This project supports a plugin system for resource generators. Plugins can be core (in `internal/generator/`) or external (in `plugins/`).

## Writing a Plugin

1. Implement the `ResourceGenerator` interface:
   ```go
   type MyPlugin struct{}
   func (p MyPlugin) Generate(props map[string]interface{}) string {
       // ...generate Terraform code...
   }
   ```
2. Register your plugin:
   ```go
   func init() {
       generator.RegisterPlugin("myresource", MyPlugin{})
   }
   ```
3. Place your plugin in `plugins/` or as a Go module.

## Using Plugins
- The orchestration layer will use a plugin if registered for a given resource type.
- You can dynamically add plugins for new resources/providers without changing core code.

## Example
See `plugins/example_plugin.go` for a scaffold.

## Future
- Dynamic loading of Go plugins (`.so` files) is planned for advanced extensibility.
