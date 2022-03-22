#!/usr/bin/env sh
# Copyright 2022 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -ex

if [ "$APG_REGISTRY_ADDRESS" == "apigeeregistry.googleapis.com:443" ]
then
  export APG_REGISTRY_TOKEN="$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)"
fi

registry get $REGISTRY_SPEC --contents > /openapi.yaml

prism mock -h "0.0.0.0" /openapi.yaml