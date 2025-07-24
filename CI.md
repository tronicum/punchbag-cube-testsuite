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

## References
- [act documentation](https://github.com/nektos/act)
- [Colima documentation](https://github.com/abiosoft/colima)
