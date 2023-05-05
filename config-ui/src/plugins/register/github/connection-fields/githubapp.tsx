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
import { FormGroup, InputGroup, TextArea, Button, Intent } from '@blueprintjs/core';

import { ExternalLink } from '@/components';

import * as API from '../api';

import * as S from './styled';

interface Props {
  endpoint?: string;
  proxy?: string;
  initialValue: any;
  value: any;
  error: string;
  setValue: (value: any) => void;
  setError: (error: any) => void;
}

interface GithubAppSettings {
  appId?: string;
  secretKey?: string;

  status: 'idle' | 'valid' | 'invalid';
  from?: string;
}

export const GithubApp = ({ endpoint, proxy, initialValue, value, error, setValue, setError }: Props) => {
  const [settings, setSettings] = useState<GithubAppSettings>({ status: 'idle' });

  useEffect(() => {
    setError({
      appId: value.appId ? '' : 'AppId is required',
      secretKey: value.secretKey ? '' : 'SecretKey is required',
    });

    return () => {
      setError({
        appId: '',
        secretKey: '',
      });
    }
  }, [value.appId, value.secretKey]);

  const testConfiguration = async (appId?: string, secretKey?: string): Promise<GithubAppSettings> => {
    if (!endpoint || !appId || !secretKey) {
      return {
        appId,
        secretKey,
        status: 'idle',
      };
    }

    try {
      const res = await API.testConnection({
        authMethod: 'AppKey',
        endpoint,
        proxy,
        appId,
        secretKey,
        token: '',
      });
      return {
        appId,
        secretKey,
        status: 'valid',
        from: res.login,
      };
    } catch {
      return {
        appId,
        secretKey,
        status: 'invalid',
      };
    }
  };

  const handleChangeAppId = (value: string) => {
    setSettings({ ...settings, appId: value });;
  };

  const handleChangeClientSecret = (value: string) => {
    setSettings({ ...settings, secretKey: value });
  };

  const handleTestConfiguration = async () => {
    const res = await testConfiguration(settings.appId, settings.secretKey);
    setSettings(res);
  };

  const checkConfig = async (appId: string, secretKey: string) => {
    const res = await testConfiguration(appId, secretKey);
    setSettings(res);
  };


  useEffect(() => {
    checkConfig(initialValue.appId, initialValue.secretKey);
  }, [initialValue.appId, initialValue.secretKey, endpoint]);

  useEffect(() => {
    setValue({ appId: settings.appId, secretKey: settings.secretKey });
  }, [settings.appId, settings.secretKey]);


  return (
    <FormGroup
      label={<S.Label>Github App settings</S.Label>}
      labelInfo={<S.LabelInfo>*</S.LabelInfo>}
      subLabel={
        <S.LabelDescription>
          Input information about your Github App{' '}
          <ExternalLink link="https://TODO">
            Learn how to create a github app
          </ExternalLink>
        </S.LabelDescription>
      }
    >
      <S.Input>
        <InputGroup
          placeholder="App Id"
          type="text"
          value={settings.appId ?? ''}
          onChange={(e) => handleChangeAppId(e.target.value)}
          onBlur={() => handleTestConfiguration()}
        />
        <div className="info">
          {settings.status === 'invalid' && <span className="error">Invalid</span>}
          {settings.status === 'valid' && <span className="success">Valid From: {settings.from}</span>}
        </div>
      </S.Input>
      <S.Input>
        <TextArea
          cols={90}
          rows={15}
          placeholder="Private key"
          value={settings.secretKey ?? ''}
          onChange={(e) => handleChangeClientSecret(e.target.value)}
          onBlur={() => handleTestConfiguration()}
        />
        <div className="info">
          {settings.status === 'invalid' && <span className="error">Invalid</span>}
          {settings.status === 'valid' && <span className="success">Valid From: {settings.from}</span>}
        </div>
      </S.Input>
    </FormGroup>
  );
};
