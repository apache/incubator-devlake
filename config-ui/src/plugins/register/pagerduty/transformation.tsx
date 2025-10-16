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

import { Form, Select } from 'antd';
import { HelpTooltip } from '@/components';

interface Props {
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
}

export const PagerDutyTransformation = ({ transformation, setTransformation }: Props) => {
  const priorityOptions = [
    { label: 'P1', value: 'P1' },
    { label: 'P2', value: 'P2' },
    { label: 'P3', value: 'P3' },
    { label: 'P4', value: 'P4' },
    { label: 'P5', value: 'P5' },
  ];

  const urgencyOptions = [
    { label: 'High', value: 'high' },
    { label: 'Low', value: 'low' },
  ];

  return (
    <>
      <Form.Item
        label={
          <>
            Priority Filter
            <HelpTooltip content="Filter incidents by priority levels. Leave empty to collect all priorities." />
          </>
        }
      >
        <Select
          mode="multiple"
          allowClear
          placeholder="Select priority levels to filter (leave empty for all)"
          options={priorityOptions}
          value={transformation.priorityFilter ?? []}
          onChange={(value) =>
            setTransformation({
              ...transformation,
              priorityFilter: value,
            })
          }
        />
      </Form.Item>

      <Form.Item
        label={
          <>
            Urgency Filter
            <HelpTooltip content="Filter incidents by urgency levels. Leave empty to collect all urgencies." />
          </>
        }
      >
        <Select
          mode="multiple"
          allowClear
          placeholder="Select urgency levels to filter (leave empty for all)"
          options={urgencyOptions}
          value={transformation.urgencyFilter ?? []}
          onChange={(value) =>
            setTransformation({
              ...transformation,
              urgencyFilter: value,
            })
          }
        />
      </Form.Item>
    </>
  );
};
