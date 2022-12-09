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

import { Intent } from '@blueprintjs/core'

import { Toast } from '@/components'

export type OperateConfig = {
  setOperating?: (success: boolean) => void
  formatReason?: (err: unknown) => string
}

/**
 * 
 * @param request -> a request
 * @param config 
 * @param config.setOperating -> Control the status of the request
 * @parma config.formatReason -> Show the reason for the failure
 * @returns 
 */
export const operator = async <T>(
  request: () => Promise<T>,
  config?: OperateConfig
) => {
  const { setOperating, formatReason } = config || {}

  try {
    setOperating?.(true)
    const res = await request()
    Toast.show({
      intent: Intent.SUCCESS,
      message: 'Operation successfully completed',
      icon: 'endorsed'
    })
    return [true, res]
  } catch (err) {
    const reason = formatReason?.(err)
    Toast.show({
      intent: Intent.DANGER,
      message: reason
        ? `Operation failed. Reason: ${reason}`
        : 'Operation failed.',
      icon: 'error'
    })

    return [false]
  } finally {
    setOperating?.(false)
  }
}
