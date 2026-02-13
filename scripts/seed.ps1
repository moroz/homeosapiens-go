#!/usr/bin/env pwsh

param(
  [string]$RemoteHost = "hs.authorizz.com",
  [string]$RemoteUser = "deploy"
)

$bastionSsh = "$RemoteUser@$RemoteHost"

function Start-SshTunnel
{
  param(
    [Parameter(Mandatory = $true)]
    [int]$Port
  )

  $processInfo = New-Object System.Diagnostics.ProcessStartInfo
  $processInfo.FileName = "ssh"
  $processInfo.Arguments = "-N -L $($Port):localhost:5432 $bastionSsh"
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

$tunnel = Start-SshTunnel -Port $port

Write-Host "Sleeping for 5 to wait for SSH tunnel..."
Start-Sleep -Seconds 5

$uri = [System.UriBuilder]$databaseUrl
$uri.Host = "127.0.0.1"
$uri.Port = $port

cd db/seeds && env SECRET_KEY_BASE="$secretKeyBase" DATABASE_URL="$uri" go run .

Stop-Process -Id $tunnel.Id
