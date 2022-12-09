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

import { operator } from '@/utils'

import { defaultConfig } from './config'
import { createTransformation } from './api'

export interface UseTransformationProps {
  plugin: string
  name: string
  initialValues?: any
  onSaveAfter?: (tid: string | number) => void
}

export const useTransformation = ({
  plugin,
  name,
  initialValues,
  onSaveAfter
}: UseTransformationProps) => {
  const [saving, setSaving] = useState(false)
  const [transformation, setTransformation] = useState({})

  useEffect(() => {
    setTransformation(initialValues ? initialValues : defaultConfig[plugin])
  }, [initialValues, plugin])

  const handleSave = async () => {
    const [success, res] = await operator(
      () =>
        createTransformation(plugin, {
          ...transformation,
          name
        }),
      {
        setOperating: setSaving
      }
    )

    if (success) {
      onSaveAfter?.(res.id)
    }
  }

  return useMemo(
    () => ({
      saving,
      transformation,
      setTransformation,
      onSave: handleSave
    }),
    [saving, transformation, name, onSaveAfter]
  )
}
