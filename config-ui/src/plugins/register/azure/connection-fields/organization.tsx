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
import { Input, Radio, type RadioChangeEvent } from 'antd';

import { Block, ExternalLink } from '@/components';
import { DOC_URL } from '@/release';

interface Props {
  initialValue: OrganizationSettings;
  value: string;
  label?: string;
  setValue: (value: string) => void;
}

interface OrganizationSettings {
  organization: string;
  scoped: boolean;
}

export const ConnectionOrganization = ({ label, initialValue, value, setValue }: Props) => {
  const [settings, setSettings] = useState<OrganizationSettings>({ scoped: false, organization: '' });

  useEffect(() => {
    const org = initialValue.organization || '';
    setValue(org);

    setSettings({ organization: initialValue.organization, scoped: org !== '' });
  }, [initialValue.organization]);

  const handleChange = (e: RadioChangeEvent) => {
    const scoped = e.target.value;
    if (scoped) {
      setValue(settings.organization);
    } else {
      setValue('');
    }
    setSettings({ ...settings, scoped });
  };

  const handleChangeValue = (e: React.ChangeEvent<HTMLInputElement>) => {
    const organization = e.target.value;
    setValue(organization);
    setSettings({ ...settings, organization });
  };

  return (
    <>
      <Block title={label || 'Personal Access Token Scope'}>
        <p>
          If you are using an organization-scoped token, please enter the organization. Otherwise make sure to create an
          unscoped token.{' '}
          {DOC_URL.PLUGIN.AZUREDEVOPS.AUTH_TOKEN !== '' && (
            <ExternalLink link={DOC_URL.PLUGIN.AZUREDEVOPS.AUTH_TOKEN}>Learn about how to create a PAT</ExternalLink>
          )}
        </p>
        <Radio.Group value={settings.scoped} onChange={handleChange}>
          <Radio value={false}>Unscoped</Radio>
          <Radio value={true}>Scoped</Radio>
        </Radio.Group>
      </Block>
      <Block>
        <Input
          style={{ width: 386 }}
          placeholder="Your organization"
          value={value}
          onChange={handleChangeValue}
          disabled={!settings.scoped}
        />
      </Block>
    </>
  );
};
