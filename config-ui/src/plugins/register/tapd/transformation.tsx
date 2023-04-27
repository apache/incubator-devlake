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

import React, { useEffect, useState } from 'react';
import { FormGroup, Intent, Tag } from '@blueprintjs/core';

import { HelpTooltip, MultiSelector, PageLoading } from '@/components';
import { useProxyPrefix, useRefreshData } from '@/hooks';

import * as API from './api';
import * as S from './styled';
import { uniqWith } from 'lodash';

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
  connectionId: ID;
  scopeId: ID;
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
}

export const TapdTransformation = ({ connectionId, scopeId, transformation, setTransformation }: Props) => {
  const [featureTypeList, setFeatureTypeList] = useState<string[]>([]);
  const [bugTypeList, setBugTypeList] = useState<string[]>([]);
  const [incidentTypeList, setIncidentTypeList] = useState<string[]>([]);
  const [todoStatusList, setTodoStatusList] = useState<string[]>([]);
  const [inProgressStatusList, setInProgressStatusList] = useState<string[]>([]);
  const [doneStatusList, setDoneStatusList] = useState<string[]>([]);

  const prefix = useProxyPrefix({ plugin: 'tapd', connectionId });

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
      API.getStoryType(prefix, scopeId),
      { BUG: 'bug' } as Record<string, string>,
      { TASK: 'task' } as Record<string, string>,
      API.getStatus(prefix, scopeId, 'story'),
      API.getStatus(prefix, scopeId, 'bug'),
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
  return (
    <S.TransformationWrapper>
      {/* Issue Tracking */}
      <div className="issue-tracking">
        <h2>Issue Tracking</h2>
        <div className="issue-type">
          <div className="title">
            <span>Issue Type Mapping</span>
            <HelpTooltip content="Standardize your issue types to the following issue types to view metrics such as `Requirement lead time` and `Bug age` in built-in dashboards." />
          </div>
          <div className="list">
            <FormGroup inline label="Requirement">
              <MultiSelector
                items={typeList}
                disabledItems={typeList.filter((v) => [...bugTypeList, ...incidentTypeList].includes(v.id))}
                getKey={(it) => it.id}
                getName={(it) => it.name}
                selectedItems={typeList.filter((v) => featureTypeList.includes(v.id))}
                onChangeItems={(selectedItems) =>
                  setTransformation({
                    ...transformation,
                    typeMappings: {
                      ...transformaType(
                        selectedItems.map((v) => v.id),
                        StandardType.Requirement,
                      ),
                      ...transformaType(bugTypeList, StandardType.Bug),
                      ...transformaType(incidentTypeList, StandardType.Incident),
                    },
                  })
                }
              />
            </FormGroup>
            <FormGroup inline label="Bug">
              <MultiSelector
                items={typeList}
                disabledItems={typeList.filter((v) => [...featureTypeList, ...incidentTypeList].includes(v.id))}
                getKey={(it) => it.id}
                getName={(it) => it.name}
                selectedItems={typeList.filter((v) => bugTypeList.includes(v.id))}
                onChangeItems={(selectedItems) =>
                  setTransformation({
                    ...transformation,
                    typeMappings: {
                      ...transformaType(featureTypeList, StandardType.Requirement),
                      ...transformaType(
                        selectedItems.map((v) => v.id),
                        StandardType.Bug,
                      ),
                      ...transformaType(incidentTypeList, StandardType.Incident),
                    },
                  })
                }
              />
            </FormGroup>
            <FormGroup
              inline
              label={
                <>
                  <span>Incident</span>
                  <Tag intent={Intent.PRIMARY} style={{ marginLeft: 4 }}>
                    DORA
                  </Tag>
                </>
              }
            >
              <MultiSelector
                items={typeList}
                disabledItems={typeList.filter((v) => [...featureTypeList, ...bugTypeList].includes(v.id))}
                getKey={(it) => it.id}
                getName={(it) => it.name}
                selectedItems={typeList.filter((v) => incidentTypeList.includes(v.id))}
                onChangeItems={(selectedItems) =>
                  setTransformation({
                    ...transformation,
                    typeMappings: {
                      ...transformaType(featureTypeList, StandardType.Requirement),
                      ...transformaType(bugTypeList, StandardType.Bug),
                      ...transformaType(
                        selectedItems.map((v) => v.id),
                        StandardType.Incident,
                      ),
                    },
                  })
                }
              />
            </FormGroup>
          </div>
        </div>
        <div className="issue-status">
          <div className="title">
            <span>Issue Status Mapping</span>
            <HelpTooltip content="Standardize your issue statuses to the following issue statuses to view metrics such as `Requirement Delivery Rate` in built-in dashboards." />
          </div>
          <div className="list">
            <FormGroup inline label="TODO">
              <MultiSelector
                items={statusList}
                disabledItems={statusList.filter((v) => [...inProgressStatusList, ...doneStatusList].includes(v.name))}
                getKey={(it) => it.id}
                getName={(it) => it.name}
                selectedItems={statusList.filter((v) => todoStatusList.includes(v.name))}
                onChangeItems={(selectedItems) =>
                  setTransformation({
                    ...transformation,
                    statusMappings: {
                      ...transformaType(
                        selectedItems.map((v) => v.name),
                        StandardStatus.Todo,
                      ),
                      ...transformaType(inProgressStatusList, StandardStatus.InProgress),
                      ...transformaType(doneStatusList, StandardStatus.Done),
                    },
                  })
                }
              />
            </FormGroup>
            <FormGroup inline label="IN-PROGRESS">
              <MultiSelector
                items={statusList}
                disabledItems={statusList.filter((v) => [...todoStatusList, ...doneStatusList].includes(v.name))}
                getKey={(it) => it.id}
                getName={(it) => it.name}
                selectedItems={statusList.filter((v) => inProgressStatusList.includes(v.name))}
                onChangeItems={(selectedItems) =>
                  setTransformation({
                    ...transformation,
                    statusMappings: {
                      ...transformaType(todoStatusList, StandardStatus.Todo),
                      ...transformaType(
                        selectedItems.map((v) => v.name),
                        StandardStatus.InProgress,
                      ),
                      ...transformaType(doneStatusList, StandardStatus.Done),
                    },
                  })
                }
              />
            </FormGroup>
            <FormGroup inline label="DONE">
              <MultiSelector
                items={statusList}
                disabledItems={statusList.filter((v) => [...todoStatusList, ...inProgressStatusList].includes(v.name))}
                getKey={(it) => it.id}
                getName={(it) => it.name}
                selectedItems={statusList.filter((v) => doneStatusList.includes(v.name))}
                onChangeItems={(selectedItems) =>
                  setTransformation({
                    ...transformation,
                    statusMappings: {
                      ...transformaType(todoStatusList, StandardStatus.Todo),
                      ...transformaType(inProgressStatusList, StandardStatus.InProgress),
                      ...transformaType(
                        selectedItems.map((v) => v.name),
                        StandardStatus.Done,
                      ),
                    },
                  })
                }
              />
            </FormGroup>
          </div>
        </div>
      </div>
    </S.TransformationWrapper>
  );
};
