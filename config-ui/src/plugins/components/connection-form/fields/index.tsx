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
  initialValues: any;
  values: any;
  errors: any;
  setValues: (values: any) => void;
  setErrors: (errors: any) => void;
}

export const Form = ({ name, fields, initialValues, values, errors, setValues, setErrors }: Props) => {
  const onValues = (values: any) => setValues((prev: any) => ({ ...prev, ...values }));
  const onErrors = (values: any) => setErrors((prev: any) => ({ ...prev, ...values }));

  const generateForm = () => {
    return fields.map((field) => {
      if (typeof field === 'function') {
        return field({
          initialValues,
          values,
          setValues: onValues,
          errors,
          setErrors: onErrors,
        });
      }

      const key = typeof field === 'string' ? field : field.key;

      switch (key) {
        case 'name':
          return (
            <ConnectionName
              key={key}
              initialValue={initialValues.name ?? ''}
              value={values.name ?? ''}
              error={errors.name ?? ''}
              setValue={(value) => onValues({ name: value })}
              setError={(value) => onErrors({ name: value })}
            />
          );
        case 'endpoint':
          return (
            <ConnectionEndpoint
              {...field}
              key={key}
              name={name}
              initialValue={initialValues.endpoint ?? ''}
              value={values.endpoint ?? ''}
              error={errors.endpoint ?? ''}
              setValue={(value) => onValues({ endpoint: value })}
              setError={(value) => onErrors({ endpoint: value })}
            />
          );
        case 'username':
          return (
            <ConnectionUsername
              key={key}
              initialValue={initialValues.username ?? ''}
              value={values.username ?? ''}
              setValue={(value) => onValues({ username: value })}
            />
          );
        case 'password':
          return (
            <ConnectionPassword
              {...field}
              key={key}
              initialValue={initialValues.password ?? ''}
              value={values.password ?? ''}
              setValue={(value) => onValues({ password: value })}
            />
          );
        case 'token':
          return (
            <ConnectionToken
              {...field}
              key={key}
              initialValue={initialValues.token ?? ''}
              value={values.token ?? ''}
              setValue={(value) => onValues({ token: value })}
            />
          );
        case 'proxy':
          return (
            <ConnectionProxy
              key={key}
              name={name}
              initialValue={initialValues.proxy ?? ''}
              value={values.proxy ?? ''}
              setValue={(value) => onValues({ proxy: value })}
            />
          );
        case 'rateLimitPerHour':
          return (
            <ConnectionRateLimit
              {...field}
              key={key}
              initialValue={initialValues.rateLimitPerHour ?? 0}
              value={values.rateLimitPerHour}
              setValue={(value) => onValues({ rateLimitPerHour: value })}
            />
          );
        default:
          return null;
      }
    });
  };

  return <div>{generateForm()}</div>;
};
