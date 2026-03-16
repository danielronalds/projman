#!/usr/bin/env bash

set -euo pipefail

main() {
    local readme="README.md"
    local help_output
    help_output=$(go run . help)

    local usage_line
    usage_line=$(grep -n "^## Usage$" "$readme" | head -1 | cut -d: -f1)

    if [[ -z "$usage_line" ]]; then
        echo "Could not find '## Usage' section in $readme" >&2
        exit 1
    fi

    local start_line
    start_line=$(tail -n +"$usage_line" "$readme" | grep -n '^\`\`\`console$' | head -1 | cut -d: -f1)
    start_line=$((usage_line + start_line - 1))

    local end_line
    end_line=$(tail -n +"$((start_line + 1))" "$readme" | grep -n '^\`\`\`$' | head -1 | cut -d: -f1)
    end_line=$((start_line + end_line))

    local head_content
    head_content=$(head -n "$start_line" "$readme")

    local tail_content
    tail_content=$(tail -n +"$end_line" "$readme")

    printf '%s\n%s\n%s\n' "$head_content" "$help_output" "$tail_content" > "$readme"

    echo "Updated help block in $readme"
}

main
