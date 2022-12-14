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

import { useState, useMemo } from 'react'

import { operator } from '@/utils'

import * as API from '../api'

export interface UseDeleteProps {
  initialValues?: {
    id: ID
  }
  onSubmitAfter?: (id: ID) => void
}

export const useDelete = ({ initialValues, onSubmitAfter }: UseDeleteProps) => {
  const [saving, setSaving] = useState(false)

  const handleDelete = async () => {
    if (!initialValues) return

    const [success] = await operator(
      () => API.deleteConnection(initialValues.id),
      {
        setOperating: setSaving
      }
    )

    if (success) {
      onSubmitAfter?.(initialValues.id)
    }
  }

  return useMemo(
    () => ({
      saving,
      onSubmit: handleDelete
    }),
    [saving]
  )
}
