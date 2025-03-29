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

export const Top = styled.div`
  cursor: pointer;

  & > span.back {
    margin-left: 8px;
    text-decoration: underline;
  }
`;

export const ProjectSelector = styled.div`
  position: relative;
  padding: 16px 8px;

  & > h1 {
    display: flex;
    align-items: center;
    justify-content: space-between;
    cursor: pointer;
  }
`;

export const Selector = styled.div`
  position: absolute;
  top: 50px;
  right: 0;
  left: 0;
  padding: 16px 8px;
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  border-radius: 6px;
  z-index: 1;

  & > ul,
  & > p {
    margin-top: 16px;
  }

  & > ul > li {
    display: flex;
    padding: 8px 16px;
    align-items: center;
    min-height: 40px;
    transition: background-color 0.3s;
    cursor: pointer;

    &:hover {
      background-color: #f5f5f5;
    }
  }
`;
