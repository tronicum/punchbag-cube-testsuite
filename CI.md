# See `TODOs.md` for the multitool CLI improvement roadmap and actionable feature tasks.
# Local CI and GitHub Actions Testing

This project supports local testing of GitHub Actions workflows using [`act`](https://github.com/nektos/act) and [Colima](https://github.com/abiosoft/colima) for container runtime compatibility on macOS (including Apple Silicon/M-series).

## Prerequisites
- Homebrew (macOS)
- Docker Desktop **or** Colima (recommended for Apple Silicon)
- act

## Setup

### 1. Install Colima (recommended for Apple Silicon)
```sh
brew install colima
colima start
```

### 2. Install act
```sh
brew install act
```

### 3. Start Colima (if not already running)
```sh
colima start
```

### 4. Run GitHub Actions workflows locally
```sh
act
```
- If you are on Apple Silicon (M-series) and encounter issues, use:
  ```sh
  act --container-architecture linux/amd64
  ```
- You may be prompted to select a default image size (Medium is usually sufficient).
- Make sure Docker or Colima is running before running `act`.


## Troubleshooting
- If you see `Cannot connect to the Docker daemon at unix:///var/run/docker.sock`, ensure Docker Desktop or Colima is running.
- For architecture issues on Apple Silicon, always try `--container-architecture linux/amd64`.

## Fast local CI with act, Colima, and Go cache (when workspace path contains colons)

If your workspace path contains a colon (`:`), Docker/Colima/act cannot use `--bind` or mount the workspace at the same path as the host. Instead, use the following command to enable Go cache reuse:

```sh
act --container-architecture linux/amd64 \
  --container-options "-v $HOME/.cache/go-build:/root/.cache/go-build -v $HOME/go/pkg/mod:/root/go/pkg/mod" -j build-and-test
```

- This mounts your Go build and module cache directories at the default root user locations inside the container.
- The workspace will be mounted at `/github/workspace` (the act default), which works even if your path contains colons.
- This avoids permission and volume errors, and enables Go cache reuse for faster local CI runs.

If you see permission errors, ensure your user owns the cache directories:

```sh
sudo chown -R $(whoami) ~/.cache/go-build ~/go/pkg/mod
```

Then restart Colima:

```sh
colima stop && colima start
```

Now re-run the act command above.

---

**Tested on macOS with Colima and act v0.2.61+**

## References
- [act documentation](https://github.com/nektos/act)
- [Colima documentation](https://github.com/abiosoft/colima)
