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

import { useMemo } from 'react';
import { useParams } from 'react-router-dom';

import { PageHeader, Workflow } from '@/components';
import { ConnectionContextProvider } from '@/store';

import { ModeEnum, FromEnum } from '../types';

import { ContextProvider, Context } from './context';
import { Step1, Step2, Step3, Step4 } from './components';

interface Props {
  from: FromEnum;
}

export const BlueprintCreatePage = ({ from }: Props) => {
  const { pname } = useParams<{ pname: string }>();

  const breadcrumbs = useMemo(
    () =>
      from === FromEnum.project
        ? [
            { name: 'Projects', path: '/projects' },
            { name: window.decodeURIComponent(pname), path: `/projects/${pname}` },
            {
              name: 'Create a Blueprint',
              path: `/projects/${pname}/create-blueprint`,
            },
          ]
        : [
            { name: 'Blueprints', path: '/blueprints' },
            { name: 'Create a Blueprint', path: '/blueprints/create' },
          ],
    [from, pname],
  );

  return (
    <ConnectionContextProvider filterBeta>
      <ContextProvider from={from} projectName={pname}>
        <Context.Consumer>
          {({ step, mode }) => (
            <PageHeader breadcrumbs={breadcrumbs}>
              <Workflow
                steps={
                  mode === ModeEnum.normal
                    ? ['Add Data Connections', 'Set Data Scope', 'Add Transformation (Optional)', 'Set Sync Policy']
                    : ['Create Advanced Configuration', 'Set Sync Policy']
                }
                activeStep={step}
              />
              {step === 1 && <Step1 from={from} />}
              {mode === ModeEnum.normal && step === 2 && <Step2 />}
              {step === 3 && <Step3 />}
              {((mode === ModeEnum.normal && step === 4) || (mode === ModeEnum.advanced && step === 2)) && <Step4 />}
            </PageHeader>
          )}
        </Context.Consumer>
      </ContextProvider>
    </ConnectionContextProvider>
  );
};
