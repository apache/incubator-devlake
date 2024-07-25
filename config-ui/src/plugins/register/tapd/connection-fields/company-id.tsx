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

import { ChangeEvent, useEffect } from 'react';
import { Input } from 'antd';

import { Block } from '@/components';

interface Props {
  initialValue: number;
  value: number;
  error: string;
  setValue: (value: number) => void;
  setError: (value: string) => void;
}

export const CompanyId = ({ initialValue, value, setValue, setError }: Props) => {
  useEffect(() => {
    setValue(initialValue);
  }, [initialValue]);

  useEffect(() => {
    let error = '';

    if (!value) {
      error = 'company id is required';
    } else if (!/\d/.test(value.toString())) {
      error = 'company id is a number';
    }

    setError(error);
  }, [value]);

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    const value = +e.target.value;

    if (typeof value === 'number' && !isNaN(value)) {
      setValue(value);
    }
  };

  return (
    <Block title="Company ID" description="" required>
      <Input style={{ width: 386 }} placeholder="Company ID" value={value} onChange={handleChange} />
    </Block>
  );
};
