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

import React, { useEffect } from 'react';
import { FormGroup, InputGroup } from '@blueprintjs/core';

import * as S from './styled';

interface Props {
  name: string;
  initialValue: string;
  value: string;
  error: string;
  setValue: (value: string) => void;
  setError: (error: string) => void;
}

export const ConnectionName = ({ initialValue, value, setValue, setError }: Props) => {
  useEffect(() => {
    setValue(initialValue);
  }, [initialValue]);

  useEffect(() => {
    setError(value ? '' : 'name is required');
  }, [value]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValue(e.target.value);
  };

  return (
    <FormGroup
      label={<S.Label>Connection Name</S.Label>}
      labelInfo={<S.LabelInfo>*</S.LabelInfo>}
      subLabel={
        <S.LabelDescription>
          Give your connection a unique name to help you identify it in the future.
        </S.LabelDescription>
      }
    >
      <InputGroup placeholder="Your Connection Name" value={value} onChange={handleChange} />
    </FormGroup>
  );
};
