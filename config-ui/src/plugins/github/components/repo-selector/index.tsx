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

import { MultiSelector } from '@/components'

import { ScopeItemType } from '../../types'

import { useRepoSelector, UseRepoSelectorProps } from './use-repo-selector'

interface Props extends UseRepoSelectorProps {
  disabledItems: ScopeItemType[]
  selectedItems: ScopeItemType[]
  onChangeItems: (selectedItems: ScopeItemType[]) => void
}

export const RepoSelector = ({
  disabledItems,
  selectedItems,
  onChangeItems,
  ...props
}: Props) => {
  const { loading, items, onSearch } = useRepoSelector(props)

  return (
    <MultiSelector
      placeholder='Search Repositories...'
      items={items}
      getKey={(it) => it.githubId}
      getName={(it) => it.name}
      disabledItems={disabledItems}
      selectedItems={selectedItems}
      onChangeItems={onChangeItems}
      loading={loading}
      noResult='No Repositories Available.'
      onQueryChange={onSearch}
    />
  )
}
