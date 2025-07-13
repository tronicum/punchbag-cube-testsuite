#!/bin/bash
# end2end_hetzner_s3_test.sh
# End-to-end test for Hetzner S3 bucket listing in direct mode
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
WORKSPACE_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
MULTITOOL_BIN="$WORKSPACE_ROOT/multitool/mt"

export HETZNER_S3_ACCESS_KEY="${HETZNER_S3_ACCESS_KEY:-dummy-access-key}"
export HETZNER_S3_SECRET_KEY="${HETZNER_S3_SECRET_KEY:-dummy-secret-key}"
export HETZNER_S3_REGION="${HETZNER_S3_REGION:-fsn1}"

OUTPUT_JSON="$SCRIPT_DIR/hetzner_s3_buckets.json"

set +e
"$MULTITOOL_BIN" objectstorage list hetzner --region "$HETZNER_S3_REGION" > "$OUTPUT_JSON"
status=$?
set -e

if [ $status -ne 0 ]; then
  echo "[ERROR] Hetzner S3 bucket listing failed."
  cat "$OUTPUT_JSON"
  exit 1
fi

echo "[SUCCESS] Hetzner S3 bucket listing output:"
cat "$OUTPUT_JSON"
rm -f "$OUTPUT_JSON"
