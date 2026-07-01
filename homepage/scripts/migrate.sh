#!/bin/sh

target_dir="/usr/bin"

# 判断操作系统类型
if [ "$(uname -s)" = "Darwin" ]; then
    target_dir="/usr/local/bin"
fi

# 检查是否有写入目标目录的权限
if [ ! -w "$target_dir" ]; then
    if [ ! -d "$target_dir" ] && [ -w "$(dirname "$target_dir")" ]; then
        mkdir -p "$target_dir" 2>/dev/null
    fi
    
    if [ ! -w "$target_dir" ]; then
        if [ ! -t 0 ]; then
            echo "非交互环境检测不到终端，将尝试使用 sudo..."
            USE_SUDO="sudo"
        else
            echo "您没有权限写入 $target_dir。请选择操作："
            echo "1) 尝试使用 sudo 提权安装 (默认)"
            echo "2) 安装到 ~/.local/bin (当前用户目录)"
            echo "3) 退出安装"
            printf "请输入 [1-3]: "
            read choice
            case "$choice" in
                2)
                    target_dir="$HOME/.local/bin"
                    mkdir -p "$target_dir"
                    echo "目标目录已更改为 $target_dir"
                    if ! echo "$PATH" | grep -q "$target_dir"; then
                        echo "⚠️ 警告: $target_dir 不在 PATH 环境变量中，安装后可能需要手动将其加入 PATH 或使用绝对路径运行。"
                    fi
                    USE_SUDO=""
                    ;;
                3)
                    echo "已退出安装"
                    exit 0
                    ;;
                *)
                    echo "尝试使用 sudo 提权..."
                    USE_SUDO="sudo"
                    ;;
            esac
        fi
    fi
else
    USE_SUDO=""
fi


echo '清理旧版二进制文件'
if [ -n "$USE_SUDO" ]; then
    $USE_SUDO rm -f $target_dir/unlock-test
    $USE_SUDO rm -f $target_dir/unlock-monitor
else
    rm -f $target_dir/unlock-test
    rm -f $target_dir/unlock-monitor
fi

echo '安装新版本文件'

if [ -f /etc/openwrt_release ]; then
    os=$(uname -s | awk '{print tolower($0)}')
else
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
fi
arch=`uname -m`

case ${arch} in
    x86) arch="386" ;;
    x86_64|amd64) arch="amd64" ;;
    aarch64|arm64) arch="arm64" ;;
    armv7*|armv8l) arch="arm7" ;;
    armv6*) arch="arm6" ;;
    armv5*) arch="arm5" ;;
    loongarch64) arch="loong64" ;;
    mips64el) arch="mips64le" ;;
    mipsel) arch="mipsle" ;;
esac
# url="https://github.com/HsukqiLee/MediaUnlockTest/releases/latest/download/unlock-test_${os}_${arch}"
url="https://unlock.icmp.ing/test/latest/unlock-test_${os}_${arch}"
wget ${url} -O unlock-test || curl ${url} -o unlock-test
chmod +x unlock-test
if [ -n "$USE_SUDO" ]; then
    $USE_SUDO mv unlock-test $target_dir/unlock-test
else
    mv unlock-test $target_dir/unlock-test
fi
unlock-test -v && echo "unlock-test 更新成功"

url="https://unlock.icmp.ing/monitor/latest/unlock-monitor_${os}_${arch}"
wget ${url} -O unlock-monitor || curl ${url} -o unlock-monitor
chmod +x unlock-monitor
if [ -n "$USE_SUDO" ]; then
    $USE_SUDO mv unlock-monitor $target_dir/unlock-monitor
else
    mv unlock-monitor $target_dir/unlock-monitor
fi
unlock-monitor -v && echo "unlock-monitor 更新成功"

if [ -n "$USE_SUDO" ]; then
    $USE_SUDO systemctl restart unlock-monitor 2>/dev/null && echo "unlock-monitor 重启成功"
else
    systemctl restart unlock-monitor 2>/dev/null && echo "unlock-monitor 重启成功"
fi