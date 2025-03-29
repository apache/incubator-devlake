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

import { Input } from 'antd';
import styled from 'styled-components';

import { ExternalLink } from '@/components';

import Picture from './assets/workflow-run.jpeg';

const Wrapper = styled.div`
  padding: 8px 16px;
  width: 100%;
  font-size: 12px;
  background-color: #fff;
  box-sizing: border-box;
  overflow: hidden;

  li {
    margin-bottom: 20px;

    &:last-child {
      margin-bottom: 0;
    }
  }

  span.blue {
    color: #7497f7;
  }

  span.yellow {
    color: #f4be55;
  }

  span.red {
    color: #ff8b8b;
  }
`;

export const WorkflowRun = () => {
  return (
    <Wrapper>
      <h5>
        Example - Convert GitLab pipeline runs that have executed the ‘build-image’ job on the ‘master’ branch to
        production deployments
      </h5>
      <ol>
        <li>Go to your GitLab/Build/Jobs page, where you will see the job run history.</li>
        <li>
          Search for the job ‘build-image’, and you will see all pipelines that have executed this job (highlighted in
          the <span className="yellow">yellow</span> rectangle).
          <img src={Picture} width="100%" alt="" />
        </li>
        <li>
          <div>
            In the first input field, enter the following regex to identify deployments (as highlighted in the{' '}
            <span className="yellow">yellow</span> rectangle):
          </div>
          <div style={{ marginTop: 10 }}>
            Its branch or <strong>one of its jobs</strong> matches
            <Input style={{ width: 240 }} size="small" disabled value="(?i)build-image" />
          </div>
        </li>
        <li>
          <div>
            In the second input field, enter the following regex to identify the production deployments (highlighted in
            the <span className="red">red</span> rectangle). If left empty, all deployments will be regarded as
            Production Deployments.
          </div>
          <div style={{ marginTop: 10 }}>
            If the branch or the job also matches <Input style={{ width: 100 }} size="small" disabled value="master" />,
            this deployment will be regarded as a ‘Production Deployment’
          </div>
        </li>
      </ol>
      <div>
        For more information, please refer to{' '}
        <ExternalLink link="https://devlake.apache.org/docs/Configuration/GitLab/#step-13---add-scope-config-optional">
          this documentation
        </ExternalLink>
        .
      </div>
    </Wrapper>
  );
};
