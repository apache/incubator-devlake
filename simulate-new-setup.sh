#!/bin/sh


cleanup() {
    # CON_PATH=$(npx which node)
    PID=$(pgrep -f "concurrently.js")
    if [ -n "$PID" ]; then
        kill -s TERM $PID
        wait $PID || true
    fi
}

# clean previous residual if any
cleanup
# in case Ctrl-C was hit during process
trap cleanup exit

# stop on error
set -e

# reset databases
docker-compose down
sudo rm -rf data postgres-data

# kick off all services
docker-compose up -d
npm i
npx sequelize-cli db:migrate
node src/plugins/index.js
nodemon concurrently.js | tee /tmp/lake.log &

# wait untill services ready
echo waiting services to be ready
while ! curl localhost:3000 2>/dev/null; do
    printf '.'
    sleep 1
done

# trigger process with a small data set
curl -X POST "http://localhost:3001/" -H 'content-type: application/json' -H "x-token: mytoken" \
    -d '{"jira":{"boardId": 29}, "gitlab": {"projectId": 24547305}}'

# wait for finishing
while ! (grep -F 'jira enriching done!' /tmp/lake.log && grep -F 'gitlab enriching done!' /tmp/lake.log); do
    sleep 1
done

# clean up normally
cleanup
sleep 1

# print success message
echo
echo
echo "CONGURATULATION!! ALL PASSED!!!"