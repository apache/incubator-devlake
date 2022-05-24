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
import { DEVLAKE_ENDPOINT } from './config.js'
import request from './request'
import PLUGINS from '../data/availablePlugins'

const ConnectionsUtil = {
  getPluginConnections: async () => {
    const pluginsToSet = []
    const errors = []
    for (const plugin of PLUGINS) {
      try {
        const res = await request.get(`${DEVLAKE_ENDPOINT}/plugins/${plugin}/connections`)
        if (res?.data?.length > 0) {
          pluginsToSet.push(plugin)
        }
      } catch (error) {
        errors.push(error)
      }
    }
    if (errors.length > 0) {
      console.log('errors', errors)
    }
    return pluginsToSet
  }
}

export default ConnectionsUtil
