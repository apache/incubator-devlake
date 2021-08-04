#!/bin/sh

docker-compose up -d

npm i
npx sequelize-cli db:migrate

export ENRICHMENT_PORT=43000
export ENRICHMENT_HOST=0.0.0.0
export COLLECTION_PORT=43001
export COLLECTION_HOST=0.0.0.0
npm run all