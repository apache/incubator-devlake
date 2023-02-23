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

export const ScopeList = styled.ul``;

export const ScopeItem = styled.li`
  margin-bottom: 4px;

  &:last-child {
    margin-bottom: 0;
  }

  .bp4-button {
    margin-left: 6px;
  }
`;

export const ScopeItemMap = styled(ScopeItem)`
  margin-bottom: 12px;

  .title {
    display: flex;
    align-items: center;
    font-size: 12px;

    span {
      font-weight: 600;
    }

    span.bp4-icon {
      margin-right: 6px;
    }
  }

  & > ul {
    padding-left: 24px;
    margin-top: 4px;

    li {
      display: flex;
      align-items: center;
      justify-content: space-between;

      span.name {
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }
    }
  }
`;
