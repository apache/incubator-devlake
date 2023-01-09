# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
version: "3"
services:
  mysql:
    image: mysql:8
    volumes:
      - mysql-storage:/var/lib/mysql
    restart: always
    ports:
      - 127.0.0.1:3306:3306
    environment:
      MYSQL_ROOT_PASSWORD: admin
      MYSQL_DATABASE: lake
      MYSQL_USER: merico
      MYSQL_PASSWORD: merico

  grafana:
    image: apache/devlake-dashboard:v0.15.0
    ports:
      - 3002:3000
    volumes:
      - grafana-storage:/var/lib/grafana
    environment:
      GF_SERVER_ROOT_URL: "http://localhost:4000/grafana"
      GF_USERS_DEFAULT_THEME: "light"
      MYSQL_URL: mysql:3306
      MYSQL_DATABASE: lake
      MYSQL_USER: merico
      MYSQL_PASSWORD: merico
    restart: always
    depends_on:
      - mysql

  devlake:
    image: apache/devlake:v0.15.0
    ports:
      - 127.0.0.1:8080:8080
    restart: always
    volumes:
      - ./.env:/app/.env
      - ./logs:/app/logs
    environment:
      LOGGING_DIR: /app/logs
    depends_on:
      - mysql

  config-ui:
    image: apache/devlake-config-ui:v0.15.0
    ports:
      - 127.0.0.1:4000:4000
    env_file:
      - ./.env
    environment:
      DEVLAKE_ENDPOINT: devlake:8080
      GRAFANA_ENDPOINT: grafana:3000
      #ADMIN_USER: devlake
      #ADMIN_PASS: merico
    depends_on:
      - devlake

volumes:
  mysql-storage:
  grafana-storage:

