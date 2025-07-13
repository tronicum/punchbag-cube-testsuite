#!/bin/bash
# end2end_multitool_test.sh
# End-to-end test for multitool: uses static multicloud test matrix for all providers
set -euo pipefail

trap 'echo "[ERROR] Script failed at line $LINENO"; exit 1' ERR

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


# Standard providers
PROVIDERS=(azure aws gcp)
for provider in "${PROVIDERS[@]}"; do
  SIM_TEST_JSON="$SCRIPT_DIR/sim_test_${provider}.json"
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
    *)
      echo "[ERROR] Unknown provider: $provider"; exit 1
      ;;
  esac
  echo "[INFO] Running simulation for $provider..."
  set +e
  "$SIMULATOR_BIN" $args > "$SIM_TEST_JSON"
  status=$?
  set -e
  if [ $status -ne 0 ] || [ ! -s "$SIM_TEST_JSON" ]; then
    echo "[ERROR] Simulation failed for $provider or output missing."
    exit 1
  fi
  echo "[SUCCESS] $provider simulation output:"
  cat "$SIM_TEST_JSON"
  rm -f "$SIM_TEST_JSON"
done

# Hetzner S3 bucket listing test
echo "[INFO] Running Hetzner S3 bucket listing test..."
HETZNER_TEST_JSON="$SCRIPT_DIR/hetzner_s3_list.json"
export HETZNER_S3_ACCESS_KEY="${HETZNER_S3_ACCESS_KEY:-dummy-access-key}"
export HETZNER_S3_SECRET_KEY="${HETZNER_S3_SECRET_KEY:-dummy-secret-key}"
set +e
"$SIMULATOR_BIN" objectstorage list hetzner > "$HETZNER_TEST_JSON"
status=$?
set -e
if [ $status -ne 0 ]; then
  echo "[ERROR] Hetzner S3 bucket listing failed."
  exit 1
fi
echo "[SUCCESS] Hetzner S3 bucket listing output:"
cat "$HETZNER_TEST_JSON"
rm -f "$HETZNER_TEST_JSON"

if [ -n "${SERVER_PID:-}" ]; then
  kill "$SERVER_PID" || true
fi

echo "[SUCCESS] Multitool end-to-end multicloud test completed."
