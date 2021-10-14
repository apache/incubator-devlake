FROM node:14-alpine

COPY package.json /src/package.json
WORKDIR /src
RUN npm i

COPY . /src

CMD ["npm", "run", "dev"]

EXPOSE 4000
