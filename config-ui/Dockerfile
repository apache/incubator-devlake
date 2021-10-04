FROM node:14-alpine

ARG REGISTRY=https://registry.npmjs.org/
COPY ./server/package.json /src/package.json

WORKDIR /src

RUN npm i --registry=${REGISTRY}

COPY out .
COPY ./server/server.js .

CMD ["npm", "run", "server"]

EXPOSE 4000
