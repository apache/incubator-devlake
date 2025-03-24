<!--
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

# Q Developer 插件

本插件用于从AWS S3获取AWS Q Developer的使用数据，并进行处理和分析。

## 功能

- 从AWS S3指定前缀下获取CSV文件
- 解析CSV文件中的用户使用数据
- 按用户聚合数据，计算各项指标

## 配置

配置项包括：

1. AWS访问密钥ID
2. AWS秘钥
3. AWS区域
4. S3桶名称
5. 速率限制(每小时)

## 数据流程

插件包含以下任务：

1. `collectQDevS3Files`: 从S3收集文件元数据信息，不下载文件内容
2. `extractQDevS3Data`: 使用S3文件元数据，下载CSV数据并解析存储到数据库
3. `convertQDevUserMetrics`: 将用户数据转换为聚合指标，计算平均值和总值

## 数据表

- `_tool_q_dev_connections`: 存储AWS S3连接信息
- `_tool_q_dev_s3_file_meta`: 存储S3文件元数据
- `_tool_q_dev_user_data`: 存储从CSV文件中解析的用户数据
- `_tool_q_dev_user_metrics`: 存储聚合后的用户指标 