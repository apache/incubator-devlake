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

import React, { useMemo, useState, useEffect } from 'react';
import { FormGroup, InputGroup, Tag, Radio, Icon, Collapse, Intent, Switch } from '@blueprintjs/core';

import { ExternalLink, HelpTooltip, Divider, MultiSelector } from '@/components';

import ExampleJpg from './assets/bitbucket-example.jpg';
import * as S from './styled';

interface Props {
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
}

const ALL_STATES = ['new', 'open', 'resolved', 'closed', 'on hold', 'wontfix', 'duplicate', 'invalid'];

export const BitbucketTransformation = ({ transformation, setTransformation }: Props) => {
  const [enableCICD, setEnableCICD] = useState(true);
  const [useCustom, setUseCustom] = useState(false);
  const [openAdditionalSettings, setOpenAdditionalSettings] = useState(false);

  useEffect(() => {
    if (transformation.refdiff) {
      setOpenAdditionalSettings(true);
    }
    if (transformation.deploymentPattern) {
      setUseCustom(true);
    } else {
      setUseCustom(false);
    }
  }, [transformation]);

  const selectedStates = useMemo(
    () => [
      ...(transformation.issueStatusTodo ? transformation.issueStatusTodo.split(',') : []),
      ...(transformation.issueStatusInProgress ? transformation.issueStatusInProgress.split(',') : []),
      ...(transformation.issueStatusDone ? transformation.issueStatusDone.split(',') : []),
      ...(transformation.issueStatusOther ? transformation.issueStatusOther.split(',') : []),
    ],
    [transformation],
  );

  const handleChangeUseCustom = (uc: boolean) => {
    if (!uc) {
      setTransformation({
        ...transformation,
        deploymentPattern: '',
        productionPattern: '',
      });
    } else {
      setTransformation({
        ...transformation,
        deploymentPattern: '(deploy|push-image)',
        productionPattern: '',
      });
    }

    setUseCustom(uc);
  };

  const handleChangeCICDEnable = (e: React.FormEvent<HTMLInputElement>) => {
    const checked = (e.target as HTMLInputElement).checked;

    if (checked) {
      setUseCustom(false);
      setTransformation({
        ...transformation,
        deploymentPattern: '',
        productionPattern: '',
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
    <S.Transformation>
      {/* Issue Tracking */}
      <div className="issue-tracking">
        <h2>Issue Tracking</h2>
        <div className="issue-type">
          <div className="title">
            <span>Issue Status Mapping</span>
            <HelpTooltip content="Standardize your issue statuses to the following issue statuses to view metrics such as `Requirement Delivery Rate` in built-in dashboards." />
          </div>
          <div className="list">
            <FormGroup inline label="TODO">
              <MultiSelector
                items={ALL_STATES}
                disabledItems={selectedStates}
                selectedItems={transformation.issueStatusTodo ? transformation.issueStatusTodo.split(',') : []}
                onChangeItems={(selectedItems) =>
                  setTransformation({
                    ...transformation,
                    issueStatusTodo: selectedItems.join(','),
                  })
                }
              />
            </FormGroup>
            <FormGroup inline label="IN-PROGRESS">
              <MultiSelector
                items={ALL_STATES}
                disabledItems={selectedStates}
                selectedItems={
                  transformation.issueStatusInProgress ? transformation.issueStatusInProgress.split(',') : []
                }
                onChangeItems={(selectedItems) =>
                  setTransformation({
                    ...transformation,
                    issueStatusInProgress: selectedItems.join(','),
                  })
                }
              />
            </FormGroup>
            <FormGroup inline label="DONE">
              <MultiSelector
                items={ALL_STATES}
                disabledItems={selectedStates}
                selectedItems={transformation.issueStatusDone ? transformation.issueStatusDone.split(',') : []}
                onChangeItems={(selectedItems) =>
                  setTransformation({
                    ...transformation,
                    issueStatusDone: selectedItems.join(','),
                  })
                }
              />
            </FormGroup>
            <FormGroup inline label="OTHER">
              <MultiSelector
                items={ALL_STATES}
                disabledItems={selectedStates}
                selectedItems={transformation.issueStatusOther ? transformation.issueStatusOther.split(',') : []}
                onChangeItems={(selectedItems) =>
                  setTransformation({
                    ...transformation,
                    issueStatusOther: selectedItems.join(','),
                  })
                }
              />
            </FormGroup>
          </div>
        </div>
      </div>
      <Divider />
      {/* CI/CD */}
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
            <div style={{ margin: '16px 0' }}>Convert a BitBucket Pipeline as a DevLake Deployment when: </div>
            <div className="text">
              <Radio checked={!useCustom} onChange={() => handleChangeUseCustom(false)} />
              <span>It has one or more BitBucket deployments. See the example.</span>
              <HelpTooltip content={<img src={ExampleJpg} alt="" width={400} />} />
            </div>
            <div className="text">
              <Radio checked={useCustom} onChange={() => handleChangeUseCustom(true)} />
              <div style={{ display: 'flex', alignItems: 'center' }}>
                <span>Its branch/tag name or one of its pipeline steps’ names matches</span>
                <InputGroup
                  style={{ width: 200, margin: '0 8px' }}
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
                <HelpTooltip content="If you leave this field empty, all DevLake Deployments will be tagged as in the Production environment. " />
              </div>
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
    </S.Transformation>
  );
};
