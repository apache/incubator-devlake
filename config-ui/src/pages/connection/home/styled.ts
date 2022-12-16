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

import styled from 'styled-components'

export const Wrapper = styled.div`
  .block + .block {
    margin-top: 48px;
  }

  h2 {
    margin: 0 0 4px;
  }

  p {
    margin: 0 0 16px;
  }

  ul {
    margin: 0;
    padding: 0;
    list-style: none;
    display: flex;
    align-items: center;
  }

  li {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 8px 16px;
    border: 2px solid transparent;
    cursor: pointer;
    transition: all 0.2s linear;

    &:hover {
      background-color: #eeeeee;
      border-color: #7497f7;
      box-shadow: 0 2px 2px 0 rgb(0 0 0 / 16%), 0 0 2px 0 rgb(0 0 0 / 12%);
    }

    img {
      width: 50px;
    }

    span {
      margin-top: 4px;
    }
  }

  li + li {
    margin-left: 24px;
  }
`
