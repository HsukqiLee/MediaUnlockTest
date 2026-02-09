# Proxy & Advanced Settings

This tool supports testing via HTTP/SOCKS5 proxies and advanced options like specifying DNS servers.

## Proxy Configuration

|Flag|Description|Example|
|---|---|---|
|`-http-proxy`|Set HTTP Proxy|`-http-proxy "http://username:password@127.0.0.1:1080"`|
|`-socks-proxy`|Set SOCKS5 Proxy|`-socks-proxy "socks5://username:password@127.0.0.1:1080"`|

**Note**: If running with `sudo` or elevated privileges, ensure proxy settings are accessible.

### Examples

```bash
# Use HTTP Proxy
./unlock-test -http-proxy "http://127.0.0.1:8080"

# Use SOCKS5 Proxy
./unlock-test -socks-proxy "socks5://127.0.0.1:1080"

# Use Proxy with Authentication
./unlock-test -http-proxy "http://user:pass@127.0.0.1:8080"
```

## Network Configuration

|Flag|Description|Example|
|---|---|---|
|`-I`|Bind specific IP/Interface|`-I 192.168.1.100` or `-I eth0`|
|`-dns-servers`|Specify DNS Servers|`-dns-servers "1.1.1.1:53"`|

### Examples

```bash
# Bind to a specific network interface
./unlock-test -I eth0

# Specify DNS server (to resolve DNS pollution)
./unlock-test -dns-servers "8.8.8.8:53"
```
