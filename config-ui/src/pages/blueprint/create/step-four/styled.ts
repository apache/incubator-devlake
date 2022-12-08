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

export * from '../styled'

export const Tips = styled.p`
  margin: 24px 0 0;

  span:last-child {
    color: #7497f7;
    cursor: pointer;
  }
`

export const Help = styled.div`
  padding: 10px;
  width: 300px;
  font-size: 12px;

  .title {
    margin-bottom: 10px;
    font-size: 14px;
    font-weight: 700px;

    span.bp3-icon {
      margin-right: 4px;
    }
  }

  img {
    width: 100%;
  }
`
