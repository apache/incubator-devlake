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

import Entity from '@/models/Entity'

/**
 * @typedef {object} GitHubProject
 * @property {number|string?} id
 * @property {number?} key
 * @property {string|number?} value
 * @property {string|number?} name
 * @property {string|number?} title
 * @property {string?} shortTitle
 * @property {string?} icon
 * @property {string?} owner
 * @property {string?} repo
 * @property {boolean?} useApi
 * @property {'project'|'board'|'job'} variant
 * @property {string?} providerId
 */
class GitHubProject extends Entity {
  constructor(data = {}) {
    super(data)
    this.id = data?.id || null
    this.key = data?.key || this.id || null
    this.owner = data?.owner || null
    this.repo = data?.repo || null
    this.name =
      data?.owner && data?.repo ? `${data?.owner}/${data?.repo}` : null
    this.value = data?.value || this.name || this.id || null
    this.title = data?.title || this.name || this.id || null
    this.shortTitle = data?.shortTitle || null
    this.icon = data?.icon || null

    this.useApi = data?.useApi || false
    this.variant = data?.variant || 'project'
    this.providerId = 'github'
  }

  getConfiguredEntityId() {
    return this.name?.toString() || this.id
  }

  getTransformationScopeOptions() {
    return {
      owner: this.owner,
      repo: this.repo
    }
  }
}

export default GitHubProject
