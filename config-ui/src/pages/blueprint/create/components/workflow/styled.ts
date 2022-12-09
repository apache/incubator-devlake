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

export const List = styled.ul`
  display: flex;
  align-items: center;
  list-style: none;
  padding: 0;
  margin: 0;
`

export const Item = styled.li<{ active?: boolean }>`
  position: relative;
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;

  &::before {
    content: '';
    position: absolute;
    top: 15px;
    left: -15px;
    width: 50%;
    height: 2px;
    background-color: #bdcefb;
  }

  &::after {
    content: '';
    position: absolute;
    top: 15px;
    right: -15px;
    width: 50%;
    height: 2px;
    background-color: #bdcefb;
  }

  &:first-child::before {
    display: none;
  }

  &:last-child::after {
    display: none;
  }

  span.step {
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 8px;
    width: 30px;
    height: 30px;
    font-size: 16px;
    color: #7497f7;
    font-weight: 600;
    background-color: #f0f4fe;
    border-radius: 50%;

    ${({ active }) =>
      active
        ? `
  color: #fff;
  background-color: #7497f7;
  `
        : `
  color: #7497f7;
  background-color: #f0f4fe;
      `}
  }

  span.name {
    color: #70727f;
    font-weight: 600;
  }
`
