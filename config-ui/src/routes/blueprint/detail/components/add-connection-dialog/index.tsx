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
import { useNavigate } from 'react-router-dom';
import { PlusOutlined } from '@ant-design/icons';
import { Modal, Select, Space, Button } from 'antd';
import styled from 'styled-components';

import { Block } from '@/components';
import { selectAllConnections } from '@/features';
import { useAppSelector } from '@/hooks';
import { PluginName, DataScopeSelect } from '@/plugins';
import { IConnection } from '@/types';

const Option = styled.div`
  display: flex;
  align-items: center;

  .icon {
    display: inline-block;
    width: 24px;
    height: 24px;

    & > svg {
      width: 100%;
      height: 100%;
    }
  }

  .name {
    margin-left: 8px;
    max-width: 90%;
    text-overflow: ellipsis;
    white-space: nowrap;
    overflow: hidden;
  }
`;

interface Props {
  disabled: string[];
  onCancel: () => void;
  onSubmit: (value: any) => void;
}

export const AddConnectionDialog = ({ disabled = [], onCancel, onSubmit }: Props) => {
  const [step, setStep] = useState(1);
  const [selectedValue, setSelectedValue] = useState<string>();

  const navigate = useNavigate();

  const connections = useAppSelector(selectAllConnections);

  const options = useMemo(
    () =>
      [{ value: '' }].concat(
        connections
          .filter((cs) => (disabled.length ? !disabled.includes(cs.unique) : true))
          .map((cs) => ({
            plugin: cs.plugin,
            label: cs.name,
            value: cs.unique,
          })),
      ),
    [connections, disabled],
  );

  const selectedConnection = useMemo(
    () => connections.find((cs) => cs.unique === selectedValue) as IConnection,
    [selectedValue],
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
    <Modal open width={820} centered title={`Add a Connection - Step ${step}`} footer={null} onCancel={onCancel}>
      {step === 1 && (
        <>
          <Block
            title="Data Connections"
            description="Select from existing Data Connections or create a new one."
            required
          >
            <Select
              style={{ width: 384 }}
              placeholder="Select..."
              options={options}
              optionRender={(option, { index }) => {
                if (index === 0) {
                  return (
                    <Button size="small" type="link" icon={<PlusOutlined />}>
                      Add New Connection
                    </Button>
                  );
                }
                return (
                  <Option>
                    <PluginName plugin={option.data.plugin} name={option.data.label} />
                  </Option>
                );
              }}
              onChange={(value) => {
                if (!value) {
                  navigate('/connections');
                }
                setSelectedValue(value);
              }}
            />
          </Block>
          <Space style={{ display: 'flex', justifyContent: 'flex-end' }}>
            <Button onClick={onCancel}>Cancel</Button>
            <Button type="primary" disabled={!selectedConnection} onClick={() => setStep(2)}>
              Next
            </Button>
          </Space>
        </>
      )}
      {step === 2 && selectedConnection && (
        <DataScopeSelect
          plugin={selectedConnection.plugin}
          connectionId={selectedConnection.id}
          onCancel={onCancel}
          onSubmit={handleSubmit}
        />
      )}
    </Modal>
  );
};
