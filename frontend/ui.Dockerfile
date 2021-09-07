FROM node:12-alpine

RUN mkdir -p /usr/src/app/frontend

WORKDIR /usr/src/app/frontend

COPY ./frontend/package.json /usr/src/app/frontend

RUN yarn

EXPOSE 4000

CMD [ "yarn", "dev" ]
