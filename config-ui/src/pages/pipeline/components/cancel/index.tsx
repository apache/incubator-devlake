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

import { IconButton } from '@/components';
import { operator } from '@/utils';

import { StatusEnum } from '../../types';
import * as API from '../../api';

import { usePipeline } from '../context';

interface Props {
  id: ID;
  status: StatusEnum;
}

export const PipelineCancel = ({ id, status }: Props) => {
  const [canceling, setCanceling] = useState(false);

  const { setVersion } = usePipeline();

  const handleSubmit = async () => {
    const [success] = await operator(() => API.deletePipeline(id), {
      setOperating: setCanceling,
    });

    if (success) {
      setVersion((v) => v + 1);
    }
  };

  if (![StatusEnum.ACTIVE, StatusEnum.RUNNING, StatusEnum.RERUN].includes(status)) {
    return null;
  }

  return <IconButton loading={canceling} icon="disable" tooltip="Cancel" onClick={handleSubmit} />;
};
