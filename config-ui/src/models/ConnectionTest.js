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
 * @type {object} ConnectionStatusCodes
 */
const ConnectionStatusCodes = {
  OFFLINE: 0,
  ONLINE: 1,
  DISCONNECTED: 2,
  TESTING: 3
}

/**
 * @type {object} ConnectionStatusLabels
 */
const ConnectionStatusLabels = {
  [ConnectionStatusCodes.OFFLINE]: 'Offline',
  [ConnectionStatusCodes.ONLINE]: 'Online',
  [ConnectionStatusCodes.DISCONNECTED]: 'Disconnected',
  [ConnectionStatusCodes.TESTING]: 'Testing'
}

/**
 * @typedef {object} ConnectionTest
 * @property {Connection?} connection
 * @property {number?} status
 * @property {string?} statusLabel
 * @property {object?} testResponse
 */

class ConnectionTest {
  constructor(data = {}) {
    this.connection = data?.connection || null
    this.status = data?.status || ConnectionStatusCodes.OFFLINE
    this.statusLabel = data?.statusLabel || ConnectionStatusLabels[this.status]
    this.testResponse = data?.testResponse || null
  }

  get(property) {
    return this[property]
  }

  set(property, value) {
    this[property] = value
    return this.property
  }

  getStatusCodes() {
    return ConnectionStatusCodes
  }

  getStatusLabels() {
    return ConnectionStatusLabels
  }

  isOnline() {
    return this.status === ConnectionStatusCodes.ONLINE
  }

  isOffline() {
    return this.status === ConnectionStatusCodes.OFFLINE
  }

  isTesting() {
    return this.status === ConnectionStatusCodes.TESTING
  }
}

export default ConnectionTest
