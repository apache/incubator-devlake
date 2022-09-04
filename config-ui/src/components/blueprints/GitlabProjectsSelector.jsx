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
import React, { useEffect, useState } from 'react'
import { Checkbox, Intent, MenuItem, } from '@blueprintjs/core'
import { MultiSelect } from '@blueprintjs/select'

const GitlabProjectsSelector = (props) => {
  const {
    fetchGitlabProjects = () => [],
    isFetchingGitlab = false,
    configuredConnection,
    placeholder = 'Select Projects',
    items = [],
    selectedItems = [],
    activeItem = null,
    disabled = false,
    isLoading = false,
    isSaving = false,
    onItemSelect = () => {},
    onRemove = () => {},
    onClear = () => {},
    itemRenderer = (item, { handleClick, modifiers }) => (
      <MenuItem
        active={modifiers.active}
        disabled={
          selectedItems.find(i => i?.id === item?.id)
        }
        key={item.value}
        onClick={handleClick}
        text={
          selectedItems.find(i => i?.id === item?.id)
            ? (
              <>
                <input type='checkbox' checked readOnly /> {item?.title}
              </>
              )
            : (
              <span style={{ fontWeight: 700 }}>
                <input type='checkbox' readOnly /> {item?.title}
              </span>
              )
        }
        style={{
          marginBottom: '2px',
          fontWeight: items.includes(item) ? 700 : 'normal',
        }}
      />
    ),
    tagRenderer = (item) => item?.title,
  } = props

  const [query, setQuery] = useState('')
  const [onlyQueryMemberRepo, setOnlyQueryMemberRepo] = useState(true)

  useEffect(() => {
    if (query.length <= 2) {
      // only search when type more than 2 char or empty
      return
    }
    // prevent request too frequently
    const timer = setTimeout(() => {
      fetchGitlabProjects(query, onlyQueryMemberRepo)
    }, 200)
    return () => clearTimeout(timer)
  }, [fetchGitlabProjects, query])

  useEffect(() => {
    if (query.length <= 2) {
      // only search when type more than 2 char or empty
      return
    }
    fetchGitlabProjects(query, onlyQueryMemberRepo)
  }, [fetchGitlabProjects, onlyQueryMemberRepo])

  return (
    <div
      className='gitlab-projects-multiselect'
      style={{ display: 'flex', marginBottom: '10px' }}
    >
      <div
        className='gitlab-projects-multiselect-selector'
        style={{ minWidth: '200px', width: '100%' }}
      >
        <MultiSelect
          disabled={disabled || isSaving || isLoading}
          // openOnKeyDown={true}
          resetOnSelect={true}
          placeholder={placeholder}
          popoverProps={{ usePortal: false, minimal: true }}
          className='multiselector-projects'
          inline={true}
          fill={true}
          items={items}
          selectedItems={selectedItems}
          activeItem={activeItem}
          onQueryChange={query => setQuery(query)}
          itemRenderer={itemRenderer}
          tagRenderer={tagRenderer}
          tagInputProps={{
            tagProps: {
              intent: Intent.PRIMARY,
              minimal: true
            },
          }}
          noResults={
            (query.length <= 2 && <MenuItem disabled={true} text='Please type more than 2 char to search.' />) ||
            (isFetchingGitlab && <MenuItem disabled={true} text='Fetching...' />) ||
            <MenuItem disabled={true} text='No Projects Available.' />
          }
          onRemove={(item) => {
            onRemove((rT) => {
              return {
                ...rT,
                [configuredConnection.id]: rT[configuredConnection.id].filter(
                  (t) => t?.id !== item.id
                ),
              }
            })
          }}
          onItemSelect={(item) => {
            onItemSelect((rT) => {
              return !rT[configuredConnection.id].includes(item)
                ? {
                    ...rT,
                    [configuredConnection.id]: [
                      ...rT[configuredConnection.id],
                      item,
                    ],
                  }
                : { ...rT }
            })
          }}
          style={{ borderRight: 0 }}
        />

        <Checkbox
          label='only search repos joined' checked={onlyQueryMemberRepo}
          onChange={e => setOnlyQueryMemberRepo(!onlyQueryMemberRepo)}
          style={{ margin: '10px 0 0 6px' }}
        />
      </div>
    </div>
  )
}

export default GitlabProjectsSelector
