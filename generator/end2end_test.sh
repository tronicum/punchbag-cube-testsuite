#!/bin/bash
# generator/end2end_test.sh
# End-to-end test for generator: uses static multicloud test matrix for all providers
set -euo pipefail

trap 'echo "[ERROR] Script failed at line $LINENO"; exit 1' ERR

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
WORKSPACE_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
GENERATOR_BIN="$SCRIPT_DIR/main.go"
TESTDATA_DIR="$WORKSPACE_ROOT/testdata"

cd "$SCRIPT_DIR"

PROVIDERS=(azure aws gcp)
for provider in "${PROVIDERS[@]}"; do
  input_json="$TESTDATA_DIR/${provider}.json"
  output_tf="test_${provider}.tf"
  echo "[INFO] Generating Terraform for $provider..."
  go run "$GENERATOR_BIN" --generate-terraform --input "$input_json" --output "$output_tf" --provider "$provider"
  if [ ! -s "$output_tf" ]; then
    echo "[ERROR] Terraform output for $provider is empty or missing."
    exit 1
  fi
  echo "--- $provider Terraform ---"
  cat "$output_tf"
  rm -f "$output_tf"
done

echo "[SUCCESS] Generator end-to-end multicloud test completed."
