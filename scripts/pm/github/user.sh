#!/bin/sh

. "$(dirname $0)/../vars/active-vars.sh"

curl -sv $GITHUB_ENDPOINT/user \
    -H "Authorization: Token $GITHUB_TOKEN" \
    | jq
