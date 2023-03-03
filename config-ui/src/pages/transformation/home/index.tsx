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

import React, { useState, useMemo } from 'react';
import { useHistory } from 'react-router-dom';
import { ButtonGroup, Button, Icon, Intent } from '@blueprintjs/core';

import { PageHeader, Table, ColumnType, Dialog, Selector, IconButton } from '@/components';
import type { PluginConfigType } from '@/plugins';
import { TransformationContextProvider, TransformationContextConsumer, TransformationItemType } from '@/store';

import * as S from './styled';

export const TransformationHomePage = () => {
  const [active, setActive] = useState('All');
  const [isOpen, setIsOpen] = useState(false);
  const [selectedPlugin, setSelectedPlugin] = useState<PluginConfigType>();

  const history = useHistory();

  const handleCreate = () => setIsOpen(true);

  const columns = useMemo(
    () =>
      [
        {
          title: 'Name',
          dataIndex: 'name',
          key: 'name',
        },
        {
          title: 'Data Source',
          dataIndex: 'plugin',
          key: 'plugin',
        },
        {
          title: '',
          key: 'action',
          width: 100,
          align: 'center',
          render: (_, row) =>
            row.plugin !== 'jira' && (
              <IconButton
                icon="cog"
                tooltip="Detail"
                onClick={() => history.push(`/transformations/${row.plugin}/${row.id}`)}
              />
            ),
        },
      ] as ColumnType<TransformationItemType>,
    [],
  );

  return (
    <TransformationContextProvider>
      <TransformationContextConsumer>
        {({ plugins, transformations }) => (
          <PageHeader breadcrumbs={[{ name: 'Transformations', path: '/transformations' }]}>
            <S.Wrapper>
              <div className="action">
                <ButtonGroup>
                  <Button
                    intent={active === 'All' ? Intent.PRIMARY : Intent.NONE}
                    text="All"
                    onClick={() => setActive('All')}
                  />
                  {plugins.map((p) => (
                    <Button
                      key={p.plugin}
                      intent={active === p.plugin ? Intent.PRIMARY : Intent.NONE}
                      text={p.name}
                      onClick={() => setActive(p.plugin)}
                    />
                  ))}
                </ButtonGroup>
                <Button icon="plus" intent={Intent.PRIMARY} text="New Transformation" onClick={handleCreate} />
              </div>
              <Table
                columns={columns}
                dataSource={transformations.filter((ts) => (active === 'All' ? true : ts.plugin === active))}
                noData={{
                  text: 'There is no transformation yet. Please add a new transformation.',
                  btnText: 'New Transformation',
                  onCreate: handleCreate,
                }}
              />
              <Dialog
                isOpen={isOpen}
                title="Select a Data Source"
                okText="Continue"
                okDisabled={!selectedPlugin || selectedPlugin.plugin === 'jira'}
                onOk={() => history.push(`/transformations/${selectedPlugin?.plugin}/create`)}
                onCancel={() => setIsOpen(false)}
              >
                <S.DialogWrapper>
                  <p>Select from the supported data sources</p>
                  <Selector
                    items={plugins}
                    getKey={(it) => it.plugin}
                    getName={(it) => it.name}
                    selectedItem={selectedPlugin}
                    onChangeItem={(selectedItem) => setSelectedPlugin(selectedItem)}
                  />
                  {selectedPlugin?.plugin === 'jira' && (
                    <div className="warning">
                      <Icon icon="error" />
                      <span>
                        Because Jira transformation is specific to every Jira connection, you can only add a
                        Transformation in Blueprints.
                      </span>
                    </div>
                  )}
                </S.DialogWrapper>
              </Dialog>
            </S.Wrapper>
          </PageHeader>
        )}
      </TransformationContextConsumer>
    </TransformationContextProvider>
  );
};
