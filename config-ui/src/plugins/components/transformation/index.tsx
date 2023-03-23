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
import { Button, Intent } from '@blueprintjs/core';

import { Table, IconButton } from '@/components';
import { TransformationSelect, getPluginId } from '@/plugins';

import * as API from './api';
import * as S from './styled';

interface Props {
  connections: MixConnection[];
  cancelBtnProps?: {
    text?: string;
  };
  submitBtnProps?: {
    text: string;
  };
  onCancel?: () => void;
  onSubmit?: () => void;
}

export const TransformationBind = ({ connections, cancelBtnProps, submitBtnProps, onCancel, onSubmit }: Props) => {
  const [selected, setSelected] = useState<Record<string, ID[]>>({});
  const [connection, setConnection] = useState<MixConnection>();

  const handleCancel = () => setConnection(undefined);

  const handleSubmit = async (tid: ID, connection: MixConnection) => {
    const { unique, plugin, connectionId } = connection;
    const scopeIds = selected[unique];
    const scopes = connection.scope.filter((sc) => scopeIds.includes(sc.id));

    await Promise.all(
      scopes.map((scope) =>
        API.updateDataScope(plugin, connectionId, scope.id, {
          ...scope,
          transformationRuleId: tid,
        }),
      ),
    );
  };

  return (
    <S.List>
      {connections.map((cs) => (
        <S.Item key={cs.unique}>
          {connections.length !== 1 && (
            <S.Title>
              <img src={cs.icon} alt="" />
              <span>{cs.name}</span>
            </S.Title>
          )}
          <S.Action>
            <Button
              intent={Intent.PRIMARY}
              icon="annotation"
              disabled={!selected[cs.unique] || !selected[cs.unique].length}
              onClick={() => setConnection(cs)}
            >
              Select Transformation
            </Button>
          </S.Action>
          <Table
            columns={[
              { title: 'Data Scope', dataIndex: 'name', key: 'name' },
              {
                title: 'Transformation',
                dataIndex: 'transformationRuleName',
                key: 'transformation',
                align: 'center',
                render: (val, row) => (
                  <div>
                    <span>{val ?? 'N/A'}</span>
                    <IconButton
                      icon="annotation"
                      tooltip="Select Transformation"
                      onClick={() => {
                        setSelected({
                          ...selected,
                          [`${cs.unique}`]: [row.id],
                        });
                        setConnection(cs);
                      }}
                    />
                  </div>
                ),
              },
            ]}
            dataSource={cs.origin}
            rowSelection={{
              rowKey: getPluginId(cs.plugin),
              selectedRowKeys: selected[cs.unique],
              onChange: (selectedRowKeys) => setSelected({ ...selected, [`${cs.unique}`]: selectedRowKeys }),
            }}
          />
        </S.Item>
      ))}
      <S.Btns>
        <Button outlined intent={Intent.PRIMARY} text="Cancel" onClick={onCancel} {...cancelBtnProps} />
        <Button outlined intent={Intent.PRIMARY} text="Save" onClick={onSubmit} {...submitBtnProps} />
      </S.Btns>
      {connection && (
        <TransformationSelect
          plugin={connection.plugin}
          connectionId={connection.connectionId}
          onCancel={handleCancel}
          onSubmit={(tid) => handleSubmit(tid, connection)}
        />
      )}
    </S.List>
  );
};
