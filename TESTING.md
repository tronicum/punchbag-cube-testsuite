# Testing and Cross-Platform Validation

This document is for maintainers and CI contributors. It describes how to test the `mt` CLI and related tools across different architectures and OSes using Docker and Colima. **End users do not need Docker or Colima to use `mt`.**

## Local and CI Testing Matrix

## Multiarch Build for mt CLI

To build the `mt` CLI for all supported OS/arch combinations (for packaging, Docker, and cross-arch testing):

1. From the project root, run:

   ```sh
   make build-multiarch
   ```

2. This will produce binaries in `testing/docker/bin/` named `mt-<os>-<arch>`, e.g.:

   - `mt-linux-amd64`
   - `mt-linux-arm64`
   - `mt-darwin-amd64`
   - `mt-darwin-arm64`

3. Use these binaries for packaging, Docker-based tests, or distribution.

## OS/Package Manager Testing with Docker

To verify that `mt` correctly detects the OS and package manager on all major Linux distributions:

1. Build the multiarch binaries as above.
2. Copy the appropriate binary to `testing/docker/mt` (e.g., `cp testing/docker/bin/mt-linux-amd64 testing/docker/mt`).
3. Build and run the test container for each OS:

   ```sh
   # Ubuntu example
   docker build -f testing/docker/Dockerfile.ubuntu -t mt-ubuntu testing/docker
   docker run --rm mt-ubuntu os-detect
   docker run --rm mt-ubuntu list-packages

   # Alpine example
   docker build -f testing/docker/Dockerfile.alpine -t mt-alpine testing/docker
   docker run --rm mt-alpine os-detect
   docker run --rm mt-alpine list-packages
   # ...repeat for other Dockerfiles
   ```

4. This ensures the CLI's OS and package manager detection logic works on all supported distros.

**Note:** These containers are for CI/maintainer testing only, not for end users.

## CI Integration

You can add these steps to your CI pipeline to ensure cross-arch and cross-distro compatibility for every commit or release.

## Troubleshooting

- If Docker/Colima is not running, start it before running the above commands.
1. **Install [Colima](https://github.com/abiosoft/colima) and Docker** (if not already):
   - On macOS: `brew install colima docker`
   - On Linux: Use your package manager for Docker, then install Colima if desired.

2. **Start Colima:**
   ```sh
   colima start
   ```

3. **Build and test in a Linux container:**
   ```sh
   docker build -t mt-test -f Dockerfile.test .
   docker run --rm -it mt-test
   ```

4. **Sample Dockerfile.test:**
   ```Dockerfile
   FROM golang:1.22-bullseye
   WORKDIR /app
   COPY . .
   RUN go build -o mt ./multitool
   RUN go test ./multitool/... ./shared/...
   ENTRYPOINT ["/bin/bash"]
   ```

## Notes
- This setup ensures the CLI works on Linux and in CI, regardless of the developer's host OS.
- For macOS-specific issues, test natively on a Mac.
- No true OSX containers exist; Colima is for Linux container emulation on Mac.
- Do not document Docker/Colima for end users—this is for maintainers and CI only.

## Troubleshooting
- If you encounter file permission issues, ensure your user has access to the Docker/Colima socket.
- For Go module cache issues, add `RUN go clean -modcache` before building.

---

For questions, see the main README or open an issue.
