#!/bin/bash
# Usage: ./cube_server_control.sh start|stop [port]
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
CUBE_SERVER_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
PID_FILE="$CUBE_SERVER_DIR/cube-server.pid"
PORT="18080"
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
  echo "[INFO] Building cube-server via Makefile..."
  make build
  chmod +x cube-server
  # Kill any process listening on the target port to avoid address-in-use errors
  EXISTING_PID=$(lsof -t -nP -iTCP:$PORT | head -n 1 || true)
  if [ -n "$EXISTING_PID" ]; then
    echo "[WARN] Killing existing process on port $PORT (PID $EXISTING_PID)"
    kill $EXISTING_PID
    sleep 1
  fi
  echo "[INFO] Starting cube-server on port $PORT..."
  if [ -n "${CUBE_SERVER_SIM_PERSIST:-}" ]; then
    echo "[INFO] Using simulation persistence file: $CUBE_SERVER_SIM_PERSIST"
    nohup env CUBE_SERVER_SIM_PERSIST="$CUBE_SERVER_SIM_PERSIST" ./cube-server --simulate hetzner --port "$PORT" > "$LOG_FILE" 2>&1 &
  else
    nohup ./cube-server --simulate hetzner --port "$PORT" > "$LOG_FILE" 2>&1 &
  fi
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
  # Wait up to 10 seconds for server to listen on the port
  LISTENED=0
  for i in {1..10}; do
    if lsof -nP -iTCP:$PORT | grep LISTEN > /dev/null; then
      LISTENED=1
      break
    fi
    sleep 1
  done
  if [ $LISTENED -ne 1 ]; then
    echo "[ERROR] cube-server not listening on port $PORT after 10s. See $LOG_FILE for details."
    kill $SERVER_PID
    exit 1
  fi
  echo $SERVER_PID > "$PID_FILE"
  echo "[INFO] cube-server started with PID $SERVER_PID on port $PORT."
}


stop_server() {
  if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    if ps -p $PID > /dev/null; then
      echo "[INFO] Stopping cube-server (PID $PID)"
      kill $PID
      # Wait for process to terminate
      for i in {1..10}; do
        if ! ps -p $PID > /dev/null; then
          break
        fi
        sleep 1
      done
      if ps -p $PID > /dev/null; then
        echo "[ERROR] cube-server did not stop after 10 seconds."
        exit 1
      fi
      rm -f "$PID_FILE"
      echo "[INFO] cube-server stopped."
    else
      echo "[WARN] No running cube-server process found for PID $PID. Removing stale PID file."
      rm -f "$PID_FILE"
    fi
  else
    echo "[WARN] No PID file found. Is cube-server running?"
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
    echo "Usage: $0 start|stop [port] [--debug]"
    exit 1
    ;;
esac
