#!/bin/bash
set -euo pipefail

PORT=8081
SERVER_LOG=server_test.log

# Start the server in the background
echo "[INFO] Starting cube-server on port $PORT..."
go run ./cmd/cube-server/main.go > "$SERVER_LOG" 2>&1 &
SERVER_PID=$!
sleep 2

# Test /health endpoint
echo "[TEST] /health"
curl -sf http://localhost:$PORT/health && echo "[PASS] /health" || { echo "[FAIL] /health"; kill $SERVER_PID; exit 1; }

# Test /api/v1/clusters endpoint
echo "[TEST] /api/v1/clusters"
curl -sf http://localhost:$PORT/api/v1/clusters && echo "[PASS] /api/v1/clusters" || { echo "[FAIL] /api/v1/clusters"; kill $SERVER_PID; exit 1; }

# Test POST /api/v1/clusters (create cluster)
echo "[TEST] POST /api/v1/clusters (valid)"
resp=$(curl -sf -X POST http://localhost:$PORT/api/v1/clusters \
  -H 'Content-Type: application/json' \
  -d '{"name":"test-cluster","provider":"azure","location":"eastus"}')
echo "$resp" | grep '"id"' && echo "[PASS] POST /api/v1/clusters (valid)" || { echo "[FAIL] POST /api/v1/clusters (valid)"; kill $SERVER_PID; exit 1; }

# Test POST /api/v1/clusters (invalid, missing name)
echo "[TEST] POST /api/v1/clusters (invalid, missing name)"
if curl -sf -X POST http://localhost:$PORT/api/v1/clusters \
  -H 'Content-Type: application/json' \
  -d '{"provider":"azure","location":"eastus"}'; then
  echo "[FAIL] POST /api/v1/clusters (invalid, missing name) should have failed"
  kill $SERVER_PID; exit 1;
else
  echo "[PASS] POST /api/v1/clusters (invalid, missing name)"
fi

# Test DELETE /api/v1/clusters/:id (delete cluster)
cluster_id=$(echo "$resp" | grep -o '"id":"[^"]*' | grep -o '[^\"]*$')
if [ -n "$cluster_id" ]; then
  echo "[TEST] DELETE /api/v1/clusters/$cluster_id"
  curl -sf -X DELETE http://localhost:$PORT/api/v1/clusters/$cluster_id && echo "[PASS] DELETE /api/v1/clusters/$cluster_id" || { echo "[FAIL] DELETE /api/v1/clusters/$cluster_id"; kill $SERVER_PID; exit 1; }
fi

# Stop the server
echo "[INFO] Stopping cube-server (PID $SERVER_PID)"
if kill -0 $SERVER_PID 2>/dev/null; then
  kill $SERVER_PID
  wait $SERVER_PID 2>/dev/null || true
fi

echo "[INFO] All API tests completed."
