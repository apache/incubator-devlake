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
import { Modal } from 'antd';
import { MillerColumns } from '@mints/miller-columns';

import { useAppSelector } from '@/hooks';
import { Block, Loading } from '@/components';
import { selectWebhooks } from '@/features/connections';
import { IWebhook } from '@/types';

import * as S from '../styled';

interface Props {
  open: boolean;
  saving: boolean;
  onCancel: () => void;
  onSubmit: (items: IWebhook[]) => void;
}

export const SelectorDialog = ({ open, saving, onCancel, onSubmit }: Props) => {
  const [selectedIds, setSelectedIds] = useState<ID[]>([]);

  const webhooks = useAppSelector(selectWebhooks);

  const handleSubmit = () => onSubmit(webhooks.filter((it) => selectedIds.includes(it.id)));

  return (
    <Modal
      open={open}
      width={820}
      centered
      title="Select Existing Webhooks"
      okText="Confirm"
      okButtonProps={{
        disabled: !selectedIds.length,
        loading: saving,
      }}
      onCancel={onCancel}
      onOk={handleSubmit}
    >
      <S.Wrapper>
        <Block title="Webhooks" description="Select an existing Webhook to import to the current project.">
          <MillerColumns
            bordered
            theme={{
              colorPrimary: '#7497f7',
              borderColor: '#dbe4fd',
            }}
            items={webhooks.map((it) => ({
              parentId: null,
              id: it.id,
              title: it.name,
            }))}
            columnHeight={160}
            renderLoading={() => <Loading size={20} style={{ padding: '4px 12px' }} />}
            selectable
            selectedIds={selectedIds}
            onSelectedIds={setSelectedIds}
          />
        </Block>
      </S.Wrapper>
    </Modal>
  );
};

export default SelectorDialog;
