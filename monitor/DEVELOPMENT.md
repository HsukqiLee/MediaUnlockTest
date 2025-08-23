# MediaUnlockTest Monitor 开发文档

**中文开发文档** | [English Development Guide](https://github.com/HsukqiLee/MediaUnlockTest/blob/main/monitor/DEVELOPMENT_EN.md)

## 目录

- [项目结构](#项目结构)
- [开发规范](#开发规范)
- [核心组件](#核心组件)
- [添加新的监控指标](#添加新的监控指标)
- [自定义检测逻辑](#自定义检测逻辑)
- [性能优化](#性能优化)
- [测试和调试](#测试和调试)
- [贡献指南](#贡献指南)

## 项目结构

```
monitor/
├── main.go           # 主程序入口
├── monitor.go        # 监控核心逻辑
├── exporter.go       # Prometheus 指标导出
├── service.go        # 服务管理
├── update.go         # 自动更新
├── build.bat         # Windows 构建脚本
├── build.sh          # Unix 构建脚本
├── README.md         # 使用说明文档
├── DEVELOPMENT.md    # 开发文档（本文件）
└── DEVELOPMENT_EN.md # 英文开发文档
```

## 开发规范

### 文件命名规范

1. **Go 文件**: 使用 snake_case 命名，以 `.go` 结尾
   - 正确: `monitor.go`, `exporter.go`, `service.go`
   - 错误: `Monitor.go`, `Exporter.go`, `monitor.go`

2. **函数名**: 使用 PascalCase 命名
   - 正确: `StartMonitoring()`, `ExportMetrics()`, `HandleRequest()`
   - 错误: `start_monitoring()`, `exportMetrics()`, `handlerequest()`

3. **变量名**: 使用 camelCase 命名
   - 正确: `monitorConfig`, `metricsExporter`, `updateInterval`
   - 错误: `monitor_config`, `MetricsExporter`, `UPDATE_INTERVAL`

### 代码结构规范

1. **包声明**: 每个文件必须以 `package main` 开头
2. **导入顺序**: 标准库 → 第三方库 → 本地包
3. **函数顺序**: main函数 → 公共函数 → 私有函数 → 辅助函数
4. **注释**: 每个公共函数必须有注释说明
5. **错误处理**: 使用统一的错误处理方式

### 监控指标规范

1. **指标命名**: 使用 `mediaunlock_` 前缀，下划线分隔
2. **标签设计**: 标签应该有意义且数量合理
3. **指标类型**: 根据用途选择合适的指标类型（Counter, Gauge, Histogram等）
4. **文档说明**: 每个指标必须有清晰的帮助信息

## 核心组件

### Monitor 核心逻辑

#### 监控器接口

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

#### 检测器接口

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

### Prometheus 指标导出

#### 指标定义

```go
import "github.com/prometheus/client_golang/prometheus"

var (
    // 检测总数
    testTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "mediaunlock_test_total",
            Help: "Total number of tests performed",
        },
        []string{"service", "region"},
    )

    // 检测成功数
    testSuccess = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "mediaunlock_test_success",
            Help: "Total number of successful tests",
        },
        []string{"service", "region"},
    )

    // 检测耗时
    testDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "mediaunlock_test_duration_seconds",
            Help:    "Test duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"service", "region"},
    )

    // 服务状态
    serviceStatus = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "mediaunlock_service_status",
            Help: "Current service status (0=failed, 1=success)",
        },
        []string{"service", "region"},
    )
)
```

#### 指标注册

```go
func init() {
    // 注册所有指标
    prometheus.MustRegister(
        testTotal,
        testSuccess,
        testDuration,
        serviceStatus,
    )
}
```

### 服务管理

#### HTTP 服务

```go
type MonitorService struct {
    monitor Monitor
    server  *http.Server
    config  *Config
}

func (s *MonitorService) Start() error {
    mux := http.NewServeMux()
    
    // 健康检查端点
    mux.HandleFunc("/health", s.handleHealth)
    
    // 指标导出端点
    mux.Handle("/metrics", promhttp.Handler())
    
    // 状态查询端点
    mux.HandleFunc("/status", s.handleStatus)
    
    s.server = &http.Server{
        Addr:    s.config.ListenAddr,
        Handler: mux,
    }
    
    return s.server.ListenAndServe()
}
```

## 添加新的监控指标

### 步骤 1: 定义指标

在 `exporter.go` 文件中定义新的指标：

```go
var (
    // 自定义指标示例
    customMetric = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "mediaunlock_custom_total",
            Help: "Custom metric description",
        },
        []string{"label1", "label2"},
    )
)
```

### 步骤 2: 注册指标

在 `init()` 函数中注册新指标：

```go
func init() {
    prometheus.MustRegister(
        // ... 现有指标
        customMetric,
    )
}
```

### 步骤 3: 更新指标

在适当的地方更新指标值：

```go
// 增加计数
customMetric.WithLabelValues("value1", "value2").Inc()

// 设置特定值
customMetric.WithLabelValues("value1", "value2").Add(5)
```

### 步骤 4: 测试验证

```bash
# 启动监控服务
go run monitor/main.go

# 检查指标端点
curl http://localhost:8080/metrics | grep custom
```

## 自定义检测逻辑

### 实现检测器接口

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
    
    // 创建 HTTP 客户端
    client := &http.Client{
        Timeout: 30 * time.Second,
    }
    
    // 发送请求
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
    
    // 设置请求头
    req.Header.Set("User-Agent", "MediaUnlockTest/1.0")
    
    // 执行请求
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
    
    // 判断状态
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

### 注册检测器

```go
func RegisterDetector(detector Detector) {
    // 添加到检测器列表
    detectors = append(detectors, detector)
}

// 使用示例
func init() {
    RegisterDetector(&CustomDetector{
        name:   "CustomService",
        region: "Global",
        url:    "https://api.customservice.com/status",
    })
}
```

### 配置检测器

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

## 性能优化

### 并发控制

#### 工作池实现

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

### 内存管理

#### 对象池

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
    // 重置字段
    r.Service = ""
    r.Status = ""
    r.Region = ""
    r.Error = nil
    r.Duration = 0
    r.Timestamp = time.Time{}
    
    resultPool.Put(r)
}
```

#### 定期清理

```go
func (m *Monitor) cleanup() {
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            // 清理过期的结果缓存
            m.cleanupExpiredResults()
            
            // 清理过期的指标
            m.cleanupExpiredMetrics()
        }
    }
}
```

### 网络优化

#### 连接复用

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

#### 请求重试

```go
func (d *CustomDetector) DetectWithRetry(ctx context.Context, maxRetries int) (*Result, error) {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        result, err := d.Detect(ctx)
        if err == nil {
            return result, nil
        }
        
        lastErr = err
        
        // 指数退避
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

## 测试和调试

### 单元测试

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

### 集成测试

```go
func TestMonitorIntegration(t *testing.T) {
    // 启动测试服务器
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/metrics" {
            w.Header().Set("Content-Type", "text/plain")
            w.Write([]byte("# Test metrics"))
        }
    }))
    defer server.Close()
    
    // 创建监控器
    monitor := NewMonitor(&Config{
        ListenAddr: ":0",
        Interval:   1,
    })
    
    // 启动监控器
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    go func() {
        if err := monitor.Start(ctx); err != nil {
            t.Logf("Monitor stopped: %v", err)
        }
    }()
    
    // 等待启动
    time.Sleep(100 * time.Millisecond)
    
    // 检查状态
    status := monitor.GetStatus()
    if !status.Running {
        t.Error("Expected monitor to be running")
    }
}
```

### 性能测试

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

### 调试技巧

#### 日志记录

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

#### 指标调试

```go
func debugMetrics() {
    // 打印当前指标值
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

## 配置管理

### 配置文件结构

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

### 配置加载

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
    
    // 设置默认值
    if config.Monitor.Interval == 0 {
        config.Monitor.Interval = 300
    }
    
    if config.Monitor.Timeout == 0 {
        config.Monitor.Timeout = 30
    }
    
    return &config, nil
}
```

## 部署和运维

### Docker 部署

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

### 系统服务

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

### 监控和告警

#### 健康检查

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

#### 告警规则

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

## 贡献指南

### 开发流程

1. **Fork 项目**: 在 GitHub 上 Fork 本项目
2. **创建分支**: 创建功能分支 `feature/new-detector`
3. **编写代码**: 按照开发规范编写代码
4. **添加测试**: 为新功能添加单元测试和集成测试
5. **更新文档**: 更新相关文档和配置示例
6. **提交 PR**: 创建 Pull Request 并描述变更

### 代码审查要点

1. **功能完整性**: 检测逻辑是否完整
2. **错误处理**: 是否正确处理各种错误情况
3. **性能考虑**: 是否考虑了并发、超时和资源管理
4. **指标设计**: 监控指标是否合理和有意义
5. **测试覆盖**: 是否有足够的测试覆盖
6. **文档更新**: 是否更新了相关文档

### 问题反馈

如果遇到问题或有建议，请：

1. 在 GitHub Issues 中搜索是否已有相关问题
2. 创建新的 Issue，详细描述问题和复现步骤
3. 提供系统环境、Go 版本等详细信息
4. 如果可能，提供相关的日志或错误信息

## 更新日志

### v1.0.0
- 初始版本发布
- 支持基本的监控功能

### v1.1.0
- 添加 Prometheus 指标导出
- 支持 Grafana 仪表板
- 优化性能

### v1.2.0
- 添加告警功能
- 支持自定义检测逻辑
- 改进配置管理

---

感谢您为 MediaUnlockTest Monitor 项目做出贡献！如有疑问，请查看 [Issues](https://github.com/HsukqiLee/MediaUnlockTest/issues) 或创建新的讨论。

## 相关文档

- [Monitor 使用文档](README.md)
- [Monitor 英文使用文档](README_en.md)
- [CLI 开发文档](../DEVELOPMENT.md)
- [CLI 英文开发文档](../DEVELOPMENT_EN.md)
