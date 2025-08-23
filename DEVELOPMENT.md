# MediaUnlockTest CLI 开发文档

**中文开发文档** | [English Development Guide](https://github.com/HsukqiLee/MediaUnlockTest/blob/main/DEVELOPMENT_EN.md)

## 目录

- [项目结构](#项目结构)
- [CLI 开发规范](#cli-开发规范)
- [函数使用指南](#函数使用指南)
- [添加新的检测项目](#添加新的检测项目)
- [CLI 功能扩展](#cli-功能扩展)
- [测试和调试](#测试和调试)
- [贡献指南](#贡献指南)

## 项目结构

```
MediaUnlockTest/
├── checks/           # 检测函数目录
│   ├── mediaunlock.go    # 核心接口和类型定义
│   ├── Netflix.go        # Netflix 检测实现
│   ├── DisneyPlus.go     # Disney+ 检测实现
│   └── ...               # 其他检测项目
├── cli/              # 命令行工具
│   ├── main.go       # 主程序入口
│   ├── build.bat     # Windows 构建脚本
│   └── build.sh      # Unix 构建脚本
└── go.mod            # Go 模块文件
```

## CLI 开发规范

### 文件命名规范

1. **检测函数文件**: 使用 PascalCase 命名，以 `.go` 结尾
   - 正确: `Netflix.go`, `DisneyPlus.go`, `BBCiPlayer.go`
   - 错误: `netflix.go`, `disney_plus.go`, `bbc_iplayer.go`

2. **函数名**: 使用 PascalCase 命名
   - 正确: `NetflixRegion()`, `DisneyPlus()`, `BBCiPlayer()`
   - 错误: `netflix_region()`, `disneyPlus()`, `bbciplayer()`

3. **变量名**: 使用 camelCase 命名
   - 正确: `httpClient`, `userAgent`, `timeout`
   - 错误: `http_client`, `UserAgent`, `TIMEOUT`

### CLI 代码结构规范

1. **包声明**: 每个文件必须以 `package main` 开头（cli目录）
2. **导入顺序**: 标准库 → 第三方库 → 本地包
3. **函数顺序**: main函数 → 公共函数 → 私有函数 → 辅助函数
4. **注释**: 每个公共函数必须有注释说明
5. **错误处理**: 使用统一的错误处理方式

### 命令行参数规范

1. **参数定义**: 使用 `flag` 包定义命令行参数
2. **参数验证**: 对关键参数进行有效性验证
3. **默认值**: 为所有参数提供合理的默认值
4. **帮助信息**: 提供清晰的参数说明

```go
// 示例：参数定义
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

// 参数绑定
flag.StringVar(&Interface, "I", "", "Source IP or network interface to use for connections")
flag.StringVar(&DNSServers, "dns-servers", "", "Custom DNS servers (format: ip:port)")
flag.BoolVar(&Debug, "debug", false, "Enable debug mode for verbose output")
flag.IntVar(&IPMode, "m", 0, "Connection mode: 0=auto (default), 4=IPv4 only, 6=IPv6 only")
```

## 函数使用指南

### 核心类型和接口

#### Result 结构体

```go
type Result struct {
    Status       Status      // 检测状态
    Region       string      // 地区信息
    Info         string      // 额外信息
    Err          error       // 错误信息
    CachedResult bool        // 是否为缓存结果
}
```

#### Status 枚举

```go
const (
    StatusOK           Status = iota // 解锁成功
    StatusNo                         // 未解锁
    StatusRestricted                 // 地区限制
    StatusNetworkErr                 // 网络错误
    StatusErr                        // 其他错误
    StatusBanned                     // 被禁止
    StatusUnexpected                 // 意外结果
    StatusFailed                     // 测试失败
)
```

### HTTP 客户端使用

#### 自动选择客户端

```go
func NetflixRegion(client http.Client) m.Result {
    // 使用传入的 client 参数
    resp, err := client.Get("https://api.netflix.com/status")
    if err != nil {
        return m.Result{
            Status: m.StatusNetworkErr,
            Err:    err,
        }
    }
    // ... 处理响应
}
```

#### 强制使用特定协议

```go
// IPv4 客户端
resp, err := m.Ipv4HttpClient.Get(url)

// IPv6 客户端  
resp, err := m.Ipv6HttpClient.Get(url)

// 自动客户端
resp, err := m.AutoHttpClient.Get(url)
```

### CLI 工具函数

#### 进度条管理

```go
import pb "github.com/schollz/progressbar/v3"

func NewBar(count int64) *pb.ProgressBar {
    return pb.NewOptions64(
        count,
        pb.OptionSetDescription("正在测试..."),
        pb.OptionSetWriter(os.Stderr),
        pb.OptionSetWidth(60),
        pb.OptionThrottle(200*time.Millisecond),
        pb.OptionShowCount(),
        pb.OptionClearOnFinish(),
        pb.OptionEnableColorCodes(true),
    )
}
```

#### 颜色输出

```go
import "github.com/fatih/color"

var (
    Red     = color.New(color.FgRed).SprintFunc()
    Green   = color.New(color.FgGreen).SprintFunc()
    Yellow  = color.New(color.FgYellow).SprintFunc()
    Blue    = color.New(color.FgBlue).SprintFunc()
    SkyBlue = color.New(color.FgCyan).SprintFunc()
)

// 使用示例
fmt.Println(Green("成功"), Red("错误"), Yellow("警告"))
```

## 添加新的检测项目

### 步骤 1: 创建检测函数

在 `checks/` 目录下创建新文件，例如 `NewService.go`:

```go
package checks

import (
    "net/http"
    "strings"
)

// NewService 检测新服务的解锁状态
func NewService(client http.Client) m.Result {
    // 1. 发送请求
    resp, err := client.Get("https://api.newservice.com/status")
    if err != nil {
        return m.Result{
            Status: m.StatusNetworkErr,
            Err:    err,
        }
    }
    defer resp.Body.Close()

    // 2. 检查状态码
    if resp.StatusCode != 200 {
        return m.Result{
            Status: m.StatusErr,
            Err:    fmt.Errorf("unexpected status code: %d", resp.StatusCode),
        }
    }

    // 3. 读取响应内容
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return m.Result{
            Status: m.StatusErr,
            Err:    err,
        }
    }

    // 4. 解析响应
    content := string(body)
    
    // 5. 判断解锁状态
    if strings.Contains(content, "unlocked") {
        return m.Result{
            Status: m.StatusOK,
            Region: "US", // 如果知道具体地区
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

### 步骤 2: 在 main.go 中注册

在 `cli/main.go` 中找到对应的测试数组，添加新项目：

```go
var NorthAmericaTests = []testItem{
    // ... 现有项目
    {"New Service", m.NewService, false}, // 最后一个参数表示是否支持 IPv6
}
```

### 步骤 3: 测试验证

```bash
# 运行测试模式
go run cli/main.go -test

# 或者直接测试特定项目
go run cli/main.go -nf
```

## CLI 功能扩展

### 添加新的命令行参数

1. **定义变量**：在 `main()` 函数前定义新的参数变量
2. **绑定参数**：使用 `flag` 包绑定参数
3. **参数处理**：在 `main()` 函数中添加参数处理逻辑

```go
// 1. 定义新参数
var (
    NewParam string
    NewFlag  bool
)

// 2. 绑定参数
flag.StringVar(&NewParam, "new-param", "", "Description of new parameter")
flag.BoolVar(&NewFlag, "new-flag", false, "Description of new flag")

// 3. 处理参数
if NewFlag {
    // 处理新标志的逻辑
}
```

### 添加新的检测模式

1. **定义测试数组**：创建新的测试项目数组
2. **添加选择逻辑**：在 `ReadSelect()` 函数中添加新的选择项
3. **执行逻辑**：在 `main()` 函数中添加执行逻辑

```go
// 1. 定义新的测试数组
var NewRegionTests = []testItem{
    {"Service 1", m.Service1, false},
    {"Service 2", m.Service2, true},
}

// 2. 在 ReadSelect() 中添加选择项
fmt.Println("[12] :   新地区平台")

// 3. 处理选择
case "12":
    NewRegion = true

// 4. 在 regions 数组中添加
{Enabled: NewRegion, Name: "NewRegion", Tests: NewRegionTests}
```

## 检测逻辑规范

### 状态判断优先级

1. **网络错误优先**: 如果无法连接，返回 `StatusNetworkErr`
2. **状态码检查**: 检查 HTTP 状态码是否正常
3. **内容解析**: 根据响应内容判断解锁状态
4. **默认处理**: 无法判断时返回 `StatusNo`

### 地区信息处理

```go
// 如果知道具体地区
if strings.Contains(content, "US") {
    return m.Result{
        Status: m.StatusOK,
        Region: "US",
    }
}

// 如果只知道解锁但不知道地区
if strings.Contains(content, "unlocked") {
    return m.Result{
        Status: m.StatusOK,
        // Region 字段留空
    }
}
```

### 错误信息处理

```go
// 网络错误
if err != nil {
    return m.Result{
        Status: m.StatusNetworkErr,
        Err:    err, // 保留原始错误
    }
}

// 业务逻辑错误
if !isValidResponse(content) {
    return m.Result{
        Status: m.StatusErr,
        Err:    fmt.Errorf("invalid response format"),
        Info:   "Response does not match expected format",
    }
}
```

## 测试和调试

### 本地测试

```bash
# 进入项目目录
cd MediaUnlockTest

# 运行所有测试
go run cli/main.go

# 仅测试特定地区
go run cli/main.go -m 4  # 仅 IPv4

# 开启调试模式
go run cli/main.go -debug

# 限制并发数量
go run cli/main.go -conc 10

# 使用代理测试
go run cli/main.go -http-proxy "http://127.0.0.1:8080"
```

### 调试技巧

1. **使用 debug 模式**: 添加 `-debug` 参数查看详细错误信息
2. **检查网络连接**: 确认目标服务是否可访问
3. **查看响应内容**: 在代码中添加日志输出响应内容
4. **使用代理测试**: 通过 `-http-proxy` 或 `-socks-proxy` 测试代理环境
5. **性能分析**: 使用 `-conc` 参数调整并发数量

### 常见问题排查

1. **超时错误**: 检查网络连接和防火墙设置
2. **解析错误**: 检查响应格式是否符合预期
3. **地区判断错误**: 确认检测逻辑是否正确
4. **并发问题**: 使用 `-conc` 参数限制并发数量
5. **代理问题**: 检查代理配置是否正确

## 构建和部署

### 构建脚本

项目提供了跨平台的构建脚本：

```bash
# Windows
cli/build.bat

# Unix/Linux/macOS
cli/build.sh
```

### 构建选项

```bash
# 基本构建
go build -o unlock-test cli/main.go

# 跨平台构建
GOOS=linux GOARCH=amd64 go build -o unlock-test-linux cli/main.go
GOOS=windows GOARCH=amd64 go build -o unlock-test.exe cli/main.go
GOOS=darwin GOARCH=amd64 go build -o unlock-test-mac cli/main.go
```

## 贡献指南

### 提交 Pull Request

1. **Fork 项目**: 在 GitHub 上 Fork 本项目
2. **创建分支**: 创建功能分支 `feature/new-service`
3. **编写代码**: 按照开发规范编写代码
4. **测试验证**: 确保代码能正常运行
5. **提交 PR**: 创建 Pull Request 并描述变更

### 代码审查要点

1. **功能完整性**: 检测逻辑是否完整
2. **错误处理**: 是否正确处理各种错误情况
3. **代码规范**: 是否符合项目的编码规范
4. **性能考虑**: 是否考虑了超时和并发
5. **文档更新**: 是否更新了相关文档
6. **CLI 体验**: 新功能是否提供了良好的用户体验

### 问题反馈

如果遇到问题或有建议，请：

1. 在 GitHub Issues 中搜索是否已有相关问题
2. 创建新的 Issue，详细描述问题和复现步骤
3. 提供系统环境、Go 版本等详细信息
4. 如果可能，提供相关的日志或错误信息

## 更新日志

### v1.0.0
- 初始版本发布
- 支持基本的流媒体检测功能

### v1.1.0
- 添加 IPv6 支持
- 优化并发处理
- 增加缓存机制

### v1.2.0
- 改进进度条显示
- 添加活动测试显示
- 优化错误处理

---

感谢您为 MediaUnlockTest CLI 项目做出贡献！如有疑问，请查看 [Issues](https://github.com/HsukqiLee/MediaUnlockTest/issues) 或创建新的讨论。

## 相关文档

- [Monitor 开发文档](https://github.com/HsukqiLee/MediaUnlockTest/blob/main/monitor/README.md)
- [Monitor 英文文档](https://github.com/HsukqiLee/MediaUnlockTest/blob/main/monitor/README_en.md)
