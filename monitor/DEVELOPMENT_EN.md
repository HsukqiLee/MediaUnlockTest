# MediaUnlockTest Monitor Development Guide

**English Development Guide** | [中文开发文档](https://github.com/HsukqiLee/MediaUnlockTest/blob/main/monitor/DEVELOPMENT.md)

## Table of Contents

- [Project Structure](#project-structure)
- [Development Standards](#development-standards)
- [Core Components](#core-components)
- [Adding New Monitoring Metrics](#adding-new-monitoring-metrics)
- [Custom Detection Logic](#custom-detection-logic)
- [Performance Optimization](#performance-optimization)
- [Testing and Debugging](#testing-and-debugging)
- [Contribution Guidelines](#contribution-guidelines)

## Project Structure

```
monitor/
├── main.go           # Main program entry
├── monitor.go        # Monitor core logic
├── exporter.go       # Prometheus metrics export
├── service.go        # Service management
├── update.go         # Auto update
├── build.bat         # Windows build script
├── build.sh          # Unix build script
├── README.md         # Usage documentation
├── DEVELOPMENT.md    # Development guide (this file)
└── DEVELOPMENT_EN.md # English development guide
```

## Development Standards

### File Naming Conventions

1. **Go Files**: Use snake_case naming with `.go` extension
   - Correct: `monitor.go`, `exporter.go`, `service.go`
   - Incorrect: `Monitor.go`, `Exporter.go`, `monitor.go`

2. **Function Names**: Use PascalCase naming
   - Correct: `StartMonitoring()`, `ExportMetrics()`, `HandleRequest()`
   - Incorrect: `start_monitoring()`, `exportMetrics()`, `handlerequest()`

3. **Variable Names**: Use camelCase naming
   - Correct: `monitorConfig`, `metricsExporter`, `updateInterval`
   - Incorrect: `monitor_config`, `MetricsExporter`, `UPDATE_INTERVAL`

### Code Structure Standards

1. **Package Declaration**: Each file must start with `package main`
2. **Import Order**: Standard library → Third-party libraries → Local packages
3. **Function Order**: main function → Public functions → Private functions → Helper functions
4. **Comments**: Each public function must have comment documentation
5. **Error Handling**: Use unified error handling approach

### Monitoring Metrics Standards

1. **Metric Naming**: Use `mediaunlock_` prefix with underscore separation
2. **Label Design**: Labels should be meaningful and reasonable in number
3. **Metric Types**: Choose appropriate metric types based on usage (Counter, Gauge, Histogram, etc.)
4. **Documentation**: Each metric must have clear help information

## Core Components

### Monitor Core Logic

#### Monitor Interface

```go
type Monitor interface {
    Start(ctx context.Context) error
    Stop() error
    GetStatus() Status
}

type Status struct {
    Running    bool
    StartTime  time.Time
    TestCount  int64
    ErrorCount int64
}
```

#### Detector Interface

```go
type Detector interface {
    Detect(ctx context.Context) (*Result, error)
    GetName() string
    GetRegion() string
}

type Result struct {
    Service     string
    Status      string
    Region      string
    Error       error
    Duration    time.Duration
    Timestamp   time.Time
}
```

### Prometheus Metrics Export

#### Metric Definitions

```go
import "github.com/prometheus/client_golang/prometheus"

var (
    // Total test count
    testTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "mediaunlock_test_total",
            Help: "Total number of tests performed",
        },
        []string{"service", "region"},
    )

    // Successful test count
    testSuccess = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "mediaunlock_test_success",
            Help: "Total number of successful tests",
        },
        []string{"service", "region"},
    )

    // Test duration
    testDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "mediaunlock_test_duration_seconds",
            Help:    "Test duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"service", "region"},
    )

    // Service status
    serviceStatus = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "mediaunlock_service_status",
            Help: "Current service status (0=failed, 1=success)",
        },
        []string{"service", "region"},
    )
)
```

#### Metric Registration

```go
func init() {
    // Register all metrics
    prometheus.MustRegister(
        testTotal,
        testSuccess,
        testDuration,
        serviceStatus,
    )
}
```

### Service Management

#### HTTP Service

```go
type MonitorService struct {
    monitor Monitor
    server  *http.Server
    config  *Config
}

func (s *MonitorService) Start() error {
    mux := http.NewServeMux()
    
    // Health check endpoint
    mux.HandleFunc("/health", s.handleHealth)
    
    // Metrics export endpoint
    mux.Handle("/metrics", promhttp.Handler())
    
    // Status query endpoint
    mux.HandleFunc("/status", s.handleStatus)
    
    s.server = &http.Server{
        Addr:    s.config.ListenAddr,
        Handler: mux,
    }
    
    return s.server.ListenAndServe()
}
```

## Adding New Monitoring Metrics

### Step 1: Define Metrics

Define new metrics in the `exporter.go` file:

```go
var (
    // Custom metric example
    customMetric = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "mediaunlock_custom_total",
            Help: "Custom metric description",
        },
        []string{"label1", "label2"},
    )
)
```

### Step 2: Register Metrics

Register new metrics in the `init()` function:

```go
func init() {
    prometheus.MustRegister(
        // ... existing metrics
        customMetric,
    )
}
```

### Step 3: Update Metrics

Update metric values in appropriate places:

```go
// Increment count
customMetric.WithLabelValues("value1", "value2").Inc()

// Set specific value
customMetric.WithLabelValues("value1", "value2").Add(5)
```

### Step 4: Test and Verify

```bash
# Start monitor service
go run monitor/main.go

# Check metrics endpoint
curl http://localhost:8080/metrics | grep custom
```

## Custom Detection Logic

### Implement Detector Interface

```go
type CustomDetector struct {
    name   string
    region string
    url    string
}

func (d *CustomDetector) GetName() string {
    return d.name
}

func (d *CustomDetector) GetRegion() string {
    return d.region
}

func (d *CustomDetector) Detect(ctx context.Context) (*Result, error) {
    start := time.Now()
    
    // Create HTTP client
    client := &http.Client{
        Timeout: 30 * time.Second,
    }
    
    // Send request
    req, err := http.NewRequestWithContext(ctx, "GET", d.url, nil)
    if err != nil {
        return &Result{
            Service:   d.name,
            Status:    "error",
            Region:    d.region,
            Error:     err,
            Duration:  time.Since(start),
            Timestamp: time.Now(),
        }, nil
    }
    
    // Set request headers
    req.Header.Set("User-Agent", "MediaUnlockTest/1.0")
    
    // Execute request
    resp, err := client.Do(req)
    if err != nil {
        return &Result{
            Service:   d.name,
            Status:    "error",
            Region:    d.region,
            Error:     err,
            Duration:  time.Since(start),
            Timestamp: time.Now(),
        }, nil
    }
    defer resp.Body.Close()
    
    // Determine status
    var status string
    if resp.StatusCode == 200 {
        status = "success"
    } else {
        status = "failed"
    }
    
    return &Result{
        Service:   d.name,
        Status:    status,
        Region:    d.region,
        Error:     nil,
        Duration:  time.Since(start),
        Timestamp: time.Now(),
    }, nil
}
```

### Register Detector

```go
func RegisterDetector(detector Detector) {
    // Add to detector list
    detectors = append(detectors, detector)
}

// Usage example
func init() {
    RegisterDetector(&CustomDetector{
        name:   "CustomService",
        region: "Global",
        url:    "https://api.customservice.com/status",
    })
}
```

### Configure Detector

```go
type DetectorConfig struct {
    Name     string `yaml:"name"`
    Region   string `yaml:"region"`
    URL      string `yaml:"url"`
    Enabled  bool   `yaml:"enabled"`
    Interval int    `yaml:"interval"`
}

func LoadDetectorConfig(configPath string) ([]DetectorConfig, error) {
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, err
    }
    
    var configs []DetectorConfig
    err = yaml.Unmarshal(data, &configs)
    if err != nil {
        return nil, err
    }
    
    return configs, nil
}
```

## Performance Optimization

### Concurrency Control

#### Worker Pool Implementation

```go
type WorkerPool struct {
    workers    int
    jobQueue   chan Job
    workerPool chan chan Job
    quit       chan bool
}

type Job struct {
    Detector Detector
    Result   chan *Result
}

func NewWorkerPool(workers int) *WorkerPool {
    return &WorkerPool{
        workers:    workers,
        jobQueue:   make(chan Job, 100),
        workerPool: make(chan chan Job, workers),
        quit:       make(chan bool),
    }
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.workers; i++ {
        go wp.worker()
    }
    
    go wp.dispatcher()
}

func (wp *WorkerPool) worker() {
    workerQueue := make(chan Job)
    
    for {
        wp.workerPool <- workerQueue
        
        select {
        case job := <-workerQueue:
            result, err := job.Detector.Detect(context.Background())
            if err != nil {
                result = &Result{
                    Service: job.Detector.GetName(),
                    Status:  "error",
                    Error:   err,
                }
            }
            job.Result <- result
            
        case <-wp.quit:
            return
        }
    }
}
```

### Memory Management

#### Object Pool

```go
var resultPool = sync.Pool{
    New: func() interface{} {
        return &Result{}
    },
}

func getResult() *Result {
    return resultPool.Get().(*Result)
}

func putResult(r *Result) {
    // Reset fields
    r.Service = ""
    r.Status = ""
    r.Region = ""
    r.Error = nil
    r.Duration = 0
    r.Timestamp = time.Time{}
    
    resultPool.Put(r)
}
```

#### Periodic Cleanup

```go
func (m *Monitor) cleanup() {
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            // Clean up expired result cache
            m.cleanupExpiredResults()
            
            // Clean up expired metrics
            m.cleanupExpiredMetrics()
        }
    }
}
```

### Network Optimization

#### Connection Reuse

```go
var httpClient = &http.Client{
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
        DisableCompression:  true,
    },
    Timeout: 30 * time.Second,
}
```

#### Request Retry

```go
func (d *CustomDetector) DetectWithRetry(ctx context.Context, maxRetries int) (*Result, error) {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        result, err := d.Detect(ctx)
        if err == nil {
            return result, nil
        }
        
        lastErr = err
        
        // Exponential backoff
        backoff := time.Duration(1<<uint(i)) * time.Second
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        case <-time.After(backoff):
            continue
        }
    }
    
    return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}
```

## Testing and Debugging

### Unit Testing

```go
func TestCustomDetector(t *testing.T) {
    detector := &CustomDetector{
        name:   "TestService",
        region: "Test",
        url:    "https://httpbin.org/status/200",
    }
    
    ctx := context.Background()
    result, err := detector.Detect(ctx)
    
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    
    if result.Service != "TestService" {
        t.Errorf("Expected service name 'TestService', got %s", result.Service)
    }
    
    if result.Status != "success" {
        t.Errorf("Expected status 'success', got %s", result.Status)
    }
}
```

### Integration Testing

```go
func TestMonitorIntegration(t *testing.T) {
    // Start test server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/metrics" {
            w.Header().Set("Content-Type", "text/plain")
            w.Write([]byte("# Test metrics"))
        }
    }))
    defer server.Close()
    
    // Create monitor
    monitor := NewMonitor(&Config{
        ListenAddr: ":0",
        Interval:   1,
    })
    
    // Start monitor
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    go func() {
        if err := monitor.Start(ctx); err != nil {
            t.Logf("Monitor stopped: %v", err)
        }
    }()
    
    // Wait for startup
    time.Sleep(100 * time.Millisecond)
    
    // Check status
    status := monitor.GetStatus()
    if !status.Running {
        t.Error("Expected monitor to be running")
    }
}
```

### Performance Testing

```go
func BenchmarkDetector(b *testing.B) {
    detector := &CustomDetector{
        name:   "BenchmarkService",
        region: "Benchmark",
        url:    "https://httpbin.org/delay/1",
    }
    
    ctx := context.Background()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := detector.Detect(ctx)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

### Debugging Tips

#### Logging

```go
import "log"

func (d *CustomDetector) Detect(ctx context.Context) (*Result, error) {
    log.Printf("Starting detection for %s in %s", d.name, d.region)
    
    start := time.Now()
    result, err := d.Detect(ctx)
    
    log.Printf("Detection completed for %s in %s, took %v, status: %s",
        d.name, d.region, time.Since(start), result.Status)
    
    return result, err
}
```

#### Metrics Debugging

```go
func debugMetrics() {
    // Print current metric values
    mfs, err := prometheus.DefaultGatherer.Gather()
    if err != nil {
        log.Printf("Error gathering metrics: %v", err)
        return
    }
    
    for _, mf := range mfs {
        log.Printf("Metric: %s", mf.GetName())
        for _, m := range mf.GetMetric() {
            log.Printf("  Labels: %v, Value: %v", m.GetLabel(), m.GetGauge().GetValue())
        }
    }
}
```

## Configuration Management

### Configuration File Structure

```yaml
# config.yaml
monitor:
  listen_addr: ":8080"
  interval: 300
  timeout: 30
  concurrent: 10
  
  detectors:
    - name: "Netflix"
      region: "US"
      url: "https://api.netflix.com/status"
      enabled: true
      interval: 300
      
    - name: "Disney+"
      region: "Global"
      url: "https://api.disneyplus.com/status"
      enabled: true
      interval: 300

prometheus:
  enabled: true
  path: "/metrics"
  labels:
    instance: "monitor-1"
    environment: "production"

logging:
  level: "info"
  format: "json"
  file: "/var/log/monitor.log"
```

### Configuration Loading

```go
type Config struct {
    Monitor   MonitorConfig   `yaml:"monitor"`
    Prometheus PrometheusConfig `yaml:"prometheus"`
    Logging   LoggingConfig   `yaml:"logging"`
}

func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    
    var config Config
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        return nil, err
    }
    
    // Set default values
    if config.Monitor.Interval == 0 {
        config.Monitor.Interval = 300
    }
    
    if config.Monitor.Timeout == 0 {
        config.Monitor.Timeout = 30
    }
    
    return &config, nil
}
```

## Deployment and Operations

### Docker Deployment

```dockerfile
FROM golang:1.19-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o monitor ./monitor

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/monitor .
COPY --from=builder /app/config.yaml .

EXPOSE 8080
CMD ["./monitor", "-config", "config.yaml"]
```

### System Service

```ini
# /etc/systemd/system/mediaunlock-monitor.service
[Unit]
Description=MediaUnlockTest Monitor
After=network.target

[Service]
Type=simple
User=monitor
WorkingDirectory=/opt/monitor
ExecStart=/opt/monitor/monitor -config /opt/monitor/config.yaml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

### Monitoring and Alerting

#### Health Check

```go
func (s *MonitorService) handleHealth(w http.ResponseWriter, r *http.Request) {
    status := s.monitor.GetStatus()
    
    if !status.Running {
        http.Error(w, "Service not running", http.StatusServiceUnavailable)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":    "healthy",
        "running":   status.Running,
        "startTime": status.StartTime,
        "testCount": status.TestCount,
    })
}
```

#### Alert Rules

```yaml
# prometheus/alerts.yml
groups:
  - name: mediaunlock
    rules:
      - alert: MonitorDown
        expr: up{job="mediaunlock"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "MediaUnlock monitor is down"
          
      - alert: HighFailureRate
        expr: rate(mediaunlock_test_failed[5m]) / rate(mediaunlock_test_total[5m]) > 0.1
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "High test failure rate detected"
```

## Contribution Guidelines

### Development Process

1. **Fork Project**: Fork this project on GitHub
2. **Create Branch**: Create feature branch `feature/new-detector`
3. **Write Code**: Write code according to development standards
4. **Add Tests**: Add unit tests and integration tests for new features
5. **Update Documentation**: Update related documentation and configuration examples
6. **Submit PR**: Create Pull Request and describe changes

### Code Review Points

1. **Function Completeness**: Is detection logic complete
2. **Error Handling**: Are various error situations handled correctly
3. **Performance Considerations**: Are concurrency, timeout, and resource management considered
4. **Metric Design**: Are monitoring metrics reasonable and meaningful
5. **Test Coverage**: Is there sufficient test coverage
6. **Documentation Updates**: Are related documents updated

### Issue Feedback

If you encounter problems or have suggestions, please:

1. Search GitHub Issues to see if similar issues already exist
2. Create new Issue with detailed problem description and reproduction steps
3. Provide system environment, Go version and other detailed information
4. If possible, provide related logs or error information

## Update Log

### v1.0.0
- Initial version release
- Support basic monitoring functionality

### v1.1.0
- Add Prometheus metrics export
- Support Grafana dashboards
- Optimize performance

### v1.2.0
- Add alert functionality
- Support custom detection logic
- Improve configuration management

---

Thank you for contributing to the MediaUnlockTest Monitor project! If you have questions, please check [Issues](https://github.com/HsukqiLee/MediaUnlockTest/issues) or create new discussions.

## Related Documentation

- [Monitor Usage Guide](README_en.md)
- [Monitor Chinese Usage Guide](README.md)
- [CLI Development Guide](../DEVELOPMENT_EN.md)
- [CLI Chinese Development Guide](../DEVELOPMENT.md)
