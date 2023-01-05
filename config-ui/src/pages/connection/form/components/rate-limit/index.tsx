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
import { NumericInput, Switch } from '@blueprintjs/core';

import * as S from './styled';

interface Props {
  initialValue?: number;
  value?: number;
  onChange?: (value: number) => void;
}

export const RateLimit = ({ initialValue, value, onChange }: Props) => {
  const [show, setShow] = useState(false);

  useEffect(() => {
    setShow(value ? true : false);
  }, [value]);

  const handleChangeValue = (value: number) => {
    onChange?.(show ? value : 0);
  };

  const handleChangeShow = (e: React.FormEvent<HTMLInputElement>) => {
    const checked = (e.target as HTMLInputElement).checked;
    setShow(checked);
    if (!checked) {
      onChange?.(0);
    } else {
      onChange?.(initialValue ?? 0);
    }
  };

  return (
    <S.Wrapper>
      {show && <NumericInput value={value} onValueChange={handleChangeValue} />}
      <Switch checked={show} onChange={handleChangeShow} />
      <span>{show ? `Enabled - ${value} Requests/hr` : 'Disabled'}</span>
    </S.Wrapper>
  );
};
