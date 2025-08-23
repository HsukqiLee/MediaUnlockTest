# MediaUnlockTest CLI Development Guide

**English Development Guide** | [中文开发文档](https://github.com/HsukqiLee/MediaUnlockTest/blob/main/DEVELOPMENT.md)

## Table of Contents

- [Project Structure](#project-structure)
- [CLI Development Standards](#cli-development-standards)
- [Function Usage Guide](#function-usage-guide)
- [Adding New Detection Items](#adding-new-detection-items)
- [CLI Feature Extensions](#cli-feature-extensions)
- [Testing and Debugging](#testing-and-debugging)
- [Contribution Guidelines](#contribution-guidelines)

## Project Structure

```
MediaUnlockTest/
├── checks/           # Detection functions directory
│   ├── mediaunlock.go    # Core interfaces and type definitions
│   ├── Netflix.go        # Netflix detection implementation
│   ├── DisneyPlus.go     # Disney+ detection implementation
│   └── ...               # Other detection items
├── cli/              # Command line tool
│   ├── main.go       # Main program entry point
│   ├── build.bat     # Windows build script
│   └── build.sh      # Unix build script
└── go.mod            # Go module file
```

## CLI Development Standards

### File Naming Conventions

1. **Detection Function Files**: Use PascalCase naming with `.go` extension
   - Correct: `Netflix.go`, `DisneyPlus.go`, `BBCiPlayer.go`
   - Incorrect: `netflix.go`, `disney_plus.go`, `bbc_iplayer.go`

2. **Function Names**: Use PascalCase naming
   - Correct: `NetflixRegion()`, `DisneyPlus()`, `BBCiPlayer()`
   - Incorrect: `netflix_region()`, `disneyPlus()`, `bbciplayer()`

3. **Variable Names**: Use camelCase naming
   - Correct: `httpClient`, `userAgent`, `timeout`
   - Incorrect: `http_client`, `UserAgent`, `TIMEOUT`

### CLI Code Structure Standards

1. **Package Declaration**: Each file must start with `package main` (in cli directory)
2. **Import Order**: Standard library → Third-party libraries → Local packages
3. **Function Order**: main function → Public functions → Private functions → Helper functions
4. **Comments**: Each public function must have comment documentation
5. **Error Handling**: Use unified error handling approach

### Command Line Parameter Standards

1. **Parameter Definition**: Use `flag` package to define command line parameters
2. **Parameter Validation**: Validate key parameters for effectiveness
3. **Default Values**: Provide reasonable default values for all parameters
4. **Help Information**: Provide clear parameter descriptions

```go
// Example: Parameter definition
var (
    Interface   string
    DNSServers  string
    HTTPProxy   string
    SocksProxy  string
    ShowVersion bool
    CheckUpdate bool
    Debug       bool
    IPMode      int
    Conc        uint64
)

// Parameter binding
flag.StringVar(&Interface, "I", "", "Source IP or network interface to use for connections")
flag.StringVar(&DNSServers, "dns-servers", "", "Custom DNS servers (format: ip:port)")
flag.BoolVar(&Debug, "debug", false, "Enable debug mode for verbose output")
flag.IntVar(&IPMode, "m", 0, "Connection mode: 0=auto (default), 4=IPv4 only, 6=IPv6 only")
```

## Function Usage Guide

### Core Types and Interfaces

#### Result Struct

```go
type Result struct {
    Status       Status      // Detection status
    Region       string      // Region information
    Info         string      // Additional information
    Err          error       // Error information
    CachedResult bool        // Whether it's a cached result
}
```

#### Status Enum

```go
const (
    StatusOK           Status = iota // Unlock successful
    StatusNo                         // Not unlocked
    StatusRestricted                 // Region restricted
    StatusNetworkErr                 // Network error
    StatusErr                        // Other errors
    StatusBanned                     // Banned
    StatusUnexpected                 // Unexpected result
    StatusFailed                     // Test failed
)
```

### HTTP Client Usage

#### Auto-select Client

```go
func NetflixRegion(client http.Client) m.Result {
    // Use the passed client parameter
    resp, err := client.Get("https://api.netflix.com/status")
    if err != nil {
        return m.Result{
            Status: m.StatusNetworkErr,
            Err:    err,
        }
    }
    // ... process response
}
```

#### Force Specific Protocol

```go
// IPv4 client
resp, err := m.Ipv4HttpClient.Get(url)

// IPv6 client  
resp, err := m.Ipv6HttpClient.Get(url)

// Auto client
resp, err := m.AutoHttpClient.Get(url)
```

### CLI Utility Functions

#### Progress Bar Management

```go
import pb "github.com/schollz/progressbar/v3"

func NewBar(count int64) *pb.ProgressBar {
    return pb.NewOptions64(
        count,
        pb.OptionSetDescription("Testing..."),
        pb.OptionSetWriter(os.Stderr),
        pb.OptionSetWidth(60),
        pb.OptionThrottle(200*time.Millisecond),
        pb.OptionShowCount(),
        pb.OptionClearOnFinish(),
        pb.OptionEnableColorCodes(true),
    )
}
```

#### Color Output

```go
import "github.com/fatih/color"

var (
    Red     = color.New(color.FgRed).SprintFunc()
    Green   = color.New(color.FgGreen).SprintFunc()
    Yellow  = color.New(color.FgYellow).SprintFunc()
    Blue    = color.New(color.FgBlue).SprintFunc()
    SkyBlue = color.New(color.FgCyan).SprintFunc()
)

// Usage example
fmt.Println(Green("Success"), Red("Error"), Yellow("Warning"))
```

## Adding New Detection Items

### Step 1: Create Detection Function

Create a new file in the `checks/` directory, for example `NewService.go`:

```go
package checks

import (
    "net/http"
    "strings"
)

// NewService detects unlock status of new service
func NewService(client http.Client) m.Result {
    // 1. Send request
    resp, err := client.Get("https://api.newservice.com/status")
    if err != nil {
        return m.Result{
            Status: m.StatusNetworkErr,
            Err:    err,
        }
    }
    defer resp.Body.Close()

    // 2. Check status code
    if resp.StatusCode != 200 {
        return m.Result{
            Status: m.StatusErr,
            Err:    fmt.Errorf("unexpected status code: %d", resp.StatusCode),
        }
    }

    // 3. Read response content
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return m.Result{
            Status: m.StatusErr,
            Err:    err,
        }
    }

    // 4. Parse response
    content := string(body)
    
    // 5. Determine unlock status
    if strings.Contains(content, "unlocked") {
        return m.Result{
            Status: m.StatusOK,
            Region: "US", // If specific region is known
        }
    } else if strings.Contains(content, "restricted") {
        return m.Result{
            Status: m.StatusRestricted,
            Info:   "Service is restricted in this region",
        }
    } else {
        return m.Result{
            Status: m.StatusNo,
        }
    }
}
```

### Step 2: Register in main.go

Find the corresponding test array in `cli/main.go` and add the new item:

```go
var NorthAmericaTests = []testItem{
    // ... existing items
    {"New Service", m.NewService, false}, // Last parameter indicates IPv6 support
}
```

### Step 3: Test and Verify

```bash
# Run test mode
go run cli/main.go -test

# Or test specific item directly
go run cli/main.go -nf
```

## CLI Feature Extensions

### Adding New Command Line Parameters

1. **Define Variables**: Define new parameter variables before `main()` function
2. **Bind Parameters**: Use `flag` package to bind parameters
3. **Parameter Processing**: Add parameter processing logic in `main()` function

```go
// 1. Define new parameters
var (
    NewParam string
    NewFlag  bool
)

// 2. Bind parameters
flag.StringVar(&NewParam, "new-param", "", "Description of new parameter")
flag.BoolVar(&NewFlag, "new-flag", false, "Description of new flag")

// 3. Process parameters
if NewFlag {
    // Logic to handle new flag
}
```

### Adding New Detection Modes

1. **Define Test Arrays**: Create new test item arrays
2. **Add Selection Logic**: Add new selection items in `ReadSelect()` function
3. **Execution Logic**: Add execution logic in `main()` function

```go
// 1. Define new test arrays
var NewRegionTests = []testItem{
    {"Service 1", m.Service1, false},
    {"Service 2", m.Service2, true},
}

// 2. Add selection items in ReadSelect()
fmt.Println("[12] :   New Region Platform")

// 3. Handle selection
case "12":
    NewRegion = true

// 4. Add to regions array
{Enabled: NewRegion, Name: "NewRegion", Tests: NewRegionTests}
```

## Detection Logic Standards

### Status Determination Priority

1. **Network Error First**: If unable to connect, return `StatusNetworkErr`
2. **Status Code Check**: Check if HTTP status code is normal
3. **Content Parsing**: Determine unlock status based on response content
4. **Default Handling**: Return `StatusNo` if unable to determine

### Region Information Handling

```go
// If specific region is known
if strings.Contains(content, "US") {
    return m.Result{
        Status: m.StatusOK,
        Region: "US",
    }
}

// If only know it's unlocked but not the region
if strings.Contains(content, "unlocked") {
    return m.Result{
        Status: m.StatusOK,
        // Region field left empty
    }
}
```

### Error Information Handling

```go
// Network error
if err != nil {
    return m.Result{
        Status: m.StatusNetworkErr,
        Err:    err, // Preserve original error
    }
}

// Business logic error
if !isValidResponse(content) {
    return m.Result{
        Status: m.StatusErr,
        Err:    fmt.Errorf("invalid response format"),
        Info:   "Response does not match expected format",
    }
}
```

## Testing and Debugging

### Local Testing

```bash
# Enter project directory
cd MediaUnlockTest

# Run all tests
go run cli/main.go

# Test specific region only
go run cli/main.go -m 4  # IPv4 only

# Enable debug mode
go run cli/main.go -debug

# Limit concurrent count
go run cli/main.go -conc 10

# Test with proxy
go run cli/main.go -http-proxy "http://127.0.0.1:8080"
```

### Debugging Tips

1. **Use Debug Mode**: Add `-debug` parameter to see detailed error information
2. **Check Network Connection**: Confirm if target service is accessible
3. **View Response Content**: Add log output in code to see response content
4. **Test with Proxy**: Test proxy environment through `-http-proxy` or `-socks-proxy`
5. **Performance Analysis**: Use `-conc` parameter to adjust concurrent count

### Common Issue Troubleshooting

1. **Timeout Errors**: Check network connection and firewall settings
2. **Parsing Errors**: Check if response format matches expectations
3. **Region Detection Errors**: Confirm if detection logic is correct
4. **Concurrency Issues**: Use `-conc` parameter to limit concurrent count
5. **Proxy Issues**: Check if proxy configuration is correct

## Building and Deployment

### Build Scripts

The project provides cross-platform build scripts:

```bash
# Windows
cli/build.bat

# Unix/Linux/macOS
cli/build.sh
```

### Build Options

```bash
# Basic build
go build -o unlock-test cli/main.go

# Cross-platform build
GOOS=linux GOARCH=amd64 go build -o unlock-test-linux cli/main.go
GOOS=windows GOARCH=amd64 go build -o unlock-test.exe cli/main.go
GOOS=darwin GOARCH=amd64 go build -o unlock-test-mac cli/main.go
```

## Contribution Guidelines

### Submitting Pull Request

1. **Fork Project**: Fork this project on GitHub
2. **Create Branch**: Create feature branch `feature/new-service`
3. **Write Code**: Write code according to development standards
4. **Test Verification**: Ensure code runs normally
5. **Submit PR**: Create Pull Request and describe changes

### Code Review Points

1. **Function Completeness**: Is detection logic complete
2. **Error Handling**: Are various error situations handled correctly
3. **Code Standards**: Does it comply with project coding standards
4. **Performance Considerations**: Are timeout and concurrency considered
5. **Documentation Updates**: Are related documents updated
6. **CLI Experience**: Do new features provide good user experience

### Issue Feedback

If you encounter problems or have suggestions, please:

1. Search GitHub Issues to see if similar issues already exist
2. Create new Issue with detailed problem description and reproduction steps
3. Provide system environment, Go version and other detailed information
4. If possible, provide related logs or error information

## Update Log

### v1.0.0
- Initial version release
- Support basic streaming media detection functionality

### v1.1.0
- Add IPv6 support
- Optimize concurrent processing
- Add caching mechanism

### v1.2.0
- Improve progress bar display
- Add active test display
- Optimize error handling

---

Thank you for contributing to the MediaUnlockTest CLI project! If you have questions, please check [Issues](https://github.com/HsukqiLee/MediaUnlockTest/issues) or create new discussions.

## Related Documentation

- [Monitor Development Guide](https://github.com/HsukqiLee/MediaUnlockTest/blob/main/monitor/README.md)
- [Monitor English Documentation](https://github.com/HsukqiLee/MediaUnlockTest/blob/main/monitor/README_en.md)
