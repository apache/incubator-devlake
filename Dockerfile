FROM mericodev/lake-builder:0.0.2 as builder

# docker build --build-arg GOPROXY=https://goproxy.io,direct -t mericodev/lake .
ARG GOPROXY=
WORKDIR /app
COPY . /app

RUN CGO_ENABLE=1 GOOS=linux go build -o lake && sh scripts/compile-plugins.sh
RUN go install ./cmd/lake-cli/

FROM alpine:edge
EXPOSE 8080
COPY --from=builder /app/lake /bin/app/lake
COPY --from=builder /app/plugins/ /bin/app/plugins
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /go/bin/lake-cli /bin/lake-cli

WORKDIR /bin/app
CMD ["/bin/sh", "-c", "./lake"]
