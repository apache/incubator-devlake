#!/bin/sh
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

# If you want to use this, you need to run `PLUGIN=github make dev`
# to compile all plugins `make dev`

set -e

echo "Usage: "
echo "  build all plugins:              $0 [golang build flags...]"
echo "  build and keep one plugin only: PLUGIN=jira $0 [golang build flags...]"

SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )"
PLUGIN_SRC_DIR=$SCRIPT_DIR/../plugins
PLUGIN_OUTPUT_DIR=$SCRIPT_DIR/../bin/plugins

if [ -z "$PLUGIN" ]; then
    PLUGINS=$(find $PLUGIN_SRC_DIR/* -maxdepth 0 -type d -not -name core -not -name helper -not -empty)
else
    PLUGINS=$PLUGIN_SRC_DIR/$PLUGIN
fi

rm -rf $PLUGIN_OUTPUT_DIR/*
for PLUG in $PLUGINS; do
    NAME=$(basename $PLUG)
    echo "Building plugin $NAME to bin/plugins/$NAME/$NAME.so"
    go build -buildmode=plugin "$@" -o $PLUGIN_OUTPUT_DIR/$NAME/$NAME.so $PLUG/*.go &
done
wait
