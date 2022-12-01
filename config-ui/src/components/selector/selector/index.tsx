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

import React, { useState, useEffect, useMemo } from 'react'
import { MenuItem, Button } from '@blueprintjs/core'
import { Select } from '@blueprintjs/select'

interface Props<T> {
  items: T[]
  disabledItems?: T[]
  getKey?: (item: T) => string
  getName?: (item: T) => string
  selectedItem?: T
  onChangeItem?: (selectedItem: T) => void
}

export const Selector = <T,>({
  items,
  disabledItems = [],
  getKey = (it) => it as string,
  getName = (it) => it as string,
  onChangeItem,
  ...props
}: Props<T>) => {
  const [selectedItem, setSelectedItem] = useState<T>()

  useEffect(() => {
    setSelectedItem(props.selectedItem)
  }, [props.selectedItem])

  const btnText = useMemo(
    () => (selectedItem ? getName(selectedItem) : 'Select...'),
    [selectedItem]
  )

  const itemPredicate = (query: string, item: T) => {
    const name = getName(item)
    return name.toLowerCase().indexOf(query.toLowerCase()) >= 0
  }

  const itemRenderer = (item: T, { handleClick }: any) => {
    const key = getKey(item)
    const name = getName(item)
    const disabled =
      !!disabledItems.find((it) => getKey(it) === getKey(item)) ||
      selectedItem === item

    return (
      <MenuItem
        key={key}
        disabled={disabled}
        text={<span>{name}</span>}
        onClick={handleClick}
      />
    )
  }

  const handleItemSelect = (item: T) => {
    if (onChangeItem) {
      onChangeItem(item)
    } else {
      setSelectedItem(item)
    }
  }

  return (
    <Select
      items={items}
      activeItem={selectedItem}
      itemPredicate={itemPredicate}
      itemRenderer={itemRenderer}
      onItemSelect={handleItemSelect}
    >
      <Button
        outlined
        small
        fill
        rightIcon='double-caret-vertical'
        text={btnText}
      />
    </Select>
  )
}
