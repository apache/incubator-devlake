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

import { ConnectionName } from './name';
import { ConnectionEndpoint } from './endpoint';
import { ConnectionUsername } from './username';
import { ConnectionPassword } from './password';
import { ConnectionToken } from './token';
import { ConnectionProxy } from './proxy';
import { ConnectionRateLimit } from './rate-limit';

interface Props {
  name: string;
  fields: any[];
  values: any;
  setValues: (values: any) => void;
  error: any;
  setError: (error: any) => void;
}

export const Form = ({ name, fields, values, setValues, error, setError }: Props) => {
  const generateForm = () => {
    return fields.map((field) => {
      if (typeof field === 'function') {
        return field({ values, setValues, error, setError });
      }

      const key = typeof field === 'string' ? field : field.key;

      switch (key) {
        case 'name':
          return (
            <ConnectionName key={key} value={values.name ?? ''} onChange={(name) => setValues({ ...values, name })} />
          );
        case 'endpoint':
          return (
            <ConnectionEndpoint
              {...field}
              key={key}
              name={name}
              value={values.endpoint ?? ''}
              onChange={(endpoint) => setValues({ ...values, endpoint })}
            />
          );
        case 'username':
          return (
            <ConnectionUsername
              key={key}
              value={values.username ?? ''}
              onChange={(username) => setValues({ ...values, username })}
            />
          );
        case 'password':
          return (
            <ConnectionPassword
              {...field}
              key={key}
              value={values.password ?? ''}
              onChange={(password) => setValues({ ...values, password })}
            />
          );
        case 'token':
          return (
            <ConnectionToken
              {...field}
              key={key}
              value={values.token ?? ''}
              onChange={(token) => setValues({ ...values, token })}
            />
          );
        case 'proxy':
          return (
            <ConnectionProxy
              key={key}
              name={name}
              value={values.proxy ?? ''}
              onChange={(proxy) => setValues({ ...values, proxy })}
            />
          );
        case 'rateLimitPerHour':
          return (
            <ConnectionRateLimit
              {...field}
              key={key}
              value={values.rateLimitPerHour}
              onChange={(rateLimitPerHour) => setValues({ ...values, rateLimitPerHour })}
            />
          );
        default:
          return null;
      }
    });
  };

  return <div>{generateForm()}</div>;
};
