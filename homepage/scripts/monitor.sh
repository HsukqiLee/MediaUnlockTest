#!/bin/sh

# 定义目标目录变量
target_dir="/usr/bin"

# 判断操作系统类型
if [ "$(uname -s)" = "Darwin" ]; then
    target_dir="/usr/local/bin"
    if [ ! -d "$target_dir" ]; then
        mkdir -p "$target_dir"
        if [ $? -ne 0 ]; then
            echo "创建文件夹失败: $target_dir"
            exit 1
        fi
    fi
fi

# 检查 unlock-monitor 是否已存在
if [ -f "$target_dir/unlock-monitor" ]; then
    $target_dir/unlock-monitor -u
    $target_dir/unlock-monitor "$@"
    exit
fi

# 获取操作系统和架构信息
if [ -f /etc/openwrt_release ]; then
    # 如果是 OpenWRT
    os=$(uname -s | awk '{print tolower($0)}')
else
    # 如果不是 OpenWRT，按默认方式处理
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
fi
arch=$(uname -m)

# 根据架构修改变量
case ${arch} in
    x86)
        arch="386"
        ;;
    x86_64)
        arch="amd64"
        ;;
    aarch64)
        arch="arm64"
        ;;
esac

# 下载并安装 unlock-monitor
url="https://unlock.icmp.ing/monitor/latest/unlock-monitor_${os}_${arch}"
wget ${url} -O unlock-monitor || curl ${url} -o unlock-monitor
chmod +x unlock-monitor
sudo mv unlock-monitor "$target_dir/unlock-monitor"

# 检查安装是否成功
if $target_dir/unlock-monitor -v; then
    echo "unlock-monitor 安装成功"
    $target_dir/unlock-monitor "$@"
else
    echo "unlock-monitor 安装失败"
fi
