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

import { useState, useMemo } from 'react';
import { Button, Intent } from '@blueprintjs/core';

import { Buttons } from '@/components';
import { getPluginConfig, getPluginScopeId } from '@/plugins';
import { operator } from '@/utils';

import { SearchLocal } from './search-local';
import { SearchRemote } from './search-remote';
import * as API from './api';

interface Props {
  plugin: string;
  connectionId: ID;
  disabledScope?: any[];
  onCancel: () => void;
  onSubmit: (origin: any) => void;
}

export const DataScopeRemote = ({ plugin, connectionId, disabledScope, onCancel, onSubmit }: Props) => {
  const [selectedScope, setSelectedScope] = useState<any[]>([]);
  const [operating, setOperating] = useState(false);

  const config = useMemo(() => getPluginConfig(plugin).dataScope, [plugin]);

  const handleSubmit = async () => {
    const [success, res] = await operator(
      () => API.updateDataScope(plugin, connectionId, { data: selectedScope.map((it) => it.data) }),
      {
        setOperating,
        formatMessage: () => 'Add data scope successful.',
      },
    );

    if (success) {
      onSubmit(res);
    }
  };

  return (
    <>
      {config.render ? (
        config.render({
          connectionId,
          disabledItems: disabledScope?.map((it) => ({ id: it.id })),
          selectedItems: selectedScope,
          onChangeSelectedItems: setSelectedScope,
        })
      ) : config.localSearch ? (
        <SearchLocal
          plugin={plugin}
          connectionId={connectionId}
          config={config}
          disabledScope={disabledScope ?? []}
          selectedScope={selectedScope}
          onChange={setSelectedScope}
        />
      ) : (
        <SearchRemote
          plugin={plugin}
          connectionId={connectionId}
          config={config}
          disabledScope={disabledScope ?? []}
          selectedScope={selectedScope}
          onChange={setSelectedScope}
        />
      )}
      <Buttons position="bottom" align="right">
        <Button outlined intent={Intent.PRIMARY} text="Cancel" disabled={operating} onClick={onCancel} />
        <Button
          outlined
          intent={Intent.PRIMARY}
          text="Save"
          loading={operating}
          disabled={!selectedScope.length}
          onClick={handleSubmit}
        />
      </Buttons>
    </>
  );
};
