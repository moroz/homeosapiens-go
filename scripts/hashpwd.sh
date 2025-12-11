#!/usr/bin/env -S bash -euo pipefail

salt="$(openssl rand 8)"
read -r password
echo $password | tr -d '\n' | argon2 "$salt" -e -id -m 16 -p 4
