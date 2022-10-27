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
import { Button, MenuItem, Icon, Colors } from '@blueprintjs/core'
import { MultiSelect } from '@blueprintjs/select'

import * as S from './styled'

const StatusList = [
  'new',
  'open',
  'resolved',
  'on hold',
  'invalid',
  'duplicate',
  'wontfix',
  'closed'
]

export const StatusSelect = ({
  name,
  saving,
  selectedItems,
  disabledItems = [],
  onItemSelect,
  onItemRemove,
  onItemClear,
  style = {}
}) => {
  return (
    <S.Container style={style}>
      <label className='label'>{name}</label>
      <MultiSelect
        className='select'
        placeholder='Select...'
        fill={true}
        resetOnSelect={true}
        disabled={saving}
        items={StatusList.filter((it) => !disabledItems.includes(it))}
        itemPredicate={(query, it) =>
          it.toLowerCase().indexOf(query.toLowerCase()) >= 0
        }
        selectedItems={selectedItems}
        itemRenderer={(it, { handleClick }) => (
          <MenuItem
            key={it}
            disabled={selectedItems.includes(it)}
            text={
              selectedItems.includes(it) ? (
                <span>
                  {it}
                  <Icon icon='small-tick' color={Colors.GREEN5} />
                </span>
              ) : (
                <span>{it}</span>
              )
            }
            onClick={handleClick}
          />
        )}
        tagRenderer={(it) => it}
        tagInputProps={{
          tagProps: {
            minimal: true
          }
        }}
        noResults={<MenuItem disabled={true} text='No results.' />}
        onRemove={(it) => onItemRemove(name, it)}
        onItemSelect={(it) => onItemSelect(name, it)}
      />
      <Button
        icon='eraser'
        disabled={selectedItems.length === 0 || saving}
        onClick={() => onItemClear(name)}
      />
    </S.Container>
  )
}
