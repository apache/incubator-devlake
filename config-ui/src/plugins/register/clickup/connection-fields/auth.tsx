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

import React, { useState, useEffect } from 'react';
import { FormGroup, RadioGroup, Radio, InputGroup } from '@blueprintjs/core';

import { ExternalLink } from '@/components';

import * as S from './styled';

interface Props {
  initialValues: any;
  values: any;
  errors: any;
  setValues: (value: any) => void;
  setErrors: (value: any) => void;
}

export const Auth = ({ initialValues, values, setValues, setErrors }: Props) => {
  useEffect(() => {
    setValues({
      token: initialValues.token,
      teamId: initialValues.teamId,
    });
  }, [
    initialValues.token,
    initialValues.teamId ? '' : 'https://api.clickup.com/api/',
  ]);

  useEffect(() => {
    setErrors({
      token: values.token ? '' : 'token is required',
      teamId: values.teamId ? '' : 'teamId is required',
    });
  }, [values]);

  const handleChangeToken = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValues({
      token: e.target.value,
    });
  };

  const handleChangeEndpoint = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValues({
      endpoint: e.target.value,
    });
  };
  const handleTeamIdChanged = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValues({
      teamId: e.target.value,
    });
  };

  return (
    <>
        <FormGroup
          style={{ marginTop: 8, marginBottom: 0 }}
          label={<S.Label>Endpoint URL</S.Label>}
          labelInfo={<S.LabelInfo>*</S.LabelInfo>}
          subLabel={
            <S.LabelDescription>
                Provide the ClickUp instance API endpoint. e.g. https://api.clickup.com/api/
            </S.LabelDescription>
          }
        >
          <InputGroup placeholder="Your Endpoint URL" value={values.endpoint} onChange={handleChangeEndpoint} />
        </FormGroup>
        <FormGroup
          style={{ marginTop: 8, marginBottom: 0 }}
          label={<S.Label>Team ID</S.Label>}
          labelInfo={<S.LabelInfo>*</S.LabelInfo>}
          subLabel={
            <S.LabelDescription>
                Provide the ClickUp TeamUp/Workspace ID
            </S.LabelDescription>
          }
        >
          <InputGroup placeholder="Team ID" value={values.teamId} onChange={handleTeamIdChanged} />
        </FormGroup>
        <FormGroup
          label={<S.Label>Personal Access Token</S.Label>}
          labelInfo={<S.LabelInfo>*</S.LabelInfo>}
          subLabel={
            <S.LabelDescription>
              <ExternalLink link="https://devlake.apache.org/docs/Configuration/Jira#personal-access-token">
                Learn about how to create a PAT
              </ExternalLink>
            </S.LabelDescription>
          }
        >
          <InputGroup type="password" placeholder="Your PAT" value={values.token} onChange={handleChangeToken} />
        </FormGroup>
    </>
  );
};
