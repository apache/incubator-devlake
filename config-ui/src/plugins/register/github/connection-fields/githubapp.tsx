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
import { Select, Input } from 'antd';

import API from '@/api';
import { Block, ExternalLink } from '@/components';

import * as S from './styled';

interface Props {
  endpoint?: string;
  proxy: string;
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
      const res = await API.connection.testOld('github', {
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
    <Block
      title="Github App settings"
      description={
        <>
          Input information about your Github App{' '}
          <ExternalLink link="https://docs.github.com/en/apps/maintaining-github-apps/modifying-a-github-app-registration#navigating-to-your-github-app-settings">
            Learn how to create a github app
          </ExternalLink>
        </>
      }
      required
    >
      <S.Input>
        <div className="input">
          <Input
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
          <Input.TextArea
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
        <Select
          style={{ width: 200 }}
          placeholder="Select App installation"
          options={
            settings.installations
              ? settings.installations.map((it) => ({
                  value: it.id,
                  label: it.account.login,
                }))
              : []
          }
          onChange={(value) => setSettings({ ...settings, installationId: value })}
        />
      </S.Input>
    </Block>
  );
};
