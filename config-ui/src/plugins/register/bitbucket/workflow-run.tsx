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

import { Flex, Input } from 'antd';
import styled from 'styled-components';

import { ExternalLink } from '@/components';

import Picture1 from './assets/workflow-run-1.jpeg';
import Picture2 from './assets/workflow-run-2.jpeg';

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
        Example - Convert Bitbucket pipelines that have executed the step ‘Deploy-to-prod’ to production deployments
      </h5>
      <ol>
        <li>Go to your Bitbucket Pipelines page (as shown in the first picture). </li>
        <li>
          Locate a successful pipeline that have executed the step ‘Deploy-to-prod’, you will see all steps executed in
          this pipeline (as shown in the second picture).
          <Flex>
            <img src={Picture1} width="100%" alt="" />
            <img src={Picture2} width="100%" alt="" />
          </Flex>
        </li>
        <li>
          <div>
            In the first input field, enter the following regex to identify this pipeline as a deployment (as
            highlighted in the <span className="yellow">yellow</span> rectangle).
          </div>
          <div style={{ marginTop: 10 }}>
            Its branch or one of its steps matches
            <Input style={{ width: 240 }} size="small" disabled value="(?i)Deploy-to.*" />
          </div>
        </li>
        <li>
          <div>
            In the second input field, enter the following regex to identify the production deployments (highlighted in
            the <span className="red">red</span> rectangle). If left empty, all pipelines containing the job in the{' '}
            <span className="yellow">yellow</span> rectangle will be regarded as Production Deployments.
          </div>
          <div style={{ marginTop: 10 }}>
            If the branch or the step also matches{' '}
            <Input style={{ width: 100 }} size="small" disabled value="(?i)prod" />, this deployment will be regarded as
            a ‘Production Deployment’
          </div>
        </li>
      </ol>
      <div>
        For more information, please refer to{' '}
        <ExternalLink link="https://devlake.apache.org/docs/Configuration/BitBucket/#cicd">
          this documentation
        </ExternalLink>
        .
      </div>
    </Wrapper>
  );
};
