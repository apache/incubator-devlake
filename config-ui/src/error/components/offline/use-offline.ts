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

import { useState, useEffect, useMemo, useRef } from 'react'
import { useHistory } from 'react-router-dom'

import { operator } from '@/utils'

import * as API from './api'

const pollTimer = 10000
const retryLimit = 10

export interface UseOfflineProps {
  onResetError: () => void
}

export const useOffline = ({ onResetError }: UseOfflineProps) => {
  const [processing, setProcessing] = useState(false)
  const [offline, setOffline] = useState(true)

  const history = useHistory()

  const timer = useRef<any>()
  const retryCount = useRef<number>(0)

  const ping = async (auto = true) => {
    const [success] = await operator(() => API.ping(), {
      setOperating: setProcessing,
      formatReason: () => 'Attempt to connect to the API failed'
    })

    if (success) {
      setOffline(false)
    }

    if (auto) {
      retryCount.current += 1
    }
  }

  useEffect(() => {
    timer.current = setInterval(() => {
      ping()
    }, pollTimer)
    return () => clearInterval(timer.current)
  }, [])

  useEffect(() => {
    if (retryCount.current >= retryLimit || !offline) {
      clearInterval(timer.current)
    }
  }, [retryCount.current, offline])

  return useMemo(
    () => ({
      processing,
      offline,
      onRefresh: () => ping(false),
      onContinue: () => {
        onResetError()
        history.push('/')
      }
    }),
    [processing, offline]
  )
}
