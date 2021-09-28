FROM node:14-alpine

ARG REGISTRY=https://registry.npmjs.org/
COPY package.json /src/package.json
WORKDIR /src
RUN npm i --registry=${REGISTRY}

COPY . /src

CMD ["npm", "run", "dev"]

EXPOSE 4000
