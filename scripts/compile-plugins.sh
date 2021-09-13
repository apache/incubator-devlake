#!/bin/sh

SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )"
PLUGIN_SRC_DIR=$SCRIPT_DIR/../plugins
PLUGIN_OUTPUT_DIR=$SCRIPT_DIR/../bin/plugins


for PLUG in $(find $PLUGIN_SRC_DIR/* -maxdepth 0 -type d -not -name core -not -empty); do
  NAME=$(basename $PLUG)
  echo "Building plugin $NAME to bin/plugins/$NAME/$NAME.so"
  go build -buildmode=plugin "$@" -o $PLUGIN_OUTPUT_DIR/$NAME/$NAME.so $PLUG/*.go
done
