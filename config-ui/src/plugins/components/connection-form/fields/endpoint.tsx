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

import { useState, useEffect } from 'react';
import type { RadioChangeEvent } from 'antd';
import { Radio, Input } from 'antd';

import { Block } from '@/components';

type VersionType = 'cloud' | 'server';

interface Props {
  subLabel?: string;
  disabled?: boolean;
  name: string;
  multipleVersions?: Record<VersionType, string>;
  cloudName?: string;
  initialValue: string;
  value: string;
  error: string;
  setValue: (value: string) => void;
  setError: (error: string) => void;
}

export const ConnectionEndpoint = ({
  subLabel,
  disabled = false,
  name,
  multipleVersions,
  cloudName,
  initialValue,
  value,
  setValue,
  setError,
}: Props) => {
  const [version, setVersion] = useState<VersionType>('cloud');

  useEffect(() => {
    setValue(initialValue);
    setVersion(initialValue === multipleVersions?.cloud ? 'cloud' : 'server');
  }, [initialValue]);

  useEffect(() => {
    setError(value ? '' : 'endpoint is required');
  }, [value]);

  const handleChange = (e: RadioChangeEvent) => {
    const version = e.target.value;
    if (version === 'cloud') {
      setValue(multipleVersions?.cloud ?? '');
    }

    if (version === 'server') {
      setValue('');
    }

    setVersion(version);
  };

  const handleChangeValue = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValue(e.target.value);
  };

  if (multipleVersions) {
    return (
      <>
        <Block title={name} required>
          <Radio.Group value={version} onChange={handleChange}>
            <Radio value="cloud">{cloudName ? cloudName : `${name} Cloud`}</Radio>
            <Radio value="server" disabled={!multipleVersions.server}>
              {name} Server {multipleVersions.server ? multipleVersions.server : '(to be supported)'}
            </Radio>
          </Radio.Group>
        </Block>
        {version === 'cloud' && (
          <Block>
            <p>
              If you are using {name} Cloud, you do not need to enter the endpoint URL, which is{' '}
              {multipleVersions.cloud}.
            </p>
          </Block>
        )}
        {version === 'server' && (
          <Block
            title="Endpoint URL"
            description={subLabel ?? `If you are using ${name} Server, please enter the endpoint URL.`}
            required
          >
            <Input style={{ width: 386 }} placeholder="Your Endpoint URL" value={value} onChange={handleChangeValue} />
          </Block>
        )}
      </>
    );
  }

  return (
    <Block title="Endpoint URL" description={subLabel ?? `Provide the ${name} instance API endpoint.`} required>
      <Input
        style={{ width: 386 }}
        disabled={disabled}
        placeholder="Your Endpoint URL"
        value={value}
        onChange={handleChangeValue}
      />
    </Block>
  );
};
