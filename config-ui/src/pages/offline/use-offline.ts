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

import * as API from './api'

const pollTimer = 10000
const retryLimit = 10

export const useOffline = () => {
  const [loading, setLoading] = useState(false)
  const [offline, setOffline] = useState(true)

  const timer = useRef<any>()
  const retryCount = useRef<number>(0)

  const ping = async (auto = true) => {
    setLoading(true)
    try {
      await API.ping()
      setOffline(false)
    } catch {
      setOffline(true)
    } finally {
      if (auto) {
        retryCount.current += 1
      }
      setLoading(false)
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
      loading,
      offline,
      onRefresh: () => ping(false)
    }),
    [loading, offline]
  )
}
