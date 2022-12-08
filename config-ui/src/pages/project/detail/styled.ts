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

export const Tabs = styled.ul`
  margin: 0;
  padding: 0;
  list-style: none;
  display: flex;
  align-items: center;
  justify-content: center;
`

export const Tab = styled.li<{ active: boolean; disabled?: boolean }>`
  margin-right: 24px;
  padding: 6px 0;
  font-size: 14px;
  border-bottom: 2px solid transparent;
  transition: all 0.3s ease;
  cursor: pointer;

  ${({ active }) =>
    active
      ? `
    color: #7497F7;
    border-color: #7497F7;
  `
      : ''}

  ${({ disabled }) =>
    disabled
      ? `
  color: #a1a1a1;
    cursor: no-drop;
  `
      : ''}

  &:last-child {
    margin-right: 0;
  }
`

export const Panel = styled.div`
  margin-top: 24px;
  background-color: #ffffff;
  box-shadow: 0px 2.4px 4.8px -0.8px rgba(0, 0, 0, 0.1),
    0px 1.6px 8px rgba(0, 0, 0, 0.07);
  border-radius: 4px;

  .blueprint,
  .webhook {
    padding: 24px;
    text-align: center;

    .logo > img {
      display: inline-block;
      width: 120px;
      height: 120px;
    }

    .desc {
      margin: 20px 0;
    }

    .action > .or {
      display: block;
      margin: 8px 0;
    }
  }

  .settings {
    padding: 24px;

    h3 {
      margin: 0;
      padding: 0;
    }

    .block {
      margin-bottom: 16px;
    }

    .bp3-input-group {
      width: 386px;
    }
  }
`
