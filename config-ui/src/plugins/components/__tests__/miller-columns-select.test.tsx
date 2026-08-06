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

import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import type { McsItem } from 'miller-columns-select';
import { MillerColumnsSelect } from 'miller-columns-select';

type Item = { name: string };

const items: McsItem<Item>[] = [
  { parentId: null, id: 'group-a', title: 'Group A', name: 'Group A' },
  { parentId: 'group-a', id: 'item-a-1', title: 'Item A-1', name: 'Item A-1' },
  { parentId: null, id: 'group-b', title: 'Group B', name: 'Group B' },
];

/**
 * Regression tests for the React 19 upgrade.
 *
 * miller-columns-select@1.4.1 ships a vendored copy of react-jsx-runtime whose
 * production and development variants both read React-18-only internals
 * (`ReactCurrentDispatcher` / `ReactCurrentOwner`). Those internals were removed
 * in React 19, so simply importing/rendering the component threw
 * "Cannot read properties of undefined (reading 'ReactCurrentDispatcher')" and
 * white-screened the whole app. We vendor-patch the package so it uses the host
 * `react/jsx-runtime`. These tests fail if that patch is dropped or the package
 * is upgraded to an again-incompatible build.
 */
describe('miller-columns-select under React 19', () => {
  let errorSpy: ReturnType<typeof vi.spyOn>;

  beforeEach(() => {
    errorSpy = vi.spyOn(console, 'error').mockImplementation(() => undefined);
  });

  afterEach(() => {
    errorSpy.mockRestore();
  });

  it('mounts without touching removed React internals', () => {
    expect(() => render(<MillerColumnsSelect<Item> items={items} mode="multiple" columnCount={2} />)).not.toThrow();

    const internalErrors = errorSpy.mock.calls
      .map((args: any[]) => args.map((a: any) => String(a)).join(' '))
      .filter((msg: string) => /ReactCurrentDispatcher|ReactCurrentOwner/.test(msg));
    expect(internalErrors).toEqual([]);
  });

  it('renders the provided root items', () => {
    render(<MillerColumnsSelect<Item> items={items} mode="multiple" />);
    expect(screen.getByText('Group A')).toBeDefined();
    expect(screen.getByText('Group B')).toBeDefined();
  });
});
