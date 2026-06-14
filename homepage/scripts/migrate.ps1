Write-Host "清理旧版二进制文件"

# 定义目标目录
$targetDir = "$env:ProgramFiles\MediaUnlockTest"

# 确保目标目录存在
if (-not (Test-Path -Path $targetDir)) {
    New-Item -ItemType Directory -Path $targetDir
}

# 删除旧版二进制文件
$unlockTestPath = "$targetDir\unlock-test.exe"
$unlockMonitorPath = "$targetDir\unlock-monitor.exe"

if (Test-Path -Path $unlockTestPath) {
    Remove-Item -Path $unlockTestPath -Force
    Write-Host "删除 unlock-test 成功"
}

if (Test-Path -Path $unlockMonitorPath) {
    Remove-Item -Path $unlockMonitorPath -Force
    Write-Host "删除 unlock-monitor 成功"
}

Write-Host "安装新版本文件"

# 判断系统架构
switch ($env:PROCESSOR_ARCHITECTURE) {
    'AMD64' { $arch = 'amd64' }
    'x86'   { $arch = '386' }
    'ARM64' { $arch = 'arm64' }
    'ARM' {
        $arch = 'arm7'
        Write-Host "当前架构为 ARM 32，默认下载 ARMv7 版本"
        Write-Host "如果无法执行可自行到 Releases 页面下载"
    }
    'ARM32' {
        $arch = 'arm7'
        Write-Host "当前架构为 ARM 32，默认下载 ARMv7 版本"
        Write-Host "如果无法执行可自行到 Releases 页面下载"
    }
    default {
        Write-Error "不支持的系统架构: $env:PROCESSOR_ARCHITECTURE"
        exit 1
    }
}


# 下载文件的函数
function Download-File {
    param (
        [string]$url,
        [string]$output
    )
    try {
        Invoke-WebRequest -Uri $url -OutFile $output
        Write-Host "下载 $output 成功"
    } catch {
        Write-Error "下载 $output 失败: $_"
    }
}

# 下载并安装 unlock-test
$unlockTestUrl = "https://unlock.icmp.ing/test/latest/unlock-test_windows_${arch}.exe"
Download-File -url $unlockTestUrl -output "$targetDir\unlock-test.exe"
if (Test-Path -Path "$targetDir\unlock-test.exe") {
    Write-Host "unlock-test 更新成功"
    & "$targetDir\unlock-test.exe" -v
}

# 下载并安装 unlock-monitor
$unlockMonitorUrl = "https://unlock.icmp.ing/monitor/latest/unlock-monitor_windows_${arch}.exe"
Download-File -url $unlockMonitorUrl -output "$targetDir\unlock-monitor.exe"
if (Test-Path -Path "$targetDir\unlock-monitor.exe") {
    Write-Host "unlock-monitor 更新成功"
    & "$targetDir\unlock-monitor.exe" -v
}

# 重启 unlock-monitor 服务
if (Get-Service -Name "unlock-monitor" -ErrorAction SilentlyContinue) {
    Restart-Service -Name "unlock-monitor"
    Write-Host "unlock-monitor 服务重启成功"
} else {
    Write-Host "unlock-monitor 服务不存在，无法重启"
}
