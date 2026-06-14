#!/bin/sh

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


echo '清理旧版二进制文件'
rm -f $target_dir/unlock-test
rm -f $target_dir/unlock-monitor

echo '安装新版本文件'

if [ -f /etc/openwrt_release ]; then
    os=$(uname -s | awk '{print tolower($0)}')
else
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
fi
arch=`uname -m`

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
# url="https://github.com/HsukqiLee/MediaUnlockTest/releases/latest/download/unlock-test_"${os}"_"${arch}
url="https://unlock.icmp.ing/test/latest/unlock-test_"${os}"_"${arch}
wget ${url} -O unlock-test || curl ${url} -o unlock-test
chmod +x unlock-test
mv unlock-test $target_dir/unlock-test
unlock-test -v && echo "unlock-test 更新成功"

url="https://unlock.icmp.ing/monitor/latest/unlock-monitor_"${os}"_"${arch}
wget ${url} -O unlock-monitor || curl ${url} -o unlock-monitor
chmod +x unlock-monitor
mv unlock-monitor $target_dir/unlock-monitor
unlock-monitor -v && echo "unlock-monitor 更新成功"

systemctl restart unlock-monitor && echo "unlock-monitor 重启成功"