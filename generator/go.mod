module github.com/tronicum/punchbag-cube-testsuite/generator

go 1.23.0

require (
	github.com/tronicum/punchbag-cube-testsuite/shared v0.0.0-00010101000000-000000000000
	gopkg.in/yaml.v3 v3.0.1
)

replace github.com/tronicum/punchbag-cube-testsuite/shared => ../shared
