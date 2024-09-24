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

import { useState, useEffect } from 'react';
import { Space, Select, Input } from 'antd';
import { useRequest } from '@mints/hooks';

import API from '@/api';
import { Loading } from '@/components';

interface Props {
  style?: React.CSSProperties;
  plugin: string;
  connectionId: ID;
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
}

export const Deployments = ({ style, plugin, connectionId, transformation, setTransformation }: Props) => {
  const [type, setType] = useState('select');

  const { loading, data } = useRequest(() => API.scopeConfig.deployments(plugin, connectionId), [plugin, connectionId]);

  useEffect(() => {
    if (transformation.envNamePattern) {
      setType('regex');
    }
  }, [transformation]);

  const handleChangeType = (t: string) => {
    if (t === 'regex') {
      setTransformation({
        ...transformation,
        envNameList: [],
      });
    }

    if (t === 'select') {
      setTransformation({
        ...transformation,
        envNamePattern: '',
      });
    }

    setType(t);
  };

  const handleChangeRegex = (e: React.ChangeEvent<HTMLInputElement>) => {
    const envNamePattern = e.target.value;
    setTransformation({
      ...transformation,
      envNamePattern,
      envNameList: [],
    });
  };

  const handleChangeSelect = (value: string[]) => {
    setTransformation({
      ...transformation,
      envNamePattern: '',
      envNameList: value,
    });
  };

  if (loading || !data) {
    return <Loading style={style} />;
  }

  return (
    <Space style={style}>
      <Select value={type} onChange={handleChangeType}>
        <Select.Option value="select">is one of</Select.Option>
        <Select.Option value="regex">matches</Select.Option>
      </Select>
      {type === 'regex' ? (
        <Input placeholder="(?i)prod(.*)" value={transformation.envNamePattern} onChange={handleChangeRegex} />
      ) : (
        <Select
          mode="tags"
          style={{ width: 180 }}
          maxTagCount={2}
          value={transformation.envNameList}
          onChange={handleChangeSelect}
        >
          {data.map((d) => (
            <Select.Option key={d} value={d}>
              {d}
            </Select.Option>
          ))}
        </Select>
      )}
    </Space>
  );
};
