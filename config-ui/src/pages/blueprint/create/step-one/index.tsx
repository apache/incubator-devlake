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

import React, { useMemo } from 'react';
import { pick } from 'lodash';
import { InputGroup, Icon } from '@blueprintjs/core';

import { useConnection, ConnectionStatusEnum } from '@/store';
import { Card, Divider, MultiSelector, Loading } from '@/components';

import { ModeEnum, FromEnum } from '../../types';
import { AdvancedEditor } from '../../components';
import { useCreateBP } from '../bp-context';

import * as S from './styled';

interface Props {
  from: FromEnum;
}

export const StepOne = ({ from }: Props) => {
  const { connections, onTest } = useConnection();

  const {
    mode,
    name,
    rawPlan,
    uniqueList,
    scopeMap,
    onChangeMode,
    onChangeName,
    onChangeRawPlan,
    onChangeUniqueList,
    onChangeScopeMap,
  } = useCreateBP();

  const fromProject = useMemo(() => from === FromEnum.project, [from]);

  return (
    <>
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
                const uniqueList = selectedItems.map((sc) => sc.unique);
                onChangeUniqueList(uniqueList);
                onChangeScopeMap(pick(scopeMap, uniqueList));
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
    </>
  );
};
