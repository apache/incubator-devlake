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

export const ConfigurationPanel = styled.div`
  .top {
    display: flex;
    align-items: flex-start;

    .block + .block {
      margin-left: 32px;
    }

    h3 {
      margin: 0 0 8px;
    }

    .detail {
      .bp4-icon {
        margin-left: 4px;
        cursor: pointer;
      }
    }
  }

  .bottom {
    margin-top: 32px;
  }
`

export const ConnectionColumn = styled.div`
  display: flex;
  align-items: center;

  img {
    margin-right: 4px;
    width: 20px;
  }
`

export const ActionColumn = styled.div`
  display: flex;
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
`
