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

import { useEffect, useState } from 'react';
import { uniqWith } from 'lodash';
import { CaretRightOutlined } from '@ant-design/icons';
import { theme, Collapse, Tag, Form, Select } from 'antd';

import API from '@/api';
import { PageLoading } from '@/components';
import { useProxyPrefix, useRefreshData } from '@/hooks';

enum StandardType {
  Requirement = 'Requirement',
  Bug = 'BUG',
  Incident = 'INCIDENT',
}

enum StandardStatus {
  Todo = 'TODO',
  InProgress = 'IN-PROGRESS',
  Done = 'DONE',
}

interface Props {
  entities: string[];
  connectionId: ID;
  scopeId: ID;
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
}

export const TapdTransformation = ({ entities, connectionId, scopeId, transformation, setTransformation }: Props) => {
  const [featureTypeList, setFeatureTypeList] = useState<string[]>([]);
  const [bugTypeList, setBugTypeList] = useState<string[]>([]);
  const [incidentTypeList, setIncidentTypeList] = useState<string[]>([]);
  const [todoStatusList, setTodoStatusList] = useState<string[]>([]);
  const [inProgressStatusList, setInProgressStatusList] = useState<string[]>([]);
  const [doneStatusList, setDoneStatusList] = useState<string[]>([]);

  const prefix = useProxyPrefix({ plugin: 'tapd', connectionId });

  const { token } = theme.useToken();

  const { ready, data } = useRefreshData<{
    statusList: Array<{
      id: string;
      name: string;
    }>;
    typeList: Array<{
      id: string;
      name: string;
    }>;
  }>(async () => {
    if (!prefix) {
      return {
        statusList: [],
        typeList: [],
      };
    }

    const [storyType, bugType, taskType, storyStatus, bugStatus, taskStatus] = await Promise.all([
      API.plugin.tapd.storyCategories(prefix, scopeId),
      { BUG: 'bug' } as Record<string, string>,
      { TASK: 'task' } as Record<string, string>,
      API.plugin.tapd.statusMap(prefix, scopeId, 'story'),
      API.plugin.tapd.statusMap(prefix, scopeId, 'bug'),
      { open: 'task-open', progressing: 'task-progressing', done: 'task-done' } as Record<string, string>,
    ]);

    const statusList: { id: string; name: string }[] = uniqWith(
      [
        { id: 'open', name: taskStatus.open },
        { id: 'progressing', name: taskStatus.progressing },
        { id: 'done', name: taskStatus.done },
        ...(Object.values(storyStatus.data) as string[]).map((it) => ({ id: it, name: it })),
        ...(Object.values(bugStatus.data) as string[]).map((it) => ({ id: it, name: it })),
      ],
      (a, b) => a.id === b.id,
    );

    const typeList: { id: string; name: string }[] = [
      ...storyType.data.map((it: any) => ({ id: it.Category.id, name: it.Category.name })),
      { id: 'BUG', name: bugType['BUG'] },
      { id: 'TASK', name: taskType['TASK'] },
    ];

    return {
      statusList,
      typeList,
    };
  }, [prefix]);

  useEffect(() => {
    const typeList = Object.entries(transformation.typeMappings ?? {}).map(([key, value]: any) => ({ key, value }));
    setFeatureTypeList(typeList.filter((it) => it.value === StandardType.Requirement).map((it) => it.key));
    setBugTypeList(typeList.filter((it) => it.value === StandardType.Bug).map((it) => it.key));
    setIncidentTypeList(typeList.filter((it) => it.value === StandardType.Incident).map((it) => it.key));

    const statusList = Object.entries(transformation.statusMappings ?? {}).map(([key, value]: any) => ({ key, value }));
    setTodoStatusList(statusList.filter((it) => it.value === StandardStatus.Todo).map((it) => it.key));
    setInProgressStatusList(statusList.filter((it) => it.value === StandardStatus.InProgress).map((it) => it.key));
    setDoneStatusList(statusList.filter((it) => it.value === StandardStatus.Done).map((it) => it.key));
  }, [transformation]);

  if (!ready || !data) {
    return <PageLoading />;
  }

  const { statusList, typeList } = data;

  const transformaType = (its: string[], standardType: string) => {
    return its.reduce((acc, cur) => {
      acc[cur] = standardType;
      return acc;
    }, {} as Record<string, string>);
  };

  const panelStyle: React.CSSProperties = {
    marginBottom: 24,
    background: token.colorFillAlter,
    borderRadius: token.borderRadiusLG,
    border: 'none',
  };

  return (
    <Collapse
      bordered={false}
      defaultActiveKey={['TICKET', 'CROSS']}
      expandIcon={({ isActive }) => <CaretRightOutlined rotate={isActive ? 90 : 0} />}
      style={{ background: token.colorBgContainer }}
      size="large"
      items={renderCollapseItems({
        entities,
        panelStyle,
        transformation,
        onChangeTransformation: setTransformation,
        typeList,
        featureTypeList,
        bugTypeList,
        incidentTypeList,
        statusList,
        todoStatusList,
        inProgressStatusList,
        doneStatusList,
        transformaType,
      })}
    />
  );
};

const renderCollapseItems = ({
  entities,
  panelStyle,
  transformation,
  onChangeTransformation,
  typeList,
  featureTypeList,
  bugTypeList,
  incidentTypeList,
  statusList,
  todoStatusList,
  inProgressStatusList,
  doneStatusList,
  transformaType,
}: {
  entities: string[];
  panelStyle: React.CSSProperties;
  transformation: any;
  onChangeTransformation: any;
  typeList: Array<{
    id: string;
    name: string;
  }>;
  featureTypeList: any;
  bugTypeList: any;
  incidentTypeList: any;
  statusList: Array<{
    id: string;
    name: string;
  }>;
  todoStatusList: any;
  inProgressStatusList: any;
  doneStatusList: any;
  transformaType: any;
}) =>
  [
    {
      key: 'TICKET',
      label: 'Issue Tracking',
      style: panelStyle,
      children: (
        <Form labelCol={{ span: 5 }}>
          <p>
            Standardize your issue types to the following issue types to view metrics such as `Requirement lead time`
            and `Bug age` in built-in dashboards.
          </p>
          <Form.Item label="Requirement">
            <Select
              mode="multiple"
              options={typeList.map((it) => ({ label: it.name, value: it.id }))}
              value={featureTypeList}
              onChange={(value) =>
                onChangeTransformation({
                  ...transformation,
                  typeMappings: {
                    ...transformaType(value, StandardType.Requirement),
                    ...transformaType(bugTypeList, StandardType.Bug),
                    ...transformaType(incidentTypeList, StandardType.Incident),
                  },
                })
              }
            />
          </Form.Item>
          <Form.Item label="Bug">
            <Select
              mode="multiple"
              options={typeList.map((it) => ({ label: it.name, value: it.id }))}
              value={bugTypeList}
              onChange={(value) =>
                onChangeTransformation({
                  ...transformation,
                  typeMappings: {
                    ...transformaType(featureTypeList, StandardType.Requirement),
                    ...transformaType(value, StandardType.Bug),
                    ...transformaType(incidentTypeList, StandardType.Incident),
                  },
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
            <Select
              mode="multiple"
              options={typeList.map((it) => ({ label: it.name, value: it.id }))}
              value={incidentTypeList}
              onChange={(value) =>
                onChangeTransformation({
                  ...transformation,
                  typeMappings: {
                    ...transformaType(featureTypeList, StandardType.Requirement),
                    ...transformaType(bugTypeList, StandardType.Bug),
                    ...transformaType(value, StandardType.Incident),
                  },
                })
              }
            />
          </Form.Item>
          <p>
            Standardize your issue statuses to the following issue statuses to view metrics such as `Requirement
            Delivery Rate` in built-in dashboards.
          </p>
          <Form.Item label="TODO">
            <Select
              mode="multiple"
              options={statusList.map((it) => ({ label: it.name, value: it.id }))}
              value={todoStatusList}
              onChange={(value) =>
                onChangeTransformation({
                  ...transformation,
                  statusMappings: {
                    ...transformaType(value, StandardStatus.Todo),
                    ...transformaType(inProgressStatusList, StandardStatus.InProgress),
                    ...transformaType(doneStatusList, StandardStatus.Done),
                  },
                })
              }
            />
          </Form.Item>
          <Form.Item label="IN-PROGRESS">
            <Select
              mode="multiple"
              options={statusList.map((it) => ({ label: it.name, value: it.id }))}
              value={inProgressStatusList}
              onChange={(value) =>
                onChangeTransformation({
                  ...transformation,
                  statusMappings: {
                    ...transformaType(todoStatusList, StandardStatus.Todo),
                    ...transformaType(value, StandardStatus.InProgress),
                    ...transformaType(doneStatusList, StandardStatus.Done),
                  },
                })
              }
            />
          </Form.Item>
          <Form.Item label="DONE">
            <Select
              mode="multiple"
              options={statusList.map((it) => ({ label: it.name, value: it.id }))}
              value={doneStatusList}
              onChange={(value) =>
                onChangeTransformation({
                  ...transformation,
                  statusMappings: {
                    ...transformaType(todoStatusList, StandardStatus.Todo),
                    ...transformaType(inProgressStatusList, StandardStatus.InProgress),
                    ...transformaType(value, StandardStatus.Done),
                  },
                })
              }
            />
          </Form.Item>
        </Form>
      ),
    },
  ].filter((it) => entities.includes(it.key));
