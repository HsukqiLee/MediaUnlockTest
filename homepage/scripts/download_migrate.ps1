function Start-Script {
    $currentUser = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($currentUser)
    if (-not $principal.IsInRole([Security.Principal.WindowsBuiltinRole]::Administrator)) {
        Write-Host "Please run this script as Administrator!" -ForegroundColor Red
    }

    if (-not (Test-Path -Path "$env:ProgramFiles\MediaUnlockTest")) {
        New-Item -ItemType Directory -Path "$env:ProgramFiles\MediaUnlockTest" -Force
    }

    [Console]::OutputEncoding = [System.Text.Encoding]::UTF8
 
    $remoteScriptUrl = "https://unlock.icmp.ing/scripts/migrate.ps1"
    $localScriptPath = "C:\Program Files\MediaUnlockTest\migrate.ps1"  # 请修改为你想要保存的本地路径

    Invoke-WebRequest -Uri $remoteScriptUrl -OutFile $localScriptPath
    
    $scriptContent = Get-Content -Path $localScriptPath -Raw -Encoding UTF8
    Invoke-Expression $scriptContent
}
Start-MonitorScript
