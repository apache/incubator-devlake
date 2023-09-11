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
import { InputGroup, Icon } from '@blueprintjs/core';

import { Dialog, FormItem, CopyText } from '@/components';
import { useConnections } from '@/hooks';
import { operator } from '@/utils';

import * as API from '../api';
import * as S from '../styled';

interface Props {
  isOpen: boolean;
  onCancel: () => void;
  onSubmitAfter?: (id: ID) => void;
}

export const CreateDialog = ({ isOpen, onCancel, onSubmitAfter }: Props) => {
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

  const { onRefresh } = useConnections();

  const prefix = useMemo(() => `${window.location.origin}/api`, []);

  const handleSubmit = async () => {
    const [success, res] = await operator(
      async () => {
        const { id, apiKey } = await API.createConnection({ name });
        const { postIssuesEndpoint, closeIssuesEndpoint, postPipelineDeployTaskEndpoint } = await API.getConnection(id);
        return {
          id,
          apiKey: apiKey.apiKey,
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
        postIssuesEndpoint: `${prefix}${res.postIssuesEndpoint}?api_key=${res.apiKey}`,
        closeIssuesEndpoint: `${prefix}${res.closeIssuesEndpoint}?api_key=${res.apiKey}`,
        postDeploymentsCurl: `curl ${prefix}${res.postPipelineDeployTaskEndpoint} -X 'POST'
            \\ -H 'Authorization: Bearer ${res.apiKey}'
            \\ -d '{
            \\"commit_sha\\":\\"the sha of deployment commit\\",
            \\"repo_url\\":\\"the repo URL of the deployment commit\\",
            \\"start_time\\":\\"eg. 2020-01-01T12:00:00+00:00\\"
          }'`,
        apiKey: res.apiKey,
      });
      onRefresh('webhook');
      onSubmitAfter?.(res.id);
    }
  };

  return (
    <Dialog
      isOpen={isOpen}
      title="Add a New Webhook"
      style={{ width: 820 }}
      footer={step === 2 ? null : undefined}
      okText={step === 1 ? 'Generate POST URL' : 'Done'}
      okDisabled={step === 1 && !name}
      okLoading={operating}
      onCancel={onCancel}
      onOk={handleSubmit}
    >
      {step === 1 && (
        <S.Wrapper>
          <FormItem
            label="Webhook Name"
            subLabel="Give your Webhook a unique name to help you identify it in the future."
            required
          >
            <InputGroup placeholder="Webhook Name" value={name} onChange={(e) => setName(e.target.value)} />
          </FormItem>
        </S.Wrapper>
      )}
      {step === 2 && (
        <S.Wrapper>
          <h2>
            <Icon icon="endorsed" size={30} />
            <span>POST URL Generated!</span>
          </h2>
          <p>
            Copy the following POST URLs to your issue tracking or CI tools to push `Incidents` and `Deployments` by
            making a POST to DevLake.
          </p>
          <p>
            An API key is automatically generated for the authentication of this webhook. This key does not expire. You
            can revoke it in the webhook page at any time.
          </p>
          <FormItem label="Incident">
            <h5>Post to register an incident</h5>
            <CopyText content={record.postIssuesEndpoint} />
            <h5>Post to close a registered incident</h5>
            <CopyText content={record.closeIssuesEndpoint} />
          </FormItem>
          <FormItem label="Deployments">
            <h5>Post to register a deployment</h5>
            <CopyText content={record.postDeploymentsCurl} />
          </FormItem>
        </S.Wrapper>
      )}
    </Dialog>
  );
};

export default CreateDialog;
