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

export const Container = styled.div`
  position: relative;
`;

export const Loading = styled.div`
  text-align: center;
`;

export const NoData = styled.div`
  text-align: center;

  img {
    display: inline-block;
  }
`;

export const Table = styled.table`
  table-layout: fixed;
  width: 100%;
  background-color: #fff;
  box-shadow: 0px 2.4px 4.8px -0.8px rgba(0, 0, 0, 0.1), 0px 1.6px 8px rgba(0, 0, 0, 0.07);
  border-radius: 8px;
  border-spacing: 0;
`;

export const THeader = styled.thead``;

export const TBody = styled.tbody``;

export const TR = styled.tr`
  &:last-child {
    td {
      border-bottom: none;
    }
  }
`;

export const TH = styled.th`
  padding: 12px 16px;
  border-bottom: 1px solid #dbe4fd;
`;

export const TD = styled.td`
  padding: 12px 16px;
  border-bottom: 1px solid #dbe4fd;
  word-break: break-word;
`;

export const TDEllipsis = styled.td`
  word-break: break-word;
`;

export const Mask = styled.div`
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  display: flex;
  align-items: center;
  justify-content: center;
`;
