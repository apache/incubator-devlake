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
 * @typedef {object} DataScopeConnection
 * @property {number?} id
 * @property {NORMAL|ADVANCED?} mode
 * @property {number|string?} name
 * @property {number?} connectionId
 * @property {number|string?} value
 * @property {object|string?} provider
 * @property {string?} providerLabel
 * @property {string?} providerId
 * @property {object?} plugin
 * @property {object?} provider
 * @property {object|string?} icon
 * @property {<Array<object>>?} projects
 * @property {<Array<string|object>>?} boards
 * @property {<Array<number>>?} boardIds
 * @property {<Array<object>>?} boardsList
 * @property {<Array<object>>?} dataDomains
 * @property {<Array<object>>?} transformations
 * @property {<Array<object>>?} transformationStates
 * @property {object?} scope
 * @property {boolean} editable
 * @property {boolean} advancedEditable
 * @property {boolean} isMultiStage
 * @property {boolean} isSingleStage
 * @property {number?} stage
 * @property {number?} totalStages
 */
class DataScopeConnection {
  constructor(data = {}) {
    this.id = parseInt(data?.id, 10) || null
    this.mode = data?.mode || 'NORMAL'
    this.name = data?.name || null
    this.connectionId = parseInt(data?.connectionId, 10) || null
    this.value = data?.value || this.connectionId
    this.provider = data?.provider || null
    this.providerLabel = data?.providerLabel || null
    this.providerId = data?.providerId || null
    this.plugin = data?.plugin || null
    this.icon = data?.icon || null
    this.projects = data?.projects || []
    this.boards = data?.boards || []
    this.boardIds = data?.boardIds || []
    this.boardsList = data?.boardsList || []
    this.dataDomains = data?.dataDomains || []
    this.transformations = data?.transformations || []
    this.transformationStates = data?.transformationStates || []
    this.scope = data?.scope || null
    this.editable = data?.editable || false
    this.advancedEditable = data?.advancedEditable || false
    this.isMultiStage = data?.isMultiStage || false
    this.isSingleStage = data?.isSingleStage || true
    this.stage = data?.stage || 1
    this.totalStages = data?.totalStages || 1
  }

  get(property) {
    return this[property]
  }

  set(property, value) {
    this[property] = value
    return this.property
  }
}

export default DataScopeConnection
