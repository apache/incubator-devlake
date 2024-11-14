FROM mcr.microsoft.com/devcontainers/go:1-1.22-bookworm

ENV PKG_CONFIG_PATH=$PKG_CONFIG_PATH:/usr/local/lib:/usr/local/lib/pkgconfig
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib
ENV DEVLAKE_PLUGINS=bamboo,bitbucket,circleci,customize,dora,gitextractor,github,github_graphql,gitlab,jenkins,jira,org,pagerduty,refdiff,slack,sonarqube,trello,webhook

RUN apt-get update -y
RUN apt-get install pkg-config python3-dev default-libmysqlclient-dev build-essential libpq-dev cmake -y
RUN git clone -b v1.3.0 https://github.com/libgit2/libgit2.git && cd libgit2 && mkdir build && cd build && cmake .. && make && make install
