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
import { FormGroup, InputGroup, Tag, Radio, Icon, Collapse, Intent, Checkbox } from '@blueprintjs/core';

import { ExternalLink, HelpTooltip, Divider, MultiSelector } from '@/components';

import * as S from './styled';

interface Props {
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
}

const ALL_STATES = ['new', 'open', 'resolved', 'closed', 'on hold', 'wontfix', 'duplicate', 'invalid'];

export const BitbucketTransformation = ({ transformation, setTransformation }: Props) => {
  const [enableCICD, setEnableCICD] = useState(false);
  const [openAdditionalSettings, setOpenAdditionalSettings] = useState(false);

  useEffect(() => {
    if (transformation.refdiff) {
      setOpenAdditionalSettings(true);
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

  const handleChangeCICDEnable = (b: boolean) => {
    if (!b) {
      setTransformation({
        ...transformation,
        deploymentPattern: '',
        productionPattern: '',
      });
    }
    setEnableCICD(b);
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
    <S.TransformationWrapper>
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
      <div className="ci-cd">
        <h2>CI/CD</h2>
        <h3>
          <span>Deployment</span>
          <Tag minimal intent={Intent.PRIMARY}>
            DORA
          </Tag>
        </h3>
        <p>
          DevLake uses BitBucket{' '}
          <ExternalLink link="https://support.atlassian.com/bitbucket-cloud/docs/set-up-and-monitor-deployments/">
            deployments
          </ExternalLink>{' '}
          as DevLake deployments. If you are NOT using BitBucket deployments, DevLake provides the option to detect
          deployments from BitBucket pipeline steps.{' '}
          <ExternalLink link="https://devlake.apache.org/docs/Configuration/BitBucket#step-3---adding-transformation-rules-optional">
            Learn more
          </ExternalLink>
        </p>
        <Checkbox
          label="Detect Deployments from Pipeline steps in BitBucket"
          checked={enableCICD}
          onChange={(e) => handleChangeCICDEnable((e.target as HTMLInputElement).checked)}
        />
        {enableCICD && (
          <div className="radio">
            <div className="input">
              <p>The Pipeline step name that matches</p>
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
              <p>The Pipeline step name that matches</p>
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
      </div>
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
    </S.TransformationWrapper>
  );
};
