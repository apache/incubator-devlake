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

export const DataScope = styled.div`
  h4 {
    margin-top: 16px;
  }
`;

export const TransformationWrapper = styled.div`
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

  .ci-cd {
    h3 {
      margin-top: 16px;

      .bp4-tag {
        margin-left: 4px;
      }
    }

    .radio {
      padding-left: 20px;
      margin-bottom: 16px;

      .input {
        display: flex;
        align-items: center;

        & + .input {
          margin-top: 8px;
        }

        p {
          color: #292b3f;
        }

        .bp4-input-group {
          margin: 0 4px;
        }
      }
    }
  }

  .additional-settings {
    h2 {
      display: flex;
      align-items: center;
      cursor: pointer;
    }

    .radio {
      display: flex;
      align-items: center;
      margin: 8px 0 16px;

      p {
        margin: 0;
      }

      .bp4-control {
        margin: 0;
      }
    }

    .refdiff {
      display: flex;
      align-items: center;
      padding-left: 20px;

      .bp4-input-group {
        margin: 0 8px;
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
