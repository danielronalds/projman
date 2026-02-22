#!/usr/bin/env bash

set -euo pipefail

main() {
    if [[ -x "$(command -v "gotestsum")" ]]; then
        gotestsum ./...
    else
        go test ./...
    fi
}

main
