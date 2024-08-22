@echo off
setlocal

set "tp=-gcflags=-trimpath=%GOPATH% -asmflags=-trimpath=%GOPATH%"
set "flags=-w -s"

echo build android 386 ...
set "CC=i686-linux-android%API_LEVEL%-clang"
set "CXX=i686-linux-android%API_LEVEL%-clang++"
set "GOOS=android"
set "GOARCH=386"
set "CGO_ENABLED=1"
go build -tags netcgo -ldflags="%flags%" %tp% -o build/unlock-test_android_386

echo build android amd64 ...
set "CC=x86_64-linux-android%API_LEVEL%-clang"
set "CXX=x86_64-linux-android%API_LEVEL%-clang++"
set "GOARCH=amd64"
go build -tags netcgo -ldflags="%flags%" %tp% -o build/unlock-test_android_amd64

echo build android arm7 ...
set "CC=armv7a-linux-androideabi%API_LEVEL%-clang"
set "CXX=armv7a-linux-androideabi%API_LEVEL%-clang++"
set "GOARCH=arm"
set "GOARM=7"
go build -tags netcgo -ldflags="%flags%" %tp% -o build/unlock-test_android_arm7

echo build android arm6 ...
set "GOARM=6"
go build -tags netcgo -ldflags="%flags%" %tp% -o build/unlock-test_android_arm6

echo build android arm5 ...
set "GOARM=5"
go build -tags netcgo -ldflags="%flags%" %tp% -o build/unlock-test_android_arm5

echo build android arm64 ...
set "CC=aarch64-linux-android%API_LEVEL%-clang"
set "CXX=aarch64-linux-android%API_LEVEL%-clang++"
set "GOARCH=arm64"
go build -tags netcgo -ldflags="%flags%" %tp% -o build/unlock-test_android_arm64

echo build darwin amd64 ...
set "GOOS=darwin"
set "GOARCH=amd64"
set "CGO_ENABLED=0"
go build -ldflags="%flags%" %tp% -o build/unlock-test_darwin_amd64

echo build darwin arm64 ...
set "GOARCH=arm64"
go build -ldflags="%flags%" %tp% -o build/unlock-test_darwin_arm64

echo build freebsd 386 ...
set "GOOS=freebsd"
set "GOARCH=386"
go build -ldflags="%flags%" %tp% -o build/unlock-test_freebsd_386

echo build freebsd amd64 ...
set "GOARCH=amd64"
go build -ldflags="%flags%" %tp% -o build/unlock-test_freebsd_amd64

echo build freebsd arm7 ...
set "GOARCH=arm"
set "GOARM=7"
go build -ldflags="%flags%" %tp% -o build/unlock-test_freebsd_arm7

echo build freebsd arm6 ...
set "GOARM=6"
go build -ldflags="%flags%" %tp% -o build/unlock-test_freebsd_arm6

echo build freebsd arm5 ...
set "GOARM=5"
go build -ldflags="%flags%" %tp% -o build/unlock-test_freebsd_arm5

echo build freebsd arm64 ...
set "GOARCH=arm64"
go build -ldflags="%flags%" %tp% -o build/unlock-test_freebsd_arm64

echo build freebsd riscv64 ...
set "GOARCH=riscv64"
go build -ldflags="%flags%" %tp% -o build/unlock-test_freebsd_riscv64

echo build linux 386 ...
set "GOOS=linux"
set "GOARCH=386"
go build -ldflags="%flags%" %tp% -o build/unlock-test_linux_386

echo build linux amd64 ...
set "GOARCH=amd64"
go build -ldflags="%flags%" %tp% -o build/unlock-test_linux_amd64

echo build linux arm7 ...
set "GOARCH=arm"
set "GOARM=7"
go build -ldflags="%flags%" %tp% -o build/unlock-test_linux_arm7

echo build linux arm6 ...
set "GOARM=6"
go build -ldflags="%flags%" %tp% -o build/unlock-test_linux_arm6

echo build linux arm5 ...
set "GOARM=5"
go build -ldflags="%flags%" %tp% -o build/unlock-test_linux_arm5

echo build linux arm64 ...
set "GOARCH=arm64"
go build -ldflags="%flags%" %tp% -o build/unlock-test_linux_arm64

echo build linux loong64 ...
set "GOARCH=loong64"
go build -ldflags="%flags%" %tp% -o build/unlock-test_linux_loong64

echo build linux mips ...
set "GOARCH=mips"
go build -ldflags="%flags%" %tp% -o build/unlock-test_linux_mips

echo build linux mips64 ...
set "GOARCH=mips64"
go build -ldflags="%flags%" %tp% -o build/unlock-test_linux_mips64

echo build linux mips64le ...
set "GOARCH=mips64le"
go build -ldflags="%flags%" %tp% -o build/unlock-test_linux_mips64le

echo build linux mipsle ...
set "GOARCH=mipsle"
go build -ldflags="%flags%" %tp% -o build/unlock-test_linux_mipsle

echo build linux ppc64 ...
set "GOARCH=ppc64"
go build -ldflags="%flags%" %tp% -o build/unlock-test_linux_ppc64

echo build linux ppc64le ...
set "GOARCH=ppc64le"
go build -ldflags="%flags%" %tp% -o build/unlock-test_linux_ppc64le

echo build linux riscv64 ...
set "GOARCH=riscv64"
go build -ldflags="%flags%" %tp% -o build/unlock-test_linux_riscv64

echo build linux s390x ...
set "GOARCH=s390x"
go build -ldflags="%flags%" %tp% -o build/unlock-test_linux_s390x

echo build netbsd 386 ...
set "GOOS=netbsd"
set "GOARCH=386"
go build -ldflags="%flags%" %tp% -o build/unlock-test_netbsd_386

echo build netbsd amd64 ...
set "GOARCH=amd64"
go build -ldflags="%flags%" %tp% -o build/unlock-test_netbsd_amd64

echo build netbsd arm7 ...
set "GOARCH=arm"
set "GOARM=7"
go build -ldflags="%flags%" %tp% -o build/unlock-test_netbsd_arm7

echo build netbsd arm6 ...
set "GOARM=6"
go build -ldflags="%flags%" %tp% -o build/unlock-test_netbsd_arm6

echo build netbsd arm5 ...
set "GOARM=5"
go build -ldflags="%flags%" %tp% -o build/unlock-test_netbsd_arm5

echo build netbsd arm64 ...
set "GOARCH=arm64"
go build -ldflags="%flags%" %tp% -o build/unlock-test_netbsd_arm64

echo build openbsd 386 ...
set "GOOS=openbsd"
set "GOARCH=386"
go build -ldflags="%flags%" %tp% -o build/unlock-test_openbsd_386

echo build openbsd amd64 ...
set "GOARCH=amd64"
go build -ldflags="%flags%" %tp% -o build/unlock-test_openbsd_amd64

echo build openbsd arm7 ...
set "GOARCH=arm"
set "GOARM=7"
go build -ldflags="%flags%" %tp% -o build/unlock-test_openbsd_arm7

echo build openbsd arm6 ...
set "GOARM=6"
go build -ldflags="%flags%" %tp% -o build/unlock-test_openbsd_arm6

echo build openbsd arm5 ...
set "GOARM=5"
go build -ldflags="%flags%" %tp% -o build/unlock-test_openbsd_arm5

echo build openbsd arm64 ...
set "GOARCH=arm64"
go build -ldflags="%flags%" %tp% -o build/unlock-test_openbsd_arm64

echo build openbsd ppc64 ...
set "GOARCH=ppc64"
go build -ldflags="%flags%" %tp% -o build/unlock-test_openbsd_ppc64

echo build windows 386 ...
set "GOOS=windows"
set "GOARCH=386"
go build -ldflags="%flags%" %tp% -o build/unlock-test_windows_386.exe

echo build windows amd64 ...
set "GOARCH=amd64"
go build -ldflags="%flags%" %tp% -o build/unlock-test_windows_amd64.exe

echo build windows arm7 ...
set "GOARCH=arm"
set "GOARM=7"
go build -ldflags="%flags%" %tp% -o build/unlock-test_windows_arm7.exe

echo build windows arm6 ...
set "GOARM=6"
go build -ldflags="%flags%" %tp% -o build/unlock-test_windows_arm6.exe

echo build windows arm5 ...
set "GOARM=5"
go build -ldflags="%flags%" %tp% -o build/unlock-test_windows_arm5.exe

echo build windows arm64 ...
set "GOARCH=arm64"
go build -ldflags="%flags%" %tp% -o build/unlock-test_windows_arm64.exe

endlocal
