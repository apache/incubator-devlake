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

import { FormGroup, InputGroup, Radio, RadioGroup } from '@blueprintjs/core';

import { useEffect } from 'react';
import * as S from './styled';

interface Props {
  initialValues: any;
  values: any;
  errors: any;
  setValues: (value: any) => void;
  setErrors: (value: any) => void;
  // this will be the original setValues function, to provide full control to the state
  setValuesDefault: React.Dispatch<React.SetStateAction<Record<string, any>>>;
}

interface AzureAuth {
  providerType: 'azure';
  subscriptionID: string;
  clientID: string;
  clientSecret: string;
  tenantID: string;
  resourceGroupName: string;
  clusterName: string;
}

interface AWSAuth {
  providerType: 'aws';
  accessKeyID: string;
  secretAccessKey: string;
  clusterName: string;
  awsRegion: string;
}

interface OpenShiftAuth {
  providerType: 'openShift';
  authenticationURLForOpenshift: string;
}

type AuthType = AzureAuth['providerType'] | AWSAuth['providerType'] | OpenShiftAuth['providerType'];

export const Auth = ({ initialValues, values, errors, setValues, setErrors, setValuesDefault }: Props) => {
  console.log('initialValues', initialValues);

  useEffect(() => {
    // all fields in the 3 auth types are required, if any fields is empty string or undefined, set error
    if (values.credentials?.providerType === 'azure') {
      if (
        values.subscriptionID === '' ||
        values.clientID === '' ||
        values.clientSecret === '' ||
        values.tenantID === '' ||
        values.resourceGroupName === '' ||
        values.clusterName === ''
      ) {
        setErrors({
          error: 'Required',
        });
        // unset errors
      } else {
        setErrors({ error: '' });
      }
    } else if (values.credentials?.providerType === 'aws') {
      if (
        values.accessKeyID === '' ||
        values.secretAccessKey === '' ||
        values.clusterName === '' ||
        values.awsRegion === ''
      ) {
        setErrors({
          error: 'Required',
        });
      } else {
        setErrors({ error: '' });
      }
    } else if (values.credentials?.providerType === 'openShift') {
      if (values.authenticationURLForOpenshift === '') {
        setErrors({
          error: 'Required',
        });
      } else {
        setErrors({ error: '' });
      }
    }
  }, [values]);

  console.log(errors);

  const defaultValues = {
    authMethod: 'AccessToken',
    endpoint: 'http://127.0.0.1:5002/v1',
    id: 23,
    proxy: '',
    rateLimitPerHour: 0,
  };

  const handleChangeAuthType = (authType: AuthType) => {
    switch (authType) {
      case 'azure':
        setValuesDefault((prev) => ({
          ...defaultValues,
          name: prev.name,

          credentials: {
            providerType: 'azure',
            subscriptionID: '',
            clientID: '',
            clientSecret: '',
            tenantID: '',
            resourceGroupName: '',
            clusterName: '',
          },
        }));
        break;
      case 'aws':
        setValuesDefault((prev) => ({
          ...defaultValues,
          name: prev.name,

          credentials: { providerType: 'aws', accessKeyID: '', clusterName: '', awsRegion: '', secretAccessKey: '' },
        }));
        break;
      case 'openShift':
        setValuesDefault((prev) => ({
          ...defaultValues,
          name: prev.name,

          credentials: { providerType: 'openShift', authenticationURLForOpenshift: '' },
        }));
        break;
    }
  };

  return (
    <>
      <RadioGroup
        inline={true}
        selectedValue={values.credentials?.providerType}
        onChange={(e) => handleChangeAuthType(e.currentTarget.value as AuthType)}
      >
        <Radio label="Azure" value="azure" />
        <Radio label="AWS" value="aws" />
        <Radio label="OpenShift" value="openShift" />
      </RadioGroup>

      {values.credentials?.providerType === 'azure' && (
        <>
          <FormGroup style={{ marginTop: 8, marginBottom: 0 }} label={<S.Label>Subscription ID</S.Label>}>
            <InputGroup
              placeholder="Your Subscription ID"
              value={(values.credentials as AzureAuth).subscriptionID}
              onChange={(e) =>
                setValuesDefault((prev) => ({
                  ...prev,
                  credentials: {
                    ...prev.credentials,
                    subscriptionID: e.target.value,
                  },
                }))
              }
            />
          </FormGroup>
          <FormGroup style={{ marginTop: 8, marginBottom: 0 }} label={<S.Label>Client ID</S.Label>}>
            <InputGroup
              placeholder="Your Client ID"
              value={(values.credentials as AzureAuth).clientID}
              onChange={(e) =>
                setValuesDefault((prev) => ({
                  ...prev,
                  credentials: {
                    ...prev.credentials,
                    clientID: e.target.value,
                  },
                }))
              }
            />
          </FormGroup>
          <FormGroup style={{ marginTop: 8, marginBottom: 0 }} label={<S.Label>Client Secret</S.Label>}>
            <InputGroup
              placeholder="Your Client Secret"
              value={(values.credentials as AzureAuth).clientSecret}
              onChange={(e) =>
                setValuesDefault((prev) => ({
                  ...prev,
                  credentials: {
                    ...prev.credentials,
                    clientSecret: e.target.value,
                  },
                }))
              }
            />
          </FormGroup>
          <FormGroup style={{ marginTop: 8, marginBottom: 0 }} label={<S.Label>Tenant ID</S.Label>}>
            <InputGroup
              placeholder="Your Tenant ID"
              value={(values.credentials as AzureAuth).tenantID}
              onChange={(e) =>
                setValuesDefault((prev) => ({
                  ...prev,
                  credentials: {
                    ...prev.credentials,
                    tenantID: e.target.value,
                  },
                }))
              }
            />
          </FormGroup>
          <FormGroup style={{ marginTop: 8, marginBottom: 0 }} label={<S.Label>Resource Group Name</S.Label>}>
            <InputGroup
              placeholder="Your Resource Group Name"
              value={(values.credentials as AzureAuth).resourceGroupName}
              onChange={(e) =>
                setValuesDefault((prev) => ({
                  ...prev,
                  credentials: {
                    ...prev.credentials,
                    resourceGroupName: e.target.value,
                  },
                }))
              }
            />
          </FormGroup>
          <FormGroup style={{ marginTop: 8, marginBottom: 0 }} label={<S.Label>Cluster name</S.Label>}>
            <InputGroup
              placeholder="Your Cluster name"
              value={(values.credentials as AzureAuth).clusterName}
              onChange={(e) =>
                setValuesDefault((prev) => ({
                  ...prev,
                  credentials: {
                    ...prev.credentials,
                    clusterName: e.target.value,
                  },
                }))
              }
            />
          </FormGroup>
        </>
      )}

      {values.credentials?.providerType === 'aws' && (
        <>
          <FormGroup style={{ marginTop: 8, marginBottom: 0 }} label={<S.Label>Access Key ID</S.Label>}>
            <InputGroup
              placeholder="Your Access Key ID"
              value={(values.credentials as AWSAuth).accessKeyID}
              onChange={(e) =>
                setValuesDefault((prev) => ({
                  ...prev,
                  credentials: {
                    ...prev.credentials,
                    accessKeyID: e.target.value,
                  },
                }))
              }
            />
          </FormGroup>
          <FormGroup style={{ marginTop: 8, marginBottom: 0 }} label={<S.Label>Secret Access Key</S.Label>}>
            <InputGroup
              placeholder="Your Secret Access Key"
              value={(values.credentials as AWSAuth).secretAccessKey}
              onChange={(e) =>
                setValuesDefault((prev) => ({
                  ...prev,
                  credentials: {
                    ...prev.credentials,
                    secretAccessKey: e.target.value,
                  },
                }))
              }
            />
          </FormGroup>
          <FormGroup style={{ marginTop: 8, marginBottom: 0 }} label={<S.Label>Cluster Name</S.Label>}>
            <InputGroup
              placeholder=" Your Cluster Name"
              value={(values.credentials as AWSAuth).clusterName}
              onChange={(e) =>
                setValuesDefault((prev) => ({
                  ...prev,
                  credentials: {
                    ...prev.credentials,
                    clusterName: e.target.value,
                  },
                }))
              }
            />
          </FormGroup>
          <FormGroup style={{ marginTop: 8, marginBottom: 0 }} label={<S.Label>AWS Region</S.Label>}>
            <InputGroup
              placeholder="Your AWS Region"
              value={(values.credentials as AWSAuth).awsRegion}
              onChange={(e) =>
                setValuesDefault((prev) => ({
                  ...prev,
                  credentials: {
                    ...prev.credentials,
                    awsRegion: e.target.value,
                  },
                }))
              }
            />
          </FormGroup>
        </>
      )}

      {values.credentials?.providerType === 'openShift' && (
        <FormGroup
          style={{ marginTop: 8, marginBottom: 0 }}
          label={<S.Label>Authentication URL for Openshift</S.Label>}
        >
          <InputGroup
            placeholder="Authentication URL for Openshift"
            value={(values.credentials as OpenShiftAuth).authenticationURLForOpenshift}
            onChange={(e) =>
              setValuesDefault((prev) => ({
                ...prev,
                credentials: {
                  ...prev.credentials,
                  authenticationURLForOpenshift: e.target.value,
                },
              }))
            }
          />
        </FormGroup>
      )}
    </>
  );
};
