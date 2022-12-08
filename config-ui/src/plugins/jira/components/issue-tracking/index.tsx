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

import React, { useState, useEffect, useMemo } from 'react'
import { Tag, Intent } from '@blueprintjs/core'

import { MultiSelector, Selector } from '@/components'

import type {
  UseIssueTrackingProps,
  IssueTypeItem,
  FieldItem
} from './use-issue-tracking'
import { useIssueTracking } from './use-issue-tracking'
import * as S from './styled'

enum StandardType {
  Requirement = 'REQUIREMENT',
  Bug = 'BUG',
  Incident = 'INCIDENT'
}

interface Props extends UseIssueTrackingProps {
  transformation: any
  setTransformation: React.Dispatch<any>
}

export const IssueTracking = ({
  connectionId,
  transformation,
  setTransformation
}: Props) => {
  const [requirements, setRequirements] = useState<IssueTypeItem['name'][]>([])
  const [bugs, setBugs] = useState<IssueTypeItem['name'][]>([])
  const [incidents, setIncidents] = useState<IssueTypeItem['name'][]>([])

  const { issueTypes, fields } = useIssueTracking({ connectionId })

  useEffect(() => {
    const types = Object.entries(transformation.typeMappings ?? {}).map(
      ([key, value]: any) => ({ name: key, ...value })
    )

    setRequirements(
      types
        .filter((it) => it.standardType === StandardType.Requirement)
        .map((it) => it.name)
    )
    setBugs(
      types
        .filter((it) => it.standardType === StandardType.Bug)
        .map((it) => it.name)
    )
    setIncidents(
      types
        .filter((it) => it.standardType === StandardType.Incident)
        .map((it) => it.name)
    )
  }, [transformation])

  const [requirementItems, bugItems, incidentItems] = useMemo(() => {
    return [
      issueTypes.filter((it) => requirements.includes(it.name)),
      issueTypes.filter((it) => bugs.includes(it.name)),
      issueTypes.filter((it) => incidents.includes(it.name))
    ]
  }, [requirements, bugs, incidents, issueTypes])

  const transformaType = (its: IssueTypeItem[], standardType: StandardType) => {
    return its.reduce((acc, cur) => {
      acc[cur.name] = {
        standardType
      }
      return acc
    }, {} as any)
  }

  return (
    <S.Container>
      <h3>Issue Tracking</h3>
      <p>
        Convert your issue labels with each category to view metrics such as
        Requirement Lead Time, Bug Count, Mean Time to Recover, etc.
      </p>
      <S.Item>
        <span>Requirement</span>
        <MultiSelector
          items={issueTypes}
          disabledItems={[...bugItems, ...incidentItems]}
          getKey={(it) => it.id}
          getName={(it) => it.name}
          getIcon={(it) => it.iconUrl}
          selectedItems={requirementItems}
          onChangeItems={(selectedItems) =>
            setTransformation({
              ...transformation,
              typeMappings: {
                ...transformaType(selectedItems, StandardType.Requirement),
                ...transformaType(bugItems, StandardType.Bug),
                ...transformaType(incidentItems, StandardType.Incident)
              }
            })
          }
        />
      </S.Item>
      <S.Item>
        <span>Bug</span>
        <MultiSelector
          items={issueTypes}
          disabledItems={[...requirementItems, ...incidentItems]}
          getKey={(it) => it.id}
          getName={(it) => it.name}
          getIcon={(it) => it.iconUrl}
          selectedItems={bugItems}
          onChangeItems={(selectedItems) =>
            setTransformation({
              ...transformation,
              typeMappings: {
                ...transformaType(requirementItems, StandardType.Requirement),
                ...transformaType(selectedItems, StandardType.Bug),
                ...transformaType(incidentItems, StandardType.Incident)
              }
            })
          }
        />
      </S.Item>
      <S.Item>
        <span>
          Incident{' '}
          <Tag intent={Intent.PRIMARY} style={{ fontSize: '10px' }} minimal>
            DORA
          </Tag>
        </span>
        <MultiSelector
          items={issueTypes}
          disabledItems={[...requirementItems, ...bugItems]}
          getKey={(it) => it.id}
          getName={(it) => it.name}
          getIcon={(it) => it.iconUrl}
          selectedItems={incidentItems}
          onChangeItems={(selectedItems) =>
            setTransformation({
              ...transformation,
              typeMappings: {
                ...transformaType(requirementItems, StandardType.Requirement),
                ...transformaType(bugItems, StandardType.Bug),
                ...transformaType(selectedItems, StandardType.Incident)
              }
            })
          }
        />
      </S.Item>
      <S.Item>
        <span>Epic Link</span>
        <Selector
          items={fields}
          disabledItems={
            transformation.storyPoint
              ? fields.filter((it) => it.id === transformation.storyPointField)
              : []
          }
          getKey={(it) => it.id}
          getName={(it) => it.name}
          selectedItem={fields.find(
            (it) => it.id === transformation.epicKeyField
          )}
          onChangeItem={(selectedItem) =>
            setTransformation({
              ...transformation,
              epicKeyField: selectedItem.id
            })
          }
        />
      </S.Item>
      <S.Item>
        <span>Story Point</span>
        <Selector
          items={fields}
          disabledItems={
            transformation.epicKeyField
              ? fields.filter((it) => it.id === transformation.epicKeyField)
              : []
          }
          getKey={(it) => it.id}
          getName={(it) => it.name}
          selectedItem={fields.find(
            (it) => it.id === transformation.storyPointField
          )}
          onChangeItem={(selectedItem) =>
            setTransformation({
              ...transformation,
              storyPointField: selectedItem.id
            })
          }
        />
      </S.Item>
    </S.Container>
  )
}
