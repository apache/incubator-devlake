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
import { useState, useEffect } from 'react'
import axios from 'axios'

import { DEVLAKE_ENDPOINT } from '@/utils/config'
import { ToastNotification } from '@/components/Toast'

export const useWebhookManager = () => {
  const [loading, setLoading] = useState(false)
  const [operating, setOperating] = useState(false)
  const [data, setData] = useState([])

  const fetch = async () => {
    setLoading(true)
    try {
      const res = await axios.get(`${DEVLAKE_ENDPOINT}/plugins/webhook/connections`)
      setData(res.data)
    } catch (err) {
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    ;(async () => {
      await fetch()
    })()
  }, [])

  const onCreate = async (payload) => {
    setOperating(true)
    try {
      const {
        data: { id },
      } = await axios.post(`${DEVLAKE_ENDPOINT}/plugins/webhook/connections`, payload)
      const { data } = await axios.get(`${DEVLAKE_ENDPOINT}/plugins/webhook/connections/${id}`)
      fetch()
      return data
    } catch (err) {
    } finally {
      setOperating(false)
    }
  }

  const onUpdate = async (id, payload) => {
    setOperating(true)
    ToastNotification.clear()
    try {
      await axios.patch(`${DEVLAKE_ENDPOINT}/plugins/webhook/connections/${id}`, payload)
      ToastNotification.show({
        message: 'Update record succeeded.',
        intent: 'success',
        icon: 'small-tick',
      })
      fetch()
    } catch (err) {
      ToastNotification.show({
        message: err.response.data.message,
        intent: 'danger',
        icon: 'error',
      })
    } finally {
      setOperating(false)
    }
  }

  const onDelete = async (id) => {
    setOperating(true)
    ToastNotification.clear()
    try {
      await axios.delete(`${DEVLAKE_ENDPOINT}/plugins/webhook/connections/${id}`)
      ToastNotification.show({
        message: 'Delete record succeeded.',
        intent: 'success',
        icon: 'small-tick',
      })
      fetch()
    } catch (err) {
    } finally {
      setOperating(false)
    }
  }

  return {
    loading,
    data,
    operating,
    onCreate,
    onUpdate,
    onDelete,
  }
}
