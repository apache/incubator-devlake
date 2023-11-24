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

import { ExternalLink, HelpTooltip } from '@/components';
import { DOC_URL } from '@/release';

interface Props {
  entities: string[];
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
}

export const BambooTransformation = ({ entities, transformation, setTransformation }: Props) => {
  const [useCustom, setUseCustom] = useState(false);

  useEffect(() => {
    if (transformation.deploymentPattern || transformation.productionPattern) {
      setUseCustom(true);
    } else {
      setUseCustom(false);
    }
  }, [transformation]);

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
  transformation,
  onChangeTransformation,
  useCustom,
  onChangeUseCustom,
}: {
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
            <ExternalLink link={DOC_URL.PLUGIN.BAMBOO.TRANSFORMATION}>Learn more</ExternalLink>
          </p>
          <Checkbox disabled checked>
            Convert a Bamboo Deployment to a DevLake Deployment
          </Checkbox>
          <div style={{ margin: '8px 0', paddingLeft: 28 }}>
            <span>If its environment name matches</span>
            <Input
              style={{ width: 180, margin: '0 8px' }}
              placeholder="(?i)prod(.*)"
              value={transformation.envNamePattern}
              onChange={(e) =>
                onChangeTransformation({
                  ...transformation,
                  envNamePattern: e.target.value,
                })
              }
            />
            <span>, this deployment is a ‘Production Deployment’</span>
          </div>
          <Checkbox checked={useCustom} onChange={onChangeUseCustom}>
            Convert a Bamboo Plan Build to a DevLake Deployment when its name or one of its job builds’ names
          </Checkbox>
          <div style={{ margin: '8px 0', paddingLeft: 28 }}>
            <span>matches</span>
            <Input
              style={{ width: 180, margin: '0 8px' }}
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
            <HelpTooltip content="View your Bamboo Builds: https://confluence.atlassian.com/bamboo/viewing-a-plan-s-build-information-289276861.html" />
          </div>
          <div style={{ margin: '8px 0', paddingLeft: 28 }}>
            <span>If the name also matches</span>
            <Input
              style={{ width: 180, margin: '0 8px' }}
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
  ].filter((it) => entities.includes(it.key));
