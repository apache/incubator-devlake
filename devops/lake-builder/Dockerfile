FROM golang:1.17.0-alpine3.13 as builder
RUN apk update
RUN apk upgrade
RUN apk add --update gcc=10.2.1_pre1-r3 g++=10.2.1_pre1-r3
RUN apk add --no-cache tzdata
