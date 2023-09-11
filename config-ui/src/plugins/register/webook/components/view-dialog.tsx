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
import { Button, Icon, Intent } from '@blueprintjs/core';

import { Dialog, FormItem, CopyText } from '@/components';
import { operator } from '@/utils';

import * as API from '../api';
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
  const [operating, setOperating] = useState(false);
  const [apiKey, setApiKey] = useState('');

  const prefix = useMemo(() => `${window.location.origin}/api`, []);

  useEffect(() => {
    (async () => {
      const res = await API.getConnection(initialId);
      setRecord({
        apiKeyId: res.apiKey.id,
        postIssuesEndpoint: `${prefix}${res.postIssuesEndpoint}`,
        closeIssuesEndpoint: `${prefix}${res.closeIssuesEndpoint}`,
        postDeploymentsCurl: `curl ${prefix}${res.postPipelineDeployTaskEndpoint} -X 'POST' -d "{
      \\"commit_sha\\":\\"the sha of deployment commit\\",
      \\"repo_url\\":\\"the repo URL of the deployment commit\\",
      \\"start_time\\":\\"Optional, eg. 2020-01-01T12:00:00+00:00\\"
      }"`,
      });
    })();
  }, [initialId]);

  const handleGenerateNewKey = async () => {
    const [success, res] = await operator(() => API.renewApiKey(record.apiKeyId), {
      setOperating,
    });

    if (success) {
      setApiKey(res.apiKey);
      setRecord({
        ...record,
        postIssuesEndpoint: `${prefix}${res.postIssuesEndpoint}?api_key=${res.apiKey}`,
        closeIssuesEndpoint: `${prefix}${res.closeIssuesEndpoint}?api_key=${res.apiKey}`,
        postDeploymentsCurl: `curl ${prefix}${res.postPipelineDeployTaskEndpoint} -X 'POST' -d "{
          \\  -H 'Authorization: Bearer ${res.apiKey}'
          \\ "commit_sha\\":\\"the sha of deployment commit\\",
          \\ "repo_url\\":\\"the repo URL of the deployment commit\\",
          \\ "start_time\\":\\"Optional, eg. 2020-01-01T12:00:00+00:00\\"
      }"`,
      });
    }
  };

  return (
    <Dialog style={{ width: 820 }} isOpen title="View Webhook" footer={null} onCancel={onCancel}>
      <S.Wrapper>
        {apiKey && (
          <h2>
            <Icon icon="endorsed" size={30} />
            <span>POST URL Generated!</span>
          </h2>
        )}
        <p>
          Copy the following POST URLs to your issue tracking or CI tools to push `Incidents` and `Deployments` by
          making a POST to DevLake. Please replace the API_key in the following URLs.
        </p>
        <FormItem label="Incidents">
          <h5>Post to register an incident</h5>
          <CopyText content={record.postIssuesEndpoint} />
          <h5>Post to close a registered incident</h5>
          <CopyText content={record.closeIssuesEndpoint} />
        </FormItem>
        <FormItem label="Deployment">
          <h5>Post to register a deployment</h5>
          <CopyText content={record.postDeploymentsCurl} />
        </FormItem>
        <FormItem
          label="API Key"
          subLabel="If you have forgotten your API key, you can revoke the previous key and generate a new one as a replacement."
        >
          {!apiKey ? (
            <Button
              intent={Intent.PRIMARY}
              loading={operating}
              text="Revoke and generate a new key"
              onClick={handleGenerateNewKey}
            />
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
    </Dialog>
  );
};
