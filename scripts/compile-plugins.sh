#!/bin/sh
# Change this if you only want to complile the plugin you are working on
PLUGIN_TO_COMPILE=all # this must be all or the name of the plugin

set -e

SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )"
PLUGIN_SRC_DIR=$SCRIPT_DIR/../plugins
PLUGIN_OUTPUT_DIR=$SCRIPT_DIR/../bin/plugins

for PLUG in $(find $PLUGIN_SRC_DIR/* -maxdepth 0 -type d -not -name core -not -empty); do
  NAME=$(basename $PLUG)

  if [ "$NAME" == "$PLUGIN_TO_COMPILE" ]; then
      echo "Building plugin $NAME to bin/plugins/$NAME/$NAME.so"
      go build -buildmode=plugin "$@" -o $PLUGIN_OUTPUT_DIR/$NAME/$NAME.so $PLUG/*.go
  fi

  if [ "$PLUGIN_TO_COMPILE" == "all" ]; then
      echo "Building plugin $NAME to bin/plugins/$NAME/$NAME.so"
      go build -buildmode=plugin "$@" -o $PLUGIN_OUTPUT_DIR/$NAME/$NAME.so $PLUG/*.go
  fi
done
