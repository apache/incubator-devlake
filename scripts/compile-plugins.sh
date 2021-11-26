#!/bin/sh
# Change this if you only want to complile the plugin you are working on

# If you want to use this, you neec to run `LOCAL=true make dev`
# to compile all plugins `make dev`

PLUGIN_TO_COMPILE=merico-analysis-engine # this must be the name of the plugin folder

set -e

SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )"
PLUGIN_SRC_DIR=$SCRIPT_DIR/../plugins
PLUGIN_OUTPUT_DIR=$SCRIPT_DIR/../bin/plugins

if [ "$LOCAL" == "true" ]; then
  for PLUG in $(find $PLUGIN_SRC_DIR/* -maxdepth 0 -type d -not -name core -not -empty); do
    NAME=$(basename $PLUG)

    if [ "$NAME" == "$PLUGIN_TO_COMPILE" ]; then
        echo "Building plugin $NAME to bin/plugins/$NAME/$NAME.so"
        go build -buildmode=plugin "$@" -o $PLUGIN_OUTPUT_DIR/$NAME/$NAME.so $PLUG/*.go
    fi
  done
else
  for PLUG in $(find $PLUGIN_SRC_DIR/* -maxdepth 0 -type d -not -name core -not -empty); do
    NAME=$(basename $PLUG)

    echo "Building plugin $NAME to bin/plugins/$NAME/$NAME.so"
    go build -buildmode=plugin "$@" -o $PLUGIN_OUTPUT_DIR/$NAME/$NAME.so $PLUG/*.go
  done
fi
