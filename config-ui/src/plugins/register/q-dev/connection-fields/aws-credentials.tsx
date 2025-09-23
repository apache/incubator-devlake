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
  accessKeyId: string;
  secretAccessKey: string;
  region: string;
  errors: {
    accessKeyId?: string;
    secretAccessKey?: string;
    region?: string;
  };
  setValues: (values: any) => void;
  setErrors: (errors: any) => void;
}

export const AwsCredentials = ({ 
  accessKeyId, 
  secretAccessKey, 
  region, 
  errors, 
  setValues, 
  setErrors 
}: Props) => {
  
  const validateAccessKeyId = (value: string) => {
    if (!value) {
      return 'AWS Access Key ID是必填项';
    }
    if (!/^[A-Z0-9]{20}$/.test(value)) {
      return 'AWS Access Key ID格式不正确';
    }
    return '';
  };

  const validateSecretAccessKey = (value: string) => {
    if (!value) {
      return 'AWS Secret Access Key是必填项';
    }
    if (value.length < 40) {
      return 'AWS Secret Access Key长度不足';
    }
    return '';
  };

  const validateRegion = (value: string) => {
    if (!value) {
      return 'AWS区域是必填项';
    }
    if (!/^[a-z0-9-]+$/.test(value)) {
      return 'AWS区域格式不正确';
    }
    return '';
  };

  const handleAccessKeyIdChange = (e: ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setValues({ accessKeyId: value });
    setErrors({ accessKeyId: validateAccessKeyId(value) });
  };

  const handleSecretAccessKeyChange = (e: ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setValues({ secretAccessKey: value });
    setErrors({ secretAccessKey: validateSecretAccessKey(value) });
  };

  const handleRegionChange = (e: ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setValues({ region: value });
    setErrors({ region: validateRegion(value) });
  };

  return (
    <>
      <Block title="AWS Access Key ID" description="请输入您的AWS Access Key ID" required>
        <Input 
          style={{ width: 386 }} 
          placeholder="AKIAIOSFODNN7EXAMPLE" 
          value={accessKeyId} 
          onChange={handleAccessKeyIdChange}
          status={errors.accessKeyId ? 'error' : ''}
        />
        {errors.accessKeyId && <div style={{ color: 'red', fontSize: '12px', marginTop: '4px' }}>{errors.accessKeyId}</div>}
      </Block>

      <Block title="AWS Secret Access Key" description="请输入您的AWS Secret Access Key" required>
        <Input.Password 
          style={{ width: 386 }} 
          placeholder="wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" 
          value={secretAccessKey} 
          onChange={handleSecretAccessKeyChange}
          status={errors.secretAccessKey ? 'error' : ''}
        />
        {errors.secretAccessKey && <div style={{ color: 'red', fontSize: '12px', marginTop: '4px' }}>{errors.secretAccessKey}</div>}
      </Block>

      <Block title="AWS区域" description="请输入AWS区域，例如：us-east-1" required>
        <Input 
          style={{ width: 386 }} 
          placeholder="us-east-1" 
          value={region} 
          onChange={handleRegionChange}
          status={errors.region ? 'error' : ''}
        />
        {errors.region && <div style={{ color: 'red', fontSize: '12px', marginTop: '4px' }}>{errors.region}</div>}
      </Block>
    </>
  );
};