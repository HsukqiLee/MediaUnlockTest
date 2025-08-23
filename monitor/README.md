# MediaUnlockTest Monitor

**中文文档** | [English Docs](https://github.com/HsukqiLee/MediaUnlockTest/blob/main/monitor/README_en.md)

> 基于 Prometheus 和 Grafana 的流媒体解锁监控工具

![](https://raw.githubusercontent.com/HsukqiLee/MediaUnlockTest/main/monitor/image.png)

## 功能特性

- 🔍 **实时监控**: 持续监控流媒体服务的解锁状态
- 📊 **数据可视化**: 通过 Grafana 提供丰富的图表和仪表板
- ⚡ **高性能**: 支持高并发检测和快速响应
- 🔧 **易于部署**: 提供 Docker 和手动部署方案
- 📈 **历史数据**: 保存历史检测数据，支持趋势分析
- 🚨 **告警通知**: 支持多种告警方式

## 快速开始

### 使用 Docker Compose（推荐）

1. **克隆项目**
```bash
git clone https://github.com/HsukqiLee/MediaUnlockTest.git
cd MediaUnlockTest/monitor
```

2. **启动服务**
```bash
docker-compose up -d
```

3. **访问服务**
- Grafana: http://localhost:3000 (默认账号: admin/admin)
- Prometheus: http://localhost:9090
- MediaUnlockTest Monitor: http://localhost:8080

### 手动部署

#### 1. 安装依赖

```bash
# 安装 Go 1.19+
go version

# 安装依赖
go mod download
```

#### 2. 构建监控工具

```bash
# Windows
monitor/build.bat

# Unix/Linux/macOS
monitor/build.sh
```

#### 3. 配置 Prometheus

创建 `prometheus.yml` 配置文件：

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

#### 4. 启动服务

```bash
# 启动 Prometheus
./prometheus --config.file=prometheus.yml

# 启动 MediaUnlockTest Monitor
./monitor

# 启动 Grafana
./grafana-server
```

## 配置说明

### 环境变量

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `MONITOR_PORT` | `8080` | 监控服务端口 |
| `MONITOR_INTERVAL` | `300` | 检测间隔（秒） |
| `MONITOR_TIMEOUT` | `30` | 单个检测超时时间（秒） |
| `MONITOR_CONCURRENT` | `10` | 并发检测数量 |
| `PROMETHEUS_ENABLED` | `true` | 是否启用 Prometheus 指标 |
| `LOG_LEVEL` | `info` | 日志级别 |

### 配置文件

创建 `config.yaml` 配置文件：

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

## 监控指标

### 核心指标

- `mediaunlock_test_total`: 总检测次数
- `mediaunlock_test_success`: 成功检测次数
- `mediaunlock_test_failed`: 失败检测次数
- `mediaunlock_service_status`: 各服务状态（0=失败, 1=成功）
- `mediaunlock_test_duration_seconds`: 检测耗时

### 自定义标签

- `service`: 服务名称
- `region`: 地区
- `status`: 状态
- `error_type`: 错误类型

### 示例查询

```promql
# 成功率
rate(mediaunlock_test_success[5m]) / rate(mediaunlock_test_total[5m])

# 平均响应时间
histogram_quantile(0.95, rate(mediaunlock_test_duration_seconds_bucket[5m]))

# 各服务状态
mediaunlock_service_status
```

## Grafana 仪表板

### 默认仪表板

项目提供了预配置的 Grafana 仪表板，包含：

- 📊 **总览面板**: 显示整体解锁状态和成功率
- 📈 **趋势图表**: 展示解锁状态随时间的变化
- 🎯 **服务详情**: 各服务的详细状态和错误信息
- ⚡ **性能监控**: 检测耗时和并发数统计
- 🚨 **告警面板**: 显示当前告警状态

### 导入仪表板

1. 在 Grafana 中点击 "+" → "Import"
2. 上传 `grafana-dashboard.json` 文件
3. 选择 Prometheus 数据源
4. 点击 "Import" 完成导入

### 自定义仪表板

可以根据需要创建自定义仪表板：

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

## 告警配置

### Prometheus 告警规则

创建 `alerts.yml` 文件：

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

### 告警通知

支持多种告警通知方式：

- **邮件**: 配置 SMTP 服务器
- **Slack**: 发送到 Slack 频道
- **Webhook**: 自定义 HTTP 回调
- **钉钉/企业微信**: 国内常用通讯工具

## 开发指南

### 项目结构

```
monitor/
├── main.go           # 主程序入口
├── monitor.go        # 监控核心逻辑
├── exporter.go       # Prometheus 指标导出
├── service.go        # 服务管理
├── update.go         # 自动更新
├── build.bat         # Windows 构建脚本
├── build.sh          # Unix 构建脚本
└── README.md         # 说明文档
```

### 添加新的监控指标

1. **定义指标**
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

2. **注册指标**
```go
func init() {
    prometheus.MustRegister(customMetric)
}
```

3. **更新指标**
```go
customMetric.WithLabelValues("value1", "value2").Inc()
```

### 自定义检测逻辑

1. **实现检测接口**
```go
type Detector interface {
    Detect(ctx context.Context) (*Result, error)
}

type CustomDetector struct {
    // 自定义字段
}

func (d *CustomDetector) Detect(ctx context.Context) (*Result, error) {
    // 实现检测逻辑
    return &Result{
        Service: "CustomService",
        Status:  "success",
        Region:  "Global",
    }, nil
}
```

2. **注册检测器**
```go
func RegisterDetector(name string, detector Detector) {
    // 注册逻辑
}
```

## 性能优化

### 并发控制

- 使用工作池控制并发数量
- 实现超时机制避免长时间阻塞
- 支持优雅关闭

### 内存管理

- 定期清理过期数据
- 使用对象池减少 GC 压力
- 限制历史数据存储量

### 网络优化

- 支持 HTTP/2
- 连接复用
- 请求重试机制

## 故障排除

### 常见问题

1. **Prometheus 无法抓取指标**
   - 检查端口是否开放
   - 确认 `/metrics` 路径可访问
   - 查看防火墙设置

2. **Grafana 显示无数据**
   - 检查数据源配置
   - 确认时间范围设置
   - 查看 Prometheus 查询语法

3. **检测失败率高**
   - 检查网络连接
   - 调整超时设置
   - 查看目标服务状态

### 日志分析

```bash
# 查看实时日志
tail -f monitor.log

# 搜索错误日志
grep "ERROR" monitor.log

# 分析性能日志
grep "duration" monitor.log | awk '{print $NF}' | sort -n
```

### 性能调优

```bash
# 调整并发数
export MONITOR_CONCURRENT=20

# 调整检测间隔
export MONITOR_INTERVAL=60

# 启用调试模式
export LOG_LEVEL=debug
```

## 部署建议

### 生产环境

- 使用反向代理（Nginx/Traefik）
- 配置 SSL 证书
- 设置监控和告警
- 定期备份数据

### 高可用部署

- 多实例部署
- 负载均衡
- 数据库集群
- 监控冗余

### 安全考虑

- 网络隔离
- 访问控制
- 日志审计
- 定期更新

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

## 相关链接

- [CLI 开发文档](../DEVELOPMENT.md)
- [Monitor 开发文档](./DEVELOPMENT.md)
- [项目主页](https://github.com/HsukqiLee/MediaUnlockTest)
- [问题反馈](https://github.com/HsukqiLee/MediaUnlockTest/issues)

## 贡献指南

欢迎提交 Issue 和 Pull Request！请参考 [开发文档](./DEVELOPMENT.md) 了解开发规范。
