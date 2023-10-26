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

import { useAppDispatch } from '@/app/hook';
import { Dialog, Message } from '@/components';
import { removeWebhook } from '@/features';
import { operator } from '@/utils';

interface Props {
  initialId: ID;
  onCancel: () => void;
  onSubmitAfter?: (id: ID) => void;
}

export const DeleteDialog = ({ initialId, onCancel, onSubmitAfter }: Props) => {
  const [operating, setOperating] = useState(false);

  const dispatch = useAppDispatch();

  const handleSubmit = async () => {
    const [success] = await operator(() => dispatch(removeWebhook(initialId)), {
      setOperating,
    });

    if (success) {
      onSubmitAfter?.(initialId);
      onCancel();
    }
  };

  return (
    <Dialog
      isOpen
      title="Delete this Webhook?"
      // style={{ width: 600 }}
      okText="Confirm"
      okLoading={operating}
      onCancel={onCancel}
      onOk={handleSubmit}
    >
      <Message content="This Webhook cannot be recovered once it’s deleted." />
    </Dialog>
  );
};
