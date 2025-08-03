# End-to-end test for Hetzner S3 simulation via multitool and cube-server
# Set -euo pipefail to ensure the script exits on errors and uninitialized variables
set -euo pipefail
#!/bin/bash
# End-to-end test for Hetzner S3 simulation via multitool and cube-server
set -euo pipefail

# Build everything before running tests
echo "[E2E] Building all binaries..."
(cd "$(git rev-parse --show-toplevel)" && make build)

# Start cube-server in simulate mode (background) using control script
echo "[E2E] Starting cube-server in simulate mode using control script..."
"$(git rev-parse --show-toplevel)/cube-server/scripts/cube_server_control.sh" start 18080 --debug
sleep 2

# Set dummy S3 credentials for simulation
echo "[E2E] Setting dummy S3 credentials for simulation..."
export SIMULATE_DUMMY_S3_CREDS=1

# Create a test bucket first
echo "[E2E] Creating test bucket..."
"$(git rev-parse --show-toplevel)/multitool/mt" objectstorage create test-bucket fsn1 --storage-provider hetzner --mode simulate --server http://localhost:18080 --automation-mode

# List buckets (should show the created bucket)
echo "[E2E] Listing buckets after creation (should show test-bucket)..."
"$(git rev-parse --show-toplevel)/multitool/mt" objectstorage list --storage-provider hetzner --mode simulate --server http://localhost:18080 --automation-mode

# Delete the test bucket
echo "[E2E] Deleting test bucket..."
"$(git rev-parse --show-toplevel)/multitool/mt" objectstorage delete test-bucket --storage-provider hetzner --mode simulate --server http://localhost:18080 --automation-mode --force

# List buckets again (should be empty)
echo "[E2E] Listing buckets after deletion (should be empty)..."
"$(git rev-parse --show-toplevel)/multitool/mt" objectstorage list --storage-provider hetzner --mode simulate --server http://localhost:18080 --automation-mode

# Stop cube-server using control script
"$(git rev-parse --show-toplevel)/cube-server/scripts/cube_server_control.sh" stop 18080

# Use a temp file for simulation persistence
SIM_PERSIST_FILE="/tmp/cube_server_sim_buckets_e2e_$$.json"
export CUBE_SERVER_SIM_PERSIST="$SIM_PERSIST_FILE"
echo "[E2E] Using simulation persistence file: $CUBE_SERVER_SIM_PERSIST"

# Show persistence file contents for debug
if [ -f "$SIM_PERSIST_FILE" ]; then
  echo "[E2E] Persistence file contents after test:"
  cat "$SIM_PERSIST_FILE"
  rm -f "$SIM_PERSIST_FILE"
fi

echo "[E2E] Hetzner S3 simulation test completed successfully."
