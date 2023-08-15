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

export const List = styled.ul`
  display: flex;
  align-items: center;
  justify-content: flex-end;
  margin-top: 16px;
`;

export const Item = styled.li<{ active?: boolean; disabled?: boolean }>`
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 8px;
  width: 30px;
  height: 30px;
  color: #7497f7;
  border: 1px solid #7497f7;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s ease-in-out;

  &:hover {
    color: #fff;
    background-color: #7497f7;
  }

  ${({ active }) =>
    active
      ? `
    color: #fff;
    background-color: #7497f7;
    cursor: no-drop;
  `
      : ''}

  ${({ disabled }) =>
    disabled
      ? `
    color: #a1a1a1;
    border-color: #a1a1a1;
    cursor: no-drop;

    &:hover {
      color: #a1a1a1;
      border-color: #a1a1a1;
      background-color: transparent;
    }
  `
      : ''}

  &:last-child {
    margin-right: 0;
  }
`;
