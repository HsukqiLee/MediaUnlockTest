# 基础使用指南

## 快速开始

**Linux (包括 iOS iSH) / macOS / Android Termux**:
(分别需要 /usr/bin、/usr/local/bin、$PREFIX/bin 目录的读写权限)

```bash
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh)
```

**Windows PowerShell** (管理员):

```ps
irm https://unlock.icmp.ing/scripts/download_test.ps1 | iex
```

## 命令行参数

### 基础选项

|参数|说明|示例|
|---|---|---|
|`-m`|连接模式：0=自动（默认），4=仅IPv4，6=仅IPv6|`-m 4` 仅测试IPv4|
|`-v`|显示版本信息并退出|`-v`|
|`-u`|检查并更新到最新版本|`-u`|

### 性能优化

|参数|说明|示例|
|---|---|---|
|`-conc`|最大并发测试数量（0=无限制）|`-conc 50` 限制最大50个并发测试|
|`-cache`|启用缓存和串行地区执行|`-cache` 启用缓存模式|
|`-show-active`|在进度条中显示正在进行的测试|`-show-active=false` 关闭活动测试显示|

### 调试与测试

|参数|说明|示例|
|---|---|---|
|`-debug`|开启调试模式（输出详细错误信息）|`-debug`|
|`-nf`|仅测试 Netflix 可用性|`-nf`|
|`-test`|运行特定测试|`-test`|

## 常见用例

```bash
# 默认检测所有项目
./unlock-test

# 仅检测 IPv4 项目
./unlock-test -m 4

# 限制并发数量为 30 (适合低配机器)
./unlock-test -conc 30

# 开启调试模式查看详细错误
./unlock-test -debug
```
