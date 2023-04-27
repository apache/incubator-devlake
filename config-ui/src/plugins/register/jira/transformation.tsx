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

import React, { useState, useEffect, useMemo } from 'react';
import { FormGroup, InputGroup, Tag, Icon, Intent } from '@blueprintjs/core';
import { Popover2 } from '@blueprintjs/popover2';
import { uniqWith } from 'lodash';

import { PageLoading, HelpTooltip, ExternalLink, MultiSelector, Selector, Divider } from '@/components';
import { useProxyPrefix, useRefreshData } from '@/hooks';

import * as API from './api';
import * as S from './styled';

enum StandardType {
  Requirement = 'REQUIREMENT',
  Bug = 'BUG',
  Incident = 'INCIDENT',
}

interface Props {
  connectionId: ID;
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
}

export const JiraTransformation = ({ connectionId, transformation, setTransformation }: Props) => {
  const [requirements, setRequirements] = useState<string[]>([]);
  const [bugs, setBugs] = useState<string[]>([]);
  const [incidents, setIncidents] = useState<string[]>([]);

  const prefix = useProxyPrefix({ plugin: 'jira', connectionId });

  const { ready, data } = useRefreshData<{
    issueTypes: Array<{
      id: string;
      name: string;
      iconUrl: string;
    }>;
    fields: Array<{
      id: string;
      name: string;
    }>;
  }>(async () => {
    if (!prefix) {
      return {
        issueTypes: [],
        fields: [],
      };
    }

    const [issueTypes, fields] = await Promise.all([API.getIssueType(prefix), API.getField(prefix)]);
    return {
      issueTypes: uniqWith(issueTypes, (it, oit) => it.name === oit.name),
      fields,
    };
  }, [prefix]);

  useEffect(() => {
    const types = Object.entries(transformation.typeMappings ?? {}).map(([key, value]: any) => ({
      name: key,
      ...value,
    }));

    setRequirements(types.filter((it) => it.standardType === StandardType.Requirement).map((it) => it.name));
    setBugs(types.filter((it) => it.standardType === StandardType.Bug).map((it) => it.name));
    setIncidents(types.filter((it) => it.standardType === StandardType.Incident).map((it) => it.name));
  }, [transformation]);

  const [requirementItems, bugItems, incidentItems] = useMemo(() => {
    return [
      (data?.issueTypes ?? []).filter((it) => requirements.includes(it.name)),
      (data?.issueTypes ?? []).filter((it) => bugs.includes(it.name)),
      (data?.issueTypes ?? []).filter((it) => incidents.includes(it.name)),
    ];
  }, [requirements, bugs, incidents, data?.issueTypes]);

  if (!ready || !data) {
    return <PageLoading />;
  }

  const { issueTypes, fields } = data;

  const transformaType = (
    its: Array<{
      id: string;
      name: string;
      iconUrl: string;
    }>,
    standardType: StandardType,
  ) => {
    return its.reduce((acc, cur) => {
      acc[cur.name] = {
        standardType,
      };
      return acc;
    }, {} as any);
  };
  return (
    <S.TransformationWrapper>
      {/* Issue Tracking */}
      <div className="issue-tracking">
        <h2>Issue Tracking</h2>
        <p>
          Tell DevLake what types of Jira issues you are using as features, bugs and incidents, and what field as `Epic
          Link` or `Story Points`.
        </p>
        <div className="issue-type">
          <div className="title">
            <span>Issue Type</span>
            <Popover2
              position="top"
              content={
                <div style={{ padding: '8px 12px', color: '#ffffff', backgroundColor: 'rgba(0,0,0,.8)' }}>
                  DevLake defines three standard types of issues: FEATURE, BUG and INCIDENT. Standardize your Jira issue
                  types to these three types so that DevLake can calculate metrics such as{' '}
                  <ExternalLink link="https://devlake.apache.org/docs/Metrics/RequirementLeadTime">
                    Requirement Lead Time
                  </ExternalLink>
                  , <ExternalLink link="https://devlake.apache.org/docs/Metrics/BugAge">Bug Age</ExternalLink>,
                  <ExternalLink link="https://devlake.apache.org/docs/Metrics/MTTR">
                    DORA - Median Time to Restore Service
                  </ExternalLink>
                  , etc.
                </div>
              }
            >
              <Icon icon="help" size={12} color="#94959f" style={{ marginLeft: 4, cursor: 'pointer' }} />
            </Popover2>
          </div>
          <div className="list">
            <FormGroup inline label="Requirement">
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
                      ...transformaType(incidentItems, StandardType.Incident),
                    },
                  })
                }
              />
            </FormGroup>
            <FormGroup inline label="Bug">
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
                      ...transformaType(incidentItems, StandardType.Incident),
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
                      ...transformaType(selectedItems, StandardType.Incident),
                    },
                  })
                }
              />
            </FormGroup>
          </div>
        </div>
        <FormGroup
          inline
          label={
            <>
              <span>Story Points</span>
              <HelpTooltip content="Choose the issue field you are using as `Story Points`." />
            </>
          }
        >
          <Selector
            items={fields}
            getKey={(it) => it.id}
            getName={(it) => it.name}
            selectedItem={fields.find((it) => it.id === transformation.storyPointField)}
            onChangeItem={(selectedItem) =>
              setTransformation({
                ...transformation,
                storyPointField: selectedItem.id,
              })
            }
          />
        </FormGroup>
      </div>
      <Divider />
      {/* Cross-domain */}
      <div>
        <h2>Cross-domain</h2>
        <p>
          Connect `commits` and `issues` to measure metrics such as{' '}
          <ExternalLink link="https://devlake.apache.org/docs/Metrics/BugCountPer1kLinesOfCode">
            Bug Count per 1k Lines of Code
          </ExternalLink>{' '}
          or man hour distribution on different work types.
        </p>
        <FormGroup
          inline
          label={
            <>
              <span>Connect GitLab Commits and Jira Issues</span>
              <HelpTooltip
                content={
                  <div>
                    If you are using GitLab’s{' '}
                    <ExternalLink link="https://docs.gitlab.com/ee/integration/jira/">Jira integration</ExternalLink>,
                    specify the commit SHA pattern. DevLake will parse the commit_sha from your Jira issues’ remote/web
                    links and store the relationship in the table `issue_commits`.
                  </div>
                }
              />
            </>
          }
        >
          <InputGroup
            fill
            placeholder="/commit/([0-9a-f]{40})$"
            value={transformation.remotelinkCommitShaPattern ?? ''}
            onChange={(e) =>
              setTransformation({
                ...transformation,
                remotelinkCommitShaPattern: e.target.value,
              })
            }
          />
        </FormGroup>
      </div>
    </S.TransformationWrapper>
  );
};
