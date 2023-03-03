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

import React, { useState, useEffect } from 'react';
import { FormGroup, RadioGroup, Radio, InputGroup } from '@blueprintjs/core';

import * as S from './styled';

type VersionType = 'cloud' | 'server';

interface Props {
  subLabel?: string;
  disabled?: boolean;
  name: string;
  multipleVersions?: Record<VersionType, string>;
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

  const handleChange = (e: React.FormEvent<HTMLInputElement>) => {
    const version = (e.target as HTMLInputElement).value as VersionType;
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
      <FormGroup label={<S.Label>{name} Version</S.Label>} labelInfo={<S.LabelInfo>*</S.LabelInfo>}>
        <RadioGroup inline selectedValue={version} onChange={handleChange}>
          <Radio value="cloud">{name} Cloud</Radio>
          <Radio value="server" disabled={!multipleVersions.server}>
            {name} Server {multipleVersions.server ? multipleVersions.server : '(to be supported)'}
          </Radio>
        </RadioGroup>
        {version === 'cloud' && (
          <p style={{ margin: 0 }}>
            If you are using {name} Cloud, you do not need to enter the endpoint URL, which is {multipleVersions.cloud}.
          </p>
        )}
        {version === 'server' && (
          <FormGroup
            label={<S.Label>Endpoint URL</S.Label>}
            labelInfo={<S.LabelInfo>*</S.LabelInfo>}
            subLabel={
              <S.LabelDescription>
                {subLabel ?? `If you are using ${name} Server, please enter the endpoint URL.`}
              </S.LabelDescription>
            }
          >
            <InputGroup placeholder="Your Endpoint URL" value={value} onChange={handleChangeValue} />
          </FormGroup>
        )}
      </FormGroup>
    );
  }

  return (
    <FormGroup
      label={<S.Label>Endpoint URL</S.Label>}
      labelInfo={<S.LabelInfo>*</S.LabelInfo>}
      subLabel={<S.LabelDescription>{subLabel ?? `Provide the ${name} instance API endpoint.`}</S.LabelDescription>}
    >
      <InputGroup disabled={disabled} placeholder="Your Endpoint URL" value={value} onChange={handleChangeValue} />
    </FormGroup>
  );
};
