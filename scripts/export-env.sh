#!/bin/sh
FILE=.env
if [ -f "$FILE" ]; then
    export $(grep -v '^#' .env | sed 's/#.*$//g' | xargs)
else 
    echo "$FILE does not exist."
fi
