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

import { ItemTypeEnum } from '../../types'

export const Wrapper = styled.div<{ selected: boolean; type: ItemTypeEnum }>`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 12px;

  ${({ type }) =>
    type === ItemTypeEnum.BRANCH
      ? `
    cursor: pointer;
    &:hover {
      background-color: #f5f5f7;
    }
    `
      : ''}

  ${({ selected }) => (selected ? 'background-color: #f5f5f7;' : '')}

  & > span.indicator {
    display: table;
    width: 6px;
    height: 6px;
    border: 1px solid #000;
    border-top: 0;
    border-left: 0;
    transform: rotate(-45deg);
  }
`
