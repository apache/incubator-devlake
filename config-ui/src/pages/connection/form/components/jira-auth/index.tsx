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

type Method = 'BasicAuth' | 'AccessToken';

type Value = {
  authMethod: string;
  username?: string;
  password?: string;
  token?: string;
};

interface Props {
  value: Value;
  onChange: (value: Value) => void;
}

export const JIRAAuth = ({ value, onChange }: Props) => {
  const [method, setMethod] = useState<Method>('BasicAuth');

  const handleChangeMethod = (e: React.FormEvent<HTMLInputElement>) => {
    const m = (e.target as HTMLInputElement).value as Method;

    setMethod(m);
    onChange({
      authMethod: m,
      username: m === 'BasicAuth' ? value.username : undefined,
      password: m === 'BasicAuth' ? value.password : undefined,
      token: m === 'AccessToken' ? value.token : undefined,
    });
  };

  const handleChangeUsername = (e: React.ChangeEvent<HTMLInputElement>) => {
    onChange({
      authMethod: 'BasicAuth',
      username: e.target.value,
    });
  };

  const handleChangePassword = (e: React.ChangeEvent<HTMLInputElement>) => {
    onChange({
      authMethod: 'BasicAuth',
      password: e.target.value,
    });
  };

  const handleChangeToken = (e: React.ChangeEvent<HTMLInputElement>) => {
    onChange({
      ...value,
      token: e.target.value,
    });
  };

  return (
    <div>
      <FormGroup inline label="Authentication Method" labelInfo="*">
        <RadioGroup selectedValue={method} onChange={handleChangeMethod}>
          <Radio value="BasicAuth">Basic Authentication</Radio>
          <Radio value="AccessToken">Using Personal Access Token</Radio>
        </RadioGroup>
      </FormGroup>
      {method === 'BasicAuth' && (
        <>
          <FormGroup inline label="Username/e-mail" labelInfo="*">
            <InputGroup
              placeholder="Your Username/e-mail"
              value={value.username || ''}
              onChange={handleChangeUsername}
            />
          </FormGroup>
          <FormGroup inline label="Password" labelInfo="*">
            <InputGroup
              type="password"
              placeholder="Your Password"
              value={value.password || ''}
              onChange={handleChangePassword}
            />
          </FormGroup>
        </>
      )}
      {method === 'AccessToken' && (
        <FormGroup inline label="Personal Access Token" labelInfo="*">
          <p>
            <a
              href="https://confluence.atlassian.com/enterprise/using-personal-access-tokens-1026032365.html"
              target="_blank"
              rel="noreferrer"
            >
              Learn about how to create PAT
            </a>
          </p>
          <InputGroup type="password" placeholder="Your PAT" value={value.token || ''} onChange={handleChangeToken} />
        </FormGroup>
      )}
    </div>
  );
};
