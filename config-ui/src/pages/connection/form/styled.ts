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

import { Colors } from '@blueprintjs/core';
import styled from 'styled-components';

export const Wrapper = styled.div`
  .bp4-form-group {
    display: flex;
    align-items: start;
    justify-content: space-between;

    .bp4-label {
      flex: 0 0 200px;
      font-weight: 600;

      .bp4-popover2-target {
        display: inline;
        margin: 0;
        line-height: 1;
        margin-left: 4px;

        & > .bp4-icon {
          display: block;
        }
      }

      .bp4-text-muted {
        color: ${Colors.RED3};
      }
    }

    .bp4-form-content {
      flex: auto;
    }
  }

  .footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-top: 32px;
  }
`;

export const Label = styled.span`
  display: inline-flex;
  align-items: center;
`;

export const SwitchWrapper = styled.div`
  display: flex;
  align-items: center;
`;
