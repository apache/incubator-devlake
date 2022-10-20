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
 * @typedef {object} MenuItem
 * @property {number?} id
 * @property {boolean?} disabled
 * @property {string?} label
 * @property {string|Object?} route
 * @property {boolean?} active
 * @property {string|Object} icon
 * @property {<Array<string>>?} classNames
 * @property {<Array<MenuItem>>?} children
 * @property {string?} target
 */
class MenuItem {
  constructor(data = {}) {
    this.id = data?.id || null
    this.disabled = !!data?.disabled
    this.label = data?.label || null
    this.route = data?.route || '#'
    this.active = !!data?.active
    this.icon = data?.icon || null
    this.classNames = Array.isArray(data?.classNames) ? data?.classNames : []
    this.target = data?.target || null
  }

  get(property) {
    return this[property]
  }

  set(property, value) {
    this[property] = value
    return this.property
  }
}

export default MenuItem
