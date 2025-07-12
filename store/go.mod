module github.com/tronicum/punchbag-cube-testsuite/store

go 1.24

toolchain go1.24.4

replace github.com/tronicum/punchbag-cube-testsuite/shared => ../shared

require (
	github.com/google/uuid v1.6.0
	github.com/tronicum/punchbag-cube-testsuite/shared v0.0.0
)
