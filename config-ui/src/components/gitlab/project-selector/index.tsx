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

import React, { useState } from 'react'
import { MultiSelect } from '@blueprintjs/select'
import { Checkbox, MenuItem, Intent } from '@blueprintjs/core'

import {
  useGitLabProjectSelector,
  UseGitLabProjectSelectorProps
} from './use-gitlab-project-selector'
import * as S from './styled'

type Item = {
  id: number
  key: number
  title: string
  shortTitle: string
}

interface Props extends UseGitLabProjectSelectorProps {
  selectedItems: Array<Item>
  onChangeSelectItems: (selectedItems: Array<Item>) => void
}

export const GitLabProjectSelector = ({
  connectionId,
  selectedItems,
  onChangeSelectItems
}: Props) => {
  const [membership, setMembership] = useState(true)
  const [search, setSearch] = useState('')

  const { loading, items } = useGitLabProjectSelector({
    connectionId,
    search,
    membership
  })

  const handleQueryChange = (query: string) => {
    setSearch(query)
  }

  const handleSelectItem = (item: Item) => {
    const newItems = [...selectedItems, item]
    onChangeSelectItems(newItems)
  }

  const handleRemoveItem = (item: Item) => {
    const newItems = selectedItems.filter((it) => item.id !== it.id)
    onChangeSelectItems(newItems)
  }

  const handleChangeChekbox = (e: React.ChangeEvent<HTMLInputElement>) => {
    setMembership(e.target.checked)
  }

  const tagRenderer = (item: any) => {
    return <span>{item.shortTitle || item.title}</span>
  }

  const itemRenderer = (item: Item, { handleClick }: any) => {
    const selected = !!selectedItems.find((it) => it.id === item.id)

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
        popoverProps={{ usePortal: false, minimal: true }}
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
        onQueryChange={handleQueryChange}
        onItemSelect={handleSelectItem}
        onRemove={handleRemoveItem}
      />
      <Checkbox
        className='checkbox'
        label='Only search my repositories'
        checked={membership}
        onChange={handleChangeChekbox}
      />
    </S.Container>
  )
}
