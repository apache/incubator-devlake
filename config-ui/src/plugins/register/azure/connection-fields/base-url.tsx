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

import React from 'react';
import { FormGroup, RadioGroup, Radio } from '@blueprintjs/core';

import * as S from './styled';

export const BaseURL = () => {
  return (
    <FormGroup label={<S.Label>Azure DevOps Version</S.Label>} labelInfo={<S.LabelInfo>*</S.LabelInfo>}>
      <RadioGroup inline selectedValue="cloud" onChange={() => {}}>
        <Radio value="cloud">Azure DevOps Cloud</Radio>
        <Radio value="server" disabled>
          Azure DevOps Server (not supported)
        </Radio>
      </RadioGroup>
      <p style={{ margin: 0 }}>If you are using Azure DevOps Cloud, you do not need to enter the endpoint URL.</p>
    </FormGroup>
  );
};
