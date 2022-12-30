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

import { useState, useEffect, useMemo, useRef } from 'react';

import { operator } from '@/utils';

import type { PipelineType, TaskType } from '../types';
import { StatusEnum } from '../types';

import * as API from './api';

const pollTimer = 5000;

export interface UseDetailProps {
  id?: ID;
}

export const useDetail = ({ id }: UseDetailProps) => {
  const [version, setVersion] = useState(0);
  const [loading, setLoading] = useState(false);
  const [operating, setOperating] = useState(false);
  const [pipeline, setPipeline] = useState<PipelineType>();
  const [tasks, setTasks] = useState<TaskType[]>([]);

  const timer = useRef<any>();

  const getPipeline = async () => {
    if (!id) return;
    setLoading(true);
    try {
      const [pipeRes, taskRes] = await Promise.all([API.getPipeline(id), API.getPipelineTasks(id)]);

      setPipeline(pipeRes);
      setTasks(taskRes.tasks);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    getPipeline();
  }, []);

  useEffect(() => {
    timer.current = setInterval(() => {
      getPipeline();
    }, pollTimer);
    return () => clearInterval(timer.current);
  }, [version]);

  useEffect(() => {
    if (pipeline && [StatusEnum.COMPLETED, StatusEnum.FAILED, StatusEnum.CANCELLED].includes(pipeline.status)) {
      clearInterval(timer.current);
    }
  }, [pipeline]);

  const handlePipelineCancel = async () => {
    if (!id) return;
    const [success] = await operator(() => API.deletePipeline(id), {
      setOperating,
    });

    if (success) {
      getPipeline();
      setVersion(version + 1);
    }
  };

  const handlePipelineRerun = async () => {
    if (!id) return;
    const [success] = await operator(() => API.pipeLineRerun(id), {
      setOperating,
    });

    if (success) {
      getPipeline();
      setVersion(version + 1);
    }
  };

  const handleTaskRertun = async (id: ID) => {
    if (!id) return;
    const [success] = await operator(() => API.taskRerun(id), {
      setOperating,
    });

    if (success) {
      getPipeline();
      setVersion(version + 1);
    }
  };

  return useMemo(
    () => ({
      loading,
      operating,
      pipeline,
      tasks,
      onCancel: handlePipelineCancel,
      onRerun: handlePipelineRerun,
      onRerunTask: handleTaskRertun,
    }),
    [loading, operating, pipeline, tasks],
  );
};
