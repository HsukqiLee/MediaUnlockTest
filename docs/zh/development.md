吗# 开发指南

## 前置要求

-   [Go](https://go.dev/dl/) 1.24+
-   [Make](https://www.gnu.org/software/make/) (可选)
-   [Docker](https://www.docker.com/) (可选)

## 项目结构

-   `checks/`: 包含所有流媒体检测逻辑的核心库。
-   `cli/`: 命令行工具入口 (`main.go` 及测试列表定义)。
-   `monitor/`: 监控服务入口。
-   `docker/`: Docker 构建文件。

## 核心库 (`checks/mediaunlock.go`)

`checks` 包提供了构建检测所需的基础工具。

### Result 结构体

所有检测函数必须返回 `Result` 结构体：

```go
type Result struct {
    Status int    // 检测状态码
    Region string // 解锁地区 (如 "US", "JP")
    Info   string // 额外信息 (如 "Only Available in ...")
    Err    error  // 错误信息
}
```

### 状态码常量

-   `StatusOK` (1): 解锁 / 支持
-   `StatusRestricted` (2): 有限制 / 部分支持 (如仅自制剧)
-   `StatusNo` (3): 不支持 / 无法观看
-   `StatusBanned` (4): IP 被封禁
-   `StatusNetworkErr` (-1): 网络错误
-   `StatusErr` (-2): 解析/逻辑错误
-   `StatusFailed` (5): 检测失败 (通常指检测过程异常)

### HTTP 客户端

请使用库提供的客户端以确保行为一致（如 User-Agent、TLS 指纹等）：

-   `AutoHttpClient`: 自动选择（根据系统/参数）。
-   `Ipv4HttpClient`: 强制 IPv4。
-   `Ipv6HttpClient`: 强制 IPv6。

### 辅助函数

-   `GET(c http.Client, url string, headers ...H) (*http.Response, error)`: 发送 GET 请求。
-   `PostJson`, `PostForm`: 发送 POST 请求。
-   `CheckGETStatus(...)`: 自动根据状态码映射结果。

## 如何添加新的检测

### 1. 新建检测文件

在 `checks/` 目录下新建 `.go` 文件，例如 `checks/MyService.go`。

```go
package mediaunlocktest

import (
    "net/http"
    "strings"
)

func MyService(c http.Client) Result {
    // 发起请求
    resp, err := GET(c, "https://api.myservice.com/check")
    if err != nil {
        return Result{Status: StatusNetworkErr, Err: err}
    }
    defer resp.Body.Close()

    // 逻辑判断
    if resp.StatusCode == 200 {
        return Result{Status: StatusOK, Region: "US"}
    } else if resp.StatusCode == 403 {
        return Result{Status: StatusNo}
    }

    return Result{Status: StatusUnexpected}
}
```

### 2. 注册检测函数

打开 `cli/main.go`，找到对应的地区列表变量（如 `NorthAmericaTests`），将你的函数加入其中：

```go
var NorthAmericaTests = []testItem{
    // ...
    {"My Service", m.MyService, true}, // true 表示支持 IPv6
}
```

## 源码构建

### CLI 工具

```bash
cd cli
go build -o ../unlock-test
```

### 监控服务

```bash
cd monitor
go build -o ../unlock-monitor
```
