module github.com/tronicum/punchbag-cube-testsuite/werfty-transformator

go 1.24.4

require (
   github.com/tronicum/punchbag-cube-testsuite/shared v0.1.0
   github.com/tronicum/punchbag-cube-testsuite/werfty-transformator/transform v0.0.0-20250712064408-7f7611779cda
)

replace github.com/tronicum/punchbag-cube-testsuite/shared => ../../shared
replace github.com/tronicum/punchbag-cube-testsuite/werfty-transformator/transform => ./transform
