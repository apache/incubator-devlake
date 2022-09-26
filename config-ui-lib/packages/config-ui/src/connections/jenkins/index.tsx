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
import { useState } from 'react';
import { Form, Input, InputNumber, Button } from 'antd';

import { operate } from '../../utils/operate';

import * as API from './api';
import * as S from './styled';

export const Jenkins = () => {
  const [operating, setOperating] = useState(false);
  const [form] = Form.useForm();

  const handleSubmit = () => {
    form.validateFields().then(async (values) => {
      const [success] = await operate(async () => await API.create(values), { setOperating });

      if (success) {
        form.resetFields();
      }
    });
  };

  return (
    <S.Container>
      <Form labelCol={{ span: 6 }} wrapperCol={{ span: 10 }} form={form}>
        <Form.Item label="Connection Name" name="name" rules={[{ required: true, message: 'Please connection name' }]}>
          <Input placeholder="eg. Jenkins" />
        </Form.Item>
        <Form.Item label="Endpoint URL" name="endpoint" rules={[{ required: true, message: 'Please endpont url' }]}>
          <Input placeholder="eg. https://api.jenkins.io/" />
        </Form.Item>
        <Form.Item label="Username" name="username" rules={[{ required: true, message: 'Please your username' }]}>
          <Input placeholder="Enter Username" />
        </Form.Item>
        <Form.Item label="Password" name="password" rules={[{ required: true, message: 'Please your password' }]}>
          <Input placeholder="Enter Password" />
        </Form.Item>
        <Form.Item label="Proxy URL" name="proxy">
          <Input placeholder="eg. http://proxy.localhost:8080" />
        </Form.Item>
        <Form.Item label="Rate Limit" name="rateLimitPerHour">
          <InputNumber min={0} />
        </Form.Item>
        <Form.Item wrapperCol={{ offset: 6, span: 16 }}>
          <Button type="primary" loading={operating} onClick={handleSubmit}>
            Submit
          </Button>
        </Form.Item>
      </Form>
    </S.Container>
  );
};
