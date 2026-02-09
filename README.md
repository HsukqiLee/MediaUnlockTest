# MediaUnlockTest

**中文文档** | [English Docs](https://github.com/HsukqiLee/MediaUnlockTest/blob/main/README_en.md)

> 更快的流媒体检测工具

修复了原作者留下的若干 Bugs，提供比原版更多的测试项目！

> 也许是目前最快的流媒体检测工具，欢迎来挑战

## 快速开始

Linux (包括 iOS iSH) / macOS / Android Termux：

（分别需要 /usr/bin、/usr/local/bin、$PREFIX/bin 目录的读写权限）

```bash
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh)
```

Windows PowerShell（需要以管理员身份启动）：

```ps
irm https://unlock.icmp.ing/scripts/download_test.ps1 | iex

## 深入了解

我们整理了详细的文档来帮助你使用本工具：

-   **[基础使用 / Basic Usage](docs/zh/usage.md)**: 快速开始、命令行参数、基础模式。
-   **[代理与高级设置 / Proxy & Advanced](docs/zh/proxy.md)**: 代理设置、网络接口绑定、DNS 配置。
-   **[Docker 使用 / Docker Guide](docs/zh/docker.md)**: 使用容器运行、多架构支持。
-   **[开发指南 / Development](docs/zh/development.md)**: 源码构建、添加新检测项目。

### English Documentation

-   [Basic Usage](docs/en/usage.md)
-   [Proxy & Advanced Settings](docs/en/proxy.md)
-   [Docker Guide](docs/en/docker.md)
-   [Development Guide](docs/en/development.md)


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
- [开发文档](docs/zh/development.md) - 参与项目开发

## 鸣谢

原项目基于 [lmc的全能检测脚本](https://github.com/lmc999/RegionRestrictionCheck) 的思路使用 Golang 重构，提供更快的检测速度。

本项目基于 [MediaUnlockTest](https://github.com/nkeonkeo/MediaUnlockTest) 二次开发，提供更丰富的测试项目。

Made with ❤️ By **Hsukqi Lee**.

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=HsukqiLee/MediaUnlockTest&type=Date)](https://star-history.com/#HsukqiLee/MediaUnlockTest&Date)