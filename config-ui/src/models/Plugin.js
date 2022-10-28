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
import DataEntity from './DataEntity'

/**
 * @typedef {object} Plugin
 * @property {string|number?} id
 * @property {string} name
 * @property {string?} description
 * @property {'integration'|'pipeline'|'plugin'?} type
 * @property {bool} enabled
 * @property {bool} isBeta
 * @property {bool} isProvider
 * @property {bool} multiConnection
 * @property {bool} private
 * @property {object?} icon
 * @property {object?} connection
 * @property {object?} settings
 * @property {number?} connectionLimit
 * @property {<Array<DataEntity>>?} entities
 * @property {object?} transformations
 */
class Plugin {
  constructor(data = {}) {
    this.id = data?.id || Math.random() * 99999
    this.name = data?.name || 'New Plugin'
    this.type = data?.type || 'integration'
    this.description = data?.description || null
    this.enabled = data?.enabled || false
    this.isBeta = data?.isBeta || false
    this.isProvider = data?.isProvider || false
    this.multiConnection = data?.multiConnection || false
    this.private = data?.private || false
    this.icon = data?.icon || null
    this.connection = data?.connection || null
    this.connectionLimit = data?.connectionLimit || 0
    this.entities = data?.entities?.map((e) => new DataEntity({ type: e })) || [
      new DataEntity({ type: 'CODE' })
    ]
    this.transformations = data?.transformations || {
      scopes: { options: {} },
      default: {}
    }
  }

  get(property) {
    return this[property]
  }

  set(property, value) {
    this[property] = value
    return this.property
  }

  getAuthenticationType() {
    return this.connection?.authentication || 'plain'
  }

  getConnectionFields() {
    return this.connection ? this.connection?.fields : {}
  }

  getConnectionFormLabels() {
    return this.connection ? this.connection?.labels : {}
  }

  getConnectionFormPlaceholders() {
    return this.connection ? this.connection?.placeholders : {}
  }

  getConnectionFormTooltips() {
    return this.connection ? this.connection?.tooltips : {}
  }

  getDefaultTransformations() {
    return this.transformations?.default || {}
  }

  getDefaultTransformationScopeOptions(entity) {
    const scopeOptions = {
      ...(this.transformations?.scopes?.options || {}),
      ...(entity && typeof entity.getTransformationScopeOptions === 'function'
        ? entity.getTransformationScopeOptions()
        : {})
    }
    return scopeOptions
  }

  getDataEntities() {
    return this.entities || []
  }
}

export default Plugin
