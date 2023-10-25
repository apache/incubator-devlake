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

import { useState, useEffect, useMemo } from 'react';
import { Button, Intent } from '@blueprintjs/core';

import API from '@/api';
import { Dialog, FormItem, CopyText, ExternalLink, Message } from '@/components';
import { operator } from '@/utils';

import * as S from '../styled';

interface Props {
  initialId: ID;
  onCancel: () => void;
}

export const ViewDialog = ({ initialId, onCancel }: Props) => {
  const [record, setRecord] = useState({
    apiKeyId: '',
    postIssuesEndpoint: '',
    closeIssuesEndpoint: '',
    postDeploymentsCurl: '',
  });
  const [isOpen, setIsOpen] = useState(false);
  const [operating, setOperating] = useState(false);
  const [apiKey, setApiKey] = useState('');

  const prefix = useMemo(() => `${window.location.origin}/api`, []);

  useEffect(() => {
    (async () => {
      const res = await API.plugin.webhook.get(initialId);
      setRecord({
        apiKeyId: res.apiKey.id,
        postIssuesEndpoint: ` curl ${prefix}${res.postIssuesEndpoint} -X 'POST' -H 'Authorization: Bearer {API_KEY}' -d '{
          "issue_key":"DLK-1234",
          "title":"a feature from DLK",
          "type":"INCIDENT",
          "original_status":"TODO",
          "status":"TODO",    
          "created_date":"2020-01-01T12:00:00+00:00",
          "updated_date":"2020-01-01T12:00:00+00:00"
       }'`,
        closeIssuesEndpoint: `curl ${prefix}${res.closeIssuesEndpoint} -X 'POST' -H 'Authorization: Bearer {API_KEY}'`,
        postDeploymentsCurl: `curl ${prefix}${res.postPipelineDeployTaskEndpoint} -X 'POST' -H 'Authorization: Bearer {API_KEY}' -d '{
           "commit_sha":"the sha of deployment commit",
           "repo_url":"the repo URL of the deployment commit",
           "start_time":"Optional, eg. 2020-01-01T12:00:00+00:00"
       }'`,
      });
    })();
  }, [initialId]);

  const handleGenerateNewKey = async () => {
    const [success, res] = await operator(() => API.apiKey.renew(record.apiKeyId), {
      setOperating,
    });

    if (success) {
      setIsOpen(false);
      setApiKey(res.apiKey);
      setRecord({
        ...record,
        postIssuesEndpoint: ` curl ${record.postIssuesEndpoint} -X 'POST' -H 'Authorization: Bearer ${res.apiKey}' -d '{
          "issue_key":"DLK-1234",
          "title":"a feature from DLK",
          "type":"INCIDENT",
          "status":"TODO",    
          "created_date":"2020-01-01T12:00:00+00:00",
          "updated_date":"2020-01-01T12:00:00+00:00"
       }'`,
        closeIssuesEndpoint: `curl ${record.closeIssuesEndpoint} -X 'POST' -H 'Authorization: Bearer ${res.apiKey}'`,
        postDeploymentsCurl: `curl ${record.postDeploymentsCurl} -X 'POST' -H 'Authorization: Bearer ${res.apiKey}' -d '{
           "commit_sha":"the sha of deployment commit",
           "repo_url":"the repo URL of the deployment commit",
           "start_time":"Optional, eg. 2020-01-01T12:00:00+00:00"
       }'`,
      });
    }
  };

  return (
    <Dialog style={{ width: 820 }} isOpen title="View Webhook" footer={null} onCancel={onCancel}>
      <S.Wrapper>
        <p>
          Copy the following CURL commands to your issue tracking or CI/CD tools to push `Incidents` and `Deployments`
          by making a POST to DevLake. Please replace the {'{'}API_KEY{'}'} in the following URLs.
        </p>
        <FormItem label="Incident">
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
        </FormItem>
        <FormItem label="Deployments">
          <h5>Post to register a deployment</h5>
          <CopyText content={record.postDeploymentsCurl} />
          <p>
            See the{' '}
            <ExternalLink link="https://devlake.apache.org/docs/Plugins/webhook#deployment">
              full payload schema
            </ExternalLink>
            .
          </p>
        </FormItem>
        <FormItem
          label="API Key"
          subLabel="If you have forgotten your API key, you can revoke the previous key and generate a new one as a replacement."
        >
          {!apiKey ? (
            <Button intent={Intent.PRIMARY} text="Revoke and generate a new key" onClick={() => setIsOpen(true)} />
          ) : (
            <>
              <S.ApiKey>
                <CopyText content={apiKey} />
                <span>No Expiration</span>
              </S.ApiKey>
              <S.Tips>
                <strong>Please copy your key now. You will not be able to see it again.</strong>
              </S.Tips>
            </>
          )}
        </FormItem>
      </S.Wrapper>
      <Dialog
        style={{ width: 820 }}
        isOpen={isOpen}
        title="Are you sure you want to revoke the previous API key and  generate a new one?"
        cancelText="Go Back"
        okText="Confirm"
        okLoading={operating}
        onCancel={() => setIsOpen(false)}
        onOk={handleGenerateNewKey}
      >
        <Message content="Once this action is done, the previous API key will become invalid and you will need to enter the new key in the application that uses this Webhook API." />
      </Dialog>
    </Dialog>
  );
};
