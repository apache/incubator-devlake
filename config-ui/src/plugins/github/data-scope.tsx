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

import type { ScopeItemType } from './types'
import { ScopeFromEnum } from './types'

import { MillerColumns, RepoSelector } from './components'

interface Props {
  connectionId: ID
  selectedItems: ScopeItemType[]
  onChangeItems: (selectedItems: ScopeItemType[]) => void
}

export const GitHubDataScope = ({
  connectionId,
  selectedItems,
  onChangeItems
}: Props) => {
  const handleChangeMillerColumnsItems = (sis: ScopeItemType[]) => {
    onChangeItems([
      ...selectedItems.filter((it) => it.from !== ScopeFromEnum.MILLER_COLUMNS),
      ...sis
    ])
  }

  const handleChangeRepoSelectorItems = (sis: ScopeItemType[]) => {
    onChangeItems([
      ...selectedItems.filter((it) => it.from !== ScopeFromEnum.REPO_SELECTOR),
      ...sis
    ])
  }

  return (
    <>
      <h3>Repositories *</h3>
      <p>Select the repositories you would like to sync.</p>
      <MillerColumns
        connectionId={connectionId}
        disabledItems={selectedItems.filter(
          (it) => it.from !== ScopeFromEnum.MILLER_COLUMNS
        )}
        selectedItems={selectedItems.filter(
          (it) => it.from === ScopeFromEnum.MILLER_COLUMNS
        )}
        onChangeItems={handleChangeMillerColumnsItems}
      />
      <h4>Add repositories outside of your organizations</h4>
      <p>Search for repositories and add to them</p>
      <RepoSelector
        connectionId={connectionId}
        disabledItems={selectedItems.filter(
          (it) => it.from !== ScopeFromEnum.REPO_SELECTOR
        )}
        selectedItems={selectedItems.filter(
          (it) => it.from === ScopeFromEnum.REPO_SELECTOR
        )}
        onChangeItems={handleChangeRepoSelectorItems}
      />
    </>
  )
}
