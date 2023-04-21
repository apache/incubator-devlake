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

export const Wrapper = styled.div`
  position: relative;
  padding-right: 36px;

  .collapse-control {
    position: absolute;
    right: 0;
    top: 0;
  }
`;

export const Inner = styled.div`
  overflow: auto;
`;

export const Header = styled.ul`
  display: flex;
  align-items: center;

  li {
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex: 0 0 30%;
    padding: 8px 12px;

    &.ready,
    &.cancel {
      color: #94959f;
      background-color: #f9f9fa;
    }

    &.loading {
      color: #7497f7;
      background-color: #e9efff;
    }

    &.success {
      color: #4db764;
      background-color: #edfbf0;
    }

    &.error {
      color: #e34040;
      background-color: #feefef;
    }
  }

  li + li {
    margin-left: 16px;
  }
`;

export const Tasks = styled.ul`
  display: flex;
  align-items: flex-start;

  li {
    flex: 0 0 30%;
    padding-bottom: 8px;
    overflow: hidden;
  }

  li + li {
    margin-left: 16px;
  }
`;
