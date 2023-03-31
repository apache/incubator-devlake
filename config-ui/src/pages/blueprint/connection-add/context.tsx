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

import React, { useState, useContext } from 'react';
import { useHistory } from 'react-router-dom';

import { PageLoading } from '@/components';
import { useRefreshData } from '@/hooks';
import { operator } from '@/utils';

import * as API from './api';

type ContextType = {
  name: string;
  step: number;
  filter: string[];
  connection?: MixConnection;
  onChangeConnection: (connection: MixConnection) => void;

  onPrev: () => void;
  onNext: () => void;
  saving: boolean;
  onSave: () => void;
  onCancel: () => void;
};

export const Context = React.createContext<ContextType>({
  name: '',
  step: 1,
  filter: [],

  onChangeConnection: () => {},
  onPrev: () => {},
  onNext: () => {},
  saving: false,
  onSave: () => {},
  onCancel: () => {},
});

interface Props {
  pname?: string;
  id: string;
  children: React.ReactNode;
}

export const ContextProvider = ({ pname, id, children }: Props) => {
  const [step, setStep] = useState(1);
  const [connection, setConnection] = useState<MixConnection>();
  const [saving, setSaving] = useState(false);

  const history = useHistory();

  const { ready, data } = useRefreshData(() => API.getBlueprint(id), [id]);

  const handlePrev = () => {
    setStep(step - 1);
  };

  const handleNext = () => {
    setStep(step + 1);
  };

  const handleSave = async () => {
    if (!connection) return null;
    const { plugin, connectionId, scope } = connection;
    const payload = {
      ...data,
      settings: {
        ...data.settings,
        connections: [
          ...data.settings.connections,
          {
            plugin,
            connectionId,
            scopes: scope.map((sc) => ({
              id: `${sc.id}`,
              entities: sc.entities,
            })),
          },
        ],
      },
    };

    const [success] = await operator(() => API.updateBlueprint(data.id, payload), {
      setOperating: setSaving,
    });

    if (success) {
      history.push(pname ? `/projects/${pname}` : `/blueprints/${id}`);
      return;
    }
  };

  const handleCancel = () => {
    history.push(`/blueprints/${id}`);
  };

  if (!ready || !data) {
    return <PageLoading />;
  }

  return (
    <Context.Provider
      value={{
        name: data.name,
        step,
        filter: data.settings.connections.map((cs: any) => `${cs.plugin}-${cs.connectionId}`),
        connection,
        onChangeConnection: setConnection,
        onPrev: handlePrev,
        onNext: handleNext,
        saving,
        onSave: handleSave,
        onCancel: handleCancel,
      }}
    >
      {children}
    </Context.Provider>
  );
};

export const useConnectionAdd = () => useContext(Context);
