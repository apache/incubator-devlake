FROM lake-builder:0.0.1 as builder

# docker build --build-arg GOPROXY=https://goproxy.io,direct -t lake .
ARG GOPROXY=
WORKDIR /app
COPY . /app

RUN CGO_ENABLE=1 GOOS=linux go build -o lake && sh scripts/compile-plugins.sh

FROM alpine:edge
EXPOSE 8080
COPY --from=builder /app/lake /bin/app/lake
COPY --from=builder /app/plugins/ /bin/app/plugins

WORKDIR /bin/app
CMD ["/bin/sh", "-c", "./lake"]
