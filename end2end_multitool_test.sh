#!/bin/bash
# end2end_multitool_test.sh
# End-to-end test for multitool: uses static multicloud test matrix for all providers
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
WORKSPACE_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
CUBE_SERVER_BIN="$WORKSPACE_ROOT/cube-server/cube-server"
CUBE_SERVER_LOG="$WORKSPACE_ROOT/cube-server.log"
SERVER_PORT=8080
SERVER_URL="http://localhost:$SERVER_PORT"
SIMULATOR_BIN="$SCRIPT_DIR/multitool"
TESTDATA_DIR="$WORKSPACE_ROOT/testdata"

# Start the cube-server if not running
if ! lsof -i :$SERVER_PORT | grep LISTEN >/dev/null 2>&1; then
  echo "Starting cube-server..."
  nohup "$CUBE_SERVER_BIN" > "$CUBE_SERVER_LOG" 2>&1 &
  SERVER_PID=$!
  # Wait for /healthz
  for i in {1..20}; do
    if curl -s "$SERVER_URL/healthz" | grep -q 'ok'; then
      echo "cube-server is up."
      break
    fi
    sleep 1
  done
else
  echo "cube-server already running."
  SERVER_PID=""
fi

PROVIDERS=(azure aws gcp)
for provider in "${PROVIDERS[@]}"; do
  SIM_TEST_JSON="$SCRIPT_DIR/sim_test_${provider}.json"
  # Example args for each provider (customize as needed)
  case "$provider" in
    azure)
      args="--resource-type aks --name test-aks --resource-group test-rg --location eastus --node-count 2"
      ;;
    aws)
      args="--resource-type eks --name test-eks --location us-west-2 --node-count 2"
      ;;
    gcp)
      args="--resource-type gke --name test-gke --location us-central1 --node-count 1"
      ;;
  esac
  eval "$SIMULATOR_BIN simulate --server $SERVER_URL $args --output $SIM_TEST_JSON"
  if [ ! -s "$SIM_TEST_JSON" ]; then
    echo "[ERROR] Simulation output for $provider is empty or missing."
    exit 1
  fi
  echo "--- $provider Simulation Output ---"
  cat "$SIM_TEST_JSON"
  rm -f "$SIM_TEST_JSON"
done

echo "[SUCCESS] Multitool end-to-end multicloud test completed."

# Clean up server if started by this script
if [ -n "$SERVER_PID" ]; then
  echo "Stopping cube-server (PID $SERVER_PID)..."
  kill "$SERVER_PID"
fi
