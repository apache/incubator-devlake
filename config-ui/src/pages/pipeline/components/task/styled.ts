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

import { Colors } from '@blueprintjs/core';
import styled from 'styled-components';

export const Wrapper = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 0;
  height: 80px;
  border-bottom: 1px solid #dbe4fd;
  box-sizing: border-box;
`;

export const Info = styled.div`
  flex: auto;
  overflow: hidden;

  .title {
    display: flex;
    align-items: center;
    margin-bottom: 8px;

    & > img {
      width: 20px;
    }

    & > strong {
      margin: 0 4px;
    }

    & > span {
      flex: auto;
      overflow: hidden;
    }
  }

  p {
    padding-left: 26px;
    margin: 0;
    font-size: 12px;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;

    &.error {
      color: ${Colors.RED3};
    }
  }
`;

export const Duration = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 0 0 80px;
  text-align: right;

  .bp4-icon {
    margin-top: 4px;
    cursor: pointer;
  }
`;
