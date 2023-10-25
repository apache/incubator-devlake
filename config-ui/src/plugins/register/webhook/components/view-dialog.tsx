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
import { Button, Intent } from '@blueprintjs/core';

import { useAppDispatch, useAppSelector } from '@/app/hook';
import { Dialog, FormItem, CopyText, ExternalLink, Message } from '@/components';
import { selectWebhook, renewWebhookApiKey } from '@/features';
import { IWebhook } from '@/types';
import { operator } from '@/utils';

import * as S from '../styled';

interface Props {
  initialId: ID;
  onCancel: () => void;
}

const transformURI = (prefix: string, webhook: IWebhook) => {
  return {
    postIssuesEndpoint: `curl ${prefix}${webhook.postIssuesEndpoint} -X 'POST' -H 'Authorization: Bearer ${
      webhook.apiKey ?? '{API_KEY}'
    }' -d '{
        "issue_key":"DLK-1234",
        "title":"a feature from DLK",
        "type":"INCIDENT",
        "original_status":"TODO",
        "status":"TODO",    
        "created_date":"2020-01-01T12:00:00+00:00",
        "updated_date":"2020-01-01T12:00:00+00:00"
     }'`,
    closeIssuesEndpoint: `curl ${prefix}${webhook.closeIssuesEndpoint} -X 'POST' -H 'Authorization: Bearer ${
      webhook.apiKey ?? '{API_KEY}'
    }'`,
    postDeploymentsCurl: `curl ${prefix}${webhook.postPipelineDeployTaskEndpoint} -X 'POST' -H 'Authorization: Bearer ${
      webhook.apiKey ?? '{API_KEY}'
    }' -d '{
         "commit_sha":"the sha of deployment commit",
         "repo_url":"the repo URL of the deployment commit",
         "start_time":"Optional, eg. 2020-01-01T12:00:00+00:00"
     }'`,
  };
};

export const ViewDialog = ({ initialId, onCancel }: Props) => {
  const [isOpen, setIsOpen] = useState(false);
  const [operating, setOperating] = useState(false);

  const dispatch = useAppDispatch();
  const webhook = useAppSelector((state) => selectWebhook(state, initialId)) as IWebhook;
  const prefix = useMemo(() => `${window.location.origin}/api`, []);

  const URI = transformURI(prefix, webhook);

  const handleGenerateNewKey = async () => {
    const [success] = await operator(() => dispatch(renewWebhookApiKey(initialId)), {
      setOperating,
    });

    if (success) {
      setIsOpen(false);
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
          <CopyText content={URI.postIssuesEndpoint} />
          <p>
            See the{' '}
            <ExternalLink link="https://devlake.apache.org/docs/Plugins/webhook#register-issues---update-or-create-issues">
              full payload schema
            </ExternalLink>
            .
          </p>
          <h5>Post to close a registered incident</h5>
          <CopyText content={URI.closeIssuesEndpoint} />
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
          <CopyText content={URI.postDeploymentsCurl} />
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
          {!webhook.apiKey ? (
            <Button intent={Intent.PRIMARY} text="Revoke and generate a new key" onClick={() => setIsOpen(true)} />
          ) : (
            <>
              <S.ApiKey>
                <CopyText content={webhook.apiKey} />
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
