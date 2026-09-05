#!/usr/bin/env pwsh

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"
$PSNativeCommandUseErrorActionPreference = $true

$orders = Resolve-Path "../../orders.csv"
$participants = Resolve-Path "../../participants.csv"
$users = Resolve-Path "../../students_export.csv"

get-content -Encoding utf8 $users | go run ./20260903-import-users
get-content -Encoding utf8 $participants | go run ./20260904-import-participants
