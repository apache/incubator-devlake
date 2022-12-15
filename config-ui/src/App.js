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

import React from 'react'
import { Switch, Route, Redirect } from 'react-router-dom'

import '@/styles/app.scss'
import '@fontsource/inter/400.css'
import '@fontsource/inter/600.css'
import '@fontsource/inter/variable-full.css'

import { BaseLayout } from '@/layouts'
import {
  ProjectHomePage,
  ProjectDetailPage,
  CreateBlueprintPage,
  BlueprintDetailPage,
  WebHookConnectionPage
} from '@/pages'
import Integration from '@/pages/configure/integration/index'
import ManageIntegration from '@/pages/configure/integration/manage'
import AddConnection from '@/pages/configure/connections/AddConnection'
import ConfigureConnection from '@/pages/configure/connections/ConfigureConnection'
import Blueprints from '@/pages/blueprints/index'

function App() {
  return (
    <BaseLayout>
      <Switch>
        <Route path='/' exact component={() => <Redirect to='/projects' />} />
        <Route exact path='/projects' component={() => <ProjectHomePage />} />
        <Route
          exact
          path='/projects/:pname'
          component={() => <ProjectDetailPage />}
        />
        <Route
          exact
          path='/projects/:pname/create-blueprint'
          component={() => <CreateBlueprintPage from='project' />}
        />
        <Route exact path='/integrations' component={() => <Integration />} />
        <Route
          path='/integrations/:providerId'
          component={() => <ManageIntegration />}
        />
        <Route
          path='/connections/add/:providerId'
          component={() => <AddConnection />}
        />
        <Route
          path='/connections/configure/:providerId/:connectionId'
          component={() => <ConfigureConnection />}
        />
        <Route
          exact
          path='/connections/incoming-webhook'
          component={() => <WebHookConnectionPage />}
        />
        <Route exact path='/blueprints' component={() => <Blueprints />} />
        <Route
          exact
          path='/blueprints/create'
          component={() => <CreateBlueprintPage from='blueprint' />}
        />
        <Route
          exact
          path='/blueprints/:id'
          component={() => <BlueprintDetailPage />}
        />
      </Switch>
    </BaseLayout>
  )
}

export default App
