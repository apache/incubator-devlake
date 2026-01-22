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
import { theme, Collapse, Tag, DatePicker, InputNumber, Alert } from 'antd';
import dayjs from 'dayjs';

import { HelpTooltip } from '@/components';

interface Props {
  entities: string[];
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
}

export const GhCopilotTransformation = ({ entities, transformation, setTransformation }: Props) => {
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
      defaultActiveKey={['COPILOT']}
      expandIcon={({ isActive }) => <CaretRightOutlined rotate={isActive ? 90 : 0} rev="" />}
      style={{ background: token.colorBgContainer }}
      size="large"
      items={renderCollapseItems({
        entities,
        panelStyle,
        transformation,
        onChangeTransformation: setTransformation,
        token,
      })}
    />
  );
};

const renderCollapseItems = ({
  entities,
  panelStyle,
  transformation,
  onChangeTransformation,
  token,
}: {
  entities: string[];
  panelStyle: React.CSSProperties;
  transformation: any;
  onChangeTransformation: any;
  token: any;
}) =>
  [
    {
      key: 'COPILOT',
      label: 'GitHub Copilot Impact Analysis',
      style: panelStyle,
      children: (
        <>
          <h3 style={{ marginBottom: 16 }}>
            <span>Rollout Milestone Configuration</span>
            <Tag style={{ marginLeft: 4 }} color="green">
              Impact Dashboard
            </Tag>
          </h3>
          <Alert
            style={{ marginBottom: 16 }}
            type="info"
            showIcon
            message="Optional: Configure a rollout milestone date to add visual annotations on the GitHub Copilot Impact Dashboard. The dashboard works without this - correlation analysis is the primary view."
          />
          <div style={{ marginTop: 16, marginBottom: 8 }}>
            <strong>GitHub Copilot Rollout Date (Optional)</strong>
            <HelpTooltip content="The date when GitHub Copilot was rolled out to your team/organization. When set, this adds a visual annotation to the correlation charts. Leave empty if you prefer pure correlation analysis." />
          </div>
          <div style={{ margin: '8px 0', paddingLeft: 28 }}>
            <DatePicker
              style={{ width: 200 }}
              value={transformation.implementationDate ? dayjs(transformation.implementationDate) : null}
              placeholder="Select implementation date"
              onChange={(date) =>
                onChangeTransformation({
                  ...transformation,
                  implementationDate: date ? date.utc().format('YYYY-MM-DD[T]HH:mm:ssZ') : null,
                })
              }
            />
          </div>
          <div style={{ marginTop: 24, marginBottom: 8 }}>
            <strong>Baseline Period (Days)</strong>
            <HelpTooltip content="The number of days before the implementation date to use as the baseline for comparison. Default is 90 days." />
          </div>
          <div style={{ margin: '8px 0', paddingLeft: 28 }}>
            <InputNumber
              style={{ width: 200 }}
              min={7}
              max={365}
              value={transformation.baselinePeriodDays ?? 90}
              placeholder="90"
              addonAfter="days"
              onChange={(value) =>
                onChangeTransformation({
                  ...transformation,
                  baselinePeriodDays: value ?? 90,
                })
              }
            />
            <span style={{ marginLeft: 8, color: token.colorTextSecondary }}>
              (Recommended: 30-90 days for meaningful comparison)
            </span>
          </div>
        </>
      ),
    },
  ].filter((it) => entities.includes(it.key) || entities.includes('COPILOT'));
