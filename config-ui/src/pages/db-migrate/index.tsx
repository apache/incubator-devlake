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

import { useState } from 'react';
import { Icon, ButtonGroup, Button, Colors, Intent } from '@blueprintjs/core';
import { useHistory } from 'react-router-dom';

import { Card } from '@/components';
import { operator } from '@/utils';

import * as API from './api';

export const DBMigratePage = () => {
  const [operating, setOperating] = useState(false);

  const history = useHistory();

  const handleSubmit = async () => {
    const [success] = await operator(() => API.migrate(), {
      setOperating: setOperating,
    });

    if (success) {
      history.push('/');
    }
  };

  return (
    <Card>
      <h2>
        <Icon icon="outdated" color={Colors.ORANGE5} size={20} />
        <span>New Migration Scripts Detected</span>
      </h2>
      <p>
        If you have already started, please wait for database migrations to complete, do <strong>NOT</strong> close your
        browser at this time.
      </p>
      <p className="warning">
        Warning: Performing migration may wipe collected data for consistency and re-collecting data may be required.
      </p>
      <ButtonGroup>
        <Button
          loading={operating}
          text="Proceed to Database Migration"
          intent={Intent.PRIMARY}
          onClick={handleSubmit}
        />
      </ButtonGroup>
    </Card>
  );
};
