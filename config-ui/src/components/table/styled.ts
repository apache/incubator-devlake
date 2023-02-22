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

export const Table = styled.ul<{ loading: number }>`
  transition: opacity 0.3s linear;
  ${({ loading }) => (loading ? 'opacity: 0.2; ' : '')}
`;

export const Row = styled.li`
  display: flex;
  align-items: center;
  padding: 12px 16px;
  border-top: 1px solid #dbe4fd;

  & > span {
    flex: 1;
    overflow: hidden;
  }
`;

export const Header = styled(Row)`
  font-size: 14px;
  font-weight: 600;
  border-top: none;
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
