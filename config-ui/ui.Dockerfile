FROM node:12-alpine

RUN mkdir -p /usr/src/app/config-ui

WORKDIR /usr/src/app/config-ui

COPY ./config-ui/package.json /usr/src/app/config-ui

RUN yarn

EXPOSE 4000

CMD [ "yarn", "dev" ]
