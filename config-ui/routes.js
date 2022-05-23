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
const routes = require('next-routes')

module.exports = routes()
// ? Name is URL route
// ? Page is page (file) name in /pages
// ? Pattern is for dynamic routes eg. '/user/:id'

// Main Setup
  .add({ name: 'configuration', page: '/' })
  .add({ name: 'triggers', page: '/triggers' })

// Plugins
  .add({ name: 'jira', page: '/plugins/jira' })
  .add({ name: 'gitlab', page: '/plugins/gitlab' })
  .add({ name: 'jenkins', page: '/plugins/jenkins' })
