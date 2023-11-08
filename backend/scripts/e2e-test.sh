#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

set -e

ROOT_DIR=$(dirname $(dirname "$0"))
EXTRA=""

if [ -n "$DEVLAKE_DEBUG" ]; then
    EXTRA="-gcflags='all=-N -l'"
fi

export TZ=UTC
export ENV_PATH=${ENV_PATH:-$ROOT_DIR/.env}

go run $ROOT_DIR/test/init.go
for m in $(go list $ROOT_DIR/test/e2e/... | grep -v manual); do \
  echo start running e2e test on $m
  go test -failfast -timeout 120s -v $m
done