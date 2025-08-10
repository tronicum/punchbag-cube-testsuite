# Fast local CI with act, Colima, and Go cache (when workspace path contains colons)

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
