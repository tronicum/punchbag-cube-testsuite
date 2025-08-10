
#!/bin/bash
# End-to-end test for Hetzner S3 simulation via multitool and cube-server
set -euo pipefail

# Always run from the repo root for robust path handling
REPO_ROOT="$(git rev-parse --show-toplevel)"
cd "$REPO_ROOT"
echo "[E2E] Running from repo root: $(pwd)"

# Use a temp file for simulation persistence (set before server start)
SIM_PERSIST_FILE="/tmp/cube_server_sim_buckets_e2e_$$.json"
export CUBE_SERVER_SIM_PERSIST="$SIM_PERSIST_FILE"
echo "[E2E] Using simulation persistence file: $CUBE_SERVER_SIM_PERSIST"

# Build everything before running tests
echo "[E2E] Building all binaries..."
make build


# Start cube-server in simulate mode (background) using control script
echo "[E2E] Starting cube-server in simulate mode using control script..."
set +e
scripts/cube_server_control.sh start 18080 --debug
SERVER_START_EXIT=$?
set -e
if [ $SERVER_START_EXIT -ne 0 ]; then
  echo "[WARN] Server control script exited with code $SERVER_START_EXIT. Checking if server is running anyway..."
  if lsof -nP -iTCP:18080 | grep LISTEN | grep -E 'cube-serv(er)?' > /dev/null; then
    echo "[E2E] Server is running and listening on port 18080. Proceeding with tests."
  else
    echo "[E2E] Server is NOT running. Aborting test."
    exit 1
  fi
fi
echo "[E2E] Server startup step complete."



# Set dummy S3 credentials for simulation (if needed by multitool)
echo "[E2E] Setting dummy S3 credentials for simulation..."
export SIMULATE_DUMMY_S3_CREDS=1
echo "[E2E] Dummy S3 credentials set."

# Run multitool CLI check to verify server is responding
echo "[E2E] Running multitool CLI check (objectstorage list)..."
multitool/mt objectstorage list --storage-provider hetzner --mode simulate --server http://localhost:18080 --automation-mode &
MT_CHECK_PID=$!
wait $MT_CHECK_PID
MT_CHECK_EXIT=$?
if [ $MT_CHECK_EXIT -eq 0 ]; then
  echo "[E2E] Multitool CLI check succeeded. Server is responding."
else
  echo "[E2E] Multitool CLI check failed with exit code $MT_CHECK_EXIT. Aborting test."
  exit 1
fi
echo "[E2E] Multitool CLI check step complete."





# --- Cross-platform timeout wrapper ---
BUCKET_NAME="test-bucket-e2e-$(date +%s)"

if command -v gtimeout >/dev/null 2>&1; then
  TIMEOUT_CMD="gtimeout"
elif command -v timeout >/dev/null 2>&1; then
  TIMEOUT_CMD="timeout"
else
  echo "[E2E] Neither 'timeout' nor 'gtimeout' is available. Attempting to install coreutils using multitool..."
  multitool/mt local packages install -p coreutils -y
  if command -v gtimeout >/dev/null 2>&1; then
    TIMEOUT_CMD="gtimeout"
  elif command -v timeout >/dev/null 2>&1; then
    TIMEOUT_CMD="timeout"
  else
    echo "[E2E] ERROR: Failed to install coreutils or provide a timeout implementation. Aborting test."
    exit 1
  fi
fi

# Create a test bucket (with timeout)
echo "[E2E] Creating test bucket: $BUCKET_NAME ..."
$TIMEOUT_CMD 15s multitool/mt objectstorage create "$BUCKET_NAME" fsn1 --storage-provider hetzner --mode simulate --server http://localhost:18080 --automation-mode
CREATE_EXIT=$?
if [ $CREATE_EXIT -eq 0 ]; then
  echo "[E2E] Test bucket creation succeeded."
elif [ $CREATE_EXIT -eq 124 ]; then
  echo "[E2E] ERROR: Test bucket creation timed out after 15s. Aborting test."
  exit 1
else
  echo "[E2E] Test bucket creation failed with exit code $CREATE_EXIT. Aborting test."
  exit 1
fi

# List buckets (should show the created bucket, with timeout)
echo "[E2E] Listing buckets after creation (should show $BUCKET_NAME)..."
$TIMEOUT_CMD 15s multitool/mt objectstorage list --storage-provider hetzner --mode simulate --server http://localhost:18080 --automation-mode
LIST1_EXIT=$?
if [ $LIST1_EXIT -eq 0 ]; then
  echo "[E2E] Bucket list after creation succeeded."
elif [ $LIST1_EXIT -eq 124 ]; then
  echo "[E2E] ERROR: Bucket list after creation timed out after 15s. Aborting test."
  exit 1
else
  echo "[E2E] Bucket list after creation failed with exit code $LIST1_EXIT. Aborting test."
  exit 1
fi

# Delete the test bucket (with timeout)
echo "[E2E] Deleting test bucket: $BUCKET_NAME ..."
$TIMEOUT_CMD 15s multitool/mt objectstorage delete "$BUCKET_NAME" --storage-provider hetzner --mode simulate --server http://localhost:18080 --automation-mode --force
DELETE_EXIT=$?
if [ $DELETE_EXIT -eq 0 ]; then
  echo "[E2E] Test bucket deletion succeeded."
elif [ $DELETE_EXIT -eq 124 ]; then
  echo "[E2E] ERROR: Test bucket deletion timed out after 15s. Aborting test."
  exit 1
else
  echo "[E2E] Test bucket deletion failed with exit code $DELETE_EXIT. Aborting test."
  exit 1
fi

# List buckets again (should be empty, with timeout)
echo "[E2E] Listing buckets after deletion (should be empty)..."
$TIMEOUT_CMD 15s multitool/mt objectstorage list --storage-provider hetzner --mode simulate --server http://localhost:18080 --automation-mode
LIST2_EXIT=$?
if [ $LIST2_EXIT -eq 0 ]; then
  echo "[E2E] Bucket list after deletion succeeded."
elif [ $LIST2_EXIT -eq 124 ]; then
  echo "[E2E] ERROR: Bucket list after deletion timed out after 15s. Aborting test."
  exit 1
else
  echo "[E2E] Bucket list after deletion failed with exit code $LIST2_EXIT. Aborting test."
  exit 1
fi
echo "[E2E] Full multitool objectstorage flow complete."
echo "[E2E] Bucket list after deletion step complete."


# Stop cube-server using control script
echo "[E2E] Stopping cube-server..."
scripts/cube_server_control.sh stop 18080
echo "[E2E] Server stop step complete."

# Show persistence file contents for debug and clean up
if [ -f "$SIM_PERSIST_FILE" ]; then
  echo "[E2E] Persistence file contents after test:"
  cat "$SIM_PERSIST_FILE"
  rm -f "$SIM_PERSIST_FILE"
fi

echo "[E2E] Hetzner S3 simulation test completed successfully."
