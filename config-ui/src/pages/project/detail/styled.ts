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

export const Wrapper = styled.div`
  & > .bp4-tabs > .bp4-tab-list {
    display: flex;
    justify-content: center;
  }

  .blueprint,
  .webhook {
    padding: 24px;
    text-align: center;

    .logo > img {
      display: inline-block;
      width: 120px;
      height: 120px;
    }

    .desc {
      margin: 20px 0;
    }

    .action > .or {
      display: block;
      margin: 8px 0;
    }
  }

  .settings {
    padding: 24px;

    h3 {
      margin: 0;
      padding: 0;
    }

    .block {
      margin-bottom: 16px;
    }

    .bp3-input-group {
      width: 386px;
    }
  }
`
