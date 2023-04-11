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
import { FormGroup, InputGroup, Button, Intent } from '@blueprintjs/core';

import { ExternalLink } from '@/components';

import * as API from '../api';

import * as S from './styled';

type TokenItem = {
  value: string;
  status: 'idle' | 'valid' | 'invalid';
  from?: string;
};

interface Props {
  endpoint?: string;
  proxy?: string;
  initialValue: string;
  value: string;
  error: string;
  setValue: (value: string) => void;
  setError: (error: string) => void;
}

export const Token = ({ endpoint, proxy, initialValue, value, error, setValue, setError }: Props) => {
  const [tokens, setTokens] = useState<TokenItem[]>([{ value: '', status: 'idle' }]);

  const testToken = async (token: string): Promise<TokenItem> => {
    if (!endpoint || !token) {
      return {
        value: token,
        status: 'idle',
      };
    }

    try {
      const res = await API.testConnection({
        endpoint,
        proxy,
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
    checkTokens(initialValue);
  }, [initialValue, endpoint]);

  useEffect(() => {
    setError(value ? '' : 'token is required');
  }, [value]);

  useEffect(() => {
    setValue(tokens.map((it) => it.value).join(','));
  }, [tokens]);

  const handleCreateToken = () => setTokens([...tokens, { value: '', status: 'idle' }]);

  const handleRemoveToken = (key: number) => setTokens(tokens.filter((_, i) => (i === key ? false : true)));

  const handleChangeToken = (key: number, value: string) =>
    setTokens(tokens.map((it, i) => (i === key ? { value, status: 'idle' } : it)));

  const handleTestToken = async (key: number) => {
    const token = tokens.find((_, i) => i === key) as TokenItem;
    if (token.status === 'idle' && token.value) {
      const res = await testToken(token.value);
      setTokens((tokens) => tokens.map((it, i) => (i === key ? res : it)));
    }
  };

  return (
    <FormGroup
      label={<S.Label>Personal Access Token(s) </S.Label>}
      labelInfo={<S.LabelInfo>*</S.LabelInfo>}
      subLabel={
        <S.LabelDescription>
          Add one or more personal token(s) for authentication from you and your organization members. Multiple tokens
          (from different GitHub accounts, NOT from one account) can help speed up the data collection process.{' '}
          <ExternalLink link="https://devlake.apache.org/docs/Configuration/GitHub/#auth-tokens">
            Learn how to create a personal access token
          </ExternalLink>
        </S.LabelDescription>
      }
    >
      {tokens.map(({ value, status, from }, i) => (
        <S.Token key={i}>
          <InputGroup
            placeholder="Token"
            type="password"
            value={value ?? ''}
            onChange={(e) => handleChangeToken(i, e.target.value)}
            onBlur={() => handleTestToken(i)}
          />
          <Button minimal icon="cross" onClick={() => handleRemoveToken(i)} />
          <div className="info">
            {status === 'invalid' && <span className="error">Invalid</span>}
            {status === 'valid' && <span className="success">Valid From: {from}</span>}
          </div>
        </S.Token>
      ))}
      <div className="action">
        <Button outlined small intent={Intent.PRIMARY} text="Another Token" icon="plus" onClick={handleCreateToken} />
      </div>
    </FormGroup>
  );
};
