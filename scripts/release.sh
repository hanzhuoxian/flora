#!/usr/bin/env bash

ROOT_DIR=$(dirname "${BASH_SOURCE[0]}")/..

# shellcheck source=/dev/null
source "${ROOT_DIR}"/scripts/common.sh
# shellcheck source=/dev/null
source "${ROOT_DIR}"/scripts/library/release.sh

golang::setup_env
build::verify_prereqs
release::verify_prereqs
# build::build_image
build::build_command
release::package_tarballs
# git push origin "${VERSION:-}"
# release::github_release
# release::generate_changelog

