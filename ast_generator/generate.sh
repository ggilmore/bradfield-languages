#!/usr/bin/env bash

cd "$(dirname "${BASH_SOURCE[0]}")"
set -euxo pipefail

go run . ../src/

cd ..
scripts/go-fmt.sh
