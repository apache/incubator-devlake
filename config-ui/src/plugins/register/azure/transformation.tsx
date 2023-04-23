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
import { Tag, Switch, Radio, InputGroup, Icon, Collapse, Intent } from '@blueprintjs/core';

import { Divider, ExternalLink, HelpTooltip } from '@/components';

import * as S from './styled';

interface Props {
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
}

export const AzureTransformation = ({ transformation, setTransformation }: Props) => {
  const [enableCICD, setEnableCICD] = useState(true);
  const [openAdditionalSettings, setOpenAdditionalSettings] = useState(false);

  useEffect(() => {
    if (transformation.refdiff) {
      setOpenAdditionalSettings(true);
    }
  }, [transformation]);

  const handleChangeCICDEnable = (e: React.FormEvent<HTMLInputElement>) => {
    const checked = (e.target as HTMLInputElement).checked;

    if (checked) {
      setTransformation({
        ...transformation,
        deploymentPattern: '(deploy|push-image)',
        productionPattern: 'production',
      });
    } else {
      setTransformation({
        ...transformation,
        deploymentPattern: undefined,
        productionPattern: undefined,
      });
    }
    setEnableCICD(checked);
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
      <S.CICD>
        <h2>CI/CD</h2>
        <h3>
          <span>Deployment</span>
          <Tag minimal intent={Intent.PRIMARY} style={{ marginLeft: 8 }}>
            DORA
          </Tag>
          <div className="switch">
            <span>Enable</span>
            <Switch alignIndicator="right" inline checked={enableCICD} onChange={handleChangeCICDEnable} />
          </div>
        </h3>
        {enableCICD && (
          <>
            <p>
              Use Regular Expression to define Deployments in DevLake in order to measure DORA metrics.{' '}
              <ExternalLink link="https://devlake.apache.org/docs/Configuration/GitHub#step-3---adding-transformation-rules-optional">
                Learn more
              </ExternalLink>
            </p>
            <div style={{ marginTop: 16 }}>Convert a Azure Pipeline Run as a DevLake Deployment when: </div>
            <div className="text">
              <span>
                The name of the <strong>Azure pipeline</strong> or <strong>one of its jobs</strong> matches
              </span>
              <InputGroup
                style={{ width: 224, margin: '0 8px' }}
                placeholder="(deploy|push-image)"
                value={transformation.deploymentPattern ?? ''}
                onChange={(e) =>
                  setTransformation({
                    ...transformation,
                    deploymentPattern: e.target.value,
                    productionPattern: !e.target.value ? '' : transformation.productionPattern,
                  })
                }
              />
              <i style={{ color: '#E34040' }}>*</i>
              <HelpTooltip content="Azure Pipelines: https://learn.microsoft.com/en-us/azure/devops/pipelines/get-started/what-is-azure-pipelines?view=azure-devops#continuous-testing" />
            </div>
            <div className="text">
              <span>If the name also matches</span>
              <InputGroup
                style={{ width: 120, margin: '0 8px' }}
                disabled={!transformation.deploymentPattern}
                placeholder="prod(.*)"
                value={transformation.productionPattern ?? ''}
                onChange={(e) =>
                  setTransformation({
                    ...transformation,
                    productionPattern: e.target.value,
                  })
                }
              />
              <span>, this Deployment is a ‘Production Deployment’</span>
              <HelpTooltip content="If you leave this field empty, all DevLake Deployments will be tagged as in the Production environment. " />
            </div>
          </>
        )}
      </S.CICD>
      <Divider />
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
