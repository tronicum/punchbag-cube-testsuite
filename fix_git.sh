#!/bin/bash
# Remove large files from git, add .terraform to .gitignore, purge all files >100MB from history, recommit, and force-push

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
git commit -m "Remove large Terraform provider binaries and .terraform directory from git, update .gitignore" || true

# Purge all files over 100MB from git history
if ! command -v git-filter-repo &> /dev/null; then
  echo "git-filter-repo is not installed. Please install it first: https://github.com/newren/git-filter-repo"
  exit 1
fi

find . -type f -size +100M | grep -v "/.git/" > large_files_to_purge.txt
if [ -s large_files_to_purge.txt ]; then
  echo "Purging the following files from git history:"
  cat large_files_to_purge.txt
  ARGS=""
  while read -r file; do
    ARGS+=" --path '$file'"
  done < large_files_to_purge.txt
  # shellcheck disable=SC2086
  eval git filter-repo $ARGS --invert-paths
else
  echo "No files over 100MB found."
fi

echo "Force-pushing cleaned history to origin..."
git push origin --force --all

echo "Force-pushing cleaned tags to origin..."
git push origin --force --tags

echo "Cleanup and push complete."
