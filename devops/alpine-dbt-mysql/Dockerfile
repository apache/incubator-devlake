FROM --platform=linux/amd64 alpine:3.15
RUN apk add --no-cache musl-dev libgit2-dev libffi-dev \
    && apk add --no-cache gcc

ENV PYTHONUNBUFFERED=1
RUN apk add --update --no-cache python3-dev && ln -sf python3 /usr/bin/python
RUN python3 -m ensurepip
RUN pip3 install --no-cache --upgrade pip setuptools
RUN pip3 install dbt-mysql
RUN apk add --no-cache tar
