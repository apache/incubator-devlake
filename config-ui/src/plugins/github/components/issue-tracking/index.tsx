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

import React from 'react'
import { Tag, FormGroup, InputGroup, Intent } from '@blueprintjs/core'

interface Props {
  transformation: any
  setTransformation: React.Dispatch<React.SetStateAction<{}>>
}

export const IssueTracking = ({
  transformation,
  setTransformation
}: Props) => {
  return (
    <>
      <h3>
        <span>Issue Tracking</span>
        <Tag minimal>RegExp</Tag>
      </h3>
      <p>
        Map your issue labels with each category to view corresponding metrics
        in the dashboard.
      </p>
      <FormGroup inline label='Severity'>
        <InputGroup
          placeholder='severity/(.*)$'
          value={transformation.issueSeverity}
          onChange={(e) =>
            setTransformation({
              ...transformation,
              issueSeverity: e.target.value
            })
          }
        />
      </FormGroup>
      <FormGroup inline label='Component'>
        <InputGroup
          placeholder='component/(.*)$'
          value={transformation.issueComponent}
          onChange={(e) =>
            setTransformation({
              ...transformation,
              issueComponent: e.target.value
            })
          }
        />
      </FormGroup>
      <FormGroup inline label='Priority'>
        <InputGroup
          placeholder='(highest|high|medium|low)$'
          value={transformation.issuePriority}
          onChange={(e) =>
            setTransformation({
              ...transformation,
              issuePriority: e.target.value
            })
          }
        />
      </FormGroup>
      <FormGroup inline label='Type/Requirement'>
        <InputGroup
          placeholder='(feat|feature|proposal|requirement)$'
          value={transformation.issueTypeRequirement}
          onChange={(e) =>
            setTransformation({
              ...transformation,
              issueTypeRequirement: e.target.value
            })
          }
        />
      </FormGroup>
      <FormGroup inline label='Type/Bug'>
        <InputGroup
          placeholder='(bug|broken)$'
          value={transformation.issueTypeBug}
          onChange={(e) =>
            setTransformation({
              ...transformation,
              issueTypeBug: e.target.value
            })
          }
        />
      </FormGroup>
      <FormGroup
        inline
        label={
          <span>
            Type/Incident
            <Tag
              minimal
              intent={Intent.PRIMARY}
              style={{ marginLeft: 4, fontSize: 10 }}
            >
              DORA
            </Tag>
          </span>
        }
      >
        <InputGroup
          placeholder='(incident|p0|p1|p2)$'
          value={transformation.issueTypeIncident}
          onChange={(e) =>
            setTransformation({
              ...transformation,
              issueTypeIncident: e.target.value
            })
          }
        />
      </FormGroup>
    </>
  )
}
