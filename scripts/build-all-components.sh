#!/usr/bin/env bash

# Copyright (c) 2021, NVIDIA CORPORATION.  All rights reserved.
# Copyright (c) 2024 Enflame Technologies, Inc. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script is used to build the packages for the components of the ENFLAME
# Container Stack. These include the enflame-container-toolkit in this repository
# as well as the components included in the third_party folder.
# All required packages are generated in the specified dist folder.

function assert_usage() {
    echo "Missing argument $1"
    echo "$(basename "${BASH_SOURCE[0]}") TARGET"
    exit 1
}

set -e

SCRIPTS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )"/../scripts && pwd )"
PROJECT_ROOT="$( cd "${SCRIPTS_DIR}"/.. && pwd )"

if [[ $# -ne 1 ]]; then
    assert_usage "TARGET"
fi

TARGET=$1

source "${SCRIPTS_DIR}"/utils.sh

: "${DIST_DIR:=${PROJECT_ROOT}/dist}"
export DIST_DIR

echo "Building ${TARGET} for all packages to ${DIST_DIR}"

: "${LIBENFLAME_CONTAINER_ROOT:=${PROJECT_ROOT}/third_party/libenflame-container}"
: "${ENFLAME_CONTAINER_TOOLKIT_ROOT:=${PROJECT_ROOT}}"


"${SCRIPTS_DIR}/get-component-versions.sh"

if [[ -z "${ENFLAME_CONTAINER_TOOLKIT_VERSION}" || -z "${LIBENFLAME_CONTAINER_VERSION}" ]]; then
eval $(${SCRIPTS_DIR}/get-component-versions.sh)
fi

# Build libenflame-container
if [[ -z ${SKIP_LIBENFLAME_CONTAINER} ]]; then
    make -C "${LIBENFLAME_CONTAINER_ROOT}" -f mk/docker.mk \
        LIB_VERSION=${ENFLAME_CONTAINER_TOOLKIT_VERSION} \
        LIB_TAG=${ENFLAME_CONTAINER_TOOLKIT_TAG} \
        "${TARGET}"
fi

if [[ -z ${SKIP_ENFLAME_CONTAINER_TOOLKIT} ]]; then
# Build enflame-container-toolkit
make -C "${ENFLAME_CONTAINER_TOOLKIT_ROOT}" \
    LIBENFLAME_CONTAINER_VERSION="${ENFLAME_CONTAINER_TOOLKIT_VERSION}" \
    LIBENFLAME_CONTAINER_TAG="${ENFLAME_CONTAINER_TOOLKIT_TAG}" \
        "${TARGET}"
fi
