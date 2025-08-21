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

import { useEffect } from 'react';
import { Input } from 'antd';

import { Block } from '@/components';

interface Props {
  type: 'create' | 'update';
  label?: string;
  subLabel?: string;
  placeholder?: string;
  initialValue: string;
  value: string;
  error: string;
  setValue: (value?: string) => void;
  setError: (value?: string) => void;
}

export const ConnectionPassword = ({
  type,
  label,
  subLabel,
  placeholder,
  initialValue,
  value,
  setValue,
  setError,
}: Props) => {
  useEffect(() => {
    setValue(type === 'create' ? initialValue : undefined);
  }, [type, initialValue]);

  useEffect(() => {
    setError(type === 'create' && !value ? 'password is required' : undefined);
  }, [type, value]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValue(e.target.value);
  };

  return (
    <Block title={label ?? 'Password'} description={subLabel ? subLabel : null} required>
      <Input.Password
        style={{ width: 386 }}
        placeholder={type === 'update' ? '********' : placeholder ? placeholder : 'Your Password'}
        value={value}
        onChange={handleChange}
      />
    </Block>
  );
};
