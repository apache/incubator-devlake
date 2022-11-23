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

import { ItemType, ItemTypeEnum, RowStatus } from '../../types'

import { Checkbox, CheckStatus } from '../checkbox'

import * as S from './styled'

interface Props {
  item: ItemType
  status?: RowStatus
  checkStatus?: CheckStatus | Array<CheckStatus>
  checkedCount?: number
  onExpandItem?: (it: ItemType) => void
  onSelectItem?: (it: ItemType) => void
}

export const Item = ({
  item,
  status = RowStatus.noselected,
  checkStatus = CheckStatus.nochecked,
  onExpandItem,
  onSelectItem
}: Props) => {
  const handleRowClick = () => {
    onExpandItem?.(item)
  }

  const handleCheckboxClick = (e: React.MouseEvent<HTMLLabelElement>) => {
    if (item.type === ItemTypeEnum.LEAF) {
      e.stopPropagation()
    }
    onSelectItem?.(item)
  }

  return (
    <S.Wrapper
      type={item.type}
      selected={status === RowStatus.selected}
      onClick={handleRowClick}
    >
      <Checkbox status={checkStatus} onClick={handleCheckboxClick}>
        {item.title}
      </Checkbox>
      {item.type === ItemTypeEnum.BRANCH && <span className='indicator' />}
    </S.Wrapper>
  )
}
