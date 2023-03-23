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
  padding-bottom: 24px;
`;

export const ConfigurationPanel = styled.div`
  .top {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;

    ul {
      display: flex;
      align-items: flex-start;

      li + li {
        margin-left: 48px;
      }
    }

    .detail {
      .bp4-icon {
        margin-left: 4px;
        cursor: pointer;
      }
    }
  }

  .bottom {
    margin-top: 24px;

    h3 {
      display: flex;
      align-items: center;
      justify-content: space-between;
    }

    .btns {
      margin-top: 16px;
      text-align: right;
    }
  }
`;

export const ConnectionColumn = styled.div`
  display: flex;
  align-items: center;

  img {
    margin-right: 4px;
    width: 20px;
  }
`;

export const ActionColumn = styled.div`
  display: inline-flex;
  flex-direction: column;
  align-items: flex-start;

  .item + .item {
    margin-top: 8px;
  }

  .item {
    cursor: pointer;

    .bp4-icon {
      margin-right: 4px;
    }
  }
`;

export const StatusPanel = styled.div`
  & > .info {
    display: flex;
    justify-content: flex-end;
    align-items: center;

    & > span {
      margin-left: 16px;
    }

    .bp4-switch {
      margin-bottom: 0;
    }
  }

  .block + .block {
    margin-top: 32px;
  }
`;

export const JenkinsTips = styled.div`
  position: fixed;
  right: 0;
  bottom: 0;
  left: 200px;
  background-color: #3c5088;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 36px;

  p {
    margin: 0;
    color: #fff;
  }
`;
