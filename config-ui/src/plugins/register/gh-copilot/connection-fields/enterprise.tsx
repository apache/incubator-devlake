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

import { useEffect, useMemo } from 'react';
import { Input } from 'antd';

import { Block } from '@/components';

import * as S from './styled';

interface Props {
  type: 'create' | 'update';
  initialValues: any;
  values: any;
  setValues: (value: any) => void;
  setErrors: (value: any) => void;
}

export const Enterprise = ({ type, initialValues, values, setValues, setErrors }: Props) => {
  useEffect(() => {
    setValues({ enterprise: initialValues.enterprise ?? '' });
  }, [initialValues.enterprise]);

  const error = useMemo(() => {
    if (type === 'update') return '';
    const hasOrg = !!values.organization?.trim();
    const hasEnt = !!values.enterprise?.trim();
    return hasOrg || hasEnt ? '' : 'At least one of Organization or Enterprise Slug is required';
  }, [type, values.organization, values.enterprise]);

  useEffect(() => {
    setErrors({ scopeIdentity: error });
  }, [error]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValues({ enterprise: e.target.value });
  };

  return (
    <Block
      title="Enterprise Slug"
      description="Enter the GitHub enterprise slug for enterprise-wide aggregate metrics and per-user data. At least one of Organization or Enterprise Slug is required. For the most complete data, provide both."
    >
      <Input
        style={{ width: 386 }}
        placeholder="e.g. my-enterprise"
        status={error ? 'error' : ''}
        value={values.enterprise ?? ''}
        onChange={handleChange}
      />
      {error && <S.ErrorText>{error}</S.ErrorText>}
    </Block>
  );
};
