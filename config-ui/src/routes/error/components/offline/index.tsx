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

import { useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { RedoOutlined, QuestionCircleOutlined } from '@ant-design/icons';
import { Card, Flex, Button } from 'antd';
import { Icon, Tag, Colors, IconName } from '@blueprintjs/core';

import API from '@/api';
import { DEVLAKE_ENDPOINT } from '@/config';
import { useAutoRefresh } from '@/hooks';

export const Offline = () => {
  const [version, setVersion] = useState(1);

  const navigate = useNavigate();

  const { loading, data } = useAutoRefresh<{ online: boolean }>(
    async () => {
      try {
        await API.ping();
        return { online: true };
      } catch {
        return { online: false };
      }
    },
    [version],
    {
      cancel: (data) => {
        return data?.online ?? false;
      },
      retryLimit: 2,
    },
  );

  const { online } = data || { online: false };

  const [icon, color, text] = useMemo(
    () => [online ? 'endorsed' : 'offline', online ? Colors.GREEN3 : Colors.RED3, data ? 'Online' : 'Offline'],
    [online],
  );

  const handleContinue = () => {
    navigate('/');
  };

  return (
    <Card>
      <h2>
        <Icon icon={icon as IconName} color={color} size={30} />
        <span>DevLake API</span>
        <strong style={{ marginLeft: 4, color }}>{text}</strong>
      </h2>
      <p>
        <Tag>DEVLAKE_ENDPOINT: {DEVLAKE_ENDPOINT}</Tag>
      </p>
      {!online ? (
        <>
          <p>
            Please wait for the&nbsp;
            <strong>Lake API</strong> to start before accessing the <strong>Configuration Interface</strong>.
          </p>
          <Flex justify="center">
            <Button
              type="primary"
              loading={loading}
              icon={<RedoOutlined rev={undefined} />}
              onClick={() => setVersion((v) => v + 1)}
            >
              Refresh
            </Button>
          </Flex>
        </>
      ) : (
        <>
          <p>Connectivity to the Lake API service was successful.</p>
          <Flex justify="center">
            <Button type="primary" onClick={handleContinue}>
              Continue
            </Button>
            <Button
              icon={<QuestionCircleOutlined rev={undefined} />}
              onClick={() =>
                window.open(
                  'https://github.com/apache/incubator-devlake/blob/main/README.md',
                  '_blank',
                  'noopener,noreferrer',
                )
              }
            >
              Read Documentation
            </Button>
          </Flex>
        </>
      )}
    </Card>
  );
};
