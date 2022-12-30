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

import React from 'react';
import { Icon, ButtonGroup, Button, Colors, Intent } from '@blueprintjs/core';

import { Card } from '@/components';

import type { UseBPUpgradeProps } from './use-bp-upgrade';
import { useBPUpgrade } from './use-bp-upgrade';

interface Props extends Pick<UseBPUpgradeProps, 'onResetError'> {}

export const BPUpgrade = ({ ...props }: Props) => {
  const bpId = window.location.pathname.split('/').pop();
  const { processing, onSubmit } = useBPUpgrade({ id: bpId, ...props });

  return (
    <Card>
      <h2>
        <Icon icon="outdated" color={Colors.ORANGE5} size={20} />
        <span>Current Blueprint Need to Upgrade</span>
      </h2>
      <p>
        If you have already started, please wait for database migrations to complete, do <strong>NOT</strong> close your
        browser at this time.
      </p>
      <ButtonGroup>
        <Button loading={processing} text="Proceed to Upgrade" intent={Intent.PRIMARY} onClick={onSubmit} />
      </ButtonGroup>
    </Card>
  );
};
