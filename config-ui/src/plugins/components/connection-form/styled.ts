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

export const Wrapper = styled.div``;

export const Tips = styled.div`
  margin-bottom: 24px;
  padding: 24px;
  color: #3c5088;
  background: #f0f4fe;
  border: 1px solid #bdcefb;
  border-radius: 4px;
`;

export const Form = styled.div`
  .bp5-form-group label.bp5-label {
    margin: 0 0 8px 0;
  }

  .bp5-form-group .bp5-form-group-sub-label {
    margin: 0 0 8px 0;
  }

  .bp5-input-group {
    width: 386px;
  }

  .bp5-input {
    border: 1px solid #dbe4fd;
    box-shadow: none;
    border-radius: 4px;

    &::placeholder {
      color: #b8b8bf;
    }
  }
`;
