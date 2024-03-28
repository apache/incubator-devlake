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

import { useNavigate } from 'react-router-dom';
import { Card, Flex, Progress, Button } from 'antd';

import API from '@/api';
import { useRefreshData } from '@/hooks';

interface Props {
  style?: React.CSSProperties;
}

export const OnboardCard = ({ style }: Props) => {
  const navigate = useNavigate();

  const { ready, data } = useRefreshData(() => API.store.get('onboard'), []);

  if (!ready || !data || data.done) {
    return null;
  }

  return (
    <Card style={style}>
      <Flex align="center" justify="space-between">
        <Flex align="center">
          <Progress type="circle" size={30} format={() => `${data.step}/4`} percent={(data.step / 4) * 100} />
          <div style={{ marginLeft: 16 }}>
            <h4>Onboard Session</h4>
            <h5 style={{ fontWeight: 400 }}>
              You are not far from connecting to your first tool. Continue to finish it.
            </h5>
          </div>
        </Flex>
        <Button type="primary" onClick={() => navigate('/onboard')}>
          Continue
        </Button>
      </Flex>
    </Card>
  );
};
