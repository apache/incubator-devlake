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

import { CaretRightOutlined } from '@ant-design/icons';
import {theme, Collapse, Tag, Input, Form} from 'antd';

import { ExternalLink, HelpTooltip } from '@/components';
import { DOC_URL } from '@/release';

interface Props {
  entities: string[];
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
}

export const AzureTransformation = ({ entities, transformation, setTransformation }: Props) => {
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
        transformation,
        onChangeTransformation: setTransformation,
      })}
    />
  );
};

const renderCollapseItems = ({
  entities,
  panelStyle,
  transformation,
  onChangeTransformation,
}: {
  entities: string[];
  panelStyle: React.CSSProperties;
  transformation: any;
  onChangeTransformation: any;
}) =>
  [
    {
      key: 'TICKET',
      label: 'Issue Tracking',
      style: panelStyle,
      children: (
        <>
          <p>
            Tell DevLake what your issue labels mean to view metrics such as{' '}
            <ExternalLink link={DOC_URL.METRICS.BUG_AGE}>Bug Age</ExternalLink>,{' '}
            <ExternalLink link={DOC_URL.METRICS.MTTR}>DORA - Median Time to Restore Service</ExternalLink>, etc.
          </p>
          <p>
            DevLake defines three standard types of issues: FEATURE, BUG and INCIDENT. Set your issues to these three
            types with issue labels that match the RegEx.
          </p>
          <Form.Item label="Requirement">
            <Input
              placeholder="(feat|feature|proposal|requirement)"
              value={transformation.issueTypeRequirement ?? ''}
              onChange={(e) =>
                onChangeTransformation({
                  ...transformation,
                  issueTypeRequirement: e.target.value,
                })
              }
            />
          </Form.Item>
          <Form.Item label="Bug">
            <Input
              placeholder="(bug|broken)"
              value={transformation.issueTypeBug ?? ''}
              onChange={(e) =>
                onChangeTransformation({
                  ...transformation,
                  issueTypeBug: e.target.value,
                })
              }
            />
          </Form.Item>
          <Form.Item
            label={
              <>
                <span>Incident</span>
                <Tag style={{ marginLeft: 4 }} color="blue">
                  DORA
                </Tag>
              </>
            }
          >
            <Input
              placeholder="(incident|failure)"
              value={transformation.issueTypeIncident ?? ''}
              onChange={(e) =>
                onChangeTransformation({
                  ...transformation,
                  issueTypeIncident: e.target.value,
                })
              }
            />
          </Form.Item>
          <Form.Item
            label={
              <>
                <span style={{ marginRight: 4 }}>Issue Priority</span>
                <HelpTooltip content="Labels that match the RegEx will be set as the priority of an issue." />
              </>
            }
          >
            <Input
              placeholder="(highest|high|medium|low|p0|p1|p2|p3)"
              value={transformation.issuePriority ?? ''}
              onChange={(e) =>
                onChangeTransformation({
                  ...transformation,
                  issuePriority: e.target.value,
                })
              }
            />
          </Form.Item>
          <Form.Item
            label={
              <>
                <span style={{ marginRight: 4 }}>Issue Component</span>
                <HelpTooltip content="Labels that match the RegEx will be set as the component of an issue." />
              </>
            }
          >
            <Input
              placeholder="component(.*)"
              value={transformation.issueComponent ?? ''}
              onChange={(e) =>
                onChangeTransformation({
                  ...transformation,
                  issueComponent: e.target.value,
                })
              }
            />
          </Form.Item>
          <Form.Item
            label={
              <>
                <span style={{ marginRight: 4 }}>Issue Severity</span>
                <HelpTooltip content="Labels that match the RegEx will be set as the serverity of an issue." />
              </>
            }
          >
            <Input
              placeholder="severity(.*)"
              value={transformation.issueSeverity ?? ''}
              onChange={(e) =>
                onChangeTransformation({
                  ...transformation,
                  issueSeverity: e.target.value,
                })
              }
            />
          </Form.Item>
        </>
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
            <ExternalLink link={DOC_URL.PLUGIN.AZUREDEVOPS.TRANSFORMATION}>Learn more</ExternalLink>
          </p>
          <div>Convert a Azure Pipeline Run as a DevLake Deployment when: </div>
          <div style={{ margin: '8px 0', paddingLeft: 28 }}>
            <span>
              The name of the <strong>Azure pipeline</strong> or <strong>one of its jobs</strong> matches
            </span>
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
            <i style={{ color: '#E34040' }}>*</i>
            <HelpTooltip content="Azure Pipelines: https://learn.microsoft.com/en-us/azure/devops/pipelines/get-started/what-is-azure-pipelines?view=azure-devops#continuous-testing" />
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
            <HelpTooltip content="If you leave this field empty, all DevLake Deployments will be tagged as in the Production environment. " />
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
