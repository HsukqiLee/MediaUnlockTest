# MediaUnlockTest

**中文文档** | [English Docs](https://github.com/HsukqiLee/MediaUnlockTest/blob/main/README_en.md)

> 更快的流媒体检测工具

修复了原作者留下的若干 Bugs，提供比原版更多的测试项目！

> 也许是目前最快的流媒体检测工具，欢迎来挑战

## 解锁检测

Linux (包括 iOS iSH) / macOS / Android Termux：

（分别需要 /usr/bin、/usr/local/bin、$PREFIX/bin 目录的读写权限）

```bash
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh)
```

Windows PowerShell（需要以管理员身份启动）：

```ps
irm https://unlock.icmp.ing/scripts/download_test.ps1 | iex
```

以下命令行参数的示例都以 Linux 为例。

只检测IPv4结果:

```bash
bash <(curl -Ls unlock.icmp.ing/test.sh) -m 4
```

只检测IPv6结果（仅测试已知支持 IPv6 的项目）：

```bash
bash <(curl -Ls unlock.icmp.ing/test.sh) -m 6
```

## 命令行参数详解

### 基本参数

|参数|说明|示例|
|-|-|-|
|`-m`|连接模式：0=自动（默认），4=仅IPv4，6=仅IPv6|`-m 4` 仅测试IPv4|
|`-I`|绑定的 IP / 网络接口|`-I 192.168.1.100` 或 `-I eth0`|
|`-v`|显示版本信息并退出|`-v`|
|`-u`|检查并更新到最新版本|`-u`|

### 代理设置

|参数|说明|示例|
|-|-|-|
|`-http-proxy`|设置 HTTP 代理|`-http-proxy "http://username:password@127.0.0.1:1080"`|
|`-socks-proxy`|设置 SOCKS5 代理|`-socks-proxy "socks5://username:password@127.0.0.1:1080"`|
|`-dns-servers`|指定 DNS 服务器|`-dns-servers "1.1.1.1:53"`|

### 性能优化

|参数|说明|示例|
|-|-|-|
|`-conc`|最大并发测试数量（0=无限制）|`-conc 50` 限制最大50个并发测试|
|`-cache`|启用缓存和串行地区执行|`-cache` 启用缓存模式|
|`-show-active`|在进度条中显示正在进行的测试|`-show-active=false` 关闭活动测试显示|

### 调试选项

|参数|说明|示例|
|-|-|-|
|`-debug`|开启调试模式（输出详细错误信息）|`-debug`|
|`-nf`|仅测试 Netflix 可用性|`-nf`|
|`-test`|运行特定测试|`-test`|

### 使用示例

#### 基本使用
```bash
# 默认检测所有项目
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh)

# 仅检测IPv4项目
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -m 4

# 仅检测IPv6项目
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -m 6
```

#### 代理使用
```bash
# 使用HTTP代理
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -http-proxy "http://127.0.0.1:8080"

# 使用SOCKS5代理
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -socks-proxy "socks5://127.0.0.1:1080"

# 使用带认证的代理
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -http-proxy "http://user:pass@127.0.0.1:8080"
```

#### 性能优化
```bash
# 限制并发数量为30
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -conc 30

# 启用缓存模式（串行执行，减少网络压力）
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -cache

# 关闭活动测试显示（减少输出）
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -show-active=false
```

#### 调试和测试
```bash
# 开启调试模式
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -debug

# 仅测试Netflix
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -nf

# 检查更新
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -u
```

#### 组合使用
```bash
# 使用代理，仅IPv4，限制并发，开启调试
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) \
  -http-proxy "http://127.0.0.1:8080" \
  -m 4 \
  -conc 20 \
  -debug

# 使用SOCKS5代理，仅IPv6，启用缓存
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) \
  -socks-proxy "socks5://127.0.0.1:1080" \
  -m 6 \
  -cache
```

## 解锁监控

使用 Prometheus 和 Grafana 搭建流媒体解锁监控，效果： [ICMPing](https://icmp.ing/service)。

[使用文档](https://github.com/HsukqiLee/MediaUnlockTest/blob/main/monitor/readme.md)

## 任务清单

- 补充更多地区的解锁检测
- 确保已有测试项目考虑到各种返回结果
- 修复已经存在/可能存在的问题

欢迎提交你的 Pull Requests

## 二次开发

```golang
import "https://github.com/HsukqiLee/MediaUnlockTest"
```

在你的 Golang 项目中导入即可使用，你可以使用它制作解锁监控等小玩具。

[开发文档](https://github.com/HsukqiLee/MediaUnlockTest/blob/main/DEVELOPMENT.md)

## 参与开发的小伙伴

<!--GAMFC_DELIMITER--><a href="https://github.com/HsukqiLee" title="Hsukqi Lee"><img src="https://avatars.githubusercontent.com/u/79034142?v=4" width="50;" alt="Hsukqi Lee"/></a>
<a href="https://github.com/nkeonkeo" title="neko"><img src="https://avatars.githubusercontent.com/u/36293036?v=4" width="50;" alt="neko"/></a>
<a href="https://github.com/xkww3n" title="xkww3n"><img src="https://avatars.githubusercontent.com/u/30206355?v=4" width="50;" alt="xkww3n"/></a>
<a href="https://github.com/oif" title="Neo Zhuo"><img src="https://avatars.githubusercontent.com/u/6374269?v=4" width="50;" alt="Neo Zhuo"/></a>
<a href="https://github.com/edgist" title="edgist"><img src="https://avatars.githubusercontent.com/u/34343603?v=4" width="50;" alt="edgist"/></a>
<a href="https://github.com/iawb-ray" title="iawb-ray"><img src="https://avatars.githubusercontent.com/u/49180084?v=4" width="50;" alt="iawb-ray"/></a><!--GAMFC_DELIMITER_END-->

## 社区准则

- [行为准则](CODE_OF_CONDUCT.md) - 了解我们的社区行为标准
- [安全策略](SECURITY.md) - 报告安全漏洞和了解安全最佳实践
- [开发文档](DEVELOPMENT.md) - 参与项目开发

## 鸣谢

原项目基于 [lmc的全能检测脚本](https://github.com/lmc999/RegionRestrictionCheck) 的思路使用 Golang 重构，提供更快的检测速度。

本项目基于 [MediaUnlockTest](https://github.com/nkeonkeo/MediaUnlockTest) 二次开发，提供更丰富的测试项目。

Made with ❤️ By **Hsukqi Lee**.

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=HsukqiLee/MediaUnlockTest&type=Date)](https://star-history.com/#HsukqiLee/MediaUnlockTest&Date)