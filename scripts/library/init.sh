#!/usr/bin/env bash

set -eu
set -o pipefail

# Default use go modules
export GO111MODULE=on

ROOT_DIR=$(cd $(dirname "${BASH_SOURCE[0]}")/../.. && pwd -P)


# shellcheck source=/dev/null
source "${ROOT_DIR}/scripts/library/color.sh"

# shellcheck source=/dev/null
source "${ROOT_DIR}"/scripts/library/util.sh

# shellcheck source=/dev/null
source "${ROOT_DIR}/scripts/library/logging.sh"

# shellcheck source=/dev/null
source "${ROOT_DIR}/scripts/library/version.sh"

# shellcheck source=/dev/null
source "${ROOT_DIR}/scripts/library/golang.sh"
