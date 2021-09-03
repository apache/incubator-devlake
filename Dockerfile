FROM golang:latest as builder

COPY . /go/src/lake

WORKDIR /go/src/lake
RUN go build -o lake && sh scripts/compile-plugins.sh

FROM alpine
EXPOSE 8080
RUN apk add --no-cache git make build-base
COPY --from=builder /go/src/lake/lake /bin/app/lake
COPY --from=builder /go/src/lake/plugins/ /bin/app/plugins

WORKDIR /bin/app
CMD ["/bin/sh", "-c", "./lake"]
