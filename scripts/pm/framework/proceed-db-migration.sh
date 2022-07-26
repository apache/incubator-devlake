
#!/bin/sh

. "$(dirname $0)/../vars/active-vars.sh"

curl -sv $LAKE_ENDPOINT/proceed-db-migration | jq
