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
import { useHistory } from 'react-router-dom'

import { Plugins, PluginConfig } from '@/plugins'
import { operator } from '@/utils'

import * as API from './api'

interface Props {
  plugin: Plugins
  id?: ID
}

export const useDetail = ({ plugin, id }: Props) => {
  const [loading, setLoading] = useState(false)
  const [operating, setOperating] = useState(false)
  const [name, setName] = useState('')
  const [transformation, setTransformation] = useState<any>(
    PluginConfig.find((pc) => pc.plugin === plugin)?.transformation
  )

  const history = useHistory()

  const getTransformation = async () => {
    if (!id) return
    setLoading(true)
    try {
      const res = await API.getTransformation(plugin, id)
      setName(res.name)
      setTransformation(res)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    getTransformation()
  }, [])

  const handleSave = async () => {
    const payload = {
      ...transformation,
      name
    }

    const [success] = await operator(
      () =>
        id
          ? API.updateTransformation(plugin, id, payload)
          : API.createTransformation(plugin, payload),
      {
        setOperating
      }
    )

    if (success) {
      history.push('/transformations')
    }
  }

  return useMemo(
    () => ({
      loading,
      operating,
      name,
      transformation,
      onChangeName: setName,
      onChangeTransformation: setTransformation,
      onSave: handleSave
    }),
    [loading, operating, name, transformation]
  )
}
