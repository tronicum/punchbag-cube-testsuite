#!/bin/bash
# Purge all files over 100MB from git history using git-filter-repo
# Requires: git-filter-repo (brew install git-filter-repo)

set -e

# Find all files over 100MB in the current repo (including .gitignored)
find . -type f -size +100M | grep -v "/.git/" > large_files_to_purge.txt

if [ ! -s large_files_to_purge.txt ]; then
  echo "No files over 100MB found. Nothing to purge."
  exit 0
fi

if ! command -v git-filter-repo &> /dev/null; then
  echo "git-filter-repo is not installed. Please install it first: https://github.com/newren/git-filter-repo"
  exit 1
fi

echo "Purging the following files from git history:"
cat large_files_to_purge.txt

# Build --path arguments for git-filter-repo
ARGS=""
while read -r file; do
  ARGS+=" --path '$file'"
done < large_files_to_purge.txt

# Run git-filter-repo to remove all large files from history
# shellcheck disable=SC2086
eval git filter-repo $ARGS --invert-paths

echo "Force-pushing cleaned history to origin..."
git push origin --force --all

echo "Force-pushing cleaned tags to origin..."
git push origin --force --tags

echo "Large files purged from history and remote updated."
