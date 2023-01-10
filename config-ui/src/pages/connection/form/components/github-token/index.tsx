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

import React, { useEffect, useState } from 'react';
import { InputGroup, Button, Intent } from '@blueprintjs/core';
import { pick } from 'lodash';

import { Plugins } from '@/plugins';

import * as API from '../../api';

import * as S from './styled';

type TokenItem = {
  value: string;
  status: 'idle' | 'valid' | 'invalid';
  from?: string;
};

interface Props {
  form: any;
  value?: string;
  onChange?: (value: string) => void;
}

export const GitHubToken = ({ form, value, onChange }: Props) => {
  const [tokens, setTokens] = useState<TokenItem[]>([{ value: '', status: 'idle' }]);

  const testToken = async (token: string): Promise<TokenItem> => {
    try {
      const res = await API.testConnection(Plugins.GitHub, {
        ...pick(form, ['endpoint', 'proxy']),
        token,
      });
      return {
        value: token,
        status: 'valid',
        from: res.login,
      };
    } catch {
      return {
        value: token,
        status: 'invalid',
      };
    }
  };

  const checkTokens = async (value: string) => {
    const res = await Promise.all((value ?? '').split(',').map((it) => testToken(it)));
    setTokens(res);
  };

  useEffect(() => {
    if (value) {
      checkTokens(value);
    }
  }, []);

  useEffect(() => {
    onChange?.(tokens.map((it) => it.value).join(','));
  }, [tokens]);

  const handleCreateToken = () => {
    setTokens([...tokens, { value: '', status: 'idle' }]);
  };

  const handleRemoveToken = (key: number) => {
    setTokens(tokens.filter((_, i) => (i === key ? false : true)));
  };

  const handleChangeToken = (key: number, value: string) => {
    setTokens(tokens.map((it, i) => (i === key ? { value, status: 'idle' } : it)));
  };

  const handleTestToken = async (key: number) => {
    const token = tokens.find((_, i) => i === key) as TokenItem;

    if (token.status === 'idle' && token.value) {
      const res = await testToken(token.value);
      setTokens((tokens) => tokens.map((it, i) => (i === key ? res : it)));
    }
  };

  return (
    <S.Wrapper>
      <p>
        Add one or more personal token(s) for authentication from you and your organization members. Multiple tokens can
        help speed up the data collection process.{' '}
      </p>
      <p>
        <a
          href="https://devlake.apache.org/docs/UserManuals/ConfigUI/GitHub/#auth-tokens"
          target="_blank"
          rel="noreferrer"
        >
          Learn about how to create a personal access token
        </a>
      </p>
      <h3>Personal Access Token(s)</h3>
      {tokens.map(({ value, status, from }, i) => (
        <div className="token" key={i}>
          <div className="input">
            <InputGroup
              placeholder="token"
              type="password"
              value={value ?? ''}
              onChange={(e) => handleChangeToken(i, e.target.value)}
              onBlur={() => handleTestToken(i)}
            />
            {status === 'invalid' && <span className="error">Invalid</span>}
            {status === 'valid' && <span className="success">Valid From: {from}</span>}
          </div>
          <Button minimal icon="cross" onClick={() => handleRemoveToken(i)} />
        </div>
      ))}
      <div className="action">
        <Button outlined small intent={Intent.PRIMARY} text="Another Token" icon="plus" onClick={handleCreateToken} />
      </div>
    </S.Wrapper>
  );
};
