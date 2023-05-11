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
import { useHistory } from 'react-router-dom';
import { FormGroup, InputGroup, Button, Intent } from '@blueprintjs/core';

import { operator } from '@/utils';

import * as API from './api';
import * as S from './styld';

const NEW_PASSWORD_REQUIRED = 'NEW_PASSWORD_REQUIRED';

export const LoginPage = () => {
  const [username, setUsername] = useState(localStorage.getItem('username') || '');
  const [password, setPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [challenge, setChallenge] = useState('');
  const [session, setSession] = useState('');

  const history = useHistory();

  // () =>
  const handleSubmit = async () => {
    var request: () => Promise<any>;

    switch (challenge) {
      case NEW_PASSWORD_REQUIRED:
        request = () => API.newPassword({ username, session, newPassword });
        break;
      default:
        request = () => API.login({ username, password });
        break;
    }

    const [success, res] = await operator(request, {
      formatReason: (error) => {
        const e = error as any;
        return e?.response?.data?.causes[0];
      },
    });
    localStorage.setItem('username', username);
    if (success) {
      if (res.challengeName) {
        setChallenge(res.challengeName);
        setSession(res.session);
      } else {
        localStorage.setItem('accessToken', res.authenticationResult.accessToken);
        document.cookie = 'access_token=' + res.authenticationResult.accessToken + '; path=/';
        setUsername('');
        setPassword('');
        setChallenge('');
        setSession('');
        history.push('/');
      }
    }
  };

  return (
    <S.Wrapper>
      <S.Inner>
        <h2>DevLake Login</h2>
        <FormGroup label="Username">
          <InputGroup
            placeholder="Username"
            value={username}
            disabled={challenge !== ''}
            onChange={(e) => setUsername((e.target as HTMLInputElement).value)}
          />
        </FormGroup>
        <FormGroup label="Password">
          <InputGroup
            type="password"
            placeholder="Password"
            value={password}
            disabled={challenge !== ''}
            onChange={(e) => setPassword((e.target as HTMLInputElement).value)}
          />
        </FormGroup>
        {challenge === 'NEW_PASSWORD_REQUIRED' && (
          <FormGroup label="Set New Password">
            <InputGroup
              type="password"
              placeholder="Please set a new Password for your account"
              value={newPassword}
              onChange={(e) => setNewPassword((e.target as HTMLInputElement).value)}
            />
          </FormGroup>
        )}
        <Button intent={Intent.PRIMARY} onClick={handleSubmit}>
          Login
        </Button>
      </S.Inner>
    </S.Wrapper>
  );
};
