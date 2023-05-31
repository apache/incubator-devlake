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
import { Button, FormGroup, InputGroup, MenuItem, TextArea } from '@blueprintjs/core';
import { Select2 } from '@blueprintjs/select';

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
  installationId?: number;

  status: 'idle' | 'valid' | 'invalid';
  from?: string;
  installations?: GithubInstallation[];
}

interface GithubInstallation {
  id: number;
  account: {
    login: string;
  };
}

export const GithubApp = ({ endpoint, proxy, initialValue, value, error, setValue, setError }: Props) => {
  const [settings, setSettings] = useState<GithubAppSettings>({ status: 'idle' });

  useEffect(() => {
    setError({
      appId: value.appId ? '' : 'AppId is required',
      secretKey: value.secretKey ? '' : 'SecretKey is required',
      installationId: value.installationId ? '' : 'InstallationId is required',
    });

    return () => {
      setError({
        appId: '',
        secretKey: '',
        installationId: '',
      });
    };
  }, [value.appId, value.secretKey, value.installationId]);

  const testConfiguration = async (
    appId?: string,
    secretKey?: string,
    installationId?: number,
  ): Promise<GithubAppSettings> => {
    if (!endpoint || !appId || !secretKey) {
      return {
        appId,
        secretKey,
        installationId,
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
        installationId,
        status: 'valid',
        from: res.login,
        installations: res.installations,
      };
    } catch {
      return {
        appId,
        secretKey,
        installationId,
        status: 'invalid',
      };
    }
  };

  const handleChangeAppId = (value: string) => {
    setSettings({ ...settings, appId: value });
  };

  const handleChangeClientSecret = (value: string) => {
    setSettings({ ...settings, secretKey: value });
  };

  const handleTestConfiguration = async () => {
    const res = await testConfiguration(settings.appId, settings.secretKey, settings.installationId);
    setSettings(res);
  };

  const checkConfig = async (appId: string, secretKey: string, installationId: number) => {
    const res = await testConfiguration(appId, secretKey, installationId);
    setSettings(res);
  };

  useEffect(() => {
    checkConfig(initialValue.appId, initialValue.secretKey, initialValue.installationId);
  }, [initialValue.appId, initialValue.secretKey, initialValue.installationId, endpoint]);

  useEffect(() => {
    setValue({ appId: settings.appId, secretKey: settings.secretKey, installationId: settings.installationId });
  }, [settings.appId, settings.secretKey, settings.installationId]);

  return (
    <FormGroup
      label={<S.Label>Github App settings</S.Label>}
      labelInfo={<S.LabelInfo>*</S.LabelInfo>}
      subLabel={
        <S.LabelDescription>
          Input information about your Github App{' '}
          <ExternalLink link="https://TODO">Learn how to create a github app</ExternalLink>
        </S.LabelDescription>
      }
    >
      <S.Input>
        <div className="input">
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
        </div>
      </S.Input>
      <S.Input>
        <div className="input">
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
        </div>
      </S.Input>
      <S.Input>
        <Select2
          items={settings.installations ?? []}
          activeItem={settings.installations?.find((e) => e.id === settings.installationId)}
          itemPredicate={(query, item) => item.account.login.toLowerCase().includes(query.toLowerCase())}
          itemRenderer={(item, { handleClick, handleFocus, modifiers }) => {
            return (
              <MenuItem
                active={modifiers.active}
                disabled={modifiers.disabled}
                key={item.id}
                label={item.id.toString()}
                onClick={handleClick}
                onFocus={handleFocus}
                roleStructure="listoption"
                text={item.account.login}
              />
            );
          }}
          onItemSelect={(item) => {
            setSettings({ ...settings, installationId: item.id });
          }}
          noResults={<option disabled={true}>No results</option>}
          popoverProps={{ minimal: true }}
        >
          <Button
            text={
              settings.installations?.find((e) => e.id === settings.installationId)?.account.login ??
              'Select App installation'
            }
            rightIcon="double-caret-vertical"
            placeholder="Select App installation"
          />
        </Select2>
      </S.Input>
    </FormGroup>
  );
};
