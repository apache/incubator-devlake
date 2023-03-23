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
import { InputGroup, Icon, Button, Intent, Position, Colors } from '@blueprintjs/core';
import { Tooltip2 } from '@blueprintjs/popover2';

import { useConnection, ConnectionStatusEnum } from '@/store';
import { Card, Divider, MultiSelector, Loading } from '@/components';

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
        return 'Blueprint Name: Enter a valid Name';
      case name.length < 3:
        return 'Blueprint Name: Name too short, 3 chars minimum.';
      case mode === ModeEnum.advanced && validRawPlan(rawPlan):
        return 'Advanced Mode: Invalid/Empty Configuration';
      case mode === ModeEnum.normal && !uniqueList.length:
        return 'Normal Mode: No Data Connections selected.';
      case mode === ModeEnum.normal &&
        !connections
          .filter((cs) => uniqueList.includes(cs.unique))
          .every((cs) => cs.status === ConnectionStatusEnum.ONLINE):
        return 'Normal Mode: Has some offline connections';
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
            <h3>Select Connections</h3>
            <p>Select from existing or create new connections</p>
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
                  selectedItems.map((sc) => ({
                    unique: sc.unique,
                    plugin: sc.plugin,
                    connectionId: sc.id,
                    name: sc.name,
                    icon: sc.icon,
                    scope: [],
                    origin: [],
                  })),
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
        <Button
          intent={Intent.PRIMARY}
          disabled={!!error}
          icon={
            error ? (
              <Tooltip2 defaultIsOpen placement={Position.TOP} content={error}>
                <Icon icon="warning-sign" color={Colors.ORANGE5} style={{ margin: 0 }} />
              </Tooltip2>
            ) : null
          }
          text="Next Step"
          onClick={onNext}
        />
      </S.Btns>
    </S.Wrapper>
  );
};
