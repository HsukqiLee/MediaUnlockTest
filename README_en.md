# MediaUnlockTest

[中文文档](https://github.com/HsukqiLee/MediaUnlockTest/blob/main/README.md) | **English Docs**

> Faster media unlock test tool

Fixed some bugs, provides more test items than the original version!

> Perhaps the fastest streaming media detection tool currently available, welcome to challenge.

## Unlock Test

Linux (Including iOS iSH) / macOS / Android Termux：

(requires read and write permissions for the /usr/bin, /usr/local/bin, and $PREFIX/bin directories respectively)

```bash
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh)
```

Windows PowerShell：

```ps
irm https://unlock.icmp.ing/scripts/download_test.ps1 | iex
```

The following command line parameter examples are based on Linux.

Only test IPv4 items:

```bash
bash <(curl -Ls unlock.icmp.ing/test.sh) -m 4
```

Only test IPv6 items (test items known to support IPv6):

```bash
bash <(curl -Ls unlock.icmp.ing/test.sh) -m 6
```

Other args:

|Args|Description|
|-|-|
|`-I`|Bind IP/Interface|
|`-f`|Force IPv6 testing (all selected items are tested using IPv6)|
|`-u`|Check for updates|
|`-v`|Output version|
|`-conc`|Number of concurrent requests|
|`-debug`|Enable debug mode (output detailed error reports of Err/Network Err)|
|`-dns-servers`|Specify DNS servers (example: "1.1.1.1:53")|
|`-http-proxy`|Set HTTP proxy (example: "http://username:password@127.0.0.1:1080")|
|`-socks-proxy`|Set SOCKS5 proxy (example: "socks5://username:password@127.0.0.1:1080")|

## Unlock Monitor

Usage Prometheus and Grafana build streaming media unlock monitoring, demo: [ICMPing](https://icmp.ing/service).


[README](https://github.com/HsukqiLee/MediaUnlockTest/blob/main/monitor/readme.md)

## Todo

- Add unlock detection for more regions
- Ensure that existing test items take into account various return results
- Fix existing/potential problems

## Secondary development

```golang
import "https://github.com/HsukqiLee/MediaUnlockTest"
```

Import it in your golang project and use it

You can use it to make unlock monitoring and other small toys

## Contributors

<!--GAMFC_DELIMITER--><a href="https://github.com/nkeonkeo" title="neko"><img src="https://avatars.githubusercontent.com/u/36293036?v=4" width="50;" alt="neko"/></a>
<a href="https://github.com/HsukqiLee" title="HsukqiLee"><img src="https://avatars.githubusercontent.com/u/79034142?v=4" width="50;" alt="HsukqiLee"/></a>
<a href="https://github.com/xkww3n" title="xkww3n"><img src="https://avatars.githubusercontent.com/u/30206355?v=4" width="50;" alt="xkww3n"/></a>
<a href="https://github.com/oif" title="Neo Zhuo"><img src="https://avatars.githubusercontent.com/u/6374269?v=4" width="50;" alt="Neo Zhuo"/></a>
<a href="https://github.com/edgist" title="edgist"><img src="https://avatars.githubusercontent.com/u/34343603?v=4" width="50;" alt="edgist"/></a>
<a href="https://github.com/iawb-ray" title="iawb-ray"><img src="https://avatars.githubusercontent.com/u/49180084?v=4" width="50;" alt="iawb-ray"/></a><!--GAMFC_DELIMITER_END-->

## Acknowledgements

The original project was refactored using Golang based on the idea of [lmc's all-round media unlock test script](https://github.com/lmc999/RegionRestrictionCheck) to provide faster detection speed.

This project is based on [MediaUnlockTest](https://github.com/nkeonkeo/MediaUnlockTest) secondary development, providing more abundant test projects.

Made with ❤️ By **Hsukqi Lee**.

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=HsukqiLee/MediaUnlockTest&type=Date)](https://star-history.com/#HsukqiLee/MediaUnlockTest&Date)