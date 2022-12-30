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

import { Colors } from '@blueprintjs/core';
import styled from 'styled-components';

export const Wrapper = styled.div`
  .card + .card {
    margin-top: 16px;
  }

  .card:last-child {
    position: relative;
    padding: 24px 48px 24px 24px;

    .collapse-control {
      position: absolute;
      right: 12px;
      top: 24px;
    }
  }

  p.message {
    margin: 8px 0 0;

    &.error {
      color: ${Colors.RED3};
    }
  }
`;

export const Pipeline = styled.ul`
  margin: 0;
  padding: 0;
  list-style: none;
  display: flex;
  align-items: center;

  li {
    flex: 1;
    display: flex;
    flex-direction: column;

    &.success {
      color: ${Colors.GREEN3};

      .bp4-icon {
        color: ${Colors.GREEN3};
      }
    }

    &.error {
      color: ${Colors.RED3};

      .bp4-icon {
        color: ${Colors.RED3};
      }
    }

    & > span {
      font-size: 12px;
      color: #94959f;
    }

    & > strong {
      display: flex;
      align-items: center;
      margin-top: 8px;
    }
  }
`;

export const Inner = styled.div`
  overflow: auto;
`;

export const Header = styled.ul`
  padding: 0;
  margin: 0;
  list-style: none;
  display: flex;
  align-items: center;

  li {
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex: 0 0 45%;
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
  padding: 0;
  margin: 0;
  list-style: none;
  display: flex;
  align-items: flex-start;

  li {
    flex: 0 0 45%;
    padding-bottom: 8px;
    overflow: hidden;
  }

  li + li {
    margin-left: 16px;
  }
`;
