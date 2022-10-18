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
import { Form, Button } from 'antd';

import config from '../config';

import type { ConnectionType, IPaylod, ItemType } from './typed';
import * as S from './styled';

interface IConnectionFormProps {
  type: ConnectionType;
  initialValues?: ItemType;
  onTest: (values: IPaylod) => void;
  onSubmit: (values: IPaylod) => void;
}

export const ConnectionForm = ({
  type,
  initialValues,
  onTest,
  onSubmit,
}: IConnectionFormProps) => {
  const [form] = Form.useForm();

  useEffect(() => {
    initialValues ? form.setFieldsValue(initialValues) : form.resetFields();
  }, [form, initialValues]);

  const handleTest = () => {
    form.validateFields().then(async (values) => {
      onTest(values);
    });
  };

  const handleSubmit = () => {
    form.validateFields().then(async (values) => {
      onSubmit(values);
      form.resetFields();
    });
  };

  return (
    <Form form={form} layout="vertical">
      {config[type].form.fields.map((field) => (
        <Form.Item
          key={field.name}
          name={field.name}
          label={field.label}
          rules={field.rule}
          tooltip={field.tooltip}
        >
          {field.render(field)}
        </Form.Item>
      ))}
      <S.BtnContainer align="right">
        <Button onClick={handleTest}>Test Connection</Button>
        <Button type="primary" onClick={handleSubmit}>
          Save Connection
        </Button>
      </S.BtnContainer>
    </Form>
  );
};
