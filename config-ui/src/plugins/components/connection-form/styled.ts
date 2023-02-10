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
  margin-bottom: 36px;
  padding: 24px;
  color: #3c5088;
  background: #f0f4fe;
  border: 1px solid #bdcefb;
  border-radius: 4px;
`;

export const Form = styled.div`
  padding: 24px;
  background: #ffffff;
  box-shadow: 0px 2.4px 4.8px -0.8px rgba(0, 0, 0, 0.1), 0px 1.6px 8px rgba(0, 0, 0, 0.07);
  border-radius: 8px;

  .bp4-form-group label.bp4-label {
    margin: 0 0 8px 0;
  }

  .bp4-form-group .bp4-form-group-sub-label {
    margin: 0 0 8px 0;
  }

  .bp4-input-group {
    width: 386px;
  }

  .bp4-input {
    border: 1px solid #dbe4fd;
    box-shadow: none;
    border-radius: 4px;

    &::placeholder {
      color: #b8b8bf;
    }
  }

  .btns {
    display: flex;
    justify-content: end;

    .bp4-button + .bp4-button {
      margin-left: 8px;
    }
  }
`;
