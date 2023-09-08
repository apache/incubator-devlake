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

import { InputGroup } from '@blueprintjs/core';

import { Dialog, FormItem, CopyText } from '@/components';

import type { UseViewOrEditProps } from './use-view-or-edit';
import { useViewOrEdit } from './use-view-or-edit';
import * as S from './styled';

interface Props extends UseViewOrEditProps {
  type: 'edit' | 'show';
  isOpen: boolean;
  onCancel: () => void;
}

export const WebhookViewOrEditDialog = ({ type, isOpen, onCancel, ...props }: Props) => {
  const { saving, name, record, onChangeName, onSubmit } = useViewOrEdit({
    ...props,
  });

  const handleSubmit = () => {
    if (type === 'edit') {
      onSubmit();
    }

    onCancel();
  };

  return (
    <Dialog
      isOpen={isOpen}
      title="View/Edit Webhook"
      style={{ width: 820 }}
      okText={type === 'edit' ? 'Save' : 'Done'}
      okDisabled={!name}
      okLoading={saving}
      onCancel={onCancel}
      onOk={handleSubmit}
    >
      <S.Wrapper>
        <h3>Webhook Name *</h3>
        <p>
          Copy the following POST URLs to your issue tracking or CI tools to push `Incidents` and `Deployments` by
          making a POST to DevLake.
        </p>
        <InputGroup disabled={type !== 'edit'} value={name} onChange={(e) => onChangeName(e.target.value)} />
        <p>
          Copy the following URLs to your issue tracking tool for Incidents and CI tool for Deployments by making a POST
          to DevLake.
        </p>
        <FormItem label="Incidents">
          <h5>Post to register an incident</h5>
          <CopyText content={record.postIssuesEndpoint} />
          <h5>Post to close a registered incident</h5>
          <CopyText content={record.closeIssuesEndpoint} />
        </FormItem>
        <FormItem label="Deployment">
          <h5>Post to register a deployment</h5>
          <CopyText content={record.postDeploymentsCurl} />
        </FormItem>
      </S.Wrapper>
    </Dialog>
  );
};
