#!/bin/sh
FILE=.env
if [ -f "$FILE" ]; then
    export $(grep -v '^#' .env | xargs)
else 
    echo "$FILE does not exist."
fi
