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

export const Transformation = styled.div`
  .issue-tracking {
    .issue-type {
      .title {
        margin-bottom: 8px;
      }

      .list {
        padding-left: 40px;
      }
    }
  }

  .bp4-form-group {
    display: flex;
    align-items: center;

    .bp4-label {
      flex: 0 0 140px;
    }

    .bp4-form-content {
      flex: auto;
    }
  }
`;

export const CICD = styled.div`
  h3 {
    margin-top: 16px;

    .switch {
      display: inline-flex;
      align-items: center;
      margin-left: 16px;
      font-size: 14px;
      font-weight: 500;

      .bp4-switch {
        margin-bottom: 0;
      }
    }
  }

  .text,
  .sub-text {
    display: flex;
    align-items: baseline;
    margin-bottom: 8px;
  }

  .sub-text {
    padding-left: 28px;
  }
`;
