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

import React from 'react';
import { Switch, Route, Redirect } from 'react-router-dom';

import { BaseLayout } from '@/layouts';
import { FromEnum } from '@/pages';
import {
  ProjectHomePage,
  ProjectDetailPage,
  ConnectionHomePage,
  ConnectionListPage,
  ConnectionFormPage,
  BlueprintHomePage,
  BlueprintCreatePage,
  BlueprintDetailPage,
} from '@/pages';

function App() {
  return (
    <BaseLayout>
      <Switch>
        <Route path="/" exact component={() => <Redirect to="/projects" />} />
        <Route exact path="/projects" component={() => <ProjectHomePage />} />
        <Route exact path="/projects/:pname" component={() => <ProjectDetailPage />} />
        <Route
          exact
          path="/projects/:pname/create-blueprint"
          component={() => <BlueprintCreatePage from={FromEnum.project} />}
        />
        <Route exact path="/connections" component={() => <ConnectionHomePage />} />
        <Route exact path="/connections/:plugin" component={() => <ConnectionListPage />} />
        <Route exact path="/connections/:plugin/create" component={() => <ConnectionFormPage />} />
        <Route exact path="/connections/:plugin/:cid" component={() => <ConnectionFormPage />} />
        <Route exact path="/blueprints" component={() => <BlueprintHomePage />} />
        <Route exact path="/blueprints/create" component={() => <BlueprintCreatePage from={FromEnum.blueprint} />} />
        <Route exact path="/blueprints/:id" component={() => <BlueprintDetailPage />} />
      </Switch>
    </BaseLayout>
  );
}

export default App;
