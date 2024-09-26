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

import { useState, useEffect } from 'react';
import { CaretRightOutlined } from '@ant-design/icons';
import { theme, Collapse, Tag, Input, Checkbox } from 'antd';

import { ShowMore, ExternalLink, HelpTooltip } from '@/components';
import { Deployments, CheckMatchedItems } from '@/plugins';
import { DOC_URL } from '@/release';

import { WorkflowRun } from './workflow-run';

interface Props {
  plugin: string;
  connectionId: ID;
  entities: string[];
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
  setHasError: React.Dispatch<React.SetStateAction<boolean>>;
}

export const GitLabTransformation = ({
  plugin,
  connectionId,
  entities,
  transformation,
  setTransformation,
  setHasError,
}: Props) => {
  const [useCustom, setUseCustom] = useState(false);

  useEffect(() => {
    if (transformation.deploymentPattern || transformation.productionPattern) {
      setUseCustom(true);
    } else {
      setUseCustom(false);
    }
  }, []);

  useEffect(() => {
    setHasError(useCustom && !transformation.deploymentPattern);
  }, [useCustom, transformation]);

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
  transformation,
  onChangeTransformation,
  useCustom,
  onChangeUseCustom,
}: {
  plugin: string;
  connectionId: ID;
  entities: string[];
  panelStyle: React.CSSProperties;
  transformation: any;
  onChangeTransformation: any;
  useCustom: boolean;
  onChangeUseCustom: any;
}) =>
  [
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
            <ExternalLink link={DOC_URL.PLUGIN.GITLAB.TRANSFORMATION}>Learn more</ExternalLink>
          </p>
          <Checkbox disabled checked>
            Convert a GitLab Deployment to a DevLake Deployment
          </Checkbox>
          <div style={{ margin: '8px 0', paddingLeft: 28 }}>
            <span>If the environment</span>
            <Deployments
              style={{ margin: '0 4px' }}
              plugin={plugin}
              connectionId={connectionId}
              transformation={transformation}
              setTransformation={onChangeTransformation}
            />
            <span>, this Deployment is a ‘Production Deployment’</span>
          </div>
          <Checkbox checked={useCustom} onChange={onChangeUseCustom}>
            Convert a GitLab Pipeline as a DevLake Deployment when:
          </Checkbox>
          {useCustom && (
            <div style={{ paddingLeft: 28 }}>
              <ShowMore
                text={<p>Select this option only if you are not enabling GitLab Deployments. </p>}
                btnText="See how to configure"
              >
                <WorkflowRun />
              </ShowMore>
              <div style={{ margin: '8px 0' }}>
                <span>
                  Its branch/tag name or <strong>one of its jobs</strong> matches
                </span>
                <Input
                  style={{ width: 180, margin: '0 8px' }}
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
                <HelpTooltip content="GitLab Pipelines: https://docs.gitlab.com/ee/ci/pipelines/" />
              </div>
              <div style={{ margin: '8px 0' }}>
                <span>If the name also matches</span>
                <Input
                  style={{ width: 180, margin: '0 8px' }}
                  placeholder="(?i)(prod|release)"
                  value={transformation.productionPattern ?? ''}
                  onChange={(e) =>
                    onChangeTransformation({
                      ...transformation,
                      productionPattern: e.target.value,
                    })
                  }
                />
                <span>, this deployment is a ‘Production Deployment’</span>
                <HelpTooltip content="If you leave this field empty, all Deployments will be tagged as in the Production environment. " />
              </div>
              <CheckMatchedItems plugin={plugin} connectionId={connectionId} transformation={transformation} />
            </div>
          )}
        </>
      ),
    },
  ].filter((it) => entities.includes(it.key));
