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
import { Form, Input, Button, message } from 'antd';
import { CheckCircleFilled, CopyOutlined } from '@ant-design/icons';
import { CopyToClipboard } from 'react-copy-to-clipboard';

import { WebhookItemType, WebhookPayloadType, WebhookFormEnmu } from './typed';
import * as S from './styled';

export interface WebhookConnectionFormProps {
  formType?: WebhookFormEnmu;
  initialValues: WebhookItemType;
  onGenerateURL: (values: WebhookPayloadType) => void;
  onDone?: () => void;
  onSubmit: (values: WebhookPayloadType) => void;
}

export const WebhookConnectionForm = ({
  formType,
  initialValues,
  onGenerateURL,
  onDone,
  onSubmit,
}: WebhookConnectionFormProps) => {
  const [form] = Form.useForm();

  useEffect(() => {
    console.log(form.getFieldsValue());
    form.setFieldsValue({ name: 'xxx' });
    initialValues ? form.setFieldsValue(initialValues) : form.resetFields();
  }, [form, initialValues]);

  const handleGenerateURL = () => {
    form.validateFields().then(async (values) => {
      await onGenerateURL(values);
      form.resetFields();
    });
  };

  const handleSubmit = () => {
    form.validateFields().then(async (values) => {
      await onSubmit(values);
      form.resetFields();
    });
  };

  switch (formType) {
    case WebhookFormEnmu.add:
      return (
        <S.FormContainer>
          <Form form={form} layout="vertical">
            <Form.Item
              label="Webhook Name"
              name="name"
              rules={[{ required: true, message: 'Please connection name' }]}
            >
              <Input placeholder="Your Webhook Name" />
            </Form.Item>
            <Form.Item>
              <Button type="primary" onClick={handleGenerateURL}>
                Generate POST URL
              </Button>
            </Form.Item>
          </Form>
        </S.FormContainer>
      );
    case WebhookFormEnmu.show:
      return (
        <S.FormContainer>
          <div className="tips">
            <CheckCircleFilled style={{ fontSize: 20, color: '#4DB764' }} />
            <span>POST URL Generated!</span>
          </div>
          <URLList record={initialValues} />
          <div className="btns">
            <Button type="primary" onClick={onDone}>
              Done
            </Button>
          </div>
        </S.FormContainer>
      );
    case WebhookFormEnmu.edit:
      return (
        <S.FormContainer>
          <Form form={form} layout="vertical">
            <Form.Item
              label="Webhook Name"
              name="name"
              rules={[{ required: true, message: 'Please connection name' }]}
            >
              <Input placeholder="Your Webhook Name" />
            </Form.Item>
            <URLList record={initialValues} />
            <Form.Item>
              <Button type="primary" onClick={handleSubmit}>
                Save
              </Button>
            </Form.Item>
          </Form>
        </S.FormContainer>
      );
    default:
      return null;
  }
};

interface URLListProps {
  record: WebhookItemType;
}

const URLList = ({ record }: URLListProps) => {
  return (
    <S.URLContainer>
      <h2>POST URL</h2>
      <p>
        Copy the following URLs to your issue tracking tool for Incidents and CI
        tool for Deployments by making a POST to DevLake.
      </p>
      <h3>Incident</h3>
      <p>Send incident opened and reopened events</p>
      <div className="block">
        <span>{record.postIssuesEndpoint}</span>
        <CopyToClipboard
          text={record.postIssuesEndpoint}
          onCopy={() => message.success('Copy successfully.')}
        >
          <CopyOutlined width={16} height={16} />
        </CopyToClipboard>
      </div>
      <p>Send incident resolved events</p>
      <div className="block">
        <span>{record.closeIssuesEndpoint}</span>
        <CopyToClipboard
          text={record.closeIssuesEndpoint}
          onCopy={() => message.success('Copy successfully.')}
        >
          <CopyOutlined width={16} height={16} />
        </CopyToClipboard>
      </div>
      <h3>Deployment</h3>
      <p>Send task started and finished events</p>
      <div className="block">
        <span>{record.postPipelineTaskEndpoint}</span>
        <CopyToClipboard
          text={record.postPipelineTaskEndpoint}
          onCopy={() => message.success('Copy successfully.')}
        >
          <CopyOutlined width={16} height={16} />
        </CopyToClipboard>
      </div>
      <p>Send deployment finished events</p>
      <div className="block">
        <span>{record.closePipelineEndpoint}</span>
        <CopyToClipboard
          text={record.closePipelineEndpoint}
          onCopy={() => message.success('Copy successfully.')}
        >
          <CopyOutlined width={16} height={16} />
        </CopyToClipboard>
      </div>
    </S.URLContainer>
  );
};
