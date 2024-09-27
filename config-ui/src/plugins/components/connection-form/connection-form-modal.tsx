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
import { theme, Modal } from 'antd';

import { getPluginConfig } from '@/plugins';

import { ConnectionForm } from '../connection-form';

import * as S from './styled';

interface Props {
  plugin: string;
  connectionId: ID;
  open: boolean;
  onCancel: () => void;
}

export const ConnectionFormModal = ({ plugin, connectionId, open, onCancel }: Props) => {
  const pluginConfig = useMemo(() => getPluginConfig(plugin), [plugin]);

  const {
    token: { colorPrimary },
  } = theme.useToken();

  return (
    <Modal
      open={open}
      width={820}
      centered
      title={
        <S.ModalTitle>
          <span className="icon">{pluginConfig.icon({ color: colorPrimary })}</span>
          <span className="name">Manage Connections: {pluginConfig.name}</span>
        </S.ModalTitle>
      }
      footer={null}
      onCancel={onCancel}
    >
      <ConnectionForm plugin={plugin} connectionId={connectionId} onSuccess={onCancel} />
    </Modal>
  );
};
