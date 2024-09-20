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
import { theme, Collapse, Tag, Form, Input, Checkbox, Select } from 'antd';

import { ShowMore, ExternalLink, HelpTooltip } from '@/components';
import { CheckMatchedItems } from '@/plugins';
import { DOC_URL } from '@/release';

import ExampleJpg from './assets/bitbucket-example.jpg';
import { WorkflowRun } from './workflow-run';

interface Props {
  plugin: string;
  connectionId: ID;
  entities: string[];
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
}

const ALL_STATES = ['new', 'open', 'resolved', 'closed', 'on hold', 'wontfix', 'duplicate', 'invalid'];

export const BitbucketTransformation = ({
  plugin,
  connectionId,
  entities,
  transformation,
  setTransformation,
}: Props) => {
  const [useCustom, setUseCustom] = useState(false);

  useEffect(() => {
    if (transformation.deploymentPattern || transformation.productionPattern) {
      setUseCustom(true);
    } else {
      setUseCustom(false);
    }
  }, []);

  const options = useMemo(() => {
    const disabledOptions = [
      ...(transformation.issueStatusTodo ? transformation.issueStatusTodo.split(',') : []),
      ...(transformation.issueStatusInProgress ? transformation.issueStatusInProgress.split(',') : []),
      ...(transformation.issueStatusDone ? transformation.issueStatusDone.split(',') : []),
      ...(transformation.issueStatusOther ? transformation.issueStatusOther.split(',') : []),
    ];
    return ALL_STATES.filter((it) => !disabledOptions.includes(it)).map((it) => ({ label: it, value: it }));
  }, [transformation]);

  const handleChangeUseCustom = (e: React.FormEvent<HTMLInputElement>) => {
    const checked = (e.target as HTMLInputElement).checked;

    if (!checked) {
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
        plugin,
        connectionId,
        entities,
        panelStyle,
        options,
        transformation,
        onChangeTransformation: setTransformation,
        useCustom,
        onChangeUseCustom: handleChangeUseCustom,
      })}
    />
  );
};

const renderCollapseItems = ({
  plugin,
  connectionId,
  entities,
  panelStyle,
  options,
  transformation,
  onChangeTransformation,
  useCustom,
  onChangeUseCustom,
}: {
  plugin: string;
  connectionId: ID;
  entities: string[];
  panelStyle: React.CSSProperties;
  options: Array<{ label: string; value: string }>;
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
            <Select
              mode="multiple"
              options={options}
              value={transformation.issueStatusTodo ? transformation.issueStatusTodo.split(',') : []}
              onChange={(value) =>
                onChangeTransformation({
                  ...transformation,
                  issueStatusTodo: value.join(','),
                })
              }
            />
          </Form.Item>
          <Form.Item label="IN-PROGRESS">
            <Select
              mode="multiple"
              options={options}
              value={transformation.issueStatusInProgress ? transformation.issueStatusInProgress.split(',') : []}
              onChange={(value) =>
                onChangeTransformation({
                  ...transformation,
                  issueStatusInProgress: value.join(','),
                })
              }
            />
          </Form.Item>
          <Form.Item label="DONE">
            <Select
              mode="multiple"
              options={options}
              value={transformation.issueStatusDone ? transformation.issueStatusDone.split(',') : []}
              onChange={(value) =>
                onChangeTransformation({
                  ...transformation,
                  issueStatusDone: value.join(','),
                })
              }
            />
          </Form.Item>
          <Form.Item label="OTHER">
            <Select
              mode="multiple"
              options={options}
              value={transformation.issueStatusOther ? transformation.issueStatusOther.split(',') : []}
              onChange={(value) =>
                onChangeTransformation({
                  ...transformation,
                  issueStatusOther: value.join(','),
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
            <span>Convert a Bitbucket Deployment to a DevLake Deployment</span>
            <HelpTooltip content={<img src={ExampleJpg} alt="" width={400} />} />
          </Checkbox>
          <Checkbox checked={useCustom} onChange={onChangeUseCustom}>
            Convert a Bitbucket Pipeline to a DevLake Deployment when its branch/tag name
          </Checkbox>
          {useCustom && (
            <div style={{ paddingLeft: 28 }}>
              <ShowMore
                text={<p>Select this option only if you are not enabling Bitbucket deployments.</p>}
                btnText="See how to configure"
              >
                <WorkflowRun />
              </ShowMore>
              <div style={{ margin: '8px 0' }}>
                <span>Its branch or one of its steps matches</span>
                <Input
                  style={{ width: 200, margin: '0 8px' }}
                  placeholder="(?i)(deploy|push-image)"
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
                <i style={{ color: '#E34040' }}>*</i>
                <HelpTooltip content="View your Bitbucket Pipelines: https://support.atlassian.com/bitbucket-cloud/docs/view-your-pipeline/" />
              </div>
              <div style={{ margin: '8px 0' }}>
                <span>If the branch or the step also matches</span>
                <Input
                  style={{ width: 200, margin: '0 8px' }}
                  disabled={!transformation.deploymentPattern}
                  placeholder="(?i)(prod|release)"
                  value={transformation.productionPattern ?? ''}
                  onChange={(e) =>
                    onChangeTransformation({
                      ...transformation,
                      productionPattern: e.target.value,
                    })
                  }
                />
                <span>, this deployment will be regarded as a ‘Production Deployment’.</span>
                <HelpTooltip content="If you leave this field empty, all Deployments will be tagged as in the Production environment. " />
              </div>
              <CheckMatchedItems plugin={plugin} connectionId={connectionId} transformation={transformation} />
            </div>
          )}
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
