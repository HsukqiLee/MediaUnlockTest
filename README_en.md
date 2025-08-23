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

## Command Line Arguments

### Basic Parameters

|Args|Description|Example|
|-|-|-|
|`-m`|Connection mode: 0=auto (default), 4=IPv4 only, 6=IPv6 only|`-m 4` IPv4 only|
|`-I`|Bind IP/Network interface|`-I 192.168.1.100` or `-I eth0`|
|`-v`|Show version information and exit|`-v`|
|`-u`|Check and update to latest version|`-u`|

### Proxy Settings

|Args|Description|Example|
|-|-|-|
|`-http-proxy`|Set HTTP proxy|`-http-proxy "http://username:password@127.0.0.1:1080"`|
|`-socks-proxy`|Set SOCKS5 proxy|`-socks-proxy "socks5://username:password@127.0.0.1:1080"`|
|`-dns-servers`|Specify DNS servers|`-dns-servers "1.1.1.1:53"`|

### Performance Optimization

|Args|Description|Example|
|-|-|-|
|`-conc`|Max concurrent tests (0=unlimited)|`-conc 50` limit to 50 concurrent tests|
|`-cache`|Enable caching and sequential region execution|`-cache` enable cache mode|
|`-show-active`|Show active tests in progress bar|`-show-active=false` disable active test display|

### Debug Options

|Args|Description|Example|
|-|-|-|
|`-debug`|Enable debug mode (output detailed error info)|`-debug`|
|`-nf`|Only test Netflix availability|`-nf`|
|`-test`|Run specific test|`-test`|

### Usage Examples

#### Basic Usage
```bash
# Test all items by default
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh)

# Test IPv4 items only
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -m 4

# Test IPv6 items only
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -m 6
```

#### Proxy Usage
```bash
# Use HTTP proxy
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -http-proxy "http://127.0.0.1:8080"

# Use SOCKS5 proxy
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -socks-proxy "socks5://127.0.0.1:1080"

# Use proxy with authentication
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -http-proxy "http://user:pass@127.0.0.1:8080"
```

#### Performance Optimization
```bash
# Limit concurrent tests to 30
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -conc 30

# Enable cache mode (sequential execution, reduce network pressure)
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -cache

# Disable active test display (reduce output)
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -show-active=false
```

#### Debug and Testing
```bash
# Enable debug mode
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -debug

# Test Netflix only
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -nf

# Check for updates
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -u
```

#### Combined Usage
```bash
# Use proxy, IPv4 only, limit concurrency, enable debug
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) \
  -http-proxy "http://127.0.0.1:8080" \
  -m 4 \
  -conc 20 \
  -debug

# Use SOCKS5 proxy, IPv6 only, enable cache
bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) \
  -socks-proxy "socks5://127.0.0.1:1080" \
  -m 6 \
  -cache
```

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

[Development Guide](https://github.com/HsukqiLee/MediaUnlockTest/blob/main/DEVELOPMENT_en.md)

## Contributors

<!--GAMFC_DELIMITER--><a href="https://github.com/nkeonkeo" title="neko"><img src="https://avatars.githubusercontent.com/u/79034142?v=4" width="50;" alt="neko"/></a>
<a href="https://github.com/HsukqiLee" title="HsukqiLee"><img src="https://avatars.githubusercontent.com/u/79034142?v=4" width="50;" alt="HsukqiLee"/></a>
<a href="https://github.com/xkww3n" title="xkww3n"><img src="https://avatars.githubusercontent.com/u/30206355?v=4" width="50;" alt="xkww3n"/></a>
<a href="https://github.com/oif" title="Neo Zhuo"><img src="https://avatars.githubusercontent.com/u/6374269?v=4" width="50;" alt="Neo Zhuo"/></a>
<a href="https://github.com/edgist" title="edgist"><img src="https://avatars.githubusercontent.com/u/6374269?v=4" width="50;" alt="edgist"/></a>
<a href="https://github.com/iawb-ray" title="iawb-ray"><img src="https://avatars.githubusercontent.com/u/49180084?v=4" width="50;" alt="iawb-ray"/></a><!--GAMFC_DELIMITER_END-->

## Community Guidelines

- [Code of Conduct](CODE_OF_CONDUCT.md) - Learn about our community behavior standards
- [Security Policy](SECURITY.md) - Report security vulnerabilities and learn about security best practices
- [Development Guide](DEVELOPMENT_en.md) - Contribute to project development

## Acknowledgements

The original project was refactored using Golang based on the idea of [lmc's all-round media unlock test script](https://github.com/lmc999/RegionRestrictionCheck) to provide faster detection speed.

This project is based on [MediaUnlockTest](https://github.com/nkeonkeo/MediaUnlockTest) secondary development, providing more abundant test projects.

Made with ❤️ By **Hsukqi Lee**.

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=HsukqiLee/MediaUnlockTest&type=Date)](https://star-history.com/#HsukqiLee/MediaUnlockTest&Date)