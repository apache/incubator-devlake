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
        Example - Convert workflow runs of ‘Build-Images-Push-Docker’ on the ‘release’ branch to production deployments
      </h5>
      <ol>
        <li>
          Go to your GitHub Actions page. Select the ‘Build-Images-Push-Docker’ workflow on the left menu (as shown in
          the <span className="blue">blue</span> rectangle.).
        </li>
        <li>
          Locate the workflow run names in the list (highlighted by the <span className="yellow">yellow</span> text).
          They are usually the same as the workflow name. The branch name are shown in the second column in the label
          format.
          <img src={Picture} width="100%" alt="" />
        </li>
        <li>
          <div>
            In the first input field, enter the following regex to identify deployments (highlighted by the{' '}
            <span className="yellow">yellow</span> rectangle):
          </div>
          <div style={{ marginTop: 10 }}>
            The name of the <strong>GitHub workflow run</strong> or <strong>one of its jobs</strong> matches{' '}
            <Input style={{ width: 240 }} size="small" disabled value="(?i)Build-Images-Push-Docker" />
          </div>
        </li>
        <li>
          <div>
            In the second input field, enter the following regex to identify production deployments (highlighted by the{' '}
            <span className="red">red</span> rectangle). If left empty, all deployments in the{' '}
            <span className="yellow">yellow</span> rectangle will be regarded as Production Deployments.
          </div>
          <div style={{ marginTop: 10 }}>
            If the name or branch matches <Input style={{ width: 100 }} size="small" disabled value="release-.*" /> ,
            this deployment will be regarded as ‘Production Deployment’
          </div>
        </li>
      </ol>
      <div>
        For more information, please refer to{' '}
        <ExternalLink link="https://devlake.apache.org/docs/Configuration/GitHub/#cicd">
          this documentation
        </ExternalLink>
        .
      </div>
    </Wrapper>
  );
};
