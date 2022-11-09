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

export const Wrapper = styled.label`
  display: inline-flex;
  align-items: center;
  cursor: pointer;

  .checkbox {
    position: relative;
    margin-right: 8px;

    &.checkbox-checked {
      .checkbox-inner {
        background-color: #7497f7;
        border-color: #7497f7;

        &::after {
          content: ' ';
          transform: rotate(45deg) scale(1) translate(-50%, -50%);
          opacity: 1;
          transition: all 0.2s cubic-bezier(0.12, 0.4, 0.29, 1.46) 0.1s;
        }
      }
    }

    &.checkbox-indeterminate {
      .checkbox-inner {
        background-color: #fff;
        border-color: #d9d9d9;

        &::after {
          content: ' ';
          left: 50%;
          width: 8px;
          height: 8px;
          background-color: #7497f7;
          border: 0;
          transform: translate(-50%, -50%) scale(1);
          opacity: 1;
        }
      }
    }

    .checkbox-input {
      position: absolute;
      z-index: 1;
      width: 100%;
      height: 100%;
      cursor: pointer;
      opacity: 0;
    }

    .checkbox-inner {
      position: relative;
      top: 0;
      left: 0;
      display: block;
      width: 14px;
      height: 14px;
      background-color: #fff;
      border: 1px solid #70727f;
      border-radius: 2px;
      transition: all 0.3s;

      &::after {
        content: ' ';
        position: absolute;
        top: 50%;
        left: 21.5%;
        display: table;
        width: 5.71428571px;
        height: 9.14285714px;
        border: 2px solid #fff;
        border-top: 0;
        border-left: 0;
        transform: rotate(45deg) scale(0) translate(-50%, -50%);
        opacity: 0;
        transition: all 0.1s cubic-bezier(0.71, -0.46, 0.88, 0.6), opacity 0.1s;
      }
    }
  }

  .text {
    font-size: 14px;
  }
`
