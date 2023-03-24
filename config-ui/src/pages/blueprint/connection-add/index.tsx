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

import { useParams } from 'react-router-dom';

import { PageHeader, Workflow } from '@/components';

import { Step1, Step2, Step3 } from './components';
import { ContextProvider, Context } from './context';

export const BlueprintConnectioAddPage = () => {
  const { pname, bid } = useParams<{ pname?: string; bid: string }>();

  return (
    <ContextProvider pname={pname} id={bid}>
      <Context.Consumer>
        {({ name, step }) => (
          <PageHeader
            breadcrumbs={[
              ...(pname
                ? [
                    {
                      name: 'Projects',
                      path: '/projects',
                    },
                    {
                      name: pname,
                      path: `/projects/${pname}`,
                    },
                  ]
                : [{ name: name, path: `/blueprints/${bid}` }]),
              { name: 'Add a New Connection', path: '' },
            ]}
          >
            <Workflow
              steps={['Select a Data Connection', 'Set Data Scope', 'Add Transformation (Optional)']}
              activeStep={step}
            />
            {step === 1 && <Step1 />}
            {step === 2 && <Step2 />}
            {step === 3 && <Step3 />}
          </PageHeader>
        )}
      </Context.Consumer>
    </ContextProvider>
  );
};
