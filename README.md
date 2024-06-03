# MediaUnlockTest

**中文文档** | [English Docs](https://github.com/HsukqiLee/MediaUnlockTest/blob/main/README_en.md)

> 更快的流媒体检测工具

修复了原作者留下的若干 Bugs，提供比原版更多的测试项目！

## 解锁检测

Linux / macOS：

```bash
bash <(curl -Ls unlock.icmp.ing/test.sh)
```

Windows PowerShell（需要以管理员身份启动）：

```ps
irm https://unlock.icmp.ing/test.ps1 | iex
```

以下命令行参数的示例都以 Linux 为例。

只检测IPv4结果:

```bash
bash <(curl -Ls unlock.icmp.ing/test.sh) -m 4
```

只检测IPv6结果：

```bash
bash <(curl -Ls unlock.icmp.ing/test.sh) -m 6
```

更多参数：

|参数|说明|
|-|-|
|`--dns-servers`|指定 DNS 服务器|
|`-I`|使用的 IP / 网卡|
|`--http-proxy`|设置代理 (示例: "http://username:password@127.0.0.1:1080")|

## 解锁监控

使用 Prometheus 和 Grafana 搭建流媒体解锁监控，效果： [ICMPing](https://icmp.ing/service)。

~~图文教程有空再写，暂时鸽了~~

[README](https://github.com/HsukqiLee/MediaUnlockTest/blob/main/monitor/readme.md)

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

开发文档孵化中！~~也是有空再写，暂时鸽了~~

## 参与开发的小伙伴

<!--GAMFC_DELIMITER--><a href="https://github.com/nkeonkeo" title="neko"><img src="https://avatars.githubusercontent.com/u/36293036?v=4" width="50;" alt="neko"/></a>
<a href="https://github.com/HsukqiLee" title="HsukqiLee"><img src="https://avatars.githubusercontent.com/u/79034142?v=4" width="50;" alt="HsukqiLee"/></a>
<a href="https://github.com/xkww3n" title="xkww3n"><img src="https://avatars.githubusercontent.com/u/30206355?v=4" width="50;" alt="xkww3n"/></a>
<a href="https://github.com/oif" title="Neo Zhuo"><img src="https://avatars.githubusercontent.com/u/6374269?v=4" width="50;" alt="Neo Zhuo"/></a>
<a href="https://github.com/edgist" title="edgist"><img src="https://avatars.githubusercontent.com/u/34343603?v=4" width="50;" alt="edgist"/></a>
<a href="https://github.com/iawb-ray" title="iawb-ray"><img src="https://avatars.githubusercontent.com/u/49180084?v=4" width="50;" alt="iawb-ray"/></a><!--GAMFC_DELIMITER_END-->

## 鸣谢

原项目基于 [lmc的全能检测脚本](https://github.com/lmc999/RegionRestrictionCheck) 的思路使用 Golang 重构，提供更快的检测速度。

本项目基于 [MediaUnlockTest](https://github.com/nkeonkeo/MediaUnlockTest) 二次开发，提供更丰富的测试项目。

Made with ❤️ By **Hsukqi Lee**.

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=HsukqiLee/MediaUnlockTest&type=Date)](https://star-history.com/#HsukqiLee/MediaUnlockTest&Date)