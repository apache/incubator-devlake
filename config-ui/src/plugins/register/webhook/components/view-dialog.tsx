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
import { Modal, Button } from 'antd';

import { useAppDispatch, useAppSelector } from '@/hooks';
import { Block, CopyText, ExternalLink, Message } from '@/components';
import { selectWebhook, renewWebhookApiKey } from '@/features';
import { IWebhook } from '@/types';
import { operator } from '@/utils';

import { transformURI } from './utils';

import * as S from '../styled';

interface Props {
  initialId: ID;
  onCancel: () => void;
}

export const ViewDialog = ({ initialId, onCancel }: Props) => {
  const [open, setOpen] = useState(false);
  const [operating, setOperating] = useState(false);
  const [apiKey, setApiKey] = useState('');

  const dispatch = useAppDispatch();
  const webhook = useAppSelector((state) => selectWebhook(state, initialId)) as IWebhook;
  const prefix = useMemo(() => `${window.location.origin}/api`, []);

  const URI = transformURI(prefix, webhook, apiKey);

  const handleGenerateNewKey = async () => {
    const [success, res] = await operator(async () => await dispatch(renewWebhookApiKey(initialId)).unwrap(), {
      setOperating,
    });

    if (success) {
      setApiKey(res.apiKey);
      setOpen(false);
    }
  };

  return (
    <Modal open width={820} centered title="View Webhook" footer={null} onCancel={onCancel}>
      <S.Wrapper>
        <p>
          Copy the following CURL commands to your issue tracking or CI/CD tools to push `Incidents` and `Deployments`
          by making a POST to DevLake. Please replace the {'{'}API_KEY{'}'} in the following URLs.
        </p>
        <Block title="Incident">
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
        </Block>
        <Block title="Deployments">
          <h5>Post to register a deployment</h5>
          <CopyText content={URI.postDeploymentsCurl} />
          <p>
            See the{' '}
            <ExternalLink link="https://devlake.apache.org/docs/Plugins/webhook#deployment">
              full payload schema
            </ExternalLink>
            .
          </p>
        </Block>
        <Block
          title="API Key"
          description="If you have forgotten your API key, you can revoke the previous key and generate a new one as a replacement."
        >
          {!apiKey ? (
            <Button type="primary" onClick={() => setOpen(true)}>
              Revoke and generate a new key
            </Button>
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
        </Block>
      </S.Wrapper>
      <Modal
        open={open}
        width={820}
        centered
        title="Are you sure you want to revoke the previous API key and  generate a new one?"
        okText="Confirm"
        cancelText="Go Back"
        okButtonProps={{
          loading: operating,
        }}
        onCancel={() => setOpen(false)}
        onOk={handleGenerateNewKey}
      >
        <Message content="Once this action is done, the previous API key will become invalid and you will need to enter the new key in the application that uses this Webhook API." />
      </Modal>
    </Modal>
  );
};
