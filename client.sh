#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo "skycoin-exchange client binary dir:" "$DIR"

pushd "$DIR" >/dev/null

go run cmd/client/client.go --gui-dir="${DIR}/src/web-app/static/" $@

popd >/dev/null