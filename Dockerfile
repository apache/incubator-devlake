FROM mericodev/lake-builder:0.0.3 as builder

# docker build --build-arg GOPROXY=https://goproxy.io,direct -t mericodev/lake .
ARG GOPROXY=
# docker build --build-arg HTTPS_PROXY=http://localhost:4780 -t mericodev/lake .
ARG HTTP_PROXY=
ARG HTTPS_PROXY=
WORKDIR /app
COPY . /app

RUN rm -rf /app/bin

ENV GOBIN=/app/bin

RUN go build -o bin/lake && sh scripts/compile-plugins.sh
RUN go install ./cmd/lake-cli/

FROM python:3.10.4-alpine3.15
RUN apk add --no-cache musl-dev libgit2-dev libffi-dev \
    && apk add --no-cache gcc
RUN pip3 install dbt-mysql

EXPOSE 8080

WORKDIR /app

COPY --from=builder /app/bin /app/bin
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

ENV PATH="/app/bin:${PATH}"

CMD ["lake"]
