# MediaUnlockTest Monitor

**English Docs** | [中文文档](https://github.com/HsukqiLee/MediaUnlockTest/blob/main/monitor/README.md)

> Streaming media unlock monitoring tool based on Prometheus and Grafana

## Features

- 🔍 **Real-time Monitoring**: Continuously monitor streaming media service unlock status
- 📊 **Data Visualization**: Rich charts and dashboards through Grafana
- ⚡ **High Performance**: Support high-concurrency detection and fast response
- 🔧 **Easy Deployment**: Provide Docker and manual deployment solutions
- 📈 **Historical Data**: Save historical detection data, support trend analysis
- 🚨 **Alert Notifications**: Support multiple alert methods

## Quick Start

### Using Docker Compose (Recommended)

1. **Clone Project**
```bash
git clone https://github.com/HsukqiLee/MediaUnlockTest.git
cd MediaUnlockTest/monitor
```

2. **Start Services**
```bash
docker-compose up -d
```

3. **Access Services**
- Grafana: http://localhost:3000 (default account: admin/admin)
- Prometheus: http://localhost:9090
- MediaUnlockTest Monitor: http://localhost:8080

### Manual Deployment

#### 1. Install Dependencies

```bash
# Install Go 1.19+
go version

# Install dependencies
go mod download
```

#### 2. Build Monitor Tool

```bash
# Windows
monitor/build.bat

# Unix/Linux/macOS
monitor/build.sh
```

#### 3. Configure Prometheus

Create `prometheus.yml` configuration file:

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'mediaunlock'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 30s
```

#### 4. Start Services

```bash
# Start Prometheus
./prometheus --config.file=prometheus.yml

# Start MediaUnlockTest Monitor
./monitor

# Start Grafana
./grafana-server
```

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `MONITOR_PORT` | `8080` | Monitor service port |
| `MONITOR_INTERVAL` | `300` | Detection interval (seconds) |
| `MONITOR_TIMEOUT` | `30` | Single detection timeout (seconds) |
| `MONITOR_CONCURRENT` | `10` | Concurrent detection count |
| `PROMETHEUS_ENABLED` | `true` | Whether to enable Prometheus metrics |
| `LOG_LEVEL` | `info` | Log level |

### Configuration File

Create `config.yaml` configuration file:

```yaml
monitor:
  port: 8080
  interval: 300
  timeout: 30
  concurrent: 10
  
  services:
    - name: "Netflix"
      enabled: true
      region: "US"
    - name: "Disney+"
      enabled: true
      region: "Global"
    - name: "BBC iPlayer"
      enabled: true
      region: "UK"

prometheus:
  enabled: true
  path: "/metrics"

logging:
  level: "info"
  format: "json"
```

## Monitoring Metrics

### Core Metrics

- `mediaunlock_test_total`: Total detection count
- `mediaunlock_test_success`: Successful detection count
- `mediaunlock_test_failed`: Failed detection count
- `mediaunlock_service_status`: Service status (0=failed, 1=success)
- `mediaunlock_test_duration_seconds`: Detection duration

### Custom Labels

- `service`: Service name
- `region`: Region
- `status`: Status
- `error_type`: Error type

### Example Queries

```promql
# Success rate
rate(mediaunlock_test_success[5m]) / rate(mediaunlock_test_total[5m])

# Average response time
histogram_quantile(0.95, rate(mediaunlock_test_duration_seconds_bucket[5m]))

# Service status
mediaunlock_service_status
```

## Grafana Dashboards

### Default Dashboard

The project provides pre-configured Grafana dashboards, including:

- 📊 **Overview Panel**: Display overall unlock status and success rate
- 📈 **Trend Charts**: Show unlock status changes over time
- 🎯 **Service Details**: Detailed status and error information for each service
- ⚡ **Performance Monitoring**: Detection duration and concurrency statistics
- 🚨 **Alert Panel**: Display current alert status

### Import Dashboard

1. Click "+" → "Import" in Grafana
2. Upload `grafana-dashboard.json` file
3. Select Prometheus data source
4. Click "Import" to complete

### Custom Dashboard

You can create custom dashboards as needed:

```json
{
  "dashboard": {
    "title": "Custom MediaUnlock Dashboard",
    "panels": [
      {
        "title": "Service Status",
        "type": "stat",
        "targets": [
          {
            "expr": "mediaunlock_service_status",
            "legendFormat": "{{service}} - {{region}}"
          }
        ]
      }
    ]
  }
}
```

## Alert Configuration

### Prometheus Alert Rules

Create `alerts.yml` file:

```yaml
groups:
  - name: mediaunlock
    rules:
      - alert: ServiceDown
        expr: mediaunlock_service_status == 0
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Service {{ $labels.service }} is down"
          description: "Service {{ $labels.service }} in {{ $labels.region }} has been down for more than 5 minutes"

      - alert: HighFailureRate
        expr: rate(mediaunlock_test_failed[5m]) / rate(mediaunlock_test_total[5m]) > 0.1
        for: 2m
        labels:
          severity: critical
        annotations:
          summary: "High failure rate detected"
          description: "Failure rate is {{ $value | humanizePercentage }}"
```

### Alert Notifications

Support multiple alert notification methods:

- **Email**: Configure SMTP server
- **Slack**: Send to Slack channel
- **Webhook**: Custom HTTP callback
- **DingTalk/WeChat Work**: Common communication tools in China

## Development Guide

### Project Structure

```
monitor/
├── main.go           # Main program entry
├── monitor.go        # Monitor core logic
├── exporter.go       # Prometheus metrics export
├── service.go        # Service management
├── update.go         # Auto update
├── build.bat         # Windows build script
├── build.sh          # Unix build script
└── README.md         # Documentation
```

### Adding New Monitoring Metrics

1. **Define Metrics**
```go
import "github.com/prometheus/client_golang/prometheus"

var (
    customMetric = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "mediaunlock_custom_total",
            Help: "Custom metric description",
        },
        []string{"label1", "label2"},
    )
)
```

2. **Register Metrics**
```go
func init() {
    prometheus.MustRegister(customMetric)
}
```

3. **Update Metrics**
```go
customMetric.WithLabelValues("value1", "value2").Inc()
```

### Custom Detection Logic

1. **Implement Detection Interface**
```go
type Detector interface {
    Detect(ctx context.Context) (*Result, error)
}

type CustomDetector struct {
    // Custom fields
}

func (d *CustomDetector) Detect(ctx context.Context) (*Result, error) {
    // Implement detection logic
    return &Result{
        Service: "CustomService",
        Status:  "success",
        Region:  "Global",
    }, nil
}
```

2. **Register Detector**
```go
func RegisterDetector(name string, detector Detector) {
    // Registration logic
}
```

## Performance Optimization

### Concurrency Control

- Use worker pool to control concurrency count
- Implement timeout mechanism to avoid long blocking
- Support graceful shutdown

### Memory Management

- Regularly clean expired data
- Use object pool to reduce GC pressure
- Limit historical data storage

### Network Optimization

- Support HTTP/2
- Connection reuse
- Request retry mechanism

## Troubleshooting

### Common Issues

1. **Prometheus Cannot Scrape Metrics**
   - Check if port is open
   - Confirm `/metrics` path is accessible
   - Check firewall settings

2. **Grafana Shows No Data**
   - Check data source configuration
   - Confirm time range settings
   - Check Prometheus query syntax

3. **High Detection Failure Rate**
   - Check network connection
   - Adjust timeout settings
   - Check target service status

### Log Analysis

```bash
# View real-time logs
tail -f monitor.log

# Search error logs
grep "ERROR" monitor.log

# Analyze performance logs
grep "duration" monitor.log | awk '{print $NF}' | sort -n
```

### Performance Tuning

```bash
# Adjust concurrency count
export MONITOR_CONCURRENT=20

# Adjust detection interval
export MONITOR_INTERVAL=60

# Enable debug mode
export LOG_LEVEL=debug
```

## Deployment Recommendations

### Production Environment

- Use reverse proxy (Nginx/Traefik)
- Configure SSL certificates
- Set up monitoring and alerts
- Regular data backup

### High Availability Deployment

- Multi-instance deployment
- Load balancing
- Database clustering
- Monitoring redundancy

### Security Considerations

- Network isolation
- Access control
- Log auditing
- Regular updates

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

## Related Links

- [CLI Development Guide](../DEVELOPMENT_en.md)
- [Monitor Development Guide](./DEVELOPMENT_en.md)
- [Project Homepage](https://github.com/HsukqiLee/MediaUnlockTest)
- [Issue Feedback](https://github.com/HsukqiLee/MediaUnlockTest/issues)

## Contribution Guidelines

Welcome to submit Issues and Pull Requests! Please refer to [Development Guide](./DEVELOPMENT_en.md) for development standards.