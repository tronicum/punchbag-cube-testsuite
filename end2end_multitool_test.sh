#!/bin/bash
# end2end_multitool_test.sh
# End-to-end test for multitool: downloads a simulated test set from the server with the simulator

set -euo pipefail

WORKSPACE_ROOT="$(cd "$(dirname "$0")" && pwd)"
CUBE_SERVER_BIN="$WORKSPACE_ROOT/cube-server/cube-server"
CUBE_SERVER_LOG="$WORKSPACE_ROOT/cube-server.log"
SERVER_PORT=8080
SERVER_URL="http://localhost:$SERVER_PORT"
SIMULATOR_BIN="$WORKSPACE_ROOT/multitool/multitool"

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

# Download a simulated test set using multitool (simulator)
SIM_TEST_JSON="$WORKSPACE_ROOT/multitool_sim_test.json"

"$SIMULATOR_BIN" simulate --server "$SERVER_URL" --resource-type aks --name test-aks --resource-group test-rg --location eastus --node-count 2 --output "$SIM_TEST_JSON"

if [ ! -f "$SIM_TEST_JSON" ]; then
  echo "Simulated test set was not downloaded."
  exit 1
fi

cat "$SIM_TEST_JSON"
echo "Simulated test set downloaded successfully."

# Clean up server if started by this script
if [ -n "$SERVER_PID" ]; then
  echo "Stopping cube-server (PID $SERVER_PID)..."
  kill "$SERVER_PID"
fi
