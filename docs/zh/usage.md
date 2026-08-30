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
|`-f`|强制执行更新（即使已是最新版本，需配合 `-u`）|`-u -f`|
|`--table`|显示紧凑的六列 IPv4 和 IPv6 表格|`--table`|
|`-region`|按菜单编号或地区名称选择检测区域，多个值用逗号分隔|`-region 0,11` 检测跨国和 AI 平台|

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
|`-test [名称]`|运行单个测试（支持显示名或函数名）|`-test Disney+` 或 `-test DisneyPlus`|

## 常见用例

```bash
# 默认检测所有项目
./unlock-test

# 非交互式检测跨国和 AI 平台
./unlock-test -region 0,11

# 也可以使用地区名称
./unlock-test -region Globe,AI

# 仅检测 IPv4 项目
./unlock-test -m 4

# 显示适应终端宽度的紧凑 IPv4 和 IPv6 表格
./unlock-test --table

# 限制并发数量为 30 (适合低配机器)
./unlock-test -conc 30

# 开启调试模式查看详细错误
./unlock-test -debug
```
