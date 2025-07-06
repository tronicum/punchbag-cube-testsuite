#!/bin/bash
# End-to-end test script for AKS simulation and Terraform generation
# Usage: ./end2end_test.sh
set -euo pipefail

# Always run from the workspace root
cd "$(dirname "$0")"

GENERATOR=generator/main.go
CUBE_SERVER=localhost:8080
CUBE_SERVER_BIN=main.go
TEST_JSON=generator/test_aks_expanded.json
TEST_TF=generator/test_aks_expanded.tf
LOG_FILE=cube-server.log

# 1. Start the cube-server in the background if not running
echo "[INFO] Checking if cube-server is running..."
if ! curl -s http://$CUBE_SERVER/healthz | grep -q 'ok'; then
  echo "[INFO] Starting cube-server in the background..."
  cd cube-server
  nohup go run $CUBE_SERVER_BIN > ../$LOG_FILE 2>&1 &
  SERVER_PID=$!
  cd ..
  # Wait for server to be ready
  for i in {1..10}; do
    sleep 1
    if curl -s http://$CUBE_SERVER/healthz | grep -q 'ok'; then
      echo "[INFO] cube-server is up."
      break
    fi
    if [ $i -eq 10 ]; then
      echo "[ERROR] cube-server did not start in time. Log output:" && tail -40 $LOG_FILE
      exit 1
    fi
  done
else
  SERVER_PID=""
  echo "[INFO] cube-server already running."
fi

echo "[INFO] Simulating AKS resource import (mock JSON)..."
mkdir -p generator
go run $GENERATOR --simulate-import --resource-type aks --name e2e-aks --resource-group e2e-rg --location eastus --node-count 2 > $TEST_JSON
cat $TEST_JSON

echo "[INFO] Generating Terraform from simulated JSON..."
go run $GENERATOR --generate-terraform --input $TEST_JSON --output $TEST_TF
cat $TEST_TF

cd examples

echo "[INFO] Initializing Terraform..."
terraform init -input=false

echo "[INFO] Validating Terraform..."
terraform validate

echo "[INFO] Running Terraform plan (should fail with dummy credentials)..."
set +e
terraform plan
PLAN_EXIT=$?
set -e
if [ $PLAN_EXIT -eq 0 ]; then
  echo "[WARNING] Terraform plan succeeded (unexpected with dummy credentials)."
else
  echo "[INFO] Terraform plan failed as expected with dummy credentials."
fi

# Cleanup: kill the server if we started it
if [ -n "${SERVER_PID:-}" ]; then
  echo "[INFO] Stopping cube-server (PID $SERVER_PID)..."
  kill $SERVER_PID || true
fi

echo "[SUCCESS] End-to-end test completed."
