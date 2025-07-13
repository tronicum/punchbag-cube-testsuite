#!/bin/bash
set -euo pipefail

API_KEY=${API_KEY:-dev-secret-key}
BASE_URL=${BASE_URL:-http://localhost:8080}

# Test AKS simulation (create, get, list, delete)
echo "[TEST] Simulate AKS create..."
resp=$(curl -s -X POST "$BASE_URL/api/simulate/azure/aks" -H "X-API-Key: $API_KEY" -H 'Content-Type: application/json' -d '{"name":"test-aks","resourceGroup":"test-rg","location":"eastus"}')
echo "$resp"
id=$(echo "$resp" | grep -o '"ID":"[^"]*' | grep -o '[^\"]*$')

if [ -z "$id" ]; then
  echo "[FAIL] AKS create did not return an ID"; exit 1
fi

echo "[TEST] Simulate AKS get..."
curl -s "$BASE_URL/api/simulate/azure/aks?id=$id" -H "X-API-Key: $API_KEY"

echo "[TEST] Simulate AKS list..."
curl -s "$BASE_URL/api/simulate/azure/aks" -H "X-API-Key: $API_KEY"

echo "[TEST] Simulate AKS delete..."
curl -s -X DELETE "$BASE_URL/api/simulate/azure/aks?id=$id" -H "X-API-Key: $API_KEY"

echo "[TEST] Validation endpoint (AKS)..."
curl -s -X POST "$BASE_URL/api/validation?provider=azure&resource=aks" -H "X-API-Key: $API_KEY" -H 'Content-Type: application/json' -d '{"name":"test-aks","resourceGroup":"test-rg","location":"eastus"}'

echo "[TEST] Proxy endpoint (AKS, dryrun)..."
curl -s -X POST "$BASE_URL/api/proxy/azure/aks?dryrun=true" -H "X-API-Key: $API_KEY" -H 'Content-Type: application/json' -d '{"name":"test-aks","resourceGroup":"test-rg","location":"eastus"}'

echo "[PASS] All end-to-end API tests completed."
