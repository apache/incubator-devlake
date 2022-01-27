FROM node:16 as builder
WORKDIR /home/node/code
COPY package.json /home/node/code
COPY package-lock.json /home/node/code
RUN npm ci
COPY . /home/node/code
RUN npm run build-production

FROM nginx:latest
#ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.9.0/wait /wait
#RUN chmod +x /wait
RUN rm /etc/nginx/conf.d/default.conf
COPY ./nginx.conf /etc/nginx/conf.d/default.conf.tpl
WORKDIR /usr/share/nginx/html
RUN rm -rf ./*
COPY --from=builder /home/node/code/dist/. ./
EXPOSE 80 443
RUN apt update && apt install -y apache2-utils
COPY ./nginx.sh /usr/bin/nginx.sh
RUN chmod +x /usr/bin/nginx.sh
CMD nginx.sh
