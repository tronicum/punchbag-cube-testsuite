#!/bin/bash
# End-to-end test for server+generator+terraform: uses static multicloud test matrix for all providers
set -euo pipefail

WORKSPACE_ROOT="$(cd "$(dirname "$0")" && pwd)"
CUBE_SERVER_BIN="$WORKSPACE_ROOT/cube-server/cube-server"
CUBE_SERVER_LOG="$WORKSPACE_ROOT/cube-server.log"
SERVER_PORT=8080
SERVER_URL="http://localhost:$SERVER_PORT"
GENERATOR_BIN="$WORKSPACE_ROOT/generator/main.go"
TESTDATA_DIR="$WORKSPACE_ROOT/testdata"

# Start the cube-server if not running
if ! lsof -i :$SERVER_PORT | grep LISTEN >/dev/null 2>&1; then
  echo "[INFO] Starting cube-server in the background..."
  nohup "$CUBE_SERVER_BIN" > "$CUBE_SERVER_LOG" 2>&1 &
  SERVER_PID=$!
  # Wait for server to be ready
  for i in {1..20}; do
    if curl -s "$SERVER_URL/healthz" | grep -q 'ok'; then
      echo "[INFO] cube-server is up."
      break
    fi
    sleep 1
  done
else
  SERVER_PID=""
  echo "[INFO] cube-server already running."
fi

PROVIDERS=(azure aws gcp)
for provider in "${PROVIDERS[@]}"; do
  input_json="$TESTDATA_DIR/${provider}.json"
  output_tf="generator/test_${provider}.tf"
  go run "$GENERATOR_BIN" --generate-terraform --input "$input_json" --output "$output_tf" --provider "$provider"
  if [ ! -s "$output_tf" ]; then
    echo "[ERROR] Terraform output for $provider is empty or missing."
    exit 1
  fi
  echo "--- $provider Terraform ---"
  cat "$output_tf"
  # Optionally run terraform validate/plan for Azure only (others require cloud credentials)
  if [ "$provider" = "azure" ]; then
    cd examples
    terraform init -input=false
    terraform validate
    cd "$WORKSPACE_ROOT"
  fi
  rm -f "$output_tf"
done

# Cleanup: kill the server if we started it
if [ -n "${SERVER_PID:-}" ]; then
  echo "[INFO] Stopping cube-server (PID $SERVER_PID)..."
  kill $SERVER_PID || true
fi

echo "[SUCCESS] End-to-end multicloud test completed."
