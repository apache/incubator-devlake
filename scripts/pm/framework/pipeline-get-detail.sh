#!/bin/sh

. "$(dirname $0)/../vars/active-vars.sh"

pipeline_id=${1-"2"}

curl -sv $LAKE_ENDPOINT/pipelines/$pipeline_id | jq
