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

import styled from 'styled-components'

export const ScopeList = styled.ul`
  list-style: none;
  margin: 0;
  padding: 0;
`

export const ScopeItem = styled.li`
  margin-bottom: 4px;

  &:last-child {
    margin-bottom: 0;
  }
`

export const ScopeItemMap = styled(ScopeItem)`
  margin-bottom: 12px;

  .name {
    display: flex;
    align-items: center;
    font-size: 12px;

    & > span {
      font-weight: 600;
    }

    .action {
      margin-left: 8px;

      & > span {
        font-size: 11px;
        color: #7497f7;
        cursor: pointer;

        &:hover {
          color: #106ba3;
        }
      }

      span + span {
        margin-left: 4px;
      }
    }
  }

  ul {
    margin: 0;
    padding: 0;
    list-style: none;
    padding-left: 24px;
    margin-top: 4px;
  }
`
