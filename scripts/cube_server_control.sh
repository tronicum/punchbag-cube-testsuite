#!/bin/bash
# Usage: ./cube_server_control.sh start|stop [port]
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
CUBE_SERVER_DIR="$SCRIPT_DIR/../cube-server"
PID_FILE="$CUBE_SERVER_DIR/cube-server.pid"
PORT="9000"
LOG_FILE="$CUBE_SERVER_DIR/cube_server.log"
DEBUG=0

# Parse flags
for arg in "$@"; do
  case $arg in
    --debug)
      DEBUG=1
      ;;
    [0-9]*)
      PORT="$arg"
      ;;
  esac
done

start_server() {
  cd "$CUBE_SERVER_DIR"
  pwd
  if [ "$DEBUG" -eq 1 ]; then
    set -x
  fi
  echo "[INFO] Building cube-server..."
  go build -v -o cube-server main.go
  chmod +x cube-server
  # Kill any process listening on the target port to avoid address-in-use errors
  EXISTING_PID=$(lsof -t -nP -iTCP:$PORT | head -n 1 || true)
  if [ -n "$EXISTING_PID" ]; then
    echo "[WARN] Killing existing process on port $PORT (PID $EXISTING_PID)"
    kill $EXISTING_PID
    sleep 1
  fi
  echo "[INFO] Starting cube-server on port $PORT..."
  nohup ./cube-server --simulate hetzner --port "$PORT" > "$LOG_FILE" 2>&1 &
  SERVER_PID=$!
  sleep 2
  if [ "$DEBUG" -eq 1 ]; then
    echo "[DEBUG] Checking process status for PID $SERVER_PID"
    ps -p $SERVER_PID || true
    echo "[DEBUG] Checking lsof for port $PORT"
    lsof -nP -iTCP:$PORT || true
    echo "[DEBUG] Last 20 lines of log:"
    tail -20 "$LOG_FILE"
  fi
  # Check if process is running
  if ! ps -p $SERVER_PID > /dev/null; then
    echo "[ERROR] cube-server process not running. See $LOG_FILE for details."
    exit 1
  fi
  # Check if server is listening on the port (macOS compatible)
  if ! lsof -nP -iTCP:$PORT | grep LISTEN | grep cube-server > /dev/null; then
    echo "[ERROR] cube-server not listening on port $PORT. See $LOG_FILE for details."
    kill $SERVER_PID
    exit 1
  fi
  echo $SERVER_PID > "$PID_FILE"
  echo "[INFO] cube-server started (PID $SERVER_PID, port $PORT)"
  if [ "$DEBUG" -eq 1 ]; then
    set +x
  fi
}

stop_server() {
  if [ ! -f "$PID_FILE" ]; then
    echo "[WARN] No PID file found. Is cube-server running?"
    exit 1
  fi
  SERVER_PID=$(cat "$PID_FILE")
  if ps -p $SERVER_PID > /dev/null; then
    echo "[INFO] Stopping cube-server (PID $SERVER_PID)"
    kill $SERVER_PID
    wait $SERVER_PID 2>/dev/null || true
    rm -f "$PID_FILE"
    echo "[INFO] cube-server stopped."
  else
    echo "[WARN] cube-server process not found. Removing stale PID file."
    rm -f "$PID_FILE"
  fi
}

case "${1:-}" in
  start)
    start_server
    ;;
  stop)
    stop_server
    ;;
  *)
    echo "Usage: $0 start|stop [port]"
    exit 1
    ;;
esac
