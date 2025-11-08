#!/bin/sh -ex

export GO_ENV=prod
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING=$DATABASE_URL
export GOOSE_MIGRATION_DIR=/app/db/migrations

/usr/local/bin/goose up
exec ./server