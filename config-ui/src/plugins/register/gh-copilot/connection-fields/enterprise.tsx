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
  initialValues: any;
  values: any;
  setValues: (value: any) => void;
}

export const Enterprise = ({ initialValues, values, setValues }: Props) => {
  useEffect(() => {
    setValues({ enterprise: initialValues.enterprise ?? '' });
  }, [initialValues.enterprise]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValues({ enterprise: e.target.value });
  };

  return (
    <Block
      title="Enterprise Slug"
      description='Optional. The GitHub enterprise slug (e.g. "my-enterprise"). When provided, enables enterprise-wide aggregate metrics and per-user daily metrics. Leave empty for organization-level metrics only.'
    >
      <Input
        style={{ width: 386 }}
        placeholder="e.g. my-enterprise"
        value={values.enterprise ?? ''}
        onChange={handleChange}
      />
    </Block>
  );
};
