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
import { Colors, Icon, Intent, MenuItem } from '@blueprintjs/core'
import { MultiSelect } from '@blueprintjs/select'
import JiraBoard from '@/models/JiraBoard'

const BoardsSelector = (props) => {
  const {
    configuredConnection,
    placeholder = 'Search and select boards',
    items = [],
    selectedItems = [],
    // eslint-disable-next-line max-len
    restrictedItems = [],
    activeItem = null,
    disabled = false,
    isLoading = false,
    isSaving = false,
    onItemSelect = () => {},
    onRemove = () => {},
    onClear = () => {},
    onQueryChange = () => {},
    itemRenderer = (item, { handleClick, modifiers }) => (
      <MenuItem
        active={modifiers.active}
        disabled={selectedItems.find((i) => i?.id === item?.id)}
        key={item.value}
        // label=
        onClick={handleClick}
        text={
          selectedItems.find((i) => i?.id === item?.id) ? (
            <>
              <img src={item.icon} width={12} height={12} /> {item?.title}{' '}
              <Icon icon='small-tick' color={Colors.GREEN5} />
            </>
          ) : (
            <span style={{ fontWeight: 700 }}>
              <img src={item.icon} width={12} height={12} /> {item?.title}
            </span>
          )
        }
        style={{
          marginBottom: '2px',
          fontWeight: items.includes(item) ? 700 : 'normal'
        }}
      />
    ),
    tagRenderer = (item) => item?.title
  } = props

  return (
    <>
      <div
        className='boards-multiselect'
        style={{ display: 'flex', marginBottom: '10px' }}
      >
        <div
          className='boards-multiselect-selector'
          style={{ minWidth: '200px', width: '100%' }}
        >
          <MultiSelect
            disabled={disabled || isSaving || isLoading}
            resetOnSelect={true}
            placeholder={placeholder}
            popoverProps={{ usePortal: false, minimal: true }}
            className='multiselector-boards'
            inline={true}
            fill={true}
            items={items}
            selectedItems={selectedItems}
            activeItem={activeItem}
            itemRenderer={itemRenderer}
            tagRenderer={tagRenderer}
            tagInputProps={{
              tagProps: {
                intent: Intent.PRIMARY,
                minimal: true
              }
            }}
            noResults={<MenuItem disabled={true} text='No Boards Available.' />}
            onQueryChange={(query) => onQueryChange(query)}
            onRemove={(item) => {
              onRemove(selectedItems.filter((t) => t?.id !== item.id))
            }}
            onItemSelect={(item) => {
              onItemSelect(
                !selectedItems.includes(item)
                  ? [...selectedItems, new JiraBoard(item)]
                  : selectedItems
              )
            }}
            style={{ borderRight: 0 }}
          />
        </div>
      </div>
    </>
  )
}

export default BoardsSelector
