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

import { useState, useEffect, useMemo } from 'react';
import { useHistory } from 'react-router-dom';
import { Icon, Button, Switch, Colors, Intent } from '@blueprintjs/core';
import dayjs from 'dayjs';

import { getCron } from '@/config';

import type { BlueprintType } from '../../types';
import { ModeEnum } from '../../types';
import { validRawPlan } from '../../utils';
import { AdvancedEditor } from '../../components';

import { UpdateNameDialog, UpdatePolicyDialog, ConnectionList } from '../components';
import * as S from '../styled';

type Type = 'name' | 'frequency' | 'scope' | 'transformation';

interface Props {
  paths: string[];
  blueprint: BlueprintType;
  operating: boolean;
  onUpdate: (bp: any) => void;
}

export const Configuration = ({ paths, blueprint, operating, onUpdate }: Props) => {
  const [type, setType] = useState<Type>();
  const [rawPlan, setRawPlan] = useState('');

  const history = useHistory();

  useEffect(() => {
    setRawPlan(JSON.stringify(blueprint.plan, null, '  '));
  }, [blueprint]);

  const cron = useMemo(() => getCron(blueprint.isManual, blueprint.cronConfig), [blueprint]);

  const handleCancel = () => {
    setType(undefined);
  };

  const handleUpdateName = async (name: string) => {
    await onUpdate({ name });
    handleCancel();
  };

  const handleUpdatePolicy = async (policy: any) => {
    await onUpdate(policy);
    handleCancel();
  };

  const handleToggleEnabled = (checked: boolean) => onUpdate({ enable: checked });

  const handleUpdatePlan = () =>
    onUpdate({
      plan: !validRawPlan(rawPlan) ? JSON.parse(rawPlan) : JSON.stringify([[]], null, '  '),
    });

  return (
    <S.ConfigurationPanel>
      <div className="top">
        <ul>
          <li>
            <h3>Name</h3>
            <div className="detail">
              <span>{blueprint.name}</span>
              <Icon icon="annotation" color={Colors.BLUE2} onClick={() => setType('name')} />
            </div>
          </li>
          <li>
            <h3>Sync Policy</h3>
            <div className="detail">
              <span>
                {cron.label} {cron.value !== 'manual' ? dayjs(cron.nextTime).format('HH:mm A') : null}
              </span>
              <Icon icon="annotation" color={Colors.BLUE2} onClick={() => setType('frequency')} />
            </div>
          </li>
        </ul>
        <Switch
          label="Blueprint Enabled"
          checked={blueprint.enable}
          onChange={(e) => handleToggleEnabled((e.target as HTMLInputElement).checked)}
        />
      </div>
      {blueprint.mode === ModeEnum.normal && (
        <div className="bottom">
          <h3>
            <span>Connections</span>
            <Button small intent={Intent.PRIMARY} onClick={() => history.push(paths[0])}>
              Add a Connection
            </Button>
          </h3>
          <ConnectionList path={paths[1]} blueprint={blueprint} />
        </div>
      )}
      {blueprint.mode === ModeEnum.advanced && (
        <div className="bottom">
          <h3>JSON Configuration</h3>
          <AdvancedEditor value={rawPlan} onChange={setRawPlan} />
          <div className="btns">
            <Button intent={Intent.PRIMARY} text="Save" onClick={handleUpdatePlan} />
          </div>
        </div>
      )}
      {type === 'name' && (
        <UpdateNameDialog
          name={blueprint.name}
          operating={operating}
          onCancel={handleCancel}
          onSubmit={handleUpdateName}
        />
      )}
      {type === 'frequency' && (
        <UpdatePolicyDialog
          blueprint={blueprint}
          isManual={blueprint.isManual}
          cronConfig={blueprint.cronConfig}
          skipOnFail={blueprint.skipOnFail}
          timeAfter={blueprint.settings?.timeAfter}
          operating={operating}
          onCancel={handleCancel}
          onSubmit={handleUpdatePolicy}
        />
      )}
    </S.ConfigurationPanel>
  );
};
