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

import { Switch, Route, Redirect, Router } from 'react-router-dom';
import { LoginPage } from './pages/login/login';
import { history } from './utils/history';
import { ErrorLayout, BaseLayout } from '@/layouts';
import { FromEnum } from '@/pages';
import {
  OfflinePage,
  DBMigratePage,
  ProjectHomePage,
  ProjectDetailPage,
  ConnectionHomePage,
  ConnectionListPage,
  ConnectionFormPage,
  BlueprintHomePage,
  BlueprintCreatePage,
  BlueprintDetailPage,
  BlueprintConnectioAddPage,
  BlueprintConnectionDetailPage,
} from '@/pages';

function App() {
  return (
    <Router history={history}>
      <Switch>
        <Route exact path="/login" component={() => <LoginPage />} />

        <Route
          exact
          path="/offline"
          component={() => (
            <ErrorLayout>
              <OfflinePage />
            </ErrorLayout>
          )}
        />

        <Route
          exact
          path="/db-migrate"
          component={() => (
            <ErrorLayout>
              <DBMigratePage />
            </ErrorLayout>
          )}
        />

        <Route
          path="/"
          component={() => (
            <BaseLayout>
              <Switch>
                <Route exact path="/" component={() => <Redirect to="/projects" />} />
                <Route exact path="/projects" component={() => <ProjectHomePage />} />
                <Route exact path="/projects/:pname" component={() => <ProjectDetailPage />} />
                <Route
                  exact
                  path="/projects/:pname/:bid/connection-add"
                  component={() => <BlueprintConnectioAddPage />}
                />
                <Route exact path="/projects/:pname/:bid/:unique" component={() => <BlueprintConnectionDetailPage />} />
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
                <Route
                  exact
                  path="/blueprints/create"
                  component={() => <BlueprintCreatePage from={FromEnum.blueprint} />}
                />
                <Route exact path="/blueprints/:id" component={() => <BlueprintDetailPage />} />
                <Route exact path="/blueprints/:bid/connection-add" component={() => <BlueprintConnectioAddPage />} />
                <Route exact path="/blueprints/:bid/:unique" component={() => <BlueprintConnectionDetailPage />} />
              </Switch>
            </BaseLayout>
          )}
        />
      </Switch>
    </Router>
  );
}

export default App;
