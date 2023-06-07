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
    if (values.providerType === 'azure') {
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
    } else if (values.providerType === 'aws') {
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
    } else if (values.providerType === 'openShift') {
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
          providerType: 'azure',
          subscriptionID: '',
          clientID: '',
          clientSecret: '',
          tenantID: '',
          resourceGroupName: '',
          clusterName: '',
        }));
        break;
      case 'aws':
        setValuesDefault((prev) => ({
          ...defaultValues,

          name: prev.name,
          providerType: 'aws',
          accessKeyID: '',
          clusterName: '',
          awsRegion: '',
          secretAccessKey: '',
        }));
        break;
      case 'openShift':
        setValuesDefault((prev) => ({
          ...defaultValues,

          name: prev.name,
          providerType: 'openShift',
          authenticationURLForOpenshift: '',
        }));
        break;
    }
  };

  return (
    <>
      <RadioGroup
        inline={true}
        selectedValue={values.providerType}
        onChange={(e) => handleChangeAuthType(e.currentTarget.value as AuthType)}
      >
        <Radio label="Azure" value="azure" />
        <Radio label="AWS" value="aws" />
        <Radio label="OpenShift" value="openShift" />
      </RadioGroup>

      {values.providerType === 'azure' && (
        <>
          <FormGroup style={{ marginTop: 8, marginBottom: 0 }} label={<S.Label>Subscription ID</S.Label>}>
            <InputGroup
              placeholder="Your Subscription ID"
              value={(values as AzureAuth).subscriptionID}
              onChange={(e) =>
                setValues({
                  subscriptionID: e.target.value,
                })
              }
            />
          </FormGroup>
          <FormGroup style={{ marginTop: 8, marginBottom: 0 }} label={<S.Label>Client ID</S.Label>}>
            <InputGroup
              placeholder="Your Client ID"
              value={(values as AzureAuth).clientID}
              onChange={(e) =>
                setValues({
                  clientID: e.target.value,
                })
              }
            />
          </FormGroup>
          <FormGroup style={{ marginTop: 8, marginBottom: 0 }} label={<S.Label>Client Secret</S.Label>}>
            <InputGroup
              placeholder="Your Client Secret"
              value={(values as AzureAuth).clientSecret}
              onChange={(e) =>
                setValues({
                  clientSecret: e.target.value,
                })
              }
            />
          </FormGroup>
          <FormGroup style={{ marginTop: 8, marginBottom: 0 }} label={<S.Label>Tenant ID</S.Label>}>
            <InputGroup
              placeholder="Your Tenant ID"
              value={(values as AzureAuth).tenantID}
              onChange={(e) =>
                setValues({
                  tenantID: e.target.value,
                })
              }
            />
          </FormGroup>
          <FormGroup style={{ marginTop: 8, marginBottom: 0 }} label={<S.Label>Resource Group Name</S.Label>}>
            <InputGroup
              placeholder="Your Resource Group Name"
              value={(values as AzureAuth).resourceGroupName}
              onChange={(e) =>
                setValues({
                  resourceGroupName: e.target.value,
                })
              }
            />
          </FormGroup>
          <FormGroup style={{ marginTop: 8, marginBottom: 0 }} label={<S.Label>Cluster name</S.Label>}>
            <InputGroup
              placeholder="Your Cluster name"
              value={(values as AzureAuth).clusterName}
              onChange={(e) =>
                setValues({
                  clusterName: e.target.value,
                })
              }
            />
          </FormGroup>
        </>
      )}

      {values.providerType === 'aws' && (
        <>
          <FormGroup style={{ marginTop: 8, marginBottom: 0 }} label={<S.Label>Access Key ID</S.Label>}>
            <InputGroup
              placeholder="Your Access Key ID"
              value={(values as AWSAuth).accessKeyID}
              onChange={(e) =>
                setValues({
                  accessKeyID: e.target.value,
                })
              }
            />
          </FormGroup>
          <FormGroup style={{ marginTop: 8, marginBottom: 0 }} label={<S.Label>Secret Access Key</S.Label>}>
            <InputGroup
              placeholder="Your Secret Access Key"
              value={(values as AWSAuth).secretAccessKey}
              onChange={(e) =>
                setValues({
                  secretAccessKey: e.target.value,
                })
              }
            />
          </FormGroup>
          <FormGroup style={{ marginTop: 8, marginBottom: 0 }} label={<S.Label>Cluster Name</S.Label>}>
            <InputGroup
              placeholder=" Your Cluster Name"
              value={(values as AWSAuth).clusterName}
              onChange={(e) =>
                setValues({
                  clusterName: e.target.value,
                })
              }
            />
          </FormGroup>
          <FormGroup style={{ marginTop: 8, marginBottom: 0 }} label={<S.Label>AWS Region</S.Label>}>
            <InputGroup
              placeholder="Your AWS Region"
              value={(values as AWSAuth).awsRegion}
              onChange={(e) =>
                setValues({
                  awsRegion: e.target.value,
                })
              }
            />
          </FormGroup>
        </>
      )}

      {values.providerType === 'openShift' && (
        <FormGroup
          style={{ marginTop: 8, marginBottom: 0 }}
          label={<S.Label>Authentication URL for Openshift</S.Label>}
        >
          <InputGroup
            placeholder="Authentication URL for Openshift"
            value={(values as OpenShiftAuth).authenticationURLForOpenshift}
            onChange={(e) =>
              setValues({
                authenticationURLForOpenshift: e.target.value,
              })
            }
          />
        </FormGroup>
      )}
    </>
  );
};
