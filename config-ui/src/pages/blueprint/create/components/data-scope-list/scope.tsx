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

import React, { useState, useEffect } from 'react'
import { groupBy } from 'lodash'

import type { BPScopeItemType } from '../../types'

import * as S from './styled'

type ScopeMap = Record<string, BPScopeItemType[]>

interface Props {
  groupByTs: boolean
  scope: BPScopeItemType[]
}

export const Scope = ({ groupByTs, scope }: Props) => {
  const [scopeTsMap, setScopeTsMap] = useState<ScopeMap>({})

  useEffect(() => {
    setScopeTsMap(
      groupBy(scope, (it) => it.transformationRuleName ?? 'No Transformation')
    )
  }, [scope])

  if (!scope.length) {
    return <span>No Data Scope Selected</span>
  }

  return (
    <S.ScopeList>
      {!groupByTs &&
        scope.map((sc) => <S.ScopeItem key={sc.id}>{sc.name}</S.ScopeItem>)}

      {groupByTs &&
        Object.keys(scopeTsMap).map((name) => (
          <S.ScopeItemMap key={name}>
            <div className='name'>
              <span>{name}</span>
              {name !== 'No Transformation' && (
                <div className='action'>
                  {/* <span onClick={() => onUpdateScope(scopeTsMap[name])}>
                    Remove
                  </span> */}
                </div>
              )}
            </div>
            <ul>
              {scopeTsMap[name].map((sc) => (
                <li key={sc.id}>{sc.name}</li>
              ))}
            </ul>
          </S.ScopeItemMap>
        ))}
    </S.ScopeList>
  )
}
