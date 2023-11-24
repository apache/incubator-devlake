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

import { useMemo, useState, useEffect } from 'react';
import { CaretRightOutlined } from '@ant-design/icons';
import { theme, Collapse, Tag, Form, Input, Checkbox } from 'antd';

import { ExternalLink, HelpTooltip, MultiSelector } from '@/components';
import { DOC_URL } from '@/release';

import ExampleJpg from './assets/bitbucket-example.jpg';

interface Props {
  entities: string[];
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
}

const ALL_STATES = ['new', 'open', 'resolved', 'closed', 'on hold', 'wontfix', 'duplicate', 'invalid'];

export const BitbucketTransformation = ({ entities, transformation, setTransformation }: Props) => {
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

  const { token } = theme.useToken();

  const panelStyle: React.CSSProperties = {
    marginBottom: 24,
    background: token.colorFillAlter,
    borderRadius: token.borderRadiusLG,
    border: 'none',
  };

  return (
    <Collapse
      bordered={false}
      defaultActiveKey={['TICKET', 'CICD']}
      expandIcon={({ isActive }) => <CaretRightOutlined rotate={isActive ? 90 : 0} rev="" />}
      style={{ background: token.colorBgContainer }}
      size="large"
      items={renderCollapseItems({
        entities,
        panelStyle,
        selectedStates,
        transformation,
        onChangeTransformation: setTransformation,
        useCustom,
        onChangeUseCustom: handleChangeUseCustom,
      })}
    />
  );
};

const renderCollapseItems = ({
  entities,
  panelStyle,
  selectedStates,
  transformation,
  onChangeTransformation,
  useCustom,
  onChangeUseCustom,
}: {
  entities: string[];
  panelStyle: React.CSSProperties;
  selectedStates: any;
  transformation: any;
  onChangeTransformation: any;
  useCustom: boolean;
  onChangeUseCustom: any;
}) =>
  [
    {
      key: 'TICKET',
      label: (
        <>
          <span>Issue Status Mapping</span>
          <HelpTooltip content="Standardize your issue statuses to the following issue statuses to view metrics such as `Requirement Delivery Rate` in built-in dashboards." />
        </>
      ),
      style: panelStyle,
      children: (
        <div className="list">
          <Form.Item label="TODO">
            <MultiSelector
              items={ALL_STATES}
              disabledItems={selectedStates}
              selectedItems={transformation.issueStatusTodo ? transformation.issueStatusTodo.split(',') : []}
              onChangeItems={(selectedItems) =>
                onChangeTransformation({
                  ...transformation,
                  issueStatusTodo: selectedItems.join(','),
                })
              }
            />
          </Form.Item>
          <Form.Item label="IN-PROGRESS">
            <MultiSelector
              items={ALL_STATES}
              disabledItems={selectedStates}
              selectedItems={
                transformation.issueStatusInProgress ? transformation.issueStatusInProgress.split(',') : []
              }
              onChangeItems={(selectedItems) =>
                onChangeTransformation({
                  ...transformation,
                  issueStatusInProgress: selectedItems.join(','),
                })
              }
            />
          </Form.Item>
          <Form.Item label="DONE">
            <MultiSelector
              items={ALL_STATES}
              disabledItems={selectedStates}
              selectedItems={transformation.issueStatusDone ? transformation.issueStatusDone.split(',') : []}
              onChangeItems={(selectedItems) =>
                onChangeTransformation({
                  ...transformation,
                  issueStatusDone: selectedItems.join(','),
                })
              }
            />
          </Form.Item>
          <Form.Item label="OTHER">
            <MultiSelector
              items={ALL_STATES}
              disabledItems={selectedStates}
              selectedItems={transformation.issueStatusOther ? transformation.issueStatusOther.split(',') : []}
              onChangeItems={(selectedItems) =>
                onChangeTransformation({
                  ...transformation,
                  issueStatusOther: selectedItems.join(','),
                })
              }
            />
          </Form.Item>
        </div>
      ),
    },
    {
      key: 'CICD',
      label: 'CI/CD',
      style: panelStyle,
      children: (
        <>
          <h3 style={{ marginBottom: 16 }}>
            <span>Deployment</span>
            <Tag style={{ marginLeft: 4 }} color="blue">
              DORA
            </Tag>
          </h3>
          <p style={{ marginBottom: 16 }}>
            Use Regular Expression to define Deployments in DevLake in order to measure DORA metrics.{' '}
            <ExternalLink link={DOC_URL.PLUGIN.BITBUCKET.TRANSFORMATION}>Learn more</ExternalLink>
          </p>
          <Checkbox disabled checked>
            <span>Convert a BitBucket Deployment to a DevLake Deployment</span>
            <HelpTooltip content={<img src={ExampleJpg} alt="" width={400} />} />
          </Checkbox>
          <Checkbox checked={useCustom} onChange={onChangeUseCustom}>
            Convert a BitBucket Pipeline to a DevLake Deployment when its branch/tag name
          </Checkbox>
          <div style={{ margin: '8px 0', paddingLeft: 28 }}>
            <span>matches</span>
            <Input
              style={{ width: 200, margin: '0 8px' }}
              placeholder="(deploy|push-image)"
              value={transformation.deploymentPattern ?? ''}
              onChange={(e) =>
                onChangeTransformation({
                  ...transformation,
                  deploymentPattern: e.target.value,
                  productionPattern: !e.target.value ? '' : transformation.productionPattern,
                })
              }
            />
            <span>.</span>
            <HelpTooltip content="View your BitBucket Pipelines: https://support.atlassian.com/bitbucket-cloud/docs/view-your-pipeline/" />
          </div>
          <div style={{ margin: '8px 0', paddingLeft: 28 }}>
            <span>If the name also matches</span>
            <Input
              style={{ width: 200, margin: '0 8px' }}
              placeholder="prod(.*)"
              value={transformation.productionPattern ?? ''}
              onChange={(e) =>
                onChangeTransformation({
                  ...transformation,
                  productionPattern: e.target.value,
                })
              }
            />
            <span>, this Deployment is a ‘Production Deployment’</span>
            <HelpTooltip content="If you leave this field empty, all Deployments will be tagged as in the Production environment. " />
          </div>
        </>
      ),
    },
    {
      key: 'ADDITIONAL',
      label: 'Additional Settings',
      style: panelStyle,
      children: (
        <>
          <p>
            Enable the <ExternalLink link={DOC_URL.PLUGIN.REFDIFF}>RefDiff</ExternalLink> plugin to pre-calculate
            version-based metrics
            <HelpTooltip content="Calculate the commits diff between two consecutive tags that match the following RegEx. Issues closed by PRs which contain these commits will also be calculated. The result will be shown in table.refs_commits_diffs and table.refs_issues_diffs." />
          </p>
          <div className="refdiff">
            Compare the last
            <Input
              style={{ margin: '0 8px', width: 60 }}
              placeholder="10"
              value={transformation.refdiff?.tagsLimit ?? ''}
              onChange={(e) =>
                onChangeTransformation({
                  ...transformation,
                  refdiff: {
                    ...transformation?.refdiff,
                    tagsLimit: +e.target.value,
                  },
                })
              }
            />
            tags that match the
            <Input
              style={{ margin: '0 8px', width: 200 }}
              placeholder="(regex)$"
              value={transformation.refdiff?.tagsPattern ?? ''}
              onChange={(e) =>
                onChangeTransformation({
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
        </>
      ),
    },
  ].filter((it) => entities.includes(it.key) || it.key === 'ADDITIONAL');
