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
  h5 {
    margin-top: 12px;
    font-weight: 400;
  }

  .block + .block {
    margin-top: 36px;
  }

  ul {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
  }

  li {
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
    margin-top: 24px;
    margin-right: 30px;
    padding: 20px 0;
    width: 160px;
    border-radius: 8px;
    box-shadow: 0px 2.4px 4.8px -0.8px rgba(0, 0, 0, 0.1), 0px 1.6px 8px rgba(0, 0, 0, 0.07);
    box-sizing: border-box;
    cursor: pointer;
    transition: all 0.2s linear;

    &:hover {
      background-color: #eeeeee;
    }

    & > img {
      width: 60px;
      margin-bottom: 8px;
    }

    & > .name {
      position: relative;
      margin-bottom: 8px;
      padding-bottom: 8px;

      &::after {
        position: absolute;
        bottom: 0;
        left: 50%;
        margin-left: -44px;
        content: '';
        width: 88px;
        height: 1px;
        background-color: #dbdcdf;
      }
    }

    & > .bp4-tag {
      position: absolute;
      top: 0;
      right: 0;
    }
  }
`;

export const Count = styled.span`
  color: #70727f;
`;

export const DialogTitle = styled.div`
  display: flex;
  align-items: center;

  img {
    margin-right: 8px;
    width: 24px;
  }
`;
