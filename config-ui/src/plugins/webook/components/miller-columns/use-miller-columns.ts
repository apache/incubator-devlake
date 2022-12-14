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

import { useState, useEffect, useMemo } from 'react'
import type { ItemType } from 'miller-columns-select'

import * as API from '../../api'

type WebhookItemType = ItemType<{
  id: ID
  name: string
}>

export const useMillerColumns = () => {
  const [items, setItems] = useState<WebhookItemType[]>([])
  const [isLast, setIsLast] = useState(false)

  const updateItems = (arr: any) =>
    arr.map((it: any) => ({
      parentId: null,
      id: it.id,
      title: it.name,
      name: it.name
    }))

  useEffect(() => {
    ;(async () => {
      const res = await API.getConnections()
      setItems([...updateItems(res)])
      setIsLast(true)
    })()
  }, [])

  return useMemo(
    () => ({
      items,
      getHasMore() {
        return !isLast
      }
    }),
    [items, isLast]
  )
}
