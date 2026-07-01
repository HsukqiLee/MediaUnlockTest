param (
    [string]$TargetDir = "$env:ProgramFiles\MediaUnlockTest"
)

$targetDir = $TargetDir

# 确保目标目录存在
if (-not (Test-Path -Path $targetDir)) {
    New-Item -ItemType Directory -Path $targetDir
}

# 判断系统架构
switch ($env:PROCESSOR_ARCHITECTURE) {
    'AMD64' { $arch = 'amd64' }
    'x86'   { $arch = '386' }
    'ARM64' { $arch = 'arm64' }
    default {
        Write-Error "不支持的系统架构: $env:PROCESSOR_ARCHITECTURE"
        exit 1
    }
}

# 下载文件的函数
function Invoke-DownloadFile {
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

# 下载并重命名 unlock-monitor
$unlockMonitorUrl = "https://unlock.icmp.ing/monitor/latest/unlock-monitor_windows_$arch.exe"
$unlockMonitorPath = "$targetDir\unlock-monitor.exe"
$unlockMonitorNewPath = "$targetDir\unlock-monitor_new.exe"

if (Test-Path -Path $unlockMonitorPath) {
    Start-Process -FilePath $unlockMonitorPath -ArgumentList '-u' -NoNewWindow -Wait
}
else {
    Invoke-DownloadFile -url $unlockMonitorUrl -output $unlockMonitorNewPath
    if (Test-Path -Path $unlockMonitorNewPath) {
        Rename-Item -Path $unlockMonitorNewPath -NewName "unlock-monitor.exe"
        Write-Host "unlock-monitor 安装成功"
    }
    $path = [System.Environment]::GetEnvironmentVariable("PATH", [System.EnvironmentVariableTarget]::Machine)
    $pathArray = $path -split ';'
    if ($pathArray -notcontains $targetDir) {
        [System.Environment]::SetEnvironmentVariable("PATH", "$path;$targetDir", [System.EnvironmentVariableTarget]::Machine)
        Write-Host "PATH 环境变量已更新，已自动尝试启动 unlock-monitor"
        $env:PATH += ";$targetDir"
        Write-Host "如果报错，请尝试重新启动 PowerShell 或命令提示符以使更改生效"
    }
}

# & unlock-monitor @args
