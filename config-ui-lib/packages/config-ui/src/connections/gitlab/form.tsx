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
import { Form, Input, InputNumber, Button } from 'antd';

import type { GitLabItemType, GitLabPayloadType } from './typed';

export interface GitLabConnectionFormProps {
  initialValues?: GitLabItemType;
  onSubmit: (values: GitLabPayloadType) => void;
}

export const GitLabConnectionForm = ({
  initialValues,
  onSubmit,
}: GitLabConnectionFormProps) => {
  const [form] = Form.useForm();

  useEffect(() => {
    initialValues ? form.setFieldsValue(initialValues) : form.resetFields();
  }, [form, initialValues]);

  const handleSubmit = () => {
    form.validateFields().then(async (values) => {
      await onSubmit(values);
      form.resetFields();
    });
  };

  return (
    <Form
      labelCol={{ span: 6 }}
      wrapperCol={{ span: 16 }}
      form={form}
      initialValues={initialValues}
    >
      <Form.Item
        label="Name"
        name="name"
        rules={[{ required: true, message: 'Please connection name' }]}
      >
        <Input placeholder="eg. GitLab" />
      </Form.Item>
      <Form.Item
        label="Endpoint URL"
        name="endpoint"
        rules={[{ required: true, message: 'Please endpont url' }]}
      >
        <Input placeholder="eg. https://gitlab.com/api/v4/" />
      </Form.Item>
      <Form.Item
        label="Access Token"
        name="token"
        rules={[{ required: true, message: 'Please your access token' }]}
      >
        <Input placeholder="eg. ff9d1ad0e5c04f1f98fa" />
      </Form.Item>
      <Form.Item label="Proxy URL" name="proxy">
        <Input placeholder="eg. http://proxy.localhost:8080" />
      </Form.Item>
      <Form.Item label="Rate Limit" name="rateLimitPerHour">
        <InputNumber min={0} />
      </Form.Item>
      <Form.Item
        wrapperCol={{ offset: 6, span: 16 }}
        style={{ marginBottom: 0 }}
      >
        <Button type="primary" onClick={handleSubmit}>
          Submit
        </Button>
      </Form.Item>
    </Form>
  );
};
