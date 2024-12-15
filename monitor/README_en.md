# MediaUnlockMonitor

[中文文档](https://github.com/HsukqiLee/MediaUnlockTest/blob/main/monitor/README.md) | **English Docs**

Use Grafana to connect to Prometheus and build a streaming media monitoring page.

### Effect

![](https://raw.githubusercontent.com/HsukqiLee/MediaUnlockTest/main/monitor/image.png)

### Installation

Linux (Including iOS iSH) / macOS / Android Termux：

(Read and write permissions are required for /usr/bin, /usr/local/bin, and $PREFIX/bin directories respectively)

```bash
bash <(curl -Ls unlock.icmp.ing/scripts/monitor.sh) -service
```

If you need to set up a service, you also need the operating permissions of Systemd / Upstart / SysV / Launchd / Windows Service.

Windows PowerShell (needs to be started as an administrator):

```ps
irm https://unlock.icmp.ing/scripts/download_monitor.ps1 | iex
```

### Usage

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

### Prometheus

When setting up services for each server, you can specify the `node` field, which will be displayed in Metrics.

Assume that the domain name of each server's Metrics is `<region>`-`<number>`.unlock.check (not necessarily in the same format), and the targets field is only filled with `<domain>`:`<port>`. You can set HTTP to redirect to HTTPS.

Add Job:

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

At this point, the `node` field of each panel in Grafana is filled with the specified one.

If the domain name format is the same and `node` is not specified, you can configure `relabel_config`:

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

At this time, `<region>`-`<number>` will be automatically extracted as `node`, and this part of the domain name can be filled in Grafana's `node`.

Note: The `scrape_interval` of this service should not be greater than ten minutes (600s), otherwise the panel may occasionally show `Data does not have a time field`.

### Grafana

Dashboard JSON template:

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


Replace some of the data in it and then import it, or edit the panel parameters in Grafana to set the parameters for yourself after importing.

Value mappings reference:

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