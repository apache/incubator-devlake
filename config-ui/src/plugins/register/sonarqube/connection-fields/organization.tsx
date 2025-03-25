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

import { Input } from 'antd';

import { Block, ExternalLink } from '@/components';
import { useEffect } from 'react';

interface Props {
  type: 'create' | 'update';
  initialValues: any;
  values: any;
  errors: any;
  setValues: (value: any) => void;
  setErrors: (value: any) => void;
}

export const Organization = ({ initialValues, values, setValues, setErrors }: Props) => {
  const { endpoint } = values;

  useEffect(() => {
    setValues({ org: initialValues.org });
  }, [initialValues.org]);

  useEffect(() => {
    setErrors({
      org: (values.endpoint != 'https://sonarcloud.io/api/' || values.org) ? '' : 'organization is required',
    });
  }, [values.org]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValues({
      org: e.target.value,
    });
  };

  if (endpoint !== 'https://sonarcloud.io/api/') {
    return null;
  }

  return (
    <Block
      title="Organization"
      description={
        <>
          Copy the organization key at{' '}
          <ExternalLink link="https://sonarcloud.io/account/organizations">here</ExternalLink>. If you have more than
          one, please create another connection.
        </>
      }
      required
    >
      <Input style={{ width: 386 }} placeholder="e.g. org-1" value={values.org} onChange={handleChange} />
    </Block>
  );
};
