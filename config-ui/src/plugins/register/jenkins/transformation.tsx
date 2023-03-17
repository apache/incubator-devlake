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

import React, { useState } from 'react';
import { Tag, Intent, RadioGroup, Radio, InputGroup } from '@blueprintjs/core';

import * as S from './styled';
import { ExternalLink, HelpTooltip } from '@/components';

interface Props {
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
}

export const JenkinsTransformation = ({ transformation, setTransformation }: Props) => {
  const [enable, setEnable] = useState(1);

  const handleChangeEnable = (e: number) => {
    if (e === 0) {
      setTransformation({
        ...transformation,
        deploymentPattern: undefined,
        productionPattern: undefined,
      });
    } else {
      setTransformation({
        ...transformation,
        deploymentPattern: '',
        productionPattern: '',
      });
    }
    setEnable(e);
  };

  return (
    <S.TransformationWrapper>
      <h2>CI/CD</h2>
      <h3>
        <span>Deployment</span>
        <Tag minimal intent={Intent.PRIMARY}>
          DORA
        </Tag>
      </h3>
      <p>Tell DevLake what CI builds are Deployments.</p>
      <RadioGroup selectedValue={enable} onChange={(e) => handleChangeEnable(+(e.target as HTMLInputElement).value)}>
        <Radio label="Detect Deployments from Jenkins Builds" value={1} />
        {enable === 1 && (
          <div className="radio">
            <p>
              Please fill in the following RegEx, as DevLake ONLY accounts for deployments in the production environment
              for DORA metrics. Not sure what a Jenkins Build is?{' '}
              <ExternalLink link="https://www.jenkins.io/doc/book/pipeline/#overview">See it here</ExternalLink>
            </p>
            <div className="input">
              <p>The Build name that matches</p>
              <InputGroup
                placeholder="(deploy|push-image)"
                value={transformation.deploymentPattern}
                onChange={(e) =>
                  setTransformation({
                    ...transformation,
                    deploymentPattern: e.target.value,
                  })
                }
              />
              <p>
                will be registered as a `Deployment` in DevLake. <span style={{ color: '#E34040' }}>*</span>
              </p>
            </div>
            <div className="input">
              <p>The Build name that matches</p>
              <InputGroup
                disabled={!transformation.deploymentPattern}
                placeholder="production"
                value={transformation.productionPattern}
                onChange={(e) =>
                  setTransformation({
                    ...transformation,
                    productionPattern: e.target.value,
                  })
                }
              />
              <p>
                will be registered as a `Deployment` to the Production environment in DevLake.
                <HelpTooltip content="If you leave this field empty, all data will be tagged as in the Production environment. " />
              </p>
            </div>
          </div>
        )}
        <Radio label="Not using Jenkins Builds as Deployments" value={0} />
      </RadioGroup>
    </S.TransformationWrapper>
  );
};
