#!/bin/sh
# Build mt for multiple platforms/architectures
set -e

# Output directory for built binaries
OUTDIR=../testing/docker/bin
mkdir -p "$OUTDIR"

# List of target platforms
PLATFORMS="linux/amd64 linux/arm64 darwin/amd64 darwin/arm64"

for PLATFORM in $PLATFORMS; do
  GOOS="${PLATFORM%%/*}"
  GOARCH="${PLATFORM##*/}"
  OUTFILE="$OUTDIR/mt-$GOOS-$GOARCH"
  echo "Building for $GOOS/$GOARCH -> $OUTFILE"
  env GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 go build -o "$OUTFILE"
done

echo "All binaries built in $OUTDIR:"
ls -lh "$OUTDIR"
