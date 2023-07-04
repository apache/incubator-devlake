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

import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ErrorLayout, BaseLayout } from '@/layouts';
import {
  OfflinePage,
  DBMigratePage,
  ConnectionHomePage,
  ConnectionDetailPage,
  ProjectHomePage,
  ProjectDetailPage,
  BlueprintHomePage,
  BlueprintDetailPage,
  BlueprintConnectionDetailPage,
} from '@/pages';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<ErrorLayout />}>
          <Route path="offline" element={<OfflinePage />} />
          <Route path="db-migrate" element={<DBMigratePage />} />
        </Route>

        <Route path="/" element={<BaseLayout />}>
          <Route index element={<Navigate to="connections" />} />
          <Route path="connections" element={<ConnectionHomePage />} />
          <Route path="connections/:plugin/:id" element={<ConnectionDetailPage />} />
          <Route path="projects" element={<ProjectHomePage />} />
          <Route path="projects/:pname" element={<ProjectDetailPage />} />
          <Route path="projects/:pname/:unique" element={<BlueprintConnectionDetailPage />} />
          <Route path="blueprints" element={<BlueprintHomePage />} />
          <Route path="blueprints/:id" element={<BlueprintDetailPage />} />
          <Route path="blueprints/:bid/:unique" element={<BlueprintConnectionDetailPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
