# MediaUnlockMonitor

**中文文档** | [English Docs](https://github.com/HsukqiLee/MediaUnlockTest/blob/main/monitor/README_en.md)

使用 Grafana 对接 Prometheus，搭建流媒体监控页面。

### 效果

![](https://raw.githubusercontent.com/HsukqiLee/MediaUnlockTest/main/monitor/image.png)

### 安装

Linux / macOS / Android Termux：

（分别需要 /usr/bin、/usr/local/bin、$PREFIX/bin 目录的读写权限）

```bash
bash <(curl -Ls unlock.icmp.ing/scripts/monitor.sh) -service
```

如果需要设置服务，还需要 Systemd / Upstart / SysV / Launchd / Windows Service 的操作权限。

Windows PowerShell（需要以管理员身份启动）：

```ps
irm https://unlock.icmp.ing/scripts/download_monitor.ps1 | iex
```

### 使用

```
root@server:~# unlock-monitor -h

Usage of build/unlock-monitor_linux_amd64:

  -listen string
        listen address (default ":9101")
  -interval uint
        check interval (s) (default 60)
  -I string
        source ip / interface
  -token string
        check token in http headers or queries
  -metrics-path string
        custom metrics path (default "/metrics")
        
  -mul
        Mutation (default true)
  -hk
        Hong Kong
  -tw
        Taiwan
  -jp
        Japan
  -kr
        Korea
  -na
        North America
  -sa
        South America
  -eu
        Europe
  -afr
        Africa
  -ocea
        Oceania

  -install
        install service
  -uninstall
        uninstall service
  -start
        start service
  -stop
        stop service
  -auto-update
        set auto update
  -update-interval uint
        update check interval (s) (default 86400)

  -u    check update
  -v    show version

  -node string
        Prometheus exporter node field
root@server:~#
```

### Prometheus:

各服务器设置服务时可以指定 `node` 字段，会在 Metrics 中显示。

假设各服务器 Metrics 域名为 `<地区>`-`<编号>`.unlock.check（不一定需要格式相同），targets 字段都只填 `<域名>`:`<端口>`，可以设置 HTTP 重定向到 HTTPS。

添加 Job:

```yaml
global:
  scrape_interval:     300s
  evaluation_interval: 300s

scrape_configs:
  - job_name: checkmedia
    scrape_interval: 300s
    static_configs:
      - targets:
        - <your ip>:9101
        - <your domain>
        - ...
```

此时，Grafana 中每个面板的 `node` 字段就填写自己指定的。

如果域名格式相同，没有指定 `node`，可以配置 `relabel_config`：

```yaml
global:
  scrape_interval:     300s
  evaluation_interval: 300s

scrape_configs:
  - job_name: checkmedia
    scrape_interval: 300s
    static_configs:
      - targets:
        - hk-1.unlock.check
        - hk-2.unlock.check
        - tw-1.unlock.check
        - ...
    relabel_configs:
      - source_labels: [__address__]
        regex: '(.*)\.unlock\.check'
        target_label: "node"
```

此时会自动提取 `<地区>`-`<编号>` 作为 `node`，把域名中的这部分填入 Grafana 的 `node` 即可。

注意：此项服务的 `scrape_interval` 最好不要大于十分钟（600s），否则面板可能会时不时出现 `Data does not have a time field`。

### Grafana

面板 JSON 模板：

```json
{
  "datasource": {
    "type": "prometheus",
    "uid": "d5dc8985-25a8-4199-921e-c64ce62fb93d"
  },
  "description": "",
  "fieldConfig": {
    "defaults": {
      "custom": {
        "lineWidth": 3,
        "fillOpacity": 100,
        "spanNulls": false,
        "insertNulls": false,
        "hideFrom": {
          "tooltip": false,
          "viz": false,
          "legend": false
        }
      },
      "color": {
        "mode": "continuous-GrYlRd"
      },
      "mappings": [],
      "thresholds": {
        "mode": "absolute",
        "steps": [
          {
            "color": "green",
            "value": null
          }
        ]
      }
    },
    "overrides": [
      {
        "matcher": {
          "id": "byValue",
          "options": {
            "op": "neq",
            "reducer": "allValues",
            "value": 0
          }
        },
        "properties": [
          {
            "id": "mappings",
            "value": [
              {
                "options": {
                  "0": {
                    "color": "dark-blue",
                    "index": 8,
                    "text": "Unknown"
                  },
                  "1": {
                    "color": "semi-dark-green",
                    "index": 0,
                    "text": "YES"
                  },
                  "2": {
                    "color": "dark-yellow",
                    "index": 1,
                    "text": "Restricted"
                  },
                  "3": {
                    "color": "dark-red",
                    "index": 2,
                    "text": "NO"
                  },
                  "4": {
                    "color": "dark-orange",
                    "index": 3,
                    "text": "BANNED"
                  },
                  "5": {
                    "color": "#494949",
                    "index": 4,
                    "text": "FAILED"
                  },
                  "6": {
                    "color": "purple",
                    "index": 5,
                    "text": "UNEXPECTED"
                  },
                  "-1": {
                    "color": "light-red",
                    "index": 6,
                    "text": "NET ERR"
                  },
                  "-2": {
                    "color": "semi-dark-red",
                    "index": 7,
                    "text": "ERROR"
                  }
                },
                "type": "value"
              }
            ]
          }
        ]
      }
    ]
  },
  "gridPos": {
    "h": 20,
    "w": 6,
    "x": 0,
    "y": 1
  },
  "id": 1,
  "options": {
    "mergeValues": true,
    "showValue": "auto",
    "alignValue": "center",
    "rowHeight": 0.9,
    "legend": {
      "showLegend": true,
      "displayMode": "list",
      "placement": "bottom"
    },
    "tooltip": {
      "mode": "single",
      "sort": "none",
      "maxHeight": 600
    }
  },
  "pluginVersion": "10.1.1",
  "targets": [
    {
      "datasource": {
        "type": "prometheus",
        "uid": "d5dc8985-25a8-4199-921e-c64ce62fb93d"
      },
      "disableTextWrap": false,
      "editorMode": "builder",
      "expr": "media_unblock_status{node=\"hk-1\"}",
      "fullMetaSearch": false,
      "includeNullMetadata": true,
      "instant": false,
      "legendFormat": "{{region}} {{mediaName}}",
      "range": true,
      "refId": "A",
      "useBackend": false
    }
  ],
  "title": "HK-1",
  "type": "state-timeline"
}
```

替换里面的若干数据后即可导入，或者导入后在 Grafana 中编辑面板参数为自己设置的参数。

Value mappings 参考：

|Status|Display Text|
|---|---|
|0|Unknown|
|1|YES|
|2|Restricted|
|3|NO|
|4|BANNED|
|5|FAILED|
|6|UNEXPECTED|
|-1|NET ERR|
|-2|ERROR|
