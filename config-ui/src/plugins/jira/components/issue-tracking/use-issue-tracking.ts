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
import { uniqWith } from 'lodash'

import { useProxyPrefix } from '../../hooks'
import * as API from '../../api'

export type IssueTypeItem = {
  id: string
  name: string
  iconUrl: string
}

export type FieldItem = {
  id: string
  name: string
}

export interface UseIssueTrackingProps {
  connectionId: ID
}

export const useIssueTracking = ({ connectionId }: UseIssueTrackingProps) => {
  const [issueTypes, setIssueTypes] = useState<IssueTypeItem[]>([])
  const [fields, setFields] = useState<FieldItem[]>([])

  const prefix = useProxyPrefix(connectionId)

  useEffect(() => {
    ;(async () => {
      const [its, fds] = await Promise.all([
        API.getIssueType(prefix),
        API.getField(prefix)
      ])
      setIssueTypes(uniqWith(its, (it, oit) => it.name === oit.name))
      setFields(fds)
    })()
  }, [prefix])

  return useMemo(
    () => ({
      issueTypes,
      fields
    }),
    [issueTypes, fields]
  )
}
