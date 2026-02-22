#!/usr/bin/env bash

set -euo pipefail

main() {
    local files=("projman")

    for file in "${files[@]}"; do
        if [[ ! -e $file ]]; then
            continue
        fi

        rm "$file"
        echo "Removed './$file'"
    done
}

main
