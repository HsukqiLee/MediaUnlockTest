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
    
    $remoteScriptUrl = "https://unlock.icmp.ing/scripts/monitor.ps1"
    $localScriptPath = "C:\Program Files\MediaUnlockTest\monitor.ps1"
    
    Invoke-WebRequest -Uri $remoteScriptUrl -OutFile $localScriptPath
    
    $scriptContent = Get-Content -Path $localScriptPath -Raw -Encoding UTF8
    Invoke-Expression $scriptContent
}
Start-Script
