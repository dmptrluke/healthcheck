# healthcheck

Lightweight healthcheck binary for distroless Docker containers. ~2 MB static binary, no dependencies.

Two modes: HTTP endpoint checks and heartbeat file-age checks.

## Usage

Copy into your image from GHCR:

```dockerfile
COPY --from=ghcr.io/dmptrluke/healthcheck:1 /healthcheck /usr/local/bin/healthcheck
```

### HTTP check

Checks that an HTTP endpoint returns 200. Timeout is 4 seconds.

```yaml
healthcheck:
  test: ["CMD", "healthcheck", "http", "127.0.0.1:8000", "/health/"]
  interval: 60s
  timeout: 5s
  retries: 3
```

### File-age check

Reads a unix timestamp from a file and checks it's within the given number of seconds. Useful for worker processes that write a heartbeat file.

```yaml
healthcheck:
  test: ["CMD", "healthcheck", "file-age", "/tmp/heartbeat", "120"]
  interval: 30s
  timeout: 5s
  retries: 3
```

The heartbeat file should contain a unix timestamp (integer or float), e.g. Python's `time.time()` or `date +%s`.
