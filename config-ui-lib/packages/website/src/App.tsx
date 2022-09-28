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

import * as Layout from './layouts';
import * as Page from './pages';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Layout.Base />}>
          <Route index element={<Navigate to="connections" replace />} />
          <Route path="connections" element={<Page.Connections />} />
          <Route path="connection/:type" element={<Page.Connection />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
