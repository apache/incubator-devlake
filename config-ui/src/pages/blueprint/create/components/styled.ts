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

import styled from 'styled-components';
import { Colors } from '@blueprintjs/core';

export const Wrapper = styled.div`
  margin-top: 36px;

  .card + .card {
    margin-top: 16px;
  }
`;

export const ConnectionList = styled.ul`
  padding: 12px;

  li {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 4px 0;
    border-bottom: 1px solid #f0f0f0;

    .name {
      font-weight: 600;
    }

    .status {
      display: flex;
      align-items: center;

      &.online {
        color: ${Colors.GREEN3};
      }

      &.offline {
        color: ${Colors.RED3};
      }
    }
  }
`;

export const Tips = styled.p`
  margin: 24px 0 0;

  span:last-child {
    color: #7497f7;
    cursor: pointer;
  }
`;

export const ConnectionColumn = styled.div`
  display: flex;
  align-items: center;

  img {
    margin-right: 4px;
    width: 20px;
  }
`;

export const ScopeColumn = styled.ul``;

export const ScopeItem = styled.li`
  margin-bottom: 4px;

  &:last-child {
    margin-bottom: 0;
  }

  .bp4-button {
    margin-left: 6px;
  }
`;

export const Btns = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 36px;

  .bp4-button + .bp4-button {
    margin-left: 8px;
  }
`;
