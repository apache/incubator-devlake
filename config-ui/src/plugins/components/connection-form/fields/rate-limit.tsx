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
import { FormGroup, Switch, NumericInput } from '@blueprintjs/core';

import { ExternalLink } from '@/components';

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

  const handleChange = (e: React.FormEvent<HTMLInputElement>) => {
    const checked = (e.target as HTMLInputElement).checked;
    if (checked) {
      setValue(defaultValue);
    } else {
      setValue(0);
    }
    setChecked(checked);
  };

  const handleChangeValue = (value: number) => {
    setValue(value);
  };

  return (
    <FormGroup
      label={<S.Label>Custom Rate Limit</S.Label>}
      subLabel={
        <S.LabelDescription>
          {subLabel} {learnMore && <ExternalLink link={learnMore}>Learn more</ExternalLink>}
        </S.LabelDescription>
      }
    >
      <S.RateLimit>
        <Switch checked={checked} onChange={handleChange} />
        {checked && (
          <>
            <NumericInput buttonPosition="none" min={0} value={value} onValueChange={handleChangeValue} />
            <span>requests/hour</span>
          </>
        )}
      </S.RateLimit>
      {checked && externalInfo && <S.RateLimitInfo dangerouslySetInnerHTML={{ __html: externalInfo }} />}
    </FormGroup>
  );
};
