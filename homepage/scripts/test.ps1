param (
    [string]$TargetDir = "$env:ProgramFiles\MediaUnlockTest"
)

$targetDir = $TargetDir


if (-not (Test-Path -Path $targetDir)) {
    New-Item -ItemType Directory -Path $targetDir -Force
}


switch ($env:PROCESSOR_ARCHITECTURE) {
    'AMD64' { $arch = 'amd64' }
    'x86'   { $arch = '386' }
    'ARM64' { $arch = 'arm64' }
    default {
        Write-Error "不支持的系统架构: $env:PROCESSOR_ARCHITECTURE"
        exit 1
    }
}


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
        exit 1
    }
}


$unlockTestUrl = "https://unlock.icmp.ing/test/latest/unlock-test_windows_$arch.exe"
$unlockTestPath = "$targetDir\unlock-test.exe"
$unlockTestNewPath = "$targetDir\unlock-test_new.exe"


if (Test-Path -Path $unlockTestPath) {
    Start-Process -FilePath $unlockTestPath -ArgumentList '-u' -NoNewWindow -Wait
}
else {
    Invoke-DownloadFile -url $unlockTestUrl -output $unlockTestNewPath
    if (Test-Path -Path $unlockTestNewPath) {
        Rename-Item -Path $unlockTestNewPath -NewName "unlock-test.exe"
        Write-Host "unlock-test 安装成功"
    }
    $path = [System.Environment]::GetEnvironmentVariable("PATH", [System.EnvironmentVariableTarget]::Machine)
    $pathArray = $path -split ';'
    if ($pathArray -notcontains $targetDir) {
        [System.Environment]::SetEnvironmentVariable("PATH", "$path;$targetDir", [System.EnvironmentVariableTarget]::Machine)
        Write-Host "PATH 环境变量已更新，已自动尝试启动 unlock-test"
        $env:PATH += ";$targetDir"
        Write-Host "如果报错，请尝试重新启动 PowerShell 或命令提示符以使更改生效"
    }
}


