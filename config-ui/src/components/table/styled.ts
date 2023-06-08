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

export const Table = styled.table`
  table-layout: fixed;
  width: 100%;
  background-color: #fff;
  border-radius: 4px;
  border-spacing: 0;
`;

export const THeader = styled.thead`
  background-color: #f0f4fe;
`;

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
  font-weight: 500;
  border-bottom: 1px solid #dbdcdf;

  label.bp4-control {
    margin-bottom: 0;
  }
`;

export const TD = styled.td`
  padding: 12px 16px;
  border-bottom: 1px solid #dbdcdf;
  word-break: break-word;

  label.bp4-control {
    margin-bottom: 0;
  }
`;
