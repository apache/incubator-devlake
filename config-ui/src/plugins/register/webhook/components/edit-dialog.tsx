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
import { InputGroup } from '@blueprintjs/core';

import API from '@/api';
import { Dialog, FormItem } from '@/components';
import { useConnections } from '@/hooks';
import { operator } from '@/utils';

interface Props {
  initialId: ID;
  onCancel: () => void;
}

export const EditDialog = ({ initialId, onCancel }: Props) => {
  const [name, setName] = useState('');
  const [operating, setOperating] = useState(false);

  const { onRefresh } = useConnections({ plugin: 'webhook' });

  useEffect(() => {
    (async () => {
      const res = await API.plugin.webhook.get(initialId);
      setName(res.name);
    })();
  }, [initialId]);

  const handleSubmit = async () => {
    const [success] = await operator(() => API.plugin.webhook.update(initialId, { name }), {
      setOperating,
    });

    if (success) {
      onRefresh('webhook');
      onCancel();
    }
  };

  return (
    <Dialog
      style={{ width: 820 }}
      isOpen
      title="Edit Webhook Name"
      okLoading={operating}
      okDisabled={!name}
      okText="Save"
      onCancel={onCancel}
      onOk={handleSubmit}
    >
      <FormItem label="Name" required>
        <InputGroup value={name} onChange={(e) => setName(e.target.value)} />
      </FormItem>
    </Dialog>
  );
};
