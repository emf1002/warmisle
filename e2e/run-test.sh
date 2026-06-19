#!/bin/bash
# ============================================================
# E2E 测试隔离运行器
# 用法:
#   bash e2e/run-test.sh tests/auth.spec.ts
#   bash e2e/run-test.sh tests/ledger.spec.ts tests/todo.spec.ts
#   bash e2e/run-test.sh --port 8090 tests/auth.spec.ts
#   bash e2e/run-test.sh --all
# ============================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
BINARY="$ROOT_DIR/dist/warmisle.exe"
PORT=8081
TEST_FILES=()

# Parse arguments
while [[ $# -gt 0 ]]; do
  case "$1" in
    --port) PORT="$2"; shift 2 ;;
    --all)  TEST_FILES=("."); shift ;;
    *)      TEST_FILES+=("$1"); shift ;;
  esac
done

if [[ ${#TEST_FILES[@]} -eq 0 ]]; then
  echo "Usage: run-test.sh [--port PORT] <test-file> [test-file2 ...]"
  echo "       run-test.sh [--port PORT] --all"
  exit 1
fi

# Auto-increment port if occupied
while netstat -an 2>/dev/null | grep -q ":${PORT} .*LISTEN"; do
  PORT=$((PORT + 1))
done

DB_PATH="$SCRIPT_DIR/e2e-data/test-${PORT}.db"
# Convert to Windows path for the .exe binary (env vars are NOT auto-converted by MSYS2)
DB_PATH_WIN=$(cygpath -w "$DB_PATH" 2>/dev/null || echo "$DB_PATH")
BASE_URL="http://localhost:${PORT}"

# Clean stale DB
rm -f "$DB_PATH" "$DB_PATH-wal" "$DB_PATH-shm" 2>/dev/null || true

# Build binary if missing
if [[ ! -f "$BINARY" ]]; then
  echo "[run-test] Binary not found, building..."
  (cd "$ROOT_DIR/frontend" && npm run build -- --emptyOutDir)
  (cd "$ROOT_DIR/backend" && go build -o ../dist/warmisle .)
fi

# Start server in background
echo "[run-test] Starting server on :${PORT} ..."
HC_PORT="$PORT" HC_DB_PATH="$DB_PATH_WIN" HC_TEST_MODE=true "$BINARY" &
SERVER_PID=$!

# Cleanup on exit
cleanup() {
  kill $SERVER_PID 2>/dev/null || true
  wait $SERVER_PID 2>/dev/null || true
  rm -f "$DB_PATH" "$DB_PATH-wal" "$DB_PATH-shm" 2>/dev/null || true
}
trap cleanup EXIT

# Wait for server ready
echo -n "[run-test] Waiting for server"
for i in $(seq 1 75); do
  if curl -s "${BASE_URL}/api/init/check" >/dev/null 2>&1; then
    echo " ready"
    break
  fi
  if [[ $i -eq 75 ]]; then
    echo " TIMEOUT"
    exit 1
  fi
  echo -n "."
  sleep 0.2
done

# Run tests
export HC_TEST_PORT="$PORT"
export HC_TEST_BASE_URL="$BASE_URL"

echo "[run-test] Running: npx playwright test ${TEST_FILES[*]}"
cd "$SCRIPT_DIR"
npx playwright test "${TEST_FILES[@]}"
TEST_RESULT=$?

echo "[run-test] Done (exit code: $TEST_RESULT)"
exit $TEST_RESULT
