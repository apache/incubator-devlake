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
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-top: 100px;
  height: 100vh;
  background-color: #f9f9fa;
  box-sizing: border-box;

  .inner {
    margin: 32px auto 0;
    width: 640px;
    text-align: center;

    h2 {
      display: flex;
      align-items: center;
      justify-content: center;
      margin: 0;

      .bp4-icon {
        margin-right: 4px;
      }

      strong {
        margin-left: 4px;
      }
    }

    .path {
      margin: 8px 0;
    }

    p {
      margin: 0 0 16px 0;
    }
  }
`
