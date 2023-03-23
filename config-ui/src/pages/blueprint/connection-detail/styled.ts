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

export const Action = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;

  .bp4-button + .bp4-button {
    margin-left: 8px;
  }
`;

export const ActionDelete = styled.div`
  padding: 16px 24px;

  .btns {
    display: flex;
    align-items: center;
    justify-content: end;
    margin-top: 16px;
  }
`;

export const Entities = styled.div`
  margin-bottom: 24px;

  h4 {
    margin-bottom: 16px;
  }

  ul {
    display: flex;
    align-items: center;

    li::after {
      content: ',';
    }

    li:last-child::after {
      content: '';
    }

    li + li {
      margin-left: 4px;
    }
  }
`;

export const SelectTransformationWrapper = styled.div`
  .action {
    margin-bottom: 24px;
  }

  .btns {
    display: flex;
    align-items: center;
    justify-content: end;
    margin-top: 24px;

    .bp4-button + .bp4-button {
      margin-left: 4px;
    }
  }
`;
