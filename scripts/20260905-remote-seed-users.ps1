#!/usr/bin/env pwsh

param(
  [string]$RemoteHost = "prod.homeosapiens.eu",
  [string]$RemoteUser = "deploy"
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"
$PSNativeCommandUseErrorActionPreference = $true

$remoteSsh = "$RemoteUser@$RemoteHost"
$usersSeed = "$env:HOME/working/hs/students_export.csv"
$participantsSeed = "$env:HOME/working/hs/participants.csv"

function Start-SshTunnel
{
  param(
    [Parameter(Mandatory = $true)]
    [int]$LocalPort,

    [Parameter(Mandatory = $true)]
    [string]$RemoteSSH
  )

  $processInfo = New-Object System.Diagnostics.ProcessStartInfo
  $processInfo.FileName = "ssh"
  $processInfo.Arguments = "-N -L $($LocalPort):localhost:5432 $RemoteSSH"
  $processInfo.UseShellExecute = $false
  $processInfo.RedirectStandardError = $true
  $processInfo.RedirectStandardOutput = $true

  Write-Host "ssh $($processInfo.Arguments)"

  $process = New-Object System.Diagnostics.Process
  $process.StartInfo = $processInfo

  if ($process.Start())
  {
    Write-Host "SSH tunnel started successfully on port $Port"
    return $process
  }
  return $null
}

$ENV_FILE="/usr/local/lib/server/homeosapiens.env"

Write-Host "Fetching remote DATABASE_URL..."
$databaseUrl = ssh "$($RemoteUser)@$($RemoteHost)" grep "DATABASE_URL" $ENV_FILE
$databaseUrl = $databaseUrl.Split('=')[1].Replace('"', '')

Write-Host "Fetching remote SECRET_KEY_BASE..."
$secretKeyBase = ssh "$($RemoteUser)@$($RemoteHost)" grep "SECRET_KEY_BASE" $ENV_FILE
$secretKeyBase = $secretKeyBase.Split('"')[1]

Write-Host $databaseUrl

$port = 6000

$tunnel = Start-SshTunnel -LocalPort $port -RemoteSSH $remoteSsh

Write-Host "Sleeping for 5 to wait for SSH tunnel..."
Start-Sleep -Seconds 5

$uri = [System.UriBuilder]$databaseUrl
$uri.Host = "127.0.0.1"
$uri.Port = $port

$env:SECRET_KEY_BASE = "$secretKeyBase"
$env:DATABASE_URL = "$uri"

Get-Content -Encoding UTF8 $usersSeed | go run ./20260903-import-users
Get-Content -Encoding utf8 $participantsSeed | go run ./20260904-import-participants

Stop-Process -Id $tunnel.Id
