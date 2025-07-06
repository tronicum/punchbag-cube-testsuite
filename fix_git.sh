#!/bin/bash
# Remove large files from git, add .terraform to .gitignore, recommit, and force-push

set -e

echo "Adding .terraform/ to .gitignore..."
echo ".terraform/" >> .gitignore

echo "Removing .terraform and provider binaries from git tracking..."
git rm -r --cached examples/.terraform || true
git rm -r --cached .terraform || true

echo "Re-adding all files except ignored..."
git add .gitignore
git add .

echo "Committing cleanup..."
git commit -m "Remove large Terraform provider binaries and .terraform directory from git, update .gitignore"

echo "Force-pushing to remote..."
git push -f

echo "Cleanup and push complete."
