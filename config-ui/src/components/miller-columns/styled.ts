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

import styled from '@emotion/styled'

export const Container = styled.div<{ height?: number }>`
  display: flex;
  width: 100%;
  ${({ height }) => `height: ${height}px;`}
  overflow-x: auto;

  .items {
    flex: 1;
    margin: 0;
    padding: 0;
    list-style: none;
    border: 1px solid #dbe4fd;
    border-right: none;
    border-radius: 4px;

    &:last-child {
      border-right: 1px solid #dbe4fd;
    }

    & > .title {
      padding: 4px 12px;
      font-weight: 700;
      color: #292b3f;
    }
  }
`
