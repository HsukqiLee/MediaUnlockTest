#!/bin/sh

# 获取操作系统和架构信息
if [ -f /etc/openwrt_release ]; then
    # 如果是 OpenWRT
    os=$(uname -s | awk '{print tolower($0)}')
else
    # 如果不是 OpenWRT，按默认方式处理
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
fi
arch=$(uname -m)

# 定义目标目录变量
target_dir="/usr/bin"

# 判断是否是 Termux
if command -v termux-setup-storage > /dev/null 2>&1; then  
    os="android"
    target_dir=$PREFIX/bin
fi

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

# 检查 unlock-test 是否已存在并且版本是否正确
if [ -f "$target_dir/unlock-test" ]; then
    version=$($target_dir/unlock-test -v)
    if [ "$version" != "1.1" ]; then
        $target_dir/unlock-test -u
        $target_dir/unlock-test "$@"
        exit
    fi
fi




# 根据架构修改变量
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

# 下载并安装 unlock-test
url="https://unlock.icmp.ing/test/latest/unlock-test_${os}_${arch}"
wget ${url} -O unlock-test || curl ${url} -o unlock-test
chmod +x unlock-test
if [ -n "$USE_SUDO" ]; then
    $USE_SUDO mv unlock-test "$target_dir/unlock-test"
else
    mv unlock-test "$target_dir/unlock-test"
fi

# 检查安装是否成功
if $target_dir/unlock-test -v; then
    echo "unlock-test 安装成功"
    $target_dir/unlock-test "$@"
else
    echo "unlock-test 安装失败"
fi
