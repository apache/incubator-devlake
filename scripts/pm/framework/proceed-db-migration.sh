
#!/bin/sh

sh "$(dirname $0)/../vars/active-vars.sh"

curl -v $LAKE_ENDPOINT/proceed-db-migration | jq
