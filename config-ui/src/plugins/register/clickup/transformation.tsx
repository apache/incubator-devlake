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
import { theme, Collapse, Form, Input, Select } from 'antd';

interface Props {
  entities: string[];
  connectionId: ID;
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
}

export const ClickUpTransformation = ({ transformation, setTransformation }: Props) => {
  const { token } = theme.useToken();

  const panelStyle: React.CSSProperties = {
    marginBottom: 24,
    background: token.colorFillAlter,
    borderRadius: token.borderRadiusLG,
    border: 'none',
  };

  const set = (patch: Record<string, any>) => setTransformation({ ...transformation, ...patch });

  return (
    <Collapse
      bordered={false}
      defaultActiveKey={['TICKET']}
      expandIcon={({ isActive }) => <CaretRightOutlined rotate={isActive ? 90 : 0} />}
      style={{ background: token.colorBgContainer }}
      size="large"
      items={[
        {
          key: 'TICKET',
          label: 'Issue Tracking',
          style: panelStyle,
          children: (
            <>
              <p style={{ marginBottom: 16 }}>
                Tell DevLake how the lists and tasks in this ClickUp folder map onto DevLake's domain model (boards,
                sprints, issues) so sprint velocity and DORA metrics are computed correctly.
              </p>

              <Form.Item
                label="Sprint list name pattern"
                extra="RegEx matched against each list's name. Matching lists become sprints; the rest are plain board issues."
              >
                <Input
                  placeholder="(?i)sprint\s*\d+"
                  value={transformation.sprintNamePattern ?? ''}
                  onChange={(e) => set({ sprintNamePattern: e.target.value })}
                />
              </Form.Item>

              <Form.Item
                label="Story point field"
                extra="Leave blank to use ClickUp's native sprint points. Set a custom-field name to read story points (e.g. Fibonacci LOE) from that field instead."
              >
                <Input
                  placeholder="points (native)"
                  value={transformation.storyPointField ?? ''}
                  onChange={(e) => set({ storyPointField: e.target.value })}
                />
              </Form.Item>

              <Form.Item
                label="Force issue type for this board"
                extra="Optional. Set to flag an entire folder as one type — e.g. INCIDENT for a Security Incidents folder (feeds DORA Change Failure Rate & MTTR). Leave as (auto) to detect per task."
              >
                <Select
                  allowClear
                  placeholder="(auto — detect per task)"
                  value={transformation.defaultIssueType || undefined}
                  onChange={(value) => set({ defaultIssueType: value ?? '' })}
                  options={[
                    { label: 'Requirement', value: 'REQUIREMENT' },
                    { label: 'Bug', value: 'BUG' },
                    { label: 'Incident', value: 'INCIDENT' },
                  ]}
                />
              </Form.Item>

              <p style={{ margin: '8px 0 16px' }}>
                Per-task type detection (used when the board type above is left on auto). DevLake classifies each task
                by matching its ClickUp task type against these RegEx patterns; precedence is Incident &gt; Bug &gt;
                Requirement.
              </p>
              <Form.Item label="Requirement">
                <Input
                  placeholder="(feat|feature|story|requirement)"
                  value={transformation.issueTypeRequirement ?? ''}
                  onChange={(e) => set({ issueTypeRequirement: e.target.value })}
                />
              </Form.Item>
              <Form.Item label="Bug">
                <Input
                  placeholder="(?i)^bug$"
                  value={transformation.issueTypeBug ?? ''}
                  onChange={(e) => set({ issueTypeBug: e.target.value })}
                />
              </Form.Item>
              <Form.Item label="Incident">
                <Input
                  placeholder="(incident|outage|sev)"
                  value={transformation.issueTypeIncident ?? ''}
                  onChange={(e) => set({ issueTypeIncident: e.target.value })}
                />
              </Form.Item>

              <p style={{ margin: '8px 0 16px' }}>
                Type by list name (optional). ClickUp often has no per-task type, so teams group bugs/incidents in a
                dedicated list. Tasks in a list whose name matches these RegEx are classified accordingly; Incident
                takes precedence over Bug. This overrides per-task detection above (but a forced board type still wins).
              </p>
              <Form.Item
                label="Bug list name pattern"
                extra="Tasks in matching lists become BUG — e.g. a QA Bugs / Bug Tracking list."
              >
                <Input
                  placeholder="(?i)bug"
                  value={transformation.bugListPattern ?? ''}
                  onChange={(e) => set({ bugListPattern: e.target.value })}
                />
              </Form.Item>
              <Form.Item
                label="Incident list name pattern"
                extra="Tasks in matching lists become INCIDENT (DORA CFR/MTTR)."
              >
                <Input
                  placeholder="(?i)incident"
                  value={transformation.incidentListPattern ?? ''}
                  onChange={(e) => set({ incidentListPattern: e.target.value })}
                />
              </Form.Item>

              <p style={{ margin: '8px 0 16px' }}>
                Status mapping (optional). ClickUp statuses are auto-mapped by their type (open/unstarted → To Do,
                custom → In Progress, done/closed → Done). List raw status names below to override.
              </p>
              <Form.Item label="To Do statuses">
                <Select
                  mode="tags"
                  placeholder="e.g. backlog, ready for dev"
                  value={transformation.issueStatusTodo ?? []}
                  onChange={(value) => set({ issueStatusTodo: value })}
                />
              </Form.Item>
              <Form.Item label="In Progress statuses">
                <Select
                  mode="tags"
                  placeholder="e.g. in development, in code review"
                  value={transformation.issueStatusInProgress ?? []}
                  onChange={(value) => set({ issueStatusInProgress: value })}
                />
              </Form.Item>
              <Form.Item label="Done statuses">
                <Select
                  mode="tags"
                  placeholder="e.g. deployed, Closed"
                  value={transformation.issueStatusDone ?? []}
                  onChange={(value) => set({ issueStatusDone: value })}
                />
              </Form.Item>
            </>
          ),
        },
      ]}
    />
  );
};
