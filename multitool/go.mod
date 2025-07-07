module github.com/tronicum/punchbag-cube-testsuite/multitool

go 1.23

require (
	github.com/spf13/cobra v1.9.1
	github.com/tronicum/punchbag-cube-testsuite/shared v0.0.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
)

replace github.com/tronicum/punchbag-cube-testsuite/shared => ../shared
