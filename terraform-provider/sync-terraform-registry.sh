#!/bin/bash
# sync-terraform-registry.sh
# Usage: ./sync-terraform-registry.sh <version-tag>
# Example: ./sync-terraform-registry.sh v0.1.0-beta

set -e

REPO_URL="https://github.com/tronicum/punchbag-cube-testsuite"
TAG="$1"

if [ -z "$TAG" ]; then
  echo "Usage: $0 <version-tag>"
  exit 1
fi

git add .
git commit -m "Release $TAG for multipass-cloud-layer Terraform provider"
git push origin main

git tag "$TAG"
git push origin "$TAG"

echo "Pushed $TAG to $REPO_URL."
echo "Check https://registry.terraform.io/publish/provider for status."
