#!/usr/bin/env bash

set -o pipefail

export CGO_ENABLED=1

if ! command -v goreleaser &> /dev/null; then
  echo "Installing goreleaser"
  brew install goreleaser/tap/goreleaser
fi

echo -n "Running goreleaser release: "
ERRS=$(goreleaser release --clean 2>&1)
if [ $? -eq 1 ]; then
    echo "FAIL"
    echo "${ERRS}"
    echo
    exit 1
fi
echo "PASS"
echo