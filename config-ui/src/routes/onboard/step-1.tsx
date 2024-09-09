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

import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { Input, Flex, Button, message } from 'antd';

import API from '@/api';
import { Block, Markdown } from '@/components';
import { PATHS } from '@/config';
import { selectOnboard, update, previous, changeProjectName, changePlugin } from '@/features/onboard';
import { useAppDispatch, useAppSelector } from '@/hooks';
import { ConnectionSelect } from '@/plugins';
import { operator } from '@/utils';

import * as S from './styled';

export const Step1 = () => {
  const [QA, setQA] = useState('');
  const [operating, setOperating] = useState(false);

  const dispatch = useAppDispatch();
  const { projectName, plugin } = useAppSelector(selectOnboard);

  useEffect(() => {
    fetch(`/onboard/step-1/${plugin ? plugin : 'default'}.md`)
      .then((res) => res.text())
      .then((text) => setQA(text));
  }, [plugin]);

  const handleSubmit = async () => {
    if (!projectName || !plugin) {
      return;
    }

    const [, res] = await operator(() => API.project.checkName(projectName), {
      setOperating,
      hideToast: true,
    });

    if (res.exist) {
      message.error(`Project name "${projectName}" already exists, please try another name.`);
      return;
    }

    dispatch(update({}));
  };

  return (
    <>
      <S.StepContent>
        <div className="content">
          <Block
            title="Project Name"
            description="Give your project a unique name with letters, numbers, -, _ or /"
            required
          >
            <Input
              style={{ width: 386 }}
              placeholder="Your Project Name"
              value={projectName}
              onChange={(e) => dispatch(changeProjectName(e.target.value))}
            />
          </Block>
          <Block
            title="Data Connection"
            description={
              <>
                For self-managed GitLab/GitHub/Bitbucket, please skip the onboarding and configure via{' '}
                <Link to={PATHS.CONNECTIONS()}>Data Connections</Link>.
              </>
            }
            required
          >
            <ConnectionSelect
              placeholder="Select a Data Connection"
              options={[
                {
                  plugin: 'github',
                  value: 'github',
                  label: 'GitHub',
                },
                {
                  plugin: 'gitlab',
                  value: 'gitlab',
                  label: 'GitLab',
                },
                {
                  plugin: 'bitbucket',
                  value: 'bitbucket',
                  label: 'Bitbucket',
                },
                {
                  plugin: 'azuredevops',
                  value: 'azuredevops',
                  label: 'Azure DevOps',
                },
              ]}
              value={plugin}
              onChange={(p) => dispatch(changePlugin(p))}
            />
          </Block>
        </div>
        <Markdown className="qa">{QA}</Markdown>
      </S.StepContent>
      <Flex style={{ marginTop: 64 }} justify="space-between">
        <Button ghost type="primary" loading={operating} onClick={() => dispatch(previous())}>
          Previous Step
        </Button>
        <Button type="primary" loading={operating} disabled={!projectName || !plugin} onClick={handleSubmit}>
          Next Step
        </Button>
      </Flex>
    </>
  );
};
