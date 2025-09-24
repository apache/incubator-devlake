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
  bucket: string;
  error?: string;
  setValue: (value: string) => void;
  setError: (error: string) => void;
}

export const S3Config = ({ bucket, error, setValue, setError }: Props) => {
  
  const validateBucket = (value: string) => {
    if (!value) {
      return 'S3存储桶名称是必填项';
    }
    if (!/^[a-z0-9.-]+$/.test(value)) {
      return 'S3存储桶名称格式不正确，只能包含小写字母、数字、点和连字符';
    }
    if (value.length < 3 || value.length > 63) {
      return 'S3存储桶名称长度必须在3-63个字符之间';
    }
    return '';
  };

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setValue(value);
    setError(validateBucket(value));
  };

  return (
    <Block title="S3存储桶名称" description="请输入存储Q Developer数据的S3存储桶名称" required>
      <Input 
        style={{ width: 386 }} 
        placeholder="my-qdev-data-bucket" 
        value={bucket} 
        onChange={handleChange}
        status={error ? 'error' : ''}
      />
      {error && <div style={{ color: 'red', fontSize: '12px', marginTop: '4px' }}>{error}</div>}
    </Block>
  );
};