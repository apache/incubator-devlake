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

import { useEffect, useState } from 'react';
import { FormGroup, Button, Icon, Intent } from '@blueprintjs/core';

import { ExternalLink, FormPassword } from '@/components';
import { DOC_URL } from '@/release';

import * as API from '../api';

import * as S from './styled';

type TokenItem = {
  value: string;
  isValid?: boolean;
  from?: string;
  status?: 'success' | 'warning' | 'error';
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
  const [tokens, setTokens] = useState<TokenItem[]>([{ value: '' }]);

  const testToken = async (token: string): Promise<TokenItem> => {
    if (!endpoint || !token) {
      return {
        value: token,
      };
    }

    try {
      const res = await API.testConnection({
        authMethod: 'AccessToken',
        endpoint,
        proxy,
        token,
      });
      return {
        value: token,
        isValid: true,
        from: res.login,
        status: !res.success ? 'error' : res.warning ? 'warning' : 'success',
      };
    } catch {
      return {
        value: token,
        isValid: false,
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

    return () => {
      setError('');
    };
  }, [value]);

  useEffect(() => {
    setValue(tokens.map((it) => it.value).join(','));
  }, [tokens]);

  const handleCreateToken = () => setTokens([...tokens, { value: '' }]);

  const handleRemoveToken = (key: number) => setTokens(tokens.filter((_, i) => (i === key ? false : true)));

  const handleChangeToken = (key: number, value: string) =>
    setTokens(tokens.map((it, i) => (i === key ? { value } : it)));

  const handleTestToken = async (key: number) => {
    const token = tokens.find((_, i) => i === key) as TokenItem;
    if (token.isValid === undefined && token.value) {
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
          <ExternalLink link={DOC_URL.PLUGIN.GITHUB.AUTH_TOKEN}>
            Learn how to create a personal access token
          </ExternalLink>
        </S.LabelDescription>
      }
    >
      {tokens.map(({ value, isValid, status, from }, i) => (
        <S.Input key={i}>
          <div className="input">
            <FormPassword
              placeholder="Token"
              onChange={(e) => handleChangeToken(i, e.target.value)}
              onBlur={() => handleTestToken(i)}
            />
            <Button minimal icon="cross" onClick={() => handleRemoveToken(i)} />
            <div className="info">
              {isValid === false && <span className="error">Invalid</span>}
              {isValid === true && <span className="success">Valid From: {from}</span>}
            </div>
          </div>
          {status && (
            <S.Alert>
              <h4>
                {status === 'success' && <Icon icon="tick-circle" color="#4DB764" />}
                {status === 'warning' && <Icon icon="warning-sign" color="#F4BE55" />}
                {status === 'error' && <Icon icon="cross-circle" color="#E34040" />}
                <span style={{ marginLeft: 8 }}>Token Permissions</span>
              </h4>
              {status === 'success' && <p>All required fields are checked.</p>}
              {status === 'warning' && (
                <p>
                  This token is able to collect public repositories. If you want to collect private repositories, please
                  check the field `repo`.
                </p>
              )}
              {status === 'error' && (
                <>
                  <p>Please check the field(s) `repo:status`, `repo_deployment`, `read:user`, `read:org`.</p>
                  <p>If you want to collect private repositories, please check the field `repo`.</p>
                </>
              )}
            </S.Alert>
          )}
        </S.Input>
      ))}
      <div className="action">
        <Button outlined small intent={Intent.PRIMARY} text="Another Token" icon="plus" onClick={handleCreateToken} />
      </div>
    </FormGroup>
  );
};
