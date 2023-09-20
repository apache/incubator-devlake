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
  .timezone {
    margin-bottom: 16px;
  }

  .quick-selection {
    margin-bottom: 8px;
  }

  .time-selection {
    display: flex;
    align-items: center;

    strong {
      margin-left: 8px;
    }
  }

  .cron {
    display: flex;
  }
`;

export const Input = styled.div`
  display: flex;
  align-items: center;

  .bp5-form-group {
    margin-bottom: 4px;
  }

  .bp5-form-group + .bp5-form-group {
    margin-left: 8px;
  }

  .bp5-input {
    width: 60px;
  }
`;

export const Error = styled.div`
  color: #e34040;
`;
