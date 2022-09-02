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

/**
 * @typedef {object} ProviderListConnection
 * @property {number?} id
 * @property {number?} key
 * @property {number?} connectionId
 * @property {string?} name
 * @property {string?} title
 * @property {number|string|object?} value
 * @property {string?} endpoint
 * @property {string?} proxy
 * @property {string?} token
 * @property {string?} username
 * @property {string?} password
 * @property {number?} rateLimitPerHour
 * @property {Date?} createdAt
 * @property {Date?} updatedAt
 * @property {plain|token?} authentication
 * @property {string|object?} plugin
 * @property {<Array<DataEntity>>} entities
 * @property {boolean} multiConnection
 * @property {string?} status
 * @property {object?} statusResponse
 * @property {string|object?} provider
 * @property {string?} providerId
 *
 */
class ProviderListConnection {
  constructor (data = {}) {
    this.id = parseInt(data?.id, 10) || null
    this.key = parseInt(data?.key, 10) || null
    this.connectionId = parseInt(data?.connectionId, 10) || null
    this.name = data?.name || ''
    this.title = data?.title || 'Default Connection'
    this.value = data?.value || null
    this.endpoint = data?.endpoint || ''
    this.proxy = data?.proxy || ''
    this.token = data?.token || ''
    this.username = data?.username || ''
    this.password = data?.password || ''
    this.rateLimitPerHour = data?.rateLimitPerHour || 0
    this.createdAt = data?.createdAt || null
    this.updatedAt = data?.updatedAt || null

    this.authentication = data?.authentication || 'plain'
    this.plugin = data?.plugin || null
    this.entities = data?.entities || []
    this.multiConnection = data?.multiConnection || true

    this.status = data?.status || 0
    this.statusResponse = data?.statusResponse || null
    this.provider = data?.provider || null
    this.providerId = data?.providerId || null
  }

  get (property) {
    return this[property]
  }

  set (property, value) {
    this[property] = value
    return this.property
  }
}

export default ProviderListConnection
