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
import { Link } from 'react-router-dom';

import * as S from './styled';

import GitLabIcon from '@/images/icons/gitlab.svg';
import JenkinsIcon from '@/images/icons/jenkins.svg';
import TapdIcon from '@/images/icons/tapd.svg';
import JiraIcon from '@/images/icons/jira.svg';
import GitHubIcon from '@/images/icons/github.svg';
import AzureIcon from '@/images/icons/azure.svg';
import BitBucketIcon from '@/images/icons/bitbucket.svg';
import GiteeIcon from '@/images/icons/gitee.svg';
import WebhookIcon from '@/images/icons/webhook.svg';

const connections = [
  {
    name: 'Gitlab',
    icon: GitLabIcon,
    link: '/gitlab',
  },
  {
    name: 'Jenkins',
    icon: JenkinsIcon,
    link: '/jenkins',
  },
  {
    name: 'TAPD',
    icon: TapdIcon,
    link: '/tapd',
  },
  {
    name: 'JIRA',
    icon: JiraIcon,
    link: '/jira',
  },
  {
    name: 'GitHub',
    icon: GitHubIcon,
    link: '/github',
  },
  {
    name: 'Azure',
    icon: AzureIcon,
    link: '/azure',
  },
  {
    name: 'BitBucket',
    icon: BitBucketIcon,
    link: '/bitbucket',
  },
  {
    name: 'Gitee',
    icon: GiteeIcon,
    link: '/gitee',
  },
];

export const Connections = () => {
  return (
    <S.Container>
      <h1>Connections</h1>
      <h4>
        Create and manage data connections from the following data sources or
        Webhooks to be used in syncing data in your Blueprints.
      </h4>
      <div className="item">
        <h2>Data Sources</h2>
        <h4>
          Data connections created for the following data sources can be used in
          your Blueprints.
        </h4>
        <ul className="list">
          {connections.map((c) => (
            <li>
              <Link to={`/connections${c.link}`}>
                <img src={c.icon} alt="" />
                <span>{c.name}</span>
              </Link>
            </li>
          ))}
        </ul>
      </div>
      <div className="item">
        <h2>Webhooks</h2>
        <h4>
          You can use Webhooks to define Issues and Deployments to be used in
          calculating DORA metrics. Please note: Webhooks cannot be created or
          managed in Blueprints.
        </h4>
        <ul className="list">
          <li>
            <Link to="/connections/webhook">
              <img src={WebhookIcon} alt="" />
              <span>Issue/Deployment Webhook</span>
            </Link>
          </li>
        </ul>
      </div>
    </S.Container>
  );
};
