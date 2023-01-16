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

import React from 'react';
import { InputGroup } from '@blueprintjs/core';

interface Props {
  placeholder?: string;
  initialValue?: string;
  value?: string;
  onChange?: (value: string) => void;
}

export const GitLabToken = ({ placeholder, initialValue, value, onChange }: Props) => {
  const handleChangeValue = (e: React.ChangeEvent<HTMLInputElement>) => {
    onChange?.(e.target.value);
  };

  return (
    <div>
      <p>
        <a
          href="https://devlake.apache.org/docs/UserManuals/ConfigUI/GitLab/#auth-tokens"
          target="_blank"
          rel="noreferrer"
        >
          Learn about how to create a personal access token
        </a>
      </p>
      <InputGroup
        placeholder={placeholder}
        type="password"
        value={value ?? initialValue}
        onChange={handleChangeValue}
      />
    </div>
  );
};
