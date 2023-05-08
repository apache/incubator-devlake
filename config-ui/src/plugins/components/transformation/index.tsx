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

import { IconButton, Table } from '@/components';
import { getPluginId, TransformationSelect } from '@/plugins';

import * as API from './api';
import * as S from './styled';

interface Props {
  connections: MixConnection[];
  cancelBtnProps?: {
    text?: string;
  };
  submitBtnProps?: {
    text?: string;
    loading?: boolean;
  };
  noFooter?: boolean;
  onCancel?: () => void;
  onSubmit?: (connections: MixConnection[]) => void;
  onNext?: () => void;
}

export const Transformation = ({
  connections,
  cancelBtnProps,
  submitBtnProps,
  noFooter,
  onCancel,
  onSubmit,
  onNext,
}: Props) => {
  const [selected, setSelected] = useState<Record<string, ID[]>>({});
  const [connection, setConnection] = useState<MixConnection>();
  const [tid, setTid] = useState<ID>();

  const handleCancel = () => {
    setConnection(undefined);
    setTid(undefined);
  };

  const handleSubmit = async (tid: ID, connection: MixConnection, connections: MixConnection[]) => {
    const { unique, plugin, connectionId } = connection;
    const scopeIds = selected[unique];
    const scopes = await Promise.all(
      scopeIds.map(async (scopeId) => {
        const scope = await API.getDataScope(plugin, connectionId, scopeId);
        return await API.updateDataScope(plugin, connectionId, scopeId, {
          ...scope,
          transformationRuleId: tid,
        });
      }),
    );
    onSubmit?.(
      connections.map((cs) => {
        if (cs.unique !== unique) {
          return cs;
        }

        const origin = cs.origin.map((sc) => {
          if (!scopeIds.includes(sc[getPluginId(cs.plugin)])) {
            return sc;
          }
          return scopes.find((it) => it[getPluginId(cs.plugin)] === sc[getPluginId(cs.plugin)]);
        });

        return { ...cs, origin };
      }),
    );
    setSelected({
      ...selected,
      [`${unique}`]: [],
    });
    handleCancel();
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
          {cs.transformationType === 'for-connection' && (
            <S.Action>
              <Button
                intent={Intent.PRIMARY}
                icon="many-to-one"
                disabled={!selected[cs.unique] || !selected[cs.unique].length}
                onClick={() => setConnection(cs)}
              >
                Associate Transformation
              </Button>
            </S.Action>
          )}
          <Table
            columns={[
              { title: 'Data Scope', dataIndex: 'name', key: 'name' },
              {
                title: 'Transformation',
                dataIndex: 'transformationRuleName',
                key: 'transformation',
                align: 'center',
                render: (val, row) =>
                  cs.transformationType === 'none' ? (
                    'N/A'
                  ) : (
                    <div>
                      <span>{val ?? 'N/A'}</span>
                      <IconButton
                        icon="one-to-one"
                        tooltip="Associate Transformation"
                        onClick={() => {
                          setSelected({
                            ...selected,
                            [`${cs.unique}`]: [row[getPluginId(cs.plugin)]],
                          });
                          setConnection(cs);
                          setTid(row.transformationRuleId);
                        }}
                      />
                    </div>
                  ),
              },
            ]}
            dataSource={cs.origin}
            rowSelection={
              cs.transformationType === 'for-connection'
                ? {
                    rowKey: getPluginId(cs.plugin),
                    selectedRowKeys: selected[cs.unique],
                    onChange: (selectedRowKeys) => setSelected({ ...selected, [`${cs.unique}`]: selectedRowKeys }),
                  }
                : undefined
            }
          />
        </S.Item>
      ))}
      {!noFooter && (
        <S.Btns>
          <Button outlined intent={Intent.PRIMARY} text="Previous Step" onClick={onCancel} {...cancelBtnProps} />
          <Button intent={Intent.PRIMARY} text="Next Step" onClick={onNext} {...submitBtnProps} />
        </S.Btns>
      )}
      {connection && (
        <TransformationSelect
          plugin={connection.plugin}
          connectionId={connection.connectionId}
          scopeId={selected[connection.unique][0]}
          transformationId={tid}
          transformationType={connection.transformationType}
          onCancel={handleCancel}
          onSubmit={(tid) => handleSubmit(tid, connection, connections)}
        />
      )}
    </S.List>
  );
};
