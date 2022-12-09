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
import {
  Tag,
  FormGroup,
  InputGroup,
  TextArea,
  Icon,
  Colors,
  Position
} from '@blueprintjs/core'
import { Popover2 } from '@blueprintjs/popover2'

import * as S from './styled'

interface Props {
  transformation: any
  setTransformation: React.Dispatch<React.SetStateAction<any>>
}

export const CodeReview = ({
  transformation,
  setTransformation
}: Props) => {
  return (
    <>
      <h3>
        <span>Code Review</span>
        <Tag minimal>RegExp</Tag>
      </h3>
      <p>
        Map your pull requests labels with each category to view corresponding
        metrics in the dashboard.
      </p>
      <FormGroup inline label='Type'>
        <InputGroup
          placeholder='type/(.*)$'
          value={transformation.prType}
          onChange={(e) =>
            setTransformation({ ...transformation, prType: e.target.value })
          }
        />
      </FormGroup>
      <FormGroup inline label='Component'>
        <InputGroup
          placeholder='component/(.*)$'
          value={transformation.prComponent}
          onChange={(e) =>
            setTransformation({
              ...transformation,
              prComponent: e.target.value
            })
          }
        />
      </FormGroup>
      <h3>
        <span> PR-Issue Mapping</span>
        <Tag minimal>RegExp</Tag>
      </h3>
      <p>
        Extract the issue numbers closed by pull requests. The issue numbers are
        parsed from PR bodies that meet the following RegEx.
      </p>
      <FormGroup
        inline
        label={
          <span>
            PR Body Pattern
            <Popover2
              position={Position.TOP}
              content={
                <S.Tips>
                  <p>
                    <Icon
                      icon='tick-circle'
                      size={10}
                      color={Colors.GREEN4}
                      style={{ marginRight: '4px' }}
                    />
                    Example 1: PR #321 body contains "
                    <strong>Closes #1234</strong>" (PR #321 and issue #1234 will
                    be mapped by the following RegEx)
                  </p>
                  <p>
                    <Icon
                      icon='delete'
                      size={10}
                      color={Colors.RED4}
                      style={{ marginRight: '4px' }}
                    />
                    Example 2: PR #321 body contains "
                    <strong>Related to #1234</strong>" (PR #321 and issue #1234
                    will NOT be mapped by the following RegEx)
                  </p>
                </S.Tips>
              }
            >
              <Icon
                icon='help'
                size={12}
                color={Colors.GRAY3}
                style={{ marginLeft: '4px', marginBottom: '4px' }}
              />
            </Popover2>
          </span>
        }
      >
        <TextArea
          value={transformation.prBodyClosePattern}
          placeholder='(?mi)(fix|close|resolve|fixes|closes|resolves|fixed|closed|resolved)[\s]*.*(((and )?(#|https:\/\/github.com\/%s\/%s\/issues\/)\d+[ ]*)+)'
          onChange={(e) =>
            setTransformation({
              ...transformation,
              prBodyClosePattern: e.target.value
            })
          }
          fill
          rows={2}
        />
      </FormGroup>
    </>
  )
}
