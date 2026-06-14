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
    if [ ! -d "$target_dir" ]; then
        sudo mkdir -p "$target_dir"
        if [ $? -ne 0 ]; then
            echo "创建文件夹失败: $target_dir"
            exit 1
        fi
    fi
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

# 下载并安装 unlock-test
url="https://unlock.icmp.ing/test/latest/unlock-test_${os}_${arch}"
wget ${url} -O unlock-test || curl ${url} -o unlock-test
chmod +x unlock-test
mv unlock-test "$target_dir/unlock-test" 2>/dev/null || sudo mv unlock-test "$target_dir/unlock-test"

# 检查安装是否成功
if $target_dir/unlock-test -v; then
    echo "unlock-test 安装成功"
    $target_dir/unlock-test "$@"
else
    echo "unlock-test 安装失败"
fi
