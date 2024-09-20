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
import { theme, Collapse, Tag, Input } from 'antd';

import { ShowMore, HelpTooltip } from '@/components';
import { CheckMatchedItems } from '@/plugins';

import { WorkflowRun } from './workflow-run';

interface Props {
  plugin: string;
  connectionId: ID;
  entities: string[];
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
}

export const JenkinsTransformation = ({ plugin, connectionId, entities, transformation, setTransformation }: Props) => {
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
        transformation,
        onChangeTransformation: setTransformation,
      })}
    />
  );
};

const renderCollapseItems = ({
  plugin,
  connectionId,
  entities,
  panelStyle,
  transformation,
  onChangeTransformation,
}: {
  plugin: string;
  connectionId: ID;
  entities: string[];
  panelStyle: React.CSSProperties;
  transformation: any;
  onChangeTransformation: any;
}) =>
  [
    {
      key: 'CICD',
      label: 'CI/CD',
      style: panelStyle,
      children: (
        <>
          <h3>
            <span>Deployment</span>
            <Tag style={{ marginLeft: 4 }} color="blue">
              DORA
            </Tag>
          </h3>
          <ShowMore
            text={<p>Use Regular Expression to define Deployments to measure DORA metrics.</p>}
            btnText="See how to configure"
          >
            <WorkflowRun />
          </ShowMore>
          <div>Convert a Jenkins Build as a DevLake Deployment when: </div>
          <div style={{ margin: '8px 0', paddingLeft: 28 }}>
            <span>
              The name of the <strong>Jenkins job</strong> or <strong>one of its stages</strong> matches
            </span>
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
            <i style={{ color: '#E34040' }}>*</i>
            <HelpTooltip content="Jenkins Builds: https://www.jenkins.io/doc/pipeline/steps/pipeline-build-step/" />
          </div>
          <div style={{ margin: '8px 0', paddingLeft: 28 }}>
            <span>If the name also matches</span>
            <Input
              style={{ width: 120, margin: '0 8px' }}
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
            <span>, this Deployment is a ‘Production Deployment’</span>
            <HelpTooltip content="If you leave this field empty, all DevLake Deployments will be tagged as in the Production environment. " />
          </div>
          <CheckMatchedItems plugin={plugin} connectionId={connectionId} transformation={transformation} />
        </>
      ),
    },
  ].filter((it) => entities.includes(it.key));
