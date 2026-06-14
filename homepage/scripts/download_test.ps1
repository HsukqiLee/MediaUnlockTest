function Start-Script {
    $currentUser = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($currentUser)
    if (-not $principal.IsInRole([Security.Principal.WindowsBuiltinRole]::Administrator)) {
        Write-Host "Please run this script as Administrator!" -ForegroundColor Red
    }

    if (-not (Test-Path -Path "$env:ProgramFiles\MediaUnlockTest")) {
        New-Item -ItemType Directory -Path "$env:ProgramFiles\MediaUnlockTest" -Force
    }
    
    # 设置控制台输出为UTF-8编码
    [Console]::OutputEncoding = [System.Text.Encoding]::UTF8
    
    # 定义远程脚本URL和本地文件路径
    $remoteScriptUrl = "https://unlock.icmp.ing/scripts/test.ps1"
    $localScriptPath = "C:\Program Files\MediaUnlockTest\test.ps1"  # 请修改为你想要保存的本地路径
    
    # 下载远程脚本文件
    Invoke-WebRequest -Uri $remoteScriptUrl -OutFile $localScriptPath
    
    # 读取和执行下载的脚本文件内容
    $scriptContent = Get-Content -Path $localScriptPath -Raw -Encoding UTF8
    Invoke-Expression $scriptContent
}
Start-Script
