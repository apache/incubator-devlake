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
const assert = require('assert')
const AVAILABLE_PLUGINS = require('../src/data/availablePlugins')
const TEST_DATA = require('./testData')
const { getCollectorJson, getCollectionJson } = require('../src/utils/triggersUtil')

describe('Json utils', () => {
  describe('getCollectionJson', function () {
    it('gets default JSON for plugins based on an array of names', function () {
      const expected = TEST_DATA.completeTriggersJson
      const actual = getCollectionJson(AVAILABLE_PLUGINS)
      assert.deepEqual(expected, actual)
    })
  })
  describe('getCollectorJson', function () {
    it('gets default JSON for a collector plugin based on the name', function () {
      const expected = TEST_DATA.gitlabTriggersJson
      const actual = getCollectorJson('gitlab')
      assert.deepEqual(expected, actual)
    })
  })
})
