module github.com/tronicum/punchbag-cube-testsuite/werfty-transformator

go 1.24.4

require (
	github.com/tronicum/punchbag-cube-testsuite/shared v0.1.2
	github.com/tronicum/punchbag-cube-testsuite/werfty-transformator/transform v0.0.0-20250712064408-7f7611779cda
)

require (
	github.com/kr/pretty v0.3.1 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/tronicum/punchbag-cube-testsuite/werfty-transformator/transform => ./transform
