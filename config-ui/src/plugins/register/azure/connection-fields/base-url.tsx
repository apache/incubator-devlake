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

import { Radio } from 'antd';

import { Block } from '@/components';

export const BaseURL = () => {
  return (
    <Block title="Azure DevOps Version" required>
      <Radio.Group value="cloud" onChange={() => {}}>
        <Radio value="cloud">Azure DevOps Cloud</Radio>
        <Radio value="server" disabled>
          Azure DevOps Server (not supported)
        </Radio>
      </Radio.Group>
      <p style={{ margin: 0 }}>If you are using Azure DevOps Cloud, you do not need to enter the endpoint URL.</p>
    </Block>
  );
};
