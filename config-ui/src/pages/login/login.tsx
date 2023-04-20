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

export const LoginPage = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  const history = useHistory();

  const handleSubmit = async () => {
    const [success, res] = await operator(() => API.login({ username, password }), {
      formatReason: (error) => 'Login failed',
    });

    if (success) {
      localStorage.setItem('accessToken', res.AuthenticationResult.AccessToken);
      document.cookie = 'access_token=' + res.AuthenticationResult.AccessToken + '; path=/';
      history.push('/');
    }

    setUsername('');
    setPassword('');
  };

  return (
    <S.Wrapper>
      <S.Inner>
        <h2>DevLake Login</h2>
        <FormGroup label="Username">
          <InputGroup
            placeholder="Username"
            value={username}
            onChange={(e) => setUsername((e.target as HTMLInputElement).value)}
          />
        </FormGroup>
        <FormGroup label="Password">
          <InputGroup
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword((e.target as HTMLInputElement).value)}
          />
        </FormGroup>
        <Button intent={Intent.PRIMARY} onClick={handleSubmit}>
          Login
        </Button>
      </S.Inner>
    </S.Wrapper>
  );
};
