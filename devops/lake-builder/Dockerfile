FROM golang:1.17-alpine3.15 as builder
#RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk update
RUN apk upgrade
#RUN apk add --update gcc=130.2.1_pre1-r3 g++=10.2.1_pre1-r3
RUN apk add --no-cache tzdata libgit2-dev gcc g++
