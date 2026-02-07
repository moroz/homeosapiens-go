#!/usr/bin/env -S bash -euo pipefail

gitroot="$(git rev-parse --show-toplevel)"

export GOOS="${GOOS:-freebsd}"
export GOARCH="${GOARCH:-amd64}"

cd $gitroot
mkdir -p rel
rm -rf rel
cd $gitroot/assets
pnpm run build
cd $gitroot
go build -a -installsuffix cgo -o rel/server -tags PROD .
mkdir -p rel/assets
cp -R $gitroot/assets/dist rel/assets
cp -R $gitroot/db/migrations rel/

# Remove hidden sql files, if any
rm $gitroot/db/migrations/._* || true

cd rel && tar czf release.tar.gz server assets/ migrations/
