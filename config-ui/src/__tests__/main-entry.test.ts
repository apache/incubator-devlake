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

/// <reference types="node" />

import { describe, it, expect } from 'vitest';
import { readFileSync } from 'fs';
import path from 'path';

/**
 * Regression guard for the React 19 upgrade.
 *
 * React 19 removed the legacy `ReactDOM.render` API; using it threw
 * "ReactDOM.render is not a function" and left a blank page. The entry point
 * must use `createRoot` from `react-dom/client` instead.
 */
describe('app entry (src/main.tsx)', () => {
  const source = readFileSync(path.join(process.cwd(), 'src/main.tsx'), 'utf-8');

  it('mounts with the React 19 createRoot API', () => {
    expect(source).toMatch(/createRoot\s*\(/);
    expect(source).toMatch(/from ['"]react-dom\/client['"]/);
  });

  it('does not use the removed legacy ReactDOM.render API', () => {
    expect(source).not.toMatch(/ReactDOM\.render\s*\(/);
    expect(source).not.toMatch(/from ['"]react-dom['"]/);
  });
});
