FROM alpine:edge as builder
RUN apk update
RUN apk upgrade
RUN apk add --update go=1.16.7-r0 gcc=10.3.1_git20210625-r1 g++=10.3.1_git20210625-r1

WORKDIR /app
COPY . /app

RUN CGO_ENABLE=1 GOOS=linux go build -o lake && sh scripts/compile-plugins.sh

FROM alpine:edge
EXPOSE 8080
COPY --from=builder /app/lake /bin/app/lake
COPY --from=builder /app/plugins/ /bin/app/plugins

WORKDIR /bin/app
CMD ["/bin/sh", "-c", "./lake"]
