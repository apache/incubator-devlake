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

import { Dialog, FormItem, Selector, Buttons } from '@/components';
import { useConnections } from '@/hooks';
import { DataScopeSelect, getPluginScopeId } from '@/plugins';
import type { ConnectionItemType } from '@/store';

interface Props {
  disabled: string[];
  onCancel: () => void;
  onSubmit: (value: any) => void;
}

export const AddConnectionDialog = ({ disabled = [], onCancel, onSubmit }: Props) => {
  const [step, setStep] = useState(1);
  const [selectedConnection, setSelectedConnection] = useState<ConnectionItemType>();

  const { connections } = useConnections({ filterPlugin: ['webhook'] });

  const disabledItems = useMemo(
    () => connections.filter((cs) => (disabled.length ? disabled.includes(cs.unique) : false)),
    [disabled],
  );

  const handleSubmit = (scopeIds: any) => {
    if (!selectedConnection) return;
    onSubmit({
      pluginName: selectedConnection.plugin,
      connectionId: selectedConnection.id,
      scopes: scopeIds.map((scopeId: any) => ({ scopeId })),
    });
  };

  return (
    <Dialog style={{ width: 820 }} isOpen title={`Add a Connection - Step ${step}`} footer={null} onCancel={onCancel}>
      {step === 1 && (
        <FormItem
          label="Data Connections"
          subLabel="Select from existing Data Connections. If you have not created any Data Connections yet, please create and manage Connections first."
          required
        >
          <Selector
            items={connections}
            disabledItems={disabledItems}
            getKey={(it) => it.unique}
            getName={(it) => it.name}
            getIcon={(it) => it.icon}
            selectedItem={selectedConnection}
            onChangeItem={(selectedItem) => setSelectedConnection(selectedItem)}
          />
          <Buttons position="bottom" align="right">
            <Button outlined intent={Intent.PRIMARY} text="Cancel" onClick={onCancel} />
            <Button disabled={!selectedConnection} intent={Intent.PRIMARY} text="Next" onClick={() => setStep(2)} />
          </Buttons>
        </FormItem>
      )}
      {step === 2 && selectedConnection && (
        <DataScopeSelect
          plugin={selectedConnection.plugin}
          connectionId={selectedConnection.id}
          onCancel={onCancel}
          onSubmit={handleSubmit}
        />
      )}
    </Dialog>
  );
};
