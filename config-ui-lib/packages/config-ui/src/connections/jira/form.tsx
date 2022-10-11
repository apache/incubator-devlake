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

import type { JiraItemType, JiraPayloadType } from './typed';

export interface JiraConnectionFormProps {
  initialValues?: JiraItemType;
  onSubmit: (values: JiraPayloadType) => void;
}

export const JiraConnectionForm = ({
  initialValues,
  onSubmit,
}: JiraConnectionFormProps) => {
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
        <Input placeholder="eg. Jira" />
      </Form.Item>
      <Form.Item
        label="Endpoint URL"
        name="endpoint"
        rules={[{ required: true, message: 'Please endpont url' }]}
      >
        <Input placeholder="eg. https://your-domain.atlassian.net/rest/" />
      </Form.Item>
      <Form.Item
        label="Username / E-mail"
        name="username"
        rules={[{ required: true, message: 'Please your username' }]}
      >
        <Input placeholder="Enter Username" />
      </Form.Item>
      <Form.Item
        label="Password"
        name="password"
        rules={[{ required: true, message: 'Please your password' }]}
      >
        <Input placeholder="Enter Password" />
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
