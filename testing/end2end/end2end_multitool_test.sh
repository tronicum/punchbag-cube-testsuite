#!/bin/bash
# end2end_multitool_test.sh
# End-to-end test for multitool: uses static multicloud test matrix for all providers
set -euo pipefail

trap 'echo "[ERROR] Script failed at line $LINENO"; exit 1' ERR

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
WORKSPACE_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
CUBE_SERVER_BIN="$WORKSPACE_ROOT/cube-server/cube-server"
CUBE_SERVER_LOG="$WORKSPACE_ROOT/cube-server.log"
SERVER_PORT=8080
SERVER_URL="http://localhost:$SERVER_PORT"
SIMULATOR_BIN="$WORKSPACE_ROOT/multitool/mt"
TESTDATA_DIR="$WORKSPACE_ROOT/testdata"


# Only start the cube-server if required (for non-direct multitool tests)
REQUIRE_CUBE_SERVER=${REQUIRE_CUBE_SERVER:-0}
if [[ "$REQUIRE_CUBE_SERVER" == "1" ]]; then
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
else
  echo "[INFO] Skipping cube-server startup (direct multitool mode)"
  SERVER_PID=""
fi

# Standard providers
PROVIDERS=(azure aws gcp)
for provider in "${PROVIDERS[@]}"; do
  SIM_TEST_JSON="$SCRIPT_DIR/sim_test_${provider}.json"
  echo "[INFO] Running simulation for $provider..."
  set +e
  "$SIMULATOR_BIN" $provider > "$SIM_TEST_JSON"
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



# Flexible object storage bucket tests (provider-agnostic)
# Accept providers as a script argument, or via OBJECT_STORAGE_PROVIDERS env var, or default to hetzner
CLEANUP=0
POSITIONAL=()
for arg in "$@"; do
  case $arg in
    --cleanup)
      CLEANUP=1
      ;;
    *)
      POSITIONAL+=("$arg")
      ;;
  esac
done
set -- "${POSITIONAL[@]}"
if [[ $# -ge 1 ]]; then
  OBJECT_STORAGE_PROVIDERS=($@)
elif [[ -n "${OBJECT_STORAGE_PROVIDERS:-}" ]]; then
  OBJECT_STORAGE_PROVIDERS=(${OBJECT_STORAGE_PROVIDERS})
else
  OBJECT_STORAGE_PROVIDERS=(hetzner)
fi
# Allow skipping object storage tests
if [[ "${SKIP_OBJECT_STORAGE_TESTS:-}" == "1" ]]; then
  echo "[INFO] Skipping object storage bucket tests (SKIP_OBJECT_STORAGE_TESTS=1)"
else
  for obj_provider in "${OBJECT_STORAGE_PROVIDERS[@]}"; do
    if [[ "$CLEANUP" == "1" ]]; then
      echo "[CLEANUP] Removing all e2e-test-bucket-* buckets for $obj_provider..."
      OBJ_TEST_JSON="$SCRIPT_DIR/${obj_provider}_s3_list.json"
      "$SIMULATOR_BIN" objectstorage list "$obj_provider" --output json > "$OBJ_TEST_JSON"
      JSON_START_LINE=$(grep -n '^\[' "$OBJ_TEST_JSON" | head -n1 | cut -d: -f1)
      JSON_END_LINE=$(awk "NR>=$JSON_START_LINE && /^]/ {print NR; exit}" "$OBJ_TEST_JSON")
      sed -n "${JSON_START_LINE},${JSON_END_LINE}p" "$OBJ_TEST_JSON" > "$OBJ_TEST_JSON.filtered"
      deleted_buckets=()
      failed_buckets=()
      for bucket_name in $(cat "$OBJ_TEST_JSON.filtered" | jq -r '.[] | select(.name | startswith("e2e-test-bucket-")) | .name'); do
        echo "> $SIMULATOR_BIN objectstorage delete $obj_provider $bucket_name --force"
        set +e
        "$SIMULATOR_BIN" objectstorage delete "$obj_provider" "$bucket_name" --force > "$OBJ_TEST_JSON.delete" 2>&1
        status=$?
        set -e
        cat "$OBJ_TEST_JSON.delete"
        if [ $status -ne 0 ]; then
          echo "[ERROR] $obj_provider S3 bucket deletion failed for $bucket_name."
          failed_buckets+=("$bucket_name")
        else
          echo "[SUCCESS] $obj_provider S3 bucket deleted: $bucket_name"
          deleted_buckets+=("$bucket_name")
        fi
        rm -f "$OBJ_TEST_JSON.delete"
        sleep 1
      done
      echo "[SUMMARY] Deleted buckets for $obj_provider: ${deleted_buckets[*]}"
      if [ ${#failed_buckets[@]} -gt 0 ]; then
        echo "[SUMMARY] Failed to delete buckets for $obj_provider: ${failed_buckets[*]}"
        exit 1
      fi
      rm -f "$OBJ_TEST_JSON" "$OBJ_TEST_JSON.filtered"
      continue
    fi
    echo "[INFO] Running object storage bucket tests for $obj_provider..."
    OBJ_TEST_JSON="$SCRIPT_DIR/${obj_provider}_s3_list.json"
    TEST_BUCKET_NAME="e2e-test-bucket-$(date +%s)-$RANDOM"
    # Set up environment for each provider if needed
    case "$obj_provider" in
      aws)
        export AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID:-dummy-access-key}"
        export AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY:-dummy-secret-key}"
        export AWS_REGION="${AWS_REGION:-us-east-1}"
        ;;
      hetzner)
        export HETZNER_S3_ACCESS_KEY="${HETZNER_S3_ACCESS_KEY:-dummy-access-key}"
        export HETZNER_S3_SECRET_KEY="${HETZNER_S3_SECRET_KEY:-dummy-secret-key}"
        export HETZNER_S3_REGION="${HETZNER_S3_REGION:-fsn1}"
        ;;
      # Add more providers here as needed
    esac
    # Create bucket
    echo "[INFO] Creating test bucket $TEST_BUCKET_NAME for $obj_provider..."
    set +e
    "$SIMULATOR_BIN" objectstorage create "$obj_provider" "$TEST_BUCKET_NAME" "${HETZNER_S3_REGION:-us-east-1}" > "$OBJ_TEST_JSON.create" 2>&1
    status=$?
    set -e
    if [ $status -ne 0 ]; then
      echo "[ERROR] $obj_provider S3 bucket creation failed."
      cat "$OBJ_TEST_JSON.create"
      exit 1
    fi
    echo "[SUCCESS] $obj_provider S3 bucket created: $TEST_BUCKET_NAME"
    # List buckets
    set +e
    "$SIMULATOR_BIN" objectstorage list "$obj_provider" --output json > "$OBJ_TEST_JSON"
    status=$?
    set -e
    if [ $status -ne 0 ]; then
      echo "[ERROR] $obj_provider S3 bucket listing failed."
      cat "$OBJ_TEST_JSON"
      exit 1
    fi
    echo "[SUCCESS] $obj_provider S3 bucket listing raw output:"
    cat "$OBJ_TEST_JSON"
    # Extract only the JSON array from the output (ignore log lines before and after)
    JSON_START_LINE=$(grep -n '^\[' "$OBJ_TEST_JSON" | head -n1 | cut -d: -f1)
    if [ -z "$JSON_START_LINE" ]; then
      echo "[ERROR] $obj_provider S3 bucket listing output does not contain a JSON array."
      cat "$OBJ_TEST_JSON"
      exit 1
    fi
    # Find the first matching closing bracket after JSON_START_LINE
    JSON_END_LINE=$(awk "NR>=$JSON_START_LINE && /^]/ {print NR; exit}" "$OBJ_TEST_JSON")
    if [ -z "$JSON_END_LINE" ]; then
      echo "[ERROR] $obj_provider S3 bucket listing output does not contain a closing bracket for the JSON array."
      cat "$OBJ_TEST_JSON"
      exit 1
    fi
    sed -n "${JSON_START_LINE},${JSON_END_LINE}p" "$OBJ_TEST_JSON" > "$OBJ_TEST_JSON.filtered"
    # Check if filtered output is valid JSON and not empty
    if [ ! -s "$OBJ_TEST_JSON.filtered" ]; then
      echo "[WARNING] $obj_provider S3 bucket listing output is empty after filtering."
    elif jq empty "$OBJ_TEST_JSON.filtered" 2>/dev/null; then
      echo "[SUCCESS] $obj_provider S3 bucket listing parsed output:"
      cat "$OBJ_TEST_JSON.filtered" | jq '.[] | {name, region, provider}'
      # Delete e2e-test buckets
      deleted_buckets=()
      failed_buckets=()
      for bucket_name in $(cat "$OBJ_TEST_JSON.filtered" | jq -r '.[] | select(.name | startswith("e2e-test-bucket-")) | .name'); do
        echo "[INFO] Deleting test bucket $bucket_name for $obj_provider..."
        echo "> $SIMULATOR_BIN objectstorage delete $obj_provider $bucket_name --force"
        set +e
        "$SIMULATOR_BIN" objectstorage delete "$obj_provider" "$bucket_name" --force > "$OBJ_TEST_JSON.delete" 2>&1
        status=$?
        set -e
        cat "$OBJ_TEST_JSON.delete"
        if [ $status -ne 0 ]; then
          echo "[ERROR] $obj_provider S3 bucket deletion failed for $bucket_name."
          failed_buckets+=("$bucket_name")
        else
          echo "[SUCCESS] $obj_provider S3 bucket deleted: $bucket_name"
          deleted_buckets+=("$bucket_name")
        fi
        rm -f "$OBJ_TEST_JSON.delete"
        sleep 1
      done
      echo "[SUMMARY] Deleted buckets for $obj_provider: ${deleted_buckets[*]}"
      if [ ${#failed_buckets[@]} -gt 0 ]; then
        echo "[SUMMARY] Failed to delete buckets for $obj_provider: ${failed_buckets[*]}"
        exit 1
      fi
    else
      echo "[ERROR] $obj_provider S3 bucket listing filtered output is not valid JSON."
      cat "$OBJ_TEST_JSON.filtered"
      exit 1
    fi
    rm -f "$OBJ_TEST_JSON" "$OBJ_TEST_JSON.create" "$OBJ_TEST_JSON.filtered"
  done
fi


# Only kill the cube-server if it was started by this script
if [[ "$REQUIRE_CUBE_SERVER" == "1" && -n "${SERVER_PID:-}" ]]; then
  kill "$SERVER_PID" || true
fi

echo "[SUCCESS] Multitool end-to-end multicloud test completed."
