# Docker Usage Guide

## Images

We provide two official Docker images:

-   `ghcr.io/hsukqilee/unlock-test`: The CLI tool for one-off tests.
-   `ghcr.io/hsukqilee/unlock-monitor`: The monitoring service.

### Tag Strategy

-   `latest`: The latest stable release.
-   `edge`: The latest build from the main branch (if applicable).
-   `vX.Y.Z`: Specific version releases.
-   `preview`: Manual builds triggered via workflow dispatch.

## Usage

### CLI Tool (`unlock-test`)

Run a quick test using the CLI image. The container acts as the executable.

```bash
docker run --rm ghcr.io/hsukqilee/unlock-test [arguments]
```

**Example:**

```bash
# Basic test (IPv4)
docker run --rm ghcr.io/hsukqilee/unlock-test -m 4

# With proxy
docker run --rm ghcr.io/hsukqilee/unlock-test -m 4 -http-proxy "http://host.docker.internal:7890"
```

### Monitor Service (`unlock-monitor`)

Run the monitor as a service.

```bash
docker run -d \
  --name unlock-monitor \
  -p 8080:8080 \
  ghcr.io/hsukqilee/unlock-monitor -listen :8080
```

## Supported Architectures

Our images support the following Linux architectures:
-   amd64
-   arm64
-   386
-   arm/v7
-   arm/v6
-   ppc64le
-   s390x
-   riscv64
