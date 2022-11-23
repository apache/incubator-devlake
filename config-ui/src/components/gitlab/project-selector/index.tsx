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
import { MultiSelect } from '@blueprintjs/select'
import { Checkbox, MenuItem, Intent } from '@blueprintjs/core'

import type { ItemType } from './use-gitlab-project-selector'
import {
  useGitLabProjectSelector,
  UseGitLabProjectSelectorProps
} from './use-gitlab-project-selector'
import * as S from './styled'

interface Props extends UseGitLabProjectSelectorProps {
  disabledItemIds?: Array<ItemType['id']>
}

export const GitLabProjectSelector = ({
  connectionId,
  disabledItemIds,
  selectedItems,
  onChangeItems
}: Props) => {
  const {
    loading,
    items,
    search,
    membership,
    onSearch,
    onChangeMembership,
    onSelect,
    onRemove
  } = useGitLabProjectSelector({
    connectionId,
    selectedItems,
    onChangeItems
  })

  const tagRenderer = (item: any) => {
    return <span>{item.shortTitle || item.title}</span>
  }

  const itemRenderer = (item: ItemType, { handleClick }: any) => {
    const selected = !![
      ...selectedItems.map((it) => it.id),
      ...(disabledItemIds ?? [])
    ].find((id) => id === item.id)

    return (
      <MenuItem
        key={item.key}
        text={
          <Checkbox label={item.title} checked={selected} disabled={selected} />
        }
        disabled={selected}
        onClick={handleClick}
      />
    )
  }

  return (
    <S.Container>
      <MultiSelect
        className='selector'
        placeholder='Select Projects'
        popoverProps={{ usePortal: false, minimal: true, isOpen: !!search }}
        resetOnSelect
        fill
        items={items}
        selectedItems={selectedItems}
        tagInputProps={{
          tagProps: {
            intent: Intent.PRIMARY,
            minimal: true
          }
        }}
        noResults={
          <MenuItem
            disabled={true}
            text={loading ? 'Fetching...' : 'No Projects Available.'}
          />
        }
        tagRenderer={tagRenderer}
        itemRenderer={itemRenderer}
        onQueryChange={onSearch}
        onItemSelect={onSelect}
        onRemove={onRemove}
      />
      <Checkbox
        className='checkbox'
        label='Only search my repositories'
        checked={membership}
        onChange={onChangeMembership}
      />
    </S.Container>
  )
}
