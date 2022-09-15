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
import { Button, Table } from 'antd';

import * as S from './styled';

export const Webhook = () => {
  return (
    <S.Container>
      <h1>
        <span>Webhook</span>
      </h1>
      <h4>
        Use Webhooks to define Incidents and Deployments for your CI tools if they are not listed in Data Sources.
      </h4>
      <div className="content">
        <div className="operate">
          <Button type="primary">Add Webhook</Button>
        </div>
        <Table />
      </div>
    </S.Container>
  );
};
