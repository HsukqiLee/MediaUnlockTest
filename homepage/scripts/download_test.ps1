function Start-Script {
    $currentUser = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($currentUser)
    $targetDir = "$env:ProgramFiles\MediaUnlockTest"
    if (-not $principal.IsInRole([Security.Principal.WindowsBuiltinRole]::Administrator)) {
        if (-not [Environment]::UserInteractive) {
            Write-Host "Non-interactive mode detected. Falling back to LocalAppData." -ForegroundColor Yellow
            $targetDir = "$env:LOCALAPPDATA\MediaUnlockTest"
        } else {
            Write-Host "You are not running as Administrator." -ForegroundColor Yellow
            Write-Host "Please select an option:"
            Write-Host "1) Restart as Administrator (Elevate)"
            Write-Host "2) Install to current user directory ($env:LOCALAPPDATA\MediaUnlockTest)"
            Write-Host "3) Exit"
            $choice = Read-Host "Enter [1-3] (Default 1)"
            
            switch ($choice) {
                "2" {
                    $targetDir = "$env:LOCALAPPDATA\MediaUnlockTest"
                    Write-Host "Target directory changed to $targetDir" -ForegroundColor Green
                }
                "3" {
                    Write-Host "Exiting installation."
                    exit
                }
                default {
                    Write-Host "Attempting to restart as Administrator..." -ForegroundColor Cyan
                    Start-Process powershell.exe -ArgumentList "-NoProfile -ExecutionPolicy Bypass -Command `"& { $PSCommandPath }`"" -Verb RunAs
                    exit
                }
            }
        }
    }

    if (-not (Test-Path -Path $targetDir)) {
        New-Item -ItemType Directory -Path $targetDir -Force | Out-Null
    }
    
    # 设置控制台输出为UTF-8编码
    [Console]::OutputEncoding = [System.Text.Encoding]::UTF8
    
    # 定义远程脚本URL和本地文件路径
    $remoteScriptUrl = "https://unlock.icmp.ing/scripts/test.ps1"
    $localScriptPath = Join-Path -Path $targetDir -ChildPath "test.ps1"
    
    # 下载远程脚本文件
    Invoke-WebRequest -Uri $remoteScriptUrl -OutFile $localScriptPath
    
    # 执行下载的脚本文件 (传递 targetDir 给内部脚本)
    & $localScriptPath -TargetDir $targetDir
}
Start-Script
