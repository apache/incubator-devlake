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

import { useState, useMemo, useEffect } from 'react';

import { operator } from '@/utils';

import * as API from '../api';

export interface UseViewOrEditProps {
  initialID: ID;
  onSubmitAfter?: (id: ID) => void;
}

export const useViewOrEdit = ({ initialID, onSubmitAfter }: UseViewOrEditProps) => {
  const [saving, setSaving] = useState(false);
  const [name, setName] = useState('');
  const [record, setRecord] = useState({
    postIssuesEndpoint: '',
    closeIssuesEndpoint: '',
    postDeploymentsCurl: '',
  });

  const prefix = useMemo(() => `${window.location.origin}/api`, []);

  useEffect(() => {
    (async () => {
      const res = await API.getConnection(initialID);
      setName(res.name);
      setRecord({
        postIssuesEndpoint: `${prefix}${res.postIssuesEndpoint}`,
        closeIssuesEndpoint: `${prefix}${res.closeIssuesEndpoint}`,
        postDeploymentsCurl: `curl ${prefix}${res.postPipelineDeployTaskEndpoint} -X 'POST' -d "{
      \\"commit_sha\\":\\"the sha of deployment commit\\",
      \\"repo_url\\":\\"the repo URL of the deployment commit\\",
      \\"start_time\\":\\"Optional, eg. 2020-01-01T12:00:00+00:00\\"
      }"`,
      });
    })();
  }, [initialID]);

  const handleUpdate = async () => {
    const [success] = await operator(() => API.updateConnection(initialID, { name }), {
      setOperating: setSaving,
    });

    if (success) {
      onSubmitAfter?.(initialID);
    }
  };

  return useMemo(
    () => ({
      saving,
      name,
      record,
      onChangeName: setName,
      onSubmit: handleUpdate,
    }),
    [saving, name, record],
  );
};
