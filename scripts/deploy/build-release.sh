#!/usr/bin/env -S bash -euo pipefail

gitroot="$(git rev-parse --show-toplevel)"

export GOOS="${GOOS:-freebsd}"
export GOARCH="${GOARCH:-amd64}"

cd $gitroot
mkdir -p rel
rm -rf rel
cd $gitroot/assets
pnpm run build
cd $gitroot/admin
pnpm run build
cd $gitroot
go build -a -installsuffix cgo -o rel/server -tags PROD .
mkdir -p rel/{assets,admin}
cp -R $gitroot/assets/dist/ rel/assets
cp -R $gitroot/admin/build/client/ rel/admin
cp -R $gitroot/db/migrations rel/

TAR_OPTS="--no-xattrs"

if [[ "$(uname)" = "Darwin" ]]; then
  TAR_OPTS="--no-xattrs --no-mac-metadata"
fi

cd rel && tar czf release.tar.gz $TAR_OPTS server assets/ admin/ migrations/
