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
import { Form, Input, Button } from 'antd';

import { RateLimit } from '../../components';

import type { GiteeItemType, GiteePayloadType } from './typed';
import * as S from './styled';

export interface GiteeConnectionFormProps {
  initialValues?: GiteeItemType;
  operating?: boolean;
  onTest: (values: GiteePayloadType) => void;
  onSubmit: (values: GiteePayloadType) => void;
}

export const GiteeConnectionForm = ({
  initialValues,
  operating,
  onTest,
  onSubmit,
}: GiteeConnectionFormProps) => {
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
      await onSubmit(values);
      form.resetFields();
    });
  };

  return (
    <Form form={form} initialValues={initialValues} layout="vertical">
      <Form.Item
        label="Connection Name"
        name="name"
        rules={[{ required: true, message: 'Please connection name' }]}
        tooltip="Give your connection a unique name to help you identify it in the future."
      >
        <Input placeholder="eg. Gitee" />
      </Form.Item>
      <Form.Item
        label="Endpoint URL"
        name="endpoint"
        rules={[{ required: true, message: 'Please endpont url' }]}
        tooltip="Provide the Gitee instance API endpoint."
      >
        <Input placeholder="eg. https://gitee.com/api/v5/" />
      </Form.Item>
      <Form.Item
        label="Auth Token"
        name="token"
        rules={[{ required: true, message: 'Please your access token' }]}
        tooltip={
          <a
            href="https://gitee.com/api/v5/oauth_doc#/"
            target="_blank"
            rel="noreferrer"
          >
            Learn about how to create a personal access token
          </a>
        }
      >
        <Input placeholder="Your Auth Token" />
      </Form.Item>
      <Form.Item label="Proxy URL" name="proxy">
        <Input placeholder="eg. http://proxy.localhost:8080" />
      </Form.Item>
      <Form.Item
        style={{ marginBottom: 0 }}
        label="Rate Limit (per Hour)"
        name="rateLimitPerHour"
      >
        <RateLimit />
      </Form.Item>
      <S.BtnContainer>
        <Button loading={operating} onClick={handleTest}>
          Test Connection
        </Button>
        <Button type="primary" loading={operating} onClick={handleSubmit}>
          Save Connection
        </Button>
      </S.BtnContainer>
    </Form>
  );
};
