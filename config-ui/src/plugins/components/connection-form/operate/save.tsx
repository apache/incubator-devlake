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

import { useState, useMemo } from 'react';
import { useHistory } from 'react-router-dom';
import { Button, Intent } from '@blueprintjs/core';

import { operator } from '@/utils';

import * as API from '../api';

interface Props {
  plugin: string;
  connectionId?: ID;
  values: any;
  errors: any;
}

export const Save = ({ plugin, connectionId, values, errors }: Props) => {
  const [saving, setSaving] = useState(false);
  const history = useHistory();

  const handleSubmit = async () => {
    const [success] = await operator(
      () => (!connectionId ? API.createConnection(plugin, values) : API.updateConnection(plugin, connectionId, values)),
      {
        setOperating: setSaving,
      },
    );

    if (success) {
      history.push(`/connections/${plugin}`);
    }
  };

  const disabled = useMemo(() => {
    return Object.values(errors).some((value) => value);
  }, [errors]);

  return (
    <Button
      loading={saving}
      disabled={disabled}
      intent={Intent.PRIMARY}
      outlined
      text="Save Connection"
      onClick={handleSubmit}
    />
  );
};
