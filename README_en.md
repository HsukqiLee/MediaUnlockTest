# MediaUnlockTest

[中文文档](https://github.com/HsukqiLee/MediaUnlockTest/blob/main/README.md) | **English Docs**

> Faster media unlock test tool

Fixed some bugs, provides more test items than the original version!

> Perhaps the fastest streaming media detection tool currently available, welcome to challenge.

## Quick Start

Linux (Including iOS iSH) / macOS / Android Termux：

(requires read and write permissions for the /usr/bin, /usr/local/bin, and $PREFIX/bin directories respectively)

```bash
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh)
```

Windows PowerShell (Run as Administrator)：

```ps
irm https://unlock.icmp.ing/scripts/download_test.ps1 | iex
```

## Documentation

We have organized detailed documentation to help you use this tool:

-   **[Basic Usage](docs/en/usage.md)**: Quick start, command line arguments, basic modes.
-   **[Proxy & Advanced Settings](docs/en/proxy.md)**: Proxy settings, network interface binding, DNS configuration.
-   **[Docker Guide](docs/en/docker.md)**: Running with Docker, multi-architecture support.
-   **[Development Guide](docs/en/development.md)**: Building from source, adding new checks.

### Chinese Documentation

-   [基础使用 / Basic Usage](docs/zh/usage.md)
-   [代理与高级设置 / Proxy & Advanced](docs/zh/proxy.md)
-   [Docker 使用 / Docker Guide](docs/zh/docker.md)
-   [开发指南 / Development](docs/zh/development.md)


## Unlock Monitor

Use Prometheus and Grafana to build streaming media unlock monitoring, demo: [ICMPing](https://icmp.ing/service).

[Documentation](https://github.com/HsukqiLee/MediaUnlockTest/blob/main/monitor/readme.md)

## Todo

- Add unlock detection for more regions
- Ensure that existing test items take into account various return results
- Fix existing/potential problems

Welcome to submit your Pull Requests.

## Secondary Development

```golang
import "https://github.com/HsukqiLee/MediaUnlockTest"
```

Import it in your golang project and use it. You can use it to create unlock monitoring tools and other small toys.

## Contributors

<!--GAMFC_DELIMITER--><a href="https://github.com/HsukqiLee" title="Hsukqi Lee"><img src="https://avatars.githubusercontent.com/u/79034142?v=4" width="50;" alt="Hsukqi Lee"/></a>
<a href="https://github.com/nkeonkeo" title="neko"><img src="https://avatars.githubusercontent.com/u/36293036?v=4" width="50;" alt="neko"/></a>
<a href="https://github.com/xkww3n" title="xkww3n"><img src="https://avatars.githubusercontent.com/u/30206355?v=4" width="50;" alt="xkww3n"/></a>
<a href="https://github.com/oif" title="Neo Zhuo"><img src="https://avatars.githubusercontent.com/u/6374269?v=4" width="50;" alt="Neo Zhuo"/></a>
<a href="https://github.com/edgist" title="edgist"><img src="https://avatars.githubusercontent.com/u/34343603?v=4" width="50;" alt="edgist"/></a>
<a href="https://github.com/iawb-ray" title="iawb-ray"><img src="https://avatars.githubusercontent.com/u/49180084?v=4" width="50;" alt="iawb-ray"/></a><!--GAMFC_DELIMITER_END-->

## Community Guidelines

- [Code of Conduct](CODE_OF_CONDUCT.md) - Learn about our community behavior standards
- [Security Policy](SECURITY.md) - Report security vulnerabilities and learn about security best practices
- [Development Guide](docs/en/development.md) - Contribute to project development

## Acknowledgements

The original project was refactored using Golang based on the idea of [lmc's all-round media unlock test script](https://github.com/lmc999/RegionRestrictionCheck) to provide faster detection speed.

This project is based on [MediaUnlockTest](https://github.com/nkeonkeo/MediaUnlockTest) secondary development, providing more abundant test projects.

Made with ❤️ By **Hsukqi Lee**.

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=HsukqiLee/MediaUnlockTest&type=Date)](https://star-history.com/#HsukqiLee/MediaUnlockTest&Date)