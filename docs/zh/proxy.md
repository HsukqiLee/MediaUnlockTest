# 代理与高级设置

本工具支持通过 HTTP、SOCKS5 代理进行测试，并提供 DNS 服务器指定等高级选项。

## 代理设置

|参数|说明|示例|
|---|---|---|
|`-http-proxy`|设置 HTTP 代理|`-http-proxy "http://username:password@127.0.0.1:1080"`|
|`-socks-proxy`|设置 SOCKS5 代理|`-socks-proxy "socks5://username:password@127.0.0.1:1080"`|

**注意**：如果是通过 `sudo` 或高权限运行，请确保代理设置（如环境变量或配置文件）对当前用户可见。

### 使用示例

```bash
# 使用HTTP代理
./unlock-test -http-proxy "http://127.0.0.1:8080"

# 使用SOCKS5代理
./unlock-test -socks-proxy "socks5://127.0.0.1:1080"

# 使用带认证的代理
./unlock-test -http-proxy "http://user:pass@127.0.0.1:8080"
```

## 网络配置

|参数|说明|示例|
|---|---|---|
|`-I`|绑定的 IP / 网络接口|`-I 192.168.1.100` 或 `-I eth0`|
|`-dns-servers`|指定 DNS 服务器|`-dns-servers "1.1.1.1:53"`|

### 使用示例

```bash
# 绑定特定出口网卡
./unlock-test -I eth0

# 指定 DNS 服务器（解决部分地区 DNS 污染问题）
./unlock-test -dns-servers "8.8.8.8:53"
```
