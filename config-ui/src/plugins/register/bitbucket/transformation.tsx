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
import { FormGroup, InputGroup, Tag, Intent, Checkbox } from '@blueprintjs/core';

import { ExternalLink, HelpTooltip, Divider, MultiSelector } from '@/components';

import ExampleJpg from './assets/bitbucket-example.jpg';
import * as S from './styled';

interface Props {
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
}

const ALL_STATES = ['new', 'open', 'resolved', 'closed', 'on hold', 'wontfix', 'duplicate', 'invalid'];

export const BitbucketTransformation = ({ transformation, setTransformation }: Props) => {
  const [useCustom, setUseCustom] = useState(false);

  useEffect(() => {
    if (transformation.deploymentPattern || transformation.productionPattern) {
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

  const handleChangeUseCustom = (e: React.FormEvent<HTMLInputElement>) => {
    const checked = (e.target as HTMLInputElement).checked;

    if (!checked) {
      setTransformation({
        ...transformation,
        deploymentPattern: '',
        productionPattern: '',
      });
    }

    setUseCustom(checked);
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
        </h3>
        <p style={{ marginBottom: 16 }}>
          Use Regular Expression to define Deployments in DevLake in order to measure DORA metrics.{' '}
          <ExternalLink link="https://devlake.apache.org/docs/Configuration/GitHub#step-3---adding-transformation-rules-optional">
            Learn more
          </ExternalLink>
        </p>
        <div className="text">
          <Checkbox disabled checked />
          <span>Convert a BitBucket Deployment to a DevLake Deployment </span>
          <HelpTooltip content={<img src={ExampleJpg} alt="" width={400} />} />
        </div>
        <div className="text">
          <Checkbox checked={useCustom} onChange={handleChangeUseCustom} />
          <span>
            Convert a BitBucket Pipeline to a DevLake Deployment when its branch/tag name or one of its pipeline steps’
            names
          </span>
        </div>
        <div className="sub-text">
          <span>matches</span>
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
          <span>.</span>
          <HelpTooltip content="View your BitBucket Pipelines: https://support.atlassian.com/bitbucket-cloud/docs/view-your-pipeline/" />
        </div>
        <div className="sub-text">
          <span>If the name also matches</span>
          <InputGroup
            style={{ width: 200, margin: '0 8px' }}
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
          <HelpTooltip content="If you leave this field empty, all Deployments will be tagged as in the Production environment. " />
        </div>
      </S.CICD>
    </S.Transformation>
  );
};
