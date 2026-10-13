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

test_repo=$1
echo "Setting up TEST repo: ${test_repo}"
sed -i -e "s#enflame\.github\.io/libenflame-container#${test_repo}/libenflame-container#g" /etc/apt/sources.list.d/enflame-container-toolkit.list
sed -i -e '/experimental/ s/^#//g' /etc/apt/sources.list.d/enflame-container-toolkit.list
