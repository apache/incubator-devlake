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

import { Modal } from 'antd';

import { ConnectionName } from '../connection-name';
import { ConnectionForm } from '../connection-form';

interface Props {
  plugin: string;
  connectionId?: ID;
  open: boolean;
  onCancel: () => void;
  onSuccess?: (id: ID) => void;
}

export const ConnectionFormModal = ({ plugin, connectionId, open, onCancel, onSuccess }: Props) => {
  const handleSuccess = (id: ID) => {
    onSuccess?.(id);
    onCancel();
  };

  return (
    <Modal
      open={open}
      width={820}
      centered
      title={
        <ConnectionName
          plugin={plugin}
          connectionId={connectionId}
          customName={(pluginName) => `Manage Connections: ${pluginName}`}
        />
      }
      footer={null}
      onCancel={() => onCancel()}
    >
      <ConnectionForm plugin={plugin} connectionId={connectionId} onSuccess={handleSuccess} />
    </Modal>
  );
};
