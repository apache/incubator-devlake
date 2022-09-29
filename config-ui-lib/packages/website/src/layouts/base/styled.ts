/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
import styled from '@emotion/styled';
import { Layout } from 'antd';

export const Container = styled(Layout)`
  min-height: 100vh;

  .logo {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 16px 0;

    img + img {
      margin-top: 8px;
    }
  }

  .trigger {
    padding: 0 24px;
    line-height: 64px;
    font-size: 18px;
    cursor: pointer;
    transition: color 0.3s;
  }

  .other-info {
    display: flex;
    align-items: center;
    margin-right: 20px;

    li {
      margin-right: 12px;

      &:last-child {
        margin-right: 0;
      }
    }
  }

  .copyright {
    display: flex;
    justify-content: center;
    color: #7c7c7c;
    font-size: 12px;
  }

  .ant-layout-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0;
    height: 50px;
    background-color: #fff;
  }

  .ant-layout-content {
    padding: 16px 24px;
  }
`;
