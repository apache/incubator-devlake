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

import { useMemo } from 'react';
import { Link } from 'react-router-dom';
import { InputGroup, Icon, Button, Intent } from '@blueprintjs/core';

import { Card, Divider, MultiSelector, Loading } from '@/components';
import { getPluginConfig } from '@/plugins';
import { useConnection, ConnectionStatusEnum } from '@/store';

import { ModeEnum, FromEnum } from '../../types';
import { AdvancedEditor } from '../../components';
import { validRawPlan } from '../../utils';

import { useCreate } from '../context';

import * as S from './styled';

interface Props {
  from: FromEnum;
}

export const Step1 = ({ from }: Props) => {
  const { connections, onTest } = useConnection();
  const { mode, name, rawPlan, onChangeMode, onChangeName, onChangeConnections, onChangeRawPlan, onNext, ...props } =
    useCreate();

  const fromProject = useMemo(() => from === FromEnum.project, [from]);
  const uniqueList = useMemo(() => props.connections.map((sc) => sc.unique), [props.connections]);

  const error = useMemo(() => {
    switch (true) {
      case !name:
        return true;
      case name.length < 3:
        return true;
      case mode === ModeEnum.advanced && validRawPlan(rawPlan):
        return true;
      case mode === ModeEnum.normal && !uniqueList.length:
        return true;
      case mode === ModeEnum.normal &&
        !connections
          .filter((cs) => uniqueList.includes(cs.unique))
          .every((cs) => cs.status === ConnectionStatusEnum.ONLINE):
        return true;
      default:
        return false;
    }
  }, [mode, name, connections, props.connections, rawPlan]);

  return (
    <S.Wrapper>
      <Card className="card">
        <h2>Blueprint Name</h2>
        <Divider />
        <p>Give your Blueprint a unique name to help you identify it in the future.</p>
        <InputGroup placeholder="Enter Blueprint Name" value={name} onChange={(e) => onChangeName(e.target.value)} />
      </Card>

      {mode === ModeEnum.normal && (
        <>
          <Card className="card">
            <h2>Add Data Connections</h2>
            <Divider />
            <p>
              If you have not created any connections yet, please <Link to="/connections">create connections</Link>{' '}
              first.
            </p>
            <MultiSelector
              placeholder="Select Connections..."
              items={connections}
              getKey={(it) => it.unique}
              getName={(it) => it.name}
              getIcon={(it) => it.icon}
              selectedItems={connections.filter((cs) => uniqueList.includes(cs.unique))}
              onChangeItems={(selectedItems) => {
                const lastItem = selectedItems[selectedItems.length - 1];
                if (lastItem) {
                  onTest(lastItem);
                }
                onChangeConnections(
                  selectedItems.map((sc) => {
                    const config = getPluginConfig(sc.plugin);
                    return {
                      unique: sc.unique,
                      plugin: sc.plugin,
                      connectionId: sc.id,
                      name: sc.name,
                      icon: sc.icon,
                      scope: [],
                      origin: [],
                      transformationType: config.transformationType,
                    };
                  }),
                );
              }}
            />
            <S.ConnectionList>
              {connections
                .filter((cs) => uniqueList.includes(cs.unique))
                .map((cs) => (
                  <li key={cs.unique}>
                    <span className="name">{cs.name}</span>
                    <span className={`status ${cs.status}`}>
                      {cs.status === ConnectionStatusEnum.TESTING && <Loading size={14} style={{ marginRight: 4 }} />}
                      {cs.status === ConnectionStatusEnum.OFFLINE && (
                        <Icon
                          size={14}
                          icon="repeat"
                          style={{ marginRight: 4, cursor: 'pointer' }}
                          onClick={() => onTest(cs)}
                        />
                      )}
                      {cs.status}
                    </span>
                  </li>
                ))}
            </S.ConnectionList>
          </Card>
          {!fromProject && (
            <S.Tips>
              <span>To customize how tasks are executed in the blueprint, please use </span>
              <span onClick={() => onChangeMode(ModeEnum.advanced)}>Advanced Mode.</span>
            </S.Tips>
          )}
        </>
      )}

      {mode === ModeEnum.advanced && !fromProject && (
        <>
          <Card className="card">
            <h2>JSON Configuration</h2>
            <Divider />
            <AdvancedEditor value={rawPlan} onChange={onChangeRawPlan} />
          </Card>
          <S.Tips>
            <span>To visually define blueprint tasks, please use </span>
            <span onClick={() => onChangeMode(ModeEnum.normal)}>Normal Mode.</span>
          </S.Tips>
        </>
      )}

      <S.Btns>
        <span></span>
        <Button intent={Intent.PRIMARY} disabled={error} text="Next Step" onClick={onNext} />
      </S.Btns>
    </S.Wrapper>
  );
};
