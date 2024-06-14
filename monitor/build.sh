# #!/bin/sh
# mkdir build
tp="-gcflags=-trimpath="$GOPATH" -asmflags=-trimpath="$GOPATH
flags="-w -s -X 'update.buildTime=$(date '+%Y-%m-%d %H:%M:%S')'"

# exit

echo "build android 386 ..."
CGO_ENABLED=1 GOOS=android GOARCH=386 CC=i686-linux-android$API_LEVEL-clang CXX=i686-linux-android$API_LEVEL-clang++ go build -tags netcgo -ldflags="$flags" $tp -o build/unlock-monitor_android_386
echo "build android amd64 ..."
CGO_ENABLED=1 GOOS=android GOARCH=amd64 CC=x86_64-linux-android$API_LEVEL-clang CXX=x86_64-linux-android$API_LEVEL-clang++ go build -tags netcgo -ldflags="$flags" $tp -o build/unlock-monitor_android_amd64
echo "build android arm7 ..."
CGO_ENABLED=1 GOOS=android GOARCH=arm GOARM=7 CC=armv7a-linux-androideabi$API_LEVEL-clang CXX=armv7a-linux-androideabi$API_LEVEL-clang++ go build -tags netcgo -ldflags="$flags" $tp -o build/unlock-monitor_android_arm7
echo "build android arm6 ..."
CGO_ENABLED=1 GOOS=android GOARCH=arm GOARM=6 CC=armv7a-linux-androideabi$API_LEVEL-clang CXX=armv7a-linux-androideabi$API_LEVEL-clang++ go build -tags netcgo -ldflags="$flags" $tp -o build/unlock-monitor_android_arm6
echo "build android arm5 ..."
CGO_ENABLED=1 GOOS=android GOARCH=arm GOARM=5 CC=armv7a-linux-androideabi$API_LEVEL-clang CXX=armv7a-linux-androideabi$API_LEVEL-clang++ go build -tags netcgo -ldflags="$flags" $tp -o build/unlock-monitor_android_arm5
echo "build android arm64 ..."
CGO_ENABLED=1 GOOS=android GOARCH=arm64 CC=aarch64-linux-android$API_LEVEL-clang CXX=aarch64-linux-android$API_LEVEL-clang++ go build -tags netcgo -ldflags="$flags" $tp -o build/unlock-monitor_android_arm64

echo "build darwin amd64 ..."
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="$flags" $tp -o build/unlock-monitor_darwin_amd64
echo "build darwin arm64 ..."
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="$flags" $tp -o build/unlock-monitor_darwin_arm64

#echo "build dragonfly amd64 ..."
#CGO_ENABLED=0 GOOS=dragonfly GOARCH=amd64 go build -ldflags="$flags" $tp -o build/unlock-monitor_dragonfly_amd64

echo "build freebsd 386 ..."
CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build -ldflags="$flags" $tp -o build/unlock-monitor_freebsd_386
echo "build freebsd amd64 ..."
CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -ldflags="$flags" $tp -o build/unlock-monitor_freebsd_amd64
echo "build freebsd arm7 ..."
CGO_ENABLED=0 GOOS=freebsd GOARCH=arm GOARM=7 go build -ldflags="$flags" $tp -o build/unlock-monitor_freebsd_arm7
echo "build freebsd arm6 ..."
CGO_ENABLED=0 GOOS=freebsd GOARCH=arm GOARM=6 go build -ldflags="$flags" $tp -o build/unlock-monitor_freebsd_arm6
echo "build freebsd arm5 ..."
CGO_ENABLED=0 GOOS=freebsd GOARCH=arm GOARM=5 go build -ldflags="$flags" $tp -o build/unlock-monitor_freebsd_arm5
echo "build freebsd arm64 ..."
CGO_ENABLED=0 GOOS=freebsd GOARCH=arm64 go build -ldflags="$flags" $tp -o build/unlock-monitor_freebsd_arm64
echo "build freebsd riscv64 ..."
CGO_ENABLED=0 GOOS=freebsd GOARCH=riscv64 go build -ldflags="$flags" $tp -o build/unlock-monitor_freebsd_riscv64

#echo "build illumos amd64 ..."
#CGO_ENABLED=0 GOOS=illumos GOARCH=amd64 go build -ldflags="$flags" $tp -o build/unlock-monitor_illumos_amd64

echo "build linux 386 ..."
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags="$flags" $tp -o build/unlock-monitor_linux_386
echo "build linux amd64 ..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$flags" $tp -o build/unlock-monitor_linux_amd64
echo "build linux arm7 ..."
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="$flags" $tp -o build/unlock-monitor_linux_arm7
echo "build linux arm6 ..."
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags="$flags" $tp -o build/unlock-monitor_linux_arm6
echo "build linux arm5 ..."
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -ldflags="$flags" $tp -o build/unlock-monitor_linux_arm5
echo "build linux arm64 ..."
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="$flags" $tp -o build/unlock-monitor_linux_arm64
echo "build linux loong64 ..."
CGO_ENABLED=0 GOOS=linux GOARCH=loong64 go build -ldflags="$flags" $tp -o build/unlock-monitor_linux_loong64
echo "build linux mips ..."
CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -ldflags="$flags" $tp -o build/unlock-monitor_linux_mips
echo "build linux mips64 ..."
CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -ldflags="$flags" $tp -o build/unlock-monitor_linux_mips64
echo "build linux mips64le ..."
CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -ldflags="$flags" $tp -o build/unlock-monitor_linux_mips64le
echo "build linux mipsle ..."
CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build -ldflags="$flags" $tp -o build/unlock-monitor_linux_mipsle
echo "build linux ppc64 ..."
CGO_ENABLED=0 GOOS=linux GOARCH=ppc64 go build -ldflags="$flags" $tp -o build/unlock-monitor_linux_ppc64
echo "build linux ppc64le ..."
CGO_ENABLED=0 GOOS=linux GOARCH=ppc64le go build -ldflags="$flags" $tp -o build/unlock-monitor_linux_ppc64le
echo "build linux riscv64 ..."
CGO_ENABLED=0 GOOS=linux GOARCH=riscv64 go build -ldflags="$flags" $tp -o build/unlock-monitor_linux_riscv64
echo "build linux s390x ..."
CGO_ENABLED=0 GOOS=linux GOARCH=s390x go build -ldflags="$flags" $tp -o build/unlock-monitor_linux_s390x

echo "build netbsd 386 ..."
CGO_ENABLED=0 GOOS=netbsd GOARCH=386 go build -ldflags="$flags" $tp -o build/unlock-monitor_netbsd_386
echo "build netbsd amd64 ..."
CGO_ENABLED=0 GOOS=netbsd GOARCH=amd64 go build -ldflags="$flags" $tp -o build/unlock-monitor_netbsd_amd64
echo "build netbsd arm7 ..."
CGO_ENABLED=0 GOOS=netbsd GOARCH=arm GOARM=7 go build -ldflags="$flags" $tp -o build/unlock-monitor_netbsd_arm7
echo "build netbsd arm6 ..."
CGO_ENABLED=0 GOOS=netbsd GOARCH=arm GOARM=6 go build -ldflags="$flags" $tp -o build/unlock-monitor_netbsd_arm6
echo "build netbsd arm5 ..."
CGO_ENABLED=0 GOOS=netbsd GOARCH=arm GOARM=5 go build -ldflags="$flags" $tp -o build/unlock-monitor_netbsd_arm5
echo "build netbsd arm64 ..."
CGO_ENABLED=0 GOOS=netbsd GOARCH=arm64 go build -ldflags="$flags" $tp -o build/unlock-monitor_netbsd_arm64

echo "build openbsd 386 ..."
CGO_ENABLED=0 GOOS=openbsd GOARCH=386 go build -ldflags="$flags" $tp -o build/unlock-monitor_openbsd_386
echo "build openbsd amd64 ..."
CGO_ENABLED=0 GOOS=openbsd GOARCH=amd64 go build -ldflags="$flags" $tp -o build/unlock-monitor_openbsd_amd64
echo "build openbsd arm7 ..."
CGO_ENABLED=0 GOOS=openbsd GOARCH=arm GOARM=7 go build -ldflags="$flags" $tp -o build/unlock-monitor_openbsd_arm7
echo "build openbsd arm6 ..."
CGO_ENABLED=0 GOOS=openbsd GOARCH=arm GOARM=6 go build -ldflags="$flags" $tp -o build/unlock-monitor_openbsd_arm6
echo "build openbsd arm5 ..."
CGO_ENABLED=0 GOOS=openbsd GOARCH=arm GOARM=5 go build -ldflags="$flags" $tp -o build/unlock-monitor_openbsd_arm5
echo "build openbsd arm64 ..."
CGO_ENABLED=0 GOOS=openbsd GOARCH=arm64 go build -ldflags="$flags" $tp -o build/unlock-monitor_openbsd_arm64
echo "build openbsd ppc64 ..."
CGO_ENABLED=0 GOOS=openbsd GOARCH=ppc64 go build -ldflags="$flags" $tp -o build/unlock-monitor_openbsd_ppc64

#echo "build plan9 386 ..."
#CGO_ENABLED=0 GOOS=plan9 GOARCH=386 go build -ldflags="$flags" $tp -o build/unlock-monitor_plan9_386
#echo "build plan9 amd64 ..."
#CGO_ENABLED=0 GOOS=plan9 GOARCH=amd64 go build -ldflags="$flags" $tp -o build/unlock-monitor_plan9_amd64
#echo "build plan9 arm7 ..."
#CGO_ENABLED=0 GOOS=plan9 GOARCH=arm GOARM=7 go build -ldflags="$flags" $tp -o build/unlock-monitor_plan9_arm7
#echo "build plan9 arm6 ..."
#CGO_ENABLED=0 GOOS=plan9 GOARCH=arm GOARM=6 go build -ldflags="$flags" $tp -o build/unlock-monitor_plan9_arm6
#echo "build plan9 arm5 ..."
#CGO_ENABLED=0 GOOS=plan9 GOARCH=arm GOARM=5 go build -ldflags="$flags" $tp -o build/unlock-monitor_plan9_arm5

echo "build solaris amd64 ..."
CGO_ENABLED=0 GOOS=solaris GOARCH=amd64 go build -ldflags="$flags" $tp -o build/unlock-monitor_solaris_amd64

echo "build windows 386 ..."
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags="$flags" $tp -o build/unlock-monitor_windows_386.exe
echo "build windows amd64 ..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="$flags" $tp -o build/unlock-monitor_windows_amd64.exe
echo "build windows arm7 ..."
CGO_ENABLED=0 GOOS=windows GOARCH=arm GOARM=7 go build -ldflags="$flags" $tp -o build/unlock-monitor_windows_arm7.exe
echo "build windows arm6 ..."
CGO_ENABLED=0 GOOS=windows GOARCH=arm GOARM=6 go build -ldflags="$flags" $tp -o build/unlock-monitor_windows_arm6.exe
echo "build windows arm5 ..."
CGO_ENABLED=0 GOOS=windows GOARCH=arm GOARM=5 go build -ldflags="$flags" $tp -o build/unlock-monitor_windows_arm5.exe
echo "build windows arm64 ..."
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -ldflags="$flags" $tp -o build/unlock-monitor_windows_arm64.exe


# upx build/unlock-monitor_linux*