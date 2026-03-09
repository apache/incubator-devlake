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

import { useEffect, type ChangeEvent } from 'react';
import { Button, Input } from 'antd';

import { Block } from '@/components';

interface Props {
  type: 'create' | 'update';
  initialValues: any;
  values: any;
  setValues: (value: any) => void;
  setErrors: (value: any) => void;
}

export const CustomHeaders = ({ type, initialValues, values, setValues }: Props) => {
  const headers: Array<{ key: string; value: string }> = values.customHeaders ?? [];

  useEffect(() => {
    setValues({ customHeaders: initialValues.customHeaders ?? [] });
  }, [type, initialValues.customHeaders]);

  const addHeader = () => {
    setValues({ customHeaders: [...headers, { key: '', value: '' }] });
  };

  const removeHeader = (index: number) => {
    setValues({ customHeaders: headers.filter((_, i) => i !== index) });
  };

  const updateHeader = (index: number, field: 'key' | 'value', newValue: string) => {
    setValues({
      customHeaders: headers.map((h, i) => (i === index ? { ...h, [field]: newValue } : h)),
    });
  };

  return (
    <Block
      title="Custom Headers"
      description="Add custom HTTP headers for middleware or proxy authentication (e.g. Ocp-Apim-Subscription-Key). Required when not using an Anthropic API Key."
    >
      {headers.map((header, index) => (
        <div key={index} style={{ display: 'flex', gap: 8, marginBottom: 8, alignItems: 'center' }}>
          <Input
            style={{ width: 180 }}
            placeholder="Header name"
            value={header.key}
            onChange={(e: ChangeEvent<HTMLInputElement>) => updateHeader(index, 'key', e.target.value)}
          />
          <Input.Password
            style={{ width: 180 }}
            placeholder={type === 'update' ? '********' : 'Header value'}
            value={header.value}
            onChange={(e: ChangeEvent<HTMLInputElement>) => updateHeader(index, 'value', e.target.value)}
          />
          <Button onClick={() => removeHeader(index)}>Remove</Button>
        </div>
      ))}
      <Button onClick={addHeader}>+ Add Header</Button>
    </Block>
  );
};
