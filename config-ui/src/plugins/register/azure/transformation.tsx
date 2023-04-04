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

import React, { useState, useEffect } from 'react';
import { Tag, RadioGroup, Radio, InputGroup, Icon, Collapse, Intent } from '@blueprintjs/core';

import { ExternalLink, HelpTooltip } from '@/components';

import * as S from './styled';

interface Props {
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
}

export const AzureTransformation = ({ transformation, setTransformation }: Props) => {
  const [enableCICD, setEnableCICD] = useState(1);
  const [openAdditionalSettings, setOpenAdditionalSettings] = useState(false);

  useEffect(() => {
    if (transformation.refdiff) {
      setOpenAdditionalSettings(true);
    }
  }, [transformation]);

  const handleChangeCICDEnable = (e: number) => {
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
    setEnableCICD(e);
  };

  const handleChangeAdditionalSettingsOpen = () => {
    setOpenAdditionalSettings(!openAdditionalSettings);
    if (!openAdditionalSettings) {
      setTransformation({
        ...transformation,
        refdiff: null,
      });
    }
  };

  return (
    <S.Transfromation>
      <div className="ci-cd">
        <h2>CI/CD</h2>
        <h3>
          <span>Deployment</span>
          <Tag minimal intent={Intent.PRIMARY} style={{ marginLeft: 4, fontWeight: 400 }}>
            DORA
          </Tag>
        </h3>
        <p>Tell DevLake what CI builds are Deployments.</p>
        <RadioGroup
          selectedValue={enableCICD}
          onChange={(e) => handleChangeCICDEnable(+(e.target as HTMLInputElement).value)}
        >
          <Radio label="Detect Deployment from Builds in Azure Pipelines" value={1} />
          {enableCICD === 1 && (
            <div className="radio">
              <p>
                Please fill in the following RegEx, as DevLake ONLY accounts for deployments in the production
                environment for DORA metrics. Not sure what an Azure Build is?
                <ExternalLink link="https://learn.microsoft.com/en-us/azure/devops/pipelines/get-started/what-is-azure-pipelines?view=azure-devops#continuous-testing">
                  See it here
                </ExternalLink>
              </p>
              <div className="input">
                <p>The Build name that matches</p>
                <InputGroup
                  placeholder="(?i)deploy"
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
                  placeholder="(?i)production"
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
          <Radio label="Not using Builds in Azure Pipelines as Deployments" value={0} />
        </RadioGroup>
      </div>
      {/* Additional Settings */}
      <div className="additional-settings">
        <h2 onClick={handleChangeAdditionalSettingsOpen}>
          <Icon icon={!openAdditionalSettings ? 'chevron-up' : 'chevron-down'} size={18} />
          <span>Additional Settings</span>
        </h2>
        <Collapse isOpen={openAdditionalSettings}>
          <div className="radio">
            <Radio defaultChecked />
            <p>
              Enable the <ExternalLink link="https://devlake.apache.org/docs/Plugins/refdiff">RefDiff</ExternalLink>{' '}
              plugin to pre-calculate version-based metrics
              <HelpTooltip content="Calculate the commits diff between two consecutive tags that match the following RegEx. Issues closed by PRs which contain these commits will also be calculated. The result will be shown in table.refs_commits_diffs and table.refs_issues_diffs." />
            </p>
          </div>
          <div className="refdiff">
            Compare the last
            <InputGroup
              style={{ width: 60 }}
              placeholder="10"
              value={transformation.refdiff?.tagsLimit}
              onChange={(e) =>
                setTransformation({
                  ...transformation,
                  refdiff: {
                    ...transformation?.refdiff,
                    tagsLimit: e.target.value,
                  },
                })
              }
            />
            tags that match the
            <InputGroup
              style={{ width: 200 }}
              placeholder="v\d+\.\d+(\.\d+(-rc)*\d*)*$"
              value={transformation.refdiff?.tagsPattern}
              onChange={(e) =>
                setTransformation({
                  ...transformation,
                  refdiff: {
                    ...transformation?.refdiff,
                    tagsPattern: e.target.value,
                  },
                })
              }
            />
            for calculation
          </div>
        </Collapse>
      </div>
    </S.Transfromation>
  );
};
