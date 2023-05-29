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
  .block + .block {
    margin-top: 24px;
  }
`;

export const Input = styled.div`
  display: flex;
  align-items: center;

  .bp4-form-group + .bp4-form-group {
    margin-left: 8px;
  }

  .bp4-input {
    width: 60px;
  }
`;

export const Error = styled.div`
  color: #e34040;
`;

export const FromTimeWrapper = styled.div`
  .quick-selection {
    margin-bottom: 16px;

    & > .bp4-tag {
      cursor: pointer;
    }
  }

  .time-selection {
    display: flex;
    align-items: center;

    & > strong {
      margin-left: 4px;
    }
  }
`;
