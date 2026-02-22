#!/usr/bin/env bash

set -euo pipefail

main() {
    if ! command -v gh &> /dev/null; then
        echo "error: gh CLI is not installed"
        exit 1
    fi

    local latest_tag
    latest_tag=$(git describe --tags --abbrev=0)

    local previous_tag
    previous_tag=$(git describe --tags --abbrev=0 "${latest_tag}^")

    local changelog
    changelog=$(git log "${previous_tag}..${latest_tag}" --pretty=format:"- %s" | grep -v "chore: bump version")

    local repo_url
    repo_url=$(gh repo view --json url --jq '.url')

    local notes
    notes=$(cat <<EOF
## What's Changed

${changelog}

**Full Changelog**: ${repo_url}/compare/${previous_tag}...${latest_tag}
EOF
)

    echo "Creating release ${latest_tag}..."
    echo ""
    echo "${notes}"
    echo ""

    gh release create "${latest_tag}" --title "${latest_tag}" --notes "${notes}"

    echo ""
    echo "Release ${latest_tag} created"
}

main
