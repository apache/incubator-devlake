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

export enum ItemTypeEnum {
  LEAF = 'leaf',
  BRANCH = 'branch'
}

export enum ItemStatusEnum {
  PENDING = 'pending',
  READY = 'ready'
}

export type ItemType = {
  id: string | number
  title: string
  type: ItemTypeEnum
  status: ItemStatusEnum
  items: ItemType[]
}

export type ItemInfoType = {
  item: ItemType
  parentId?: ItemType['id']
  childLoaded: boolean
}

export type ItemMapType = {
  getItem: (id: ItemType['id']) => ItemType
  getItemParent: (id: ItemType['id']) => ItemType | null
  getItemChildLoaded: (id: ItemType['id']) => boolean
}

export type ColumnType = {
  parentId: ItemType['id'] | null
  items: ItemType[]
  activeId: ItemType['id'] | null
}

export enum RowStatus {
  selected = 'selected',
  noselected = 'noselected'
}

export enum CheckedStatus {
  selected = 'selected',
  noselected = 'noselected',
  indeterminate = 'indeterminate'
}
