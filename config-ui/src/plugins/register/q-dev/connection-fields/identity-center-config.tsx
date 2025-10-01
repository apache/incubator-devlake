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

import { ChangeEvent } from 'react';
import { Input } from 'antd';

import { Block } from '@/components';

interface Props {
  identityStoreId: string;
  identityStoreRegion: string;
  errors: {
    identityStoreId?: string;
    identityStoreRegion?: string;
  };
  setValues: (values: any) => void;
  setErrors: (errors: any) => void;
}

export const IdentityCenterConfig = ({ 
  identityStoreId, 
  identityStoreRegion, 
  errors, 
  setValues, 
  setErrors 
}: Props) => {
  
  const validateIdentityStoreId = (value: string) => {
    if (!value) {
      return 'Identity Store ID是必填项';
    }
    if (!/^d-[a-z0-9]{10}$/.test(value)) {
      return 'Identity Store ID格式不正确，应为：d-xxxxxxxxxx';
    }
    return '';
  };

  const validateIdentityStoreRegion = (value: string) => {
    if (!value) {
      return 'Identity Center区域是必填项';
    }
    if (!/^[a-z0-9-]+$/.test(value)) {
      return 'Identity Center区域格式不正确';
    }
    return '';
  };

  const handleIdentityStoreIdChange = (e: ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setValues({ identityStoreId: value });
    setErrors({ identityStoreId: validateIdentityStoreId(value) });
  };

  const handleIdentityStoreRegionChange = (e: ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setValues({ identityStoreRegion: value });
    setErrors({ identityStoreRegion: validateIdentityStoreRegion(value) });
  };

  return (
    <>
      <Block title="IAM Identity Store ID" description="请输入Identity Store ID，格式：d-xxxxxxxxxx" required>
        <Input 
          style={{ width: 386 }} 
          placeholder="d-1234567890" 
          value={identityStoreId} 
          onChange={handleIdentityStoreIdChange}
          status={errors.identityStoreId ? 'error' : ''}
        />
        {errors.identityStoreId && <div style={{ color: 'red', fontSize: '12px', marginTop: '4px' }}>{errors.identityStoreId}</div>}
      </Block>

      <Block title="IAM Identity Center区域" description="请输入IAM Identity Center所在的AWS区域" required>
        <Input 
          style={{ width: 386 }} 
          placeholder="us-east-1" 
          value={identityStoreRegion} 
          onChange={handleIdentityStoreRegionChange}
          status={errors.identityStoreRegion ? 'error' : ''}
        />
        {errors.identityStoreRegion && <div style={{ color: 'red', fontSize: '12px', marginTop: '4px' }}>{errors.identityStoreRegion}</div>}
      </Block>
    </>
  );
};