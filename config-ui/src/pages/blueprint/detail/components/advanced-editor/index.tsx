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

import { ClearOutlined, CaretDownOutlined } from '@ant-design/icons';
import { Space, Input, Button, Dropdown } from 'antd';
import styled from 'styled-components';

import { ExternalLink } from '@/components';
import { DOC_URL } from '@/release';

import { EXAMPLE_CONFIG } from './example';

const Wrapper = styled.div`
  h2 {
    margin: 0;
    padding: 0;
    font-size: 16px;
    font-weight: 600;
  }

  h3 {
    margin: 0 0 8px;
    padding: 0;
    font-size: 14px;
    font-weight: 600;
  }

  p {
    margin: 0 0 8px;
  }

  textarea {
    margin-bottom: 8px;
    min-height: 240px;
    font-size: 12px;
    font-family: 'JetBrains Mono', monospace;
  }
`;

interface Props {
  value: string;
  onChange: (value: string) => void;
}

export const AdvancedEditor = ({ value, onChange }: Props) => {
  return (
    <Wrapper>
      <h3>Task Editor</h3>
      <p>
        <span>Enter JSON Configuration or preload from a template.</span>
        <ExternalLink link={DOC_URL.ADVANCED_MODE.EXAMPLES}>See examples</ExternalLink>
      </p>
      <Input.TextArea rows={6} value={value} onChange={(e) => onChange(e.target.value)} />
      <Space>
        <Button size="small" icon={<ClearOutlined />} onClick={() => onChange(JSON.stringify([[]], null, '  '))}>
          Reset
        </Button>
        <Dropdown
          menu={{
            items: EXAMPLE_CONFIG.map((it) => ({
              key: it.id,
              label: it.name,
            })),
            onClick: ({ key }) => {
              const config = EXAMPLE_CONFIG.find((it) => it.id === key)?.config;
              if (config) {
                onChange(JSON.stringify(config, null, '  '));
              }
            },
          }}
        >
          <Button size="small" icon={<CaretDownOutlined />}>
            Load Templates
          </Button>
        </Dropdown>
      </Space>
    </Wrapper>
  );
};
