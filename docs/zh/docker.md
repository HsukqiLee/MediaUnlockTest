# Docker 使用指南

## 镜像

我们提供两个官方 Docker 镜像：

-   `ghcr.io/hsukqilee/unlock-test`: 用于一次性测试的命令行工具。
-   `ghcr.io/hsukqilee/unlock-monitor`: 用于持续运行的监控服务。

### 标签策略

-   `latest`: 最新稳定版本。
-   `edge`: 来自 main 分支的最新构建（如果有）。
-   `vX.Y.Z`: 特定版本。
-   `preview`: 手动触发构建的版本。
-   `prerelease`: 最新预发布版本。

## 使用方法

### 命令行工具 (`unlock-test`)

使用 CLI 镜像运行快速测试。容器直接作为可执行程序运行。

```bash
docker run --rm ghcr.io/hsukqilee/unlock-test [参数]
```

**示例:**

```bash
# 基本测试 (IPv4)
docker run --rm ghcr.io/hsukqilee/unlock-test -m 4

# 使用代理 (注意网络环境)
# 如果使用宿主机代理，Linux下可以使用 --network host
docker run --rm --network host ghcr.io/hsukqilee/unlock-test -m 4 -http-proxy "http://127.0.0.1:7890"
```

### 监控服务 (`unlock-monitor`)

作为服务运行监控程序。

```bash
docker run -d \
  --name unlock-monitor \
  -p 8080:8080 \
  ghcr.io/hsukqilee/unlock-monitor -listen :8080
```

## 支持的架构

我们的镜像支持以下 Linux 架构：
-   amd64
-   arm64
-   386
-   arm/v7
-   arm/v6
-   ppc64le
-   s390x
-   riscv64
