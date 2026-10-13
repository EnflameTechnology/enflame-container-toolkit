#!/bin/bash
#
# Copyright 2024 Enflame. All Rights Reserved.
#
set -eu -o pipefail

ARCH_NAME=$(uname -m)

function usage() {
  cat <<EOF
Usage: ./build.sh [command]
Command:
    all     Build the binaries and make the release package.
    clean   Clean the binaries and the release package.
Example:
    ./build.sh all
    ./build.sh clean
EOF
}

# Version: read from the monorepo g_version.txt when present, otherwise default.
if [ -f "../../g_version.txt" ]; then
    source ../../g_version.txt
fi
PKG_VER=${ENFLAME_CONTAINER_TOOLKIT_VERSION:-"2.0.0"}

BINS="enflame-container-runtime enflame-container-runtime-hook enflame-container-runtime.cdi enflame-container-runtime.legacy enflame-ctk"

function clean_all() {
    make clean >/dev/null 2>&1 || true
    if [ -d "dist" ]; then
        rm -rf dist
    fi
}

function build_all() {
    echo -e "\033[33mBuilding enflame-container-toolkit_${PKG_VER} (${ARCH_NAME}).\033[0m"
    make binaries
}

function package_all() {
    PKG_DIR="dist/enflame-container-toolkit_${PKG_VER}"
    rm -rf "${PKG_DIR}"
    mkdir -p "${PKG_DIR}/bin"

    for b in ${BINS}; do
        if [ -f "${b}" ]; then
            cp -f "${b}" "${PKG_DIR}/bin/"
        fi
    done
    cp -rf LICENSE "${PKG_DIR}/"
    cp -rf README.md "${PKG_DIR}/"
    if [ -d "deployments" ]; then
        cp -rf deployments "${PKG_DIR}/"
    fi

    echo -e "\033[33menflame-container-toolkit_${PKG_VER} packaged under ${PKG_DIR}.\033[0m"
}

function main() {
    [[ "$#" -lt 1 ]] && { usage >&2; exit 1; }
    case $1 in
    "clean")
        clean_all
        ;;
    "all")
        clean_all
        build_all
        package_all
        ;;
    *)
        usage
        exit 1
        ;;
    esac
}

main "$@"
