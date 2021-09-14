FROM mericodev/lake-builder:0.0.2 as builder

# docker build --build-arg GOPROXY=https://goproxy.io,direct -t mericodev/lake .
ARG GOPROXY=
WORKDIR /app
COPY . /app

RUN rm -rf /app/bin

ENV GOBIN=/app/bin

RUN CGO_ENABLE=1 GOOS=linux go build -o bin/lake && sh scripts/compile-plugins.sh
RUN go install ./cmd/lake-cli/

FROM alpine:edge

EXPOSE 8080
WORKDIR /app

COPY --from=builder /app/bin /app/bin
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

ENV PATH="/app/bin:${PATH}"

CMD ["lake"]
