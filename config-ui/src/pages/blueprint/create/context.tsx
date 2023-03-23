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

import React, { useState, useMemo, useContext } from 'react';
import { useHistory } from 'react-router-dom';
import dayjs from 'dayjs';

import { operator, formatTime } from '@/utils';

import { ModeEnum, FromEnum } from '../types';
import { validRawPlan } from '../utils';

import type { ContextType } from './types';
import * as API from './api';

export const Context = React.createContext<ContextType>({
  step: 1,

  name: 'MY BLUEPRINT',
  mode: ModeEnum.normal,
  connections: [],
  rawPlan: JSON.stringify([[]], null, '  '),
  cronConfig: '0 0 * * *',
  isManual: false,
  skipOnFail: false,
  timeAfter: null,

  onPrev: () => {},
  onNext: () => {},
  onSave: () => {},
  onSaveAndRun: () => {},

  onChangeMode: () => {},
  onChangeName: () => {},
  onChangeConnections: () => {},
  onChangeRawPlan: () => {},
  onChangeCronConfig: () => {},
  onChangeIsManual: () => {},
  onChangeSkipOnFail: () => {},
  onChangeTimeAfter: () => {},
});

interface Props {
  from: FromEnum;
  projectName: string;
  children: React.ReactNode;
}

export const ContextProvider = ({ from, projectName, children }: Props) => {
  const [step, setStep] = useState(1);

  const [name, setName] = useState(
    from === FromEnum.project ? `${window.decodeURIComponent(projectName)}-BLUEPRINT` : 'MY BLUEPRINT',
  );
  const [mode, setMode] = useState<ModeEnum>(ModeEnum.normal);
  const [connections, setConnections] = useState<MixConnection[]>([]);
  const [rawPlan, setRawPlan] = useState(JSON.stringify([[]], null, '  '));
  const [cronConfig, setCronConfig] = useState('0 0 * * *');
  const [isManual, setIsManual] = useState(false);
  const [skipOnFail, setSkipOnFail] = useState(true);
  const [timeAfter, setTimeAfter] = useState<string | null>(
    formatTime(dayjs().subtract(6, 'month').startOf('day').toDate(), 'YYYY-MM-DD[T]HH:mm:ssZ'),
  );

  const history = useHistory();

  const payload = useMemo(() => {
    const params: any = {
      name,
      projectName: projectName ? window.decodeURIComponent(projectName) : '',
      mode,
      enable: true,
      cronConfig,
      isManual,
      skipOnFail,
    };

    if (mode === ModeEnum.normal) {
      params.settings = {
        version: '2.0.0',
        timeAfter,
        connections: connections.map((cs) => {
          return {
            plugin: cs.plugin,
            connectionId: cs.connectionId,
            scopes: cs.scope.map((sc) => ({
              id: `${sc.id}`,
              entities: sc.entities,
            })),
          };
        }),
      };
    }

    if (mode === ModeEnum.advanced) {
      params.plan = !validRawPlan(rawPlan) ? JSON.parse(rawPlan) : JSON.stringify([[]], null, '  ');
      params.settings = null;
    }

    return params;
  }, [projectName, name, mode, connections, rawPlan, cronConfig, isManual, skipOnFail, timeAfter]);

  const handleSaveAfter = (id: ID) => {
    const path =
      from === FromEnum.blueprint ? `/blueprints/${id}` : `/projects/${window.encodeURIComponent(projectName)}`;

    history.push(path);
  };

  const handleSave = async () => {
    const [success, res] = await operator(() => API.createBlueprint(payload));

    if (success) {
      handleSaveAfter(res.id);
    }
  };

  const hanldeSaveAndRun = async () => {
    const [success, res] = await operator(async () => {
      const res = await API.createBlueprint(payload);
      return await API.runBlueprint(res.id);
    });

    if (success) {
      handleSaveAfter(res.blueprintId);
    }
  };

  const handlePrev = () => setStep(step - 1);
  const handleNext = () => setStep(step + 1);

  return (
    <Context.Provider
      value={{
        step,

        mode,
        name,
        connections,
        rawPlan,
        cronConfig,
        isManual,
        skipOnFail,
        timeAfter,

        onPrev: handlePrev,
        onNext: handleNext,
        onSave: handleSave,
        onSaveAndRun: hanldeSaveAndRun,

        onChangeMode: setMode,
        onChangeName: setName,
        onChangeConnections: setConnections,
        onChangeRawPlan: setRawPlan,
        onChangeCronConfig: setCronConfig,
        onChangeIsManual: setIsManual,
        onChangeSkipOnFail: setSkipOnFail,
        onChangeTimeAfter: setTimeAfter,
      }}
    >
      {children}
    </Context.Provider>
  );
};

export const useCreate = () => {
  return useContext(Context);
};
