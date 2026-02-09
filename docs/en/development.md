# Development Guide

## Prerequisites

-   [Go](https://go.dev/dl/) 1.24+
-   [Make](https://www.gnu.org/software/make/) (Optional)
-   [Docker](https://www.docker.com/) (Optional)

## Project Structure

-   `checks/`: Core library containing all streaming service check logic.
-   `cli/`: CLI tool entry point (`main.go` and test definitions).
-   `monitor/`: Monitor service entry point.
-   `docker/`: Docker build files.

## Core Library (`checks/mediaunlock.go`)

The `checks` package provides tools for building checks.

### Result Struct

All check functions must return a `Result` struct:

```go
type Result struct {
    Status int    // Status Code
    Region string // Unlocked Region (e.g., "US", "JP")
    Info   string // Extra Info (e.g., "Only Available in ...")
    Err    error  // Error
}
```

### Status Constants

-   `StatusOK` (1): Unlocked / Supported
-   `StatusRestricted` (2): Restricted / Partially Supported
-   `StatusNo` (3): Not Supported / Banned
-   `StatusBanned` (4): IP Banned
-   `StatusNetworkErr` (-1): Network Error
-   `StatusErr` (-2): Parse/Logic Error
-   `StatusFailed` (5): Check Failed (Exception)

### HTTP Clients

Use the provided clients to ensure consistent behavior (User-Agent, TLS Fingerprint, etc.):

-   `AutoHttpClient`: Auto-selection.
-   `Ipv4HttpClient`: Force IPv4.
-   `Ipv6HttpClient`: Force IPv6.

### Helper Functions

-   `GET(c http.Client, url string, headers ...H) (*http.Response, error)`: Send GET request.
-   `PostJson`, `PostForm`: Send POST request.
-   `CheckGETStatus(...)`: Map status codes to results automatically.

## How to Add a New Check

### 1. Create a Check File

Create a new `.go` file in `checks/`, e.g., `checks/MyService.go`.

```go
package mediaunlocktest

import (
    "net/http"
    "strings"
)

func MyService(c http.Client) Result {
    // Send Request
    resp, err := GET(c, "https://api.myservice.com/check")
    if err != nil {
        return Result{Status: StatusNetworkErr, Err: err}
    }
    defer resp.Body.Close()

    // Logic
    if resp.StatusCode == 200 {
        return Result{Status: StatusOK, Region: "US"}
    } else if resp.StatusCode == 403 {
        return Result{Status: StatusNo}
    }

    return Result{Status: StatusUnexpected}
}
```

### 2. Register the Check

Open `cli/main.go`, find the corresponding region list (e.g., `NorthAmericaTests`), and add your function:

```go
var NorthAmericaTests = []testItem{
    // ...
    {"My Service", m.MyService, true}, // true indicates IPv6 support
}
```

## Build from Source

### CLI Tool

```bash
cd cli
go build -o ../unlock-test
```

### Monitor Service

```bash
cd monitor
go build -o ../unlock-monitor
```
