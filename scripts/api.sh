#!/bin/sh

set -e

notes() {
    curl -v "$GITLAB_ENDPOINT/projects/8967944/merge_requests/1349/notes?system=false&per_page=100&page=1" \
        -H "Authorization: Bearer $GITLAB_AUTH"
}

commits() {
    SIZE=${1-100}
    PAGE=${2-1}
    PROJ=${3-8967944}
    curl -v "$GITLAB_ENDPOINT/projects/$PROJ/repository/commits?with_stats=true&per_page=$SIZE&page=$PAGE" \
        -H "Authorization: Bearer $GITLAB_AUTH"
}

"$@"
