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

import React, { useContext } from 'react'

import { PageLoading } from '@/components'
import type { PluginConfigType } from '@/plugins'
import { Plugins } from '@/plugins'

import type { TransformationItemType } from './types'
import { useContextValue } from './use-context-value'

const TransformationContext = React.createContext<{
  plugins: PluginConfigType[]
  transformations: TransformationItemType[]
}>({
  plugins: [],
  transformations: []
})

interface Props {
  children?: React.ReactNode
}

export const TransformationContextProvider = ({ children }: Props) => {
  const { loading, plugins, transformations } = useContextValue()

  if (loading) {
    return <PageLoading />
  }

  return (
    <TransformationContext.Provider
      value={{
        plugins,
        transformations
      }}
    >
      {children}
    </TransformationContext.Provider>
  )
}

export const TransformationContextConsumer = TransformationContext.Consumer

export const useTransformation = () => useContext(TransformationContext)
