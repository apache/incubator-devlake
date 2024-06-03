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

import { useState, useMemo } from 'react';
import { CheckCircleOutlined } from '@ant-design/icons';
import { Modal, Input } from 'antd';

import { useAppDispatch } from '@/hooks';
import { Block, CopyText, ExternalLink } from '@/components';
import { addWebhook } from '@/features';
import { operator } from '@/utils';

import { transformURI } from './utils';

import * as S from '../styled';

interface Props {
  open: boolean;
  onCancel: () => void;
  onSubmitAfter?: (id: ID) => void;
}

export const CreateDialog = ({ open, onCancel, onSubmitAfter }: Props) => {
  const [operating, setOperating] = useState(false);
  const [step, setStep] = useState(1);
  const [name, setName] = useState('');
  const [record, setRecord] = useState({
    id: 0,
    postIssuesEndpoint: '',
    closeIssuesEndpoint: '',
    postDeploymentsCurl: '',
    apiKey: '',
  });

  const dispatch = useAppDispatch();

  const prefix = useMemo(() => `${window.location.origin}/api`, []);

  const handleSubmit = async () => {
    const [success, res] = await operator(
      async () => {
        const {
          webhook: { id, postIssuesEndpoint, closeIssuesEndpoint, postPipelineDeployTaskEndpoint },
          apiKey,
        } = await dispatch(addWebhook({ name })).unwrap();

        return {
          id,
          apiKey,
          postIssuesEndpoint,
          closeIssuesEndpoint,
          postPipelineDeployTaskEndpoint,
        };
      },
      {
        setOperating,
        hideToast: true,
      },
    );

    if (success) {
      setStep(2);
      setRecord({
        id: res.id,
        apiKey: res.apiKey,
        ...transformURI(prefix, res, res.apiKey),
      });
      onSubmitAfter?.(res.id);
    }
  };

  return (
    <Modal
      open={open}
      width={820}
      centered
      title="Add a New Webhook"
      footer={step === 2 ? null : undefined}
      okText={step === 1 ? 'Generate POST URL' : 'Done'}
      okButtonProps={{
        disabled: step === 1 && !name,
        loading: operating,
      }}
      onCancel={onCancel}
      onOk={handleSubmit}
    >
      {step === 1 && (
        <S.Wrapper>
          <Block
            title="Webhook Name"
            description="Give your Webhook a unique name to help you identify it in the future."
            required
          >
            <Input placeholder="Webhook Name" value={name} onChange={(e) => setName(e.target.value)} />
          </Block>
        </S.Wrapper>
      )}
      {step === 2 && (
        <S.Wrapper>
          <h2>
            <CheckCircleOutlined size={30} />
            <span>CURL commands generated. Please copy them now.</span>
          </h2>
          <p>
            A non-expired API key is automatically generated for the authentication of the webhook. This key will only
            show now. You can revoke it in the webhook page at any time.
          </p>
          <Block title="Incident">
            <h5>Post to register/update an incident</h5>
            <CopyText content={record.postIssuesEndpoint} />
            <p>
              See the{' '}
              <ExternalLink link="https://devlake.apache.org/docs/Plugins/webhook#register-issues---update-or-create-issues">
                full payload schema
              </ExternalLink>
              .
            </p>
            <h5>Post to close a registered incident</h5>
            <CopyText content={record.closeIssuesEndpoint} />
            <p>
              See the{' '}
              <ExternalLink link="https://devlake.apache.org/docs/Plugins/webhook#register-issues---close-issues-optional">
                full payload schema
              </ExternalLink>
              .
            </p>
          </Block>
          <Block title="Deployments">
            <h5>Post to register a deployment</h5>
            <CopyText content={record.postDeploymentsCurl} />
            <p>
              See the{' '}
              <ExternalLink link="https://devlake.apache.org/docs/Plugins/webhook#deployment">
                full payload schema
              </ExternalLink>
              .
            </p>
          </Block>
        </S.Wrapper>
      )}
    </Modal>
  );
};

export default CreateDialog;
