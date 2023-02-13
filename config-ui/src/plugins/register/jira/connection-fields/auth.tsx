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

import React, { useState } from 'react';
import { FormGroup, RadioGroup, Radio, InputGroup } from '@blueprintjs/core';

import { ExternalLink } from '@/components';

import * as S from './styled';

type Method = 'BasicAuth' | 'AccessToken';

interface Props {
  values: any;
  setValues: (value: any) => void;
}

export const Auth = ({ values, setValues }: Props) => {
  const [method, setMethod] = useState<Method>('BasicAuth');

  const handleChangeMethod = (e: React.FormEvent<HTMLInputElement>) => {
    const m = (e.target as HTMLInputElement).value as Method;

    setMethod(m);
    setValues({
      ...values,
      authMethod: m,
      username: m === 'BasicAuth' ? values.username : undefined,
      password: m === 'BasicAuth' ? values.password : undefined,
      token: m === 'AccessToken' ? values.token : undefined,
    });
  };

  const handleChangeUsername = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValues({
      ...values,
      authMethod: 'BasicAuth',
      username: e.target.value,
    });
  };

  const handleChangePassword = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValues({
      ...values,
      authMethod: 'BasicAuth',
      password: e.target.value,
    });
  };

  const handleChangeToken = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValues({
      ...values,
      token: e.target.value,
    });
  };

  return (
    <FormGroup label={<S.Label>Authentication Method</S.Label>} labelInfo={<S.LabelInfo>*</S.LabelInfo>}>
      <RadioGroup inline selectedValue={method} onChange={handleChangeMethod}>
        <Radio value="BasicAuth">Basic Authentication</Radio>
        <Radio value="AccessToken">Using Personal Access Token</Radio>
      </RadioGroup>
      {method === 'BasicAuth' && (
        <>
          <FormGroup label={<S.Label>Username/e-mail</S.Label>} labelInfo={<S.LabelInfo>*</S.LabelInfo>}>
            <InputGroup
              placeholder="Your Username/e-mail"
              value={values.username || ''}
              onChange={handleChangeUsername}
            />
          </FormGroup>
          <FormGroup
            label={<S.Label>Password</S.Label>}
            labelInfo={<S.LabelInfo>*</S.LabelInfo>}
            subLabel={
              <S.LabelDescription>
                For Jira Cloud, please enter your{' '}
                <ExternalLink link="https://confluence.atlassian.com/enterprise/using-personal-access-tokens-1026032365.html">
                  Personal Access Token
                </ExternalLink>{' '}
                For Jira Server v8+, please enter the password of your Jira account.
              </S.LabelDescription>
            }
          >
            <InputGroup
              type="password"
              placeholder="Your Token/Password"
              value={values.password || ''}
              onChange={handleChangePassword}
            />
          </FormGroup>
        </>
      )}
      {method === 'AccessToken' && (
        <FormGroup
          label={<S.Label>Personal Access Token</S.Label>}
          labelInfo={<S.LabelInfo>*</S.LabelInfo>}
          subLabel={
            <S.LabelDescription>
              <ExternalLink link="https://confluence.atlassian.com/enterprise/using-personal-access-tokens-1026032365.html">
                Learn about how to create PAT
              </ExternalLink>
            </S.LabelDescription>
          }
        >
          <InputGroup type="password" placeholder="Your PAT" value={values.token || ''} onChange={handleChangeToken} />
        </FormGroup>
      )}
    </FormGroup>
  );
};
