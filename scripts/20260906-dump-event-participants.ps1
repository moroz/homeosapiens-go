#!/usr/bin/env pwsh

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"
$PSNativeCommandUseErrorActionPreference = $true

$eventIds = @(5177, 6228, 4976, 4403, 3101, 3293, 1753, 2375, 1897, 1452)
$outfile = [System.IO.Path]::GetFullPath("../../participants.csv")

if (Test-Path $outfile)
{
    Remove-Item $outfile
}

$cookie = $env:EZYCOURSE_COOKIE

if ($cookie -eq "")
{
    write-host "Fatal: Environment variable EZYCOURSE_COOKIE is not set!"
    exit 1
}

foreach ($event in $eventIds)
{
    $url = "https://www.homeosapiens.eu/api/teacher/events/getSingleEventParticipantsById/$($event)?type=all&page=1&pageSize=250"
    $data = curl --url "$url" -b "$cookie"
    $data | jq --arg eid $event -r '.data[] | [$eid, .student.email, .event_registered_at] | @csv' | Add-Content -Path $outfile
}
