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
  type: 'pipeline' | 'task';
  id: ID;
  status: StatusEnum;
}

export const PipelineRerun = ({ type, id, status }: Props) => {
  const [reruning, setReruning] = useState(false);

  const { setVersion } = usePipeline();

  const handleSubmit = async () => {
    const [success] = await operator(() => (type === 'task' ? API.taskRerun(id) : API.pipelineRerun(id)), {
      setOperating: setReruning,
    });

    if (success) {
      setVersion((v) => v + 1);
    }
  };

  if (![StatusEnum.COMPLETED, StatusEnum.PARTIAL, StatusEnum.FAILED, StatusEnum.CANCELLED].includes(status)) {
    return null;
  }

  if (type === 'task') {
    return <IconButton loading={reruning} icon="repeat" tooltip="Rerun task" onClick={handleSubmit} />;
  }

  return <IconButton loading={reruning} icon="repeat" tooltip="Rerun failed tasks" onClick={handleSubmit} />;
};
