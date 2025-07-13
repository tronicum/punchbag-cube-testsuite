#!/bin/bash
# end2end_multitool_terraform.sh
# End-to-end test: simulate all supported applications, generate Terraform, and print results

set -euo pipefail

WORKSPACE_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CUBE_SERVER_BIN="$WORKSPACE_ROOT/cube-server/cube-server"
CUBE_SERVER_LOG="$WORKSPACE_ROOT/cube-server.log"
SERVER_PORT=8080
SERVER_URL="http://localhost:$SERVER_PORT"
SIMULATOR_BIN="$WORKSPACE_ROOT/multitool/mt"
GENERATOR_BIN="$WORKSPACE_ROOT/generator/main.go"

# Parse --provider option (default: azure)
PROVIDER="azure"
while [[ $# -gt 0 ]]; do
  case $1 in
    --provider)
      PROVIDER="$2"
      shift 2
      ;;
    *)
      echo "Unknown option: $1" >&2
      exit 1
      ;;
  esac
done

echo "Using provider: $PROVIDER"

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

# List of applications to simulate (expand as needed)
APPS=(
  "aks $PROVIDER --resource-group test-rg --location eastus"
  "monitor $PROVIDER --resource-group test-rg --location eastus"
  "loganalytics $PROVIDER --resource-group test-rg --location eastus"
  "appinsights $PROVIDER --resource-group test-rg --location eastus"
)

for APP in "${APPS[@]}"; do
  APP_NAME=$(echo $APP | awk '{print $1}')
  APP_PROVIDER=$(echo $APP | awk '{print $2}')
  ARGS=$(echo $APP | cut -d' ' -f3-)
  JSON_FILE="$WORKSPACE_ROOT/multitool_sim_${APP_NAME}.json"
  TF_FILE="$WORKSPACE_ROOT/multitool_sim_${APP_NAME}.tf"

  echo "\n=== Simulating $APP_NAME ($APP_PROVIDER) ==="
  "$SIMULATOR_BIN" simulate cluster create test-$APP_NAME $APP_PROVIDER $ARGS -o json | awk '/^{/,/^}/' > "$JSON_FILE.raw" || {
    echo "Simulation failed for $APP_NAME"; continue;
  }
  # Provider-specific flattening
  if [ "$APP_PROVIDER" = "azure" ]; then
    if [ "$APP_NAME" = "aks" ]; then
      jq '{properties: {
        name: .name,
        location: .location,
        resourceGroup: .resource_group,
        nodeCount: .node_count,
        networkPlugin: .network_plugin,
        networkPolicy: .network_policy,
        identity: .identity,
        tags: .tags
      }}' "$JSON_FILE.raw" > "$JSON_FILE"
    else
      jq '{properties: (.properties // .)}' "$JSON_FILE.raw" > "$JSON_FILE"
    fi
  else
    # For other providers, just wrap the whole object for now (to be extended)
    jq '{properties: .}' "$JSON_FILE.raw" > "$JSON_FILE"
  fi
  rm "$JSON_FILE.raw"
  cat "$JSON_FILE"

  echo "Generating Terraform for $APP_NAME..."
  go run "$GENERATOR_BIN" --generate-terraform --input "$JSON_FILE" --output "$TF_FILE" || {
    echo "Terraform generation failed for $APP_NAME"; continue;
  }
  cat "$TF_FILE"
  echo "=== Done $APP_NAME ===\n"
done

# Clean up server if started by this script
if [ -n "$SERVER_PID" ]; then
  echo "Stopping cube-server (PID $SERVER_PID)..."
  kill "$SERVER_PID"
fi
