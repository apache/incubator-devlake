#!/bin/sh

SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )"
FILE=$SCRIPT_DIR/../.env

if [ -f "$FILE" ]; then
    export $(grep -v '^#' $FILE | sed 's/#.*$//g' | xargs)
else 
    echo "$FILE does not exist."
fi
