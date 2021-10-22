FROM nginx:latest
RUN rm /etc/nginx/conf.d/default.conf
COPY ./nginx.conf /etc/nginx/conf.d/default.conf
WORKDIR /usr/share/nginx/html
RUN rm -rf ./*
COPY ./dist/* ./
EXPOSE 80 443
ENTRYPOINT ["nginx", "-g", "daemon off;"]
