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

import { useState, useEffect } from 'react';
import { Modal, Input } from 'antd';

import { useAppDispatch, useAppSelector } from '@/hooks';
import { Block } from '@/components';
import { updateWebhook, selectWebhook } from '@/features';
import { operator } from '@/utils';

interface Props {
  initialId: ID;
  onCancel: () => void;
}

export const EditDialog = ({ initialId, onCancel }: Props) => {
  const [name, setName] = useState('');
  const [operating, setOperating] = useState(false);

  const dispatch = useAppDispatch();
  const webhook = useAppSelector((state) => selectWebhook(state, initialId));

  useEffect(() => {
    setName(webhook?.name ?? '');
  }, [webhook]);

  const handleSubmit = async () => {
    const [success] = await operator(() => dispatch(updateWebhook({ id: initialId, name })), {
      setOperating,
    });

    if (success) {
      onCancel();
    }
  };

  return (
    <Modal
      open
      width={820}
      centered
      title="Edit Webhook Name"
      okText="Save"
      okButtonProps={{
        disabled: !name,
        loading: operating,
      }}
      onCancel={onCancel}
      onOk={handleSubmit}
    >
      <Block title="Name" required>
        <Input value={name} onChange={(e) => setName(e.target.value)} />
      </Block>
    </Modal>
  );
};
