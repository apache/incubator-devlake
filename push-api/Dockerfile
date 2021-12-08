FROM mericodev/lake-builder:0.0.2 as builder

# TODO - THIS DOESN'T WORK YET

# docker build --build-arg GOPROXY=https://goproxy.io,direct -t mericodev/lake .
ARG GOPROXY=
WORKDIR /app
COPY . /app
ADD ./../go.mod /app
ADD ./../go.sum /app

RUN rm -rf /app/bin

ENV GOBIN=/app/bin

RUN CGO_ENABLE=1 GOOS=linux go build -o ./bin/push-api

FROM alpine:edge

EXPOSE 9123
WORKDIR /app

COPY --from=builder /app/bin /app/bin
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

ENV PATH="/app/bin:${PATH}"

CMD ["push-api"]
