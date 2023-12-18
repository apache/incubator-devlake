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
import { Switch, InputNumber } from 'antd';

import { Block, ExternalLink } from '@/components';

import * as S from './styled';

interface Props {
  subLabel: string;
  learnMore: string;
  externalInfo: string;
  defaultValue: number;
  name: string;
  initialValue: number;
  value: number;
  error: string;
  setValue: (value: number) => void;
  setError: (value: string) => void;
}

export const ConnectionRateLimit = ({
  subLabel,
  learnMore,
  externalInfo,
  defaultValue,
  initialValue,
  value,
  setValue,
}: Props) => {
  const [checked, setChecked] = useState(true);

  useEffect(() => {
    setValue(initialValue);
  }, [initialValue]);

  useEffect(() => {
    setChecked(value ? true : false);
  }, [value]);

  const handleChange = (checked: boolean) => {
    if (checked) {
      setValue(defaultValue);
    } else {
      setValue(0);
    }
    setChecked(checked);
  };

  const handleChangeValue = (value: number | null) => {
    setValue(value ?? 0);
  };

  return (
    <Block
      title="Custom Rate Limit"
      description={
        <>
          {subLabel} {learnMore && <ExternalLink link={learnMore}>Learn more</ExternalLink>}
        </>
      }
    >
      <S.RateLimit>
        <Switch checked={checked} onChange={handleChange} />
        {checked && (
          <>
            <InputNumber min={0} value={value} onChange={handleChangeValue} />
            <span>requests/hour</span>
          </>
        )}
      </S.RateLimit>
      {checked && externalInfo && <S.RateLimitInfo dangerouslySetInnerHTML={{ __html: externalInfo }} />}
    </Block>
  );
};
