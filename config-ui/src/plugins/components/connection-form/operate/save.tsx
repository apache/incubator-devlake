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

import React, { useMemo } from 'react';
import { useHistory } from 'react-router-dom';
import { Button, Intent } from '@blueprintjs/core';

import { useOperator } from '@/hooks';

import * as API from '../api';

interface Props {
  plugin: string;
  connectionId?: ID;
  form: any;
  error: any;
}

export const Save = ({ plugin, connectionId, form, error }: Props) => {
  const history = useHistory();

  const { operating, onSubmit } = useOperator(
    (paylaod) =>
      !connectionId ? API.createConnection(plugin, paylaod) : API.updateConnection(plugin, connectionId, paylaod),
    {
      callback: () => history.push(`/connections/${plugin}`),
    },
  );

  const disabled = useMemo(() => {
    return Object.values(error).some((value) => value);
  }, [error]);

  const handleSubmit = () => {
    onSubmit(form);
  };

  return (
    <Button
      loading={operating}
      disabled={disabled}
      intent={Intent.PRIMARY}
      outlined
      text="Save Connection"
      onClick={handleSubmit}
    />
  );
};
