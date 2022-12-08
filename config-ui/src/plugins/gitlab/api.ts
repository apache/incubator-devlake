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

import request from '@/components/utils/request'

type PaginationParams = {
  page: number
  per_page: number
}

export const getUser = (prefix: string) => request(`${prefix}/user`)

export const getUserGroups = (prefix: string, params: PaginationParams) =>
  request(`${prefix}/groups`, {
    data: { top_level_only: 1, ...params }
  })

export const getUserProjects = (
  prefix: string,
  uid: ID,
  params: PaginationParams
) =>
  request(`${prefix}/users/${uid}/projects`, {
    data: params
  })

export const getGroupSubgroups = (
  prefix: string,
  gid: ID,
  params: PaginationParams
) =>
  request(`${prefix}/groups/${gid}/subgroups`, {
    data: params
  })

export const getGroupProjects = (
  prefix: string,
  gid: ID,
  params: PaginationParams
) =>
  request(`${prefix}/groups/${gid}/projects`, {
    data: params
  })

type SearchProjectParams = {
  search: string
  membership: boolean
}

export const searchProject = (prefix: string, params: SearchProjectParams) =>
  request(`${prefix}/projects`, {
    data: params
  })
