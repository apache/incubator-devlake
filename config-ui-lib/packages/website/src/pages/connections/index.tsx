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

import JenkinsIcon from '@/images/icons/jenkins.svg';
import WebhookIcon from '@/images/icons/webhook.svg';

export const Connections = () => {
  return (
    <S.Container>
      <h1>Connections</h1>
      <h4>
        Create and manage data connections from the following data sources or Webhooks to be used in syncing data in
        your Blueprints.
      </h4>
      <div className="item">
        <h2>Data Sources</h2>
        <h4>Data connections created for the following data sources can be used in your Blueprints.</h4>
        <ul className="list">
          <li>
            <Link to="/connection/jenkins">
              <img src={JenkinsIcon} width={60} alt="" />
              <span>Jenkins</span>
            </Link>
          </li>
        </ul>
      </div>
      <div className="item">
        <h2>Webhooks</h2>
        <h4>
          You can use Webhooks to define Issues and Deployments to be used in calculating DORA metrics. Please note:
          Webhooks cannot be created or managed in Blueprints.
        </h4>
        <ul className="list">
          <li>
            <Link to="/connection/webhook">
              <img src={WebhookIcon} width={60} alt="" />
              <span>Issue/Deployment Webhook</span>
            </Link>
          </li>
        </ul>
      </div>
    </S.Container>
  );
};
