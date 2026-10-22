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

import { beforeEach, describe, expect, it, vi } from 'vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { ThemeProvider } from 'styled-components';

import { getTheme } from '@/theme/tokens';
import { IPipeline, IPipelineStatus } from '@/types';

import { PipelineInfo } from './info';

const state = vi.hoisted(() => ({ data: undefined as IPipeline | undefined, rerun: vi.fn() }));
vi.mock('@/hooks', () => ({ useAutoRefresh: () => ({ data: state.data }) }));
vi.mock('@/api', () => ({ default: { pipeline: { rerun: state.rerun } } }));
vi.mock('@/components', () => ({ Loading: () => <span>Loading</span> }));
vi.mock('@/utils', async () => ({
  ...(await import('@/utils/operator')),
  ...(await import('@/utils/time')),
}));

const terminalStatuses = [
  IPipelineStatus.COMPLETED,
  IPipelineStatus.PARTIAL,
  IPipelineStatus.FAILED,
  IPipelineStatus.CANCELLED,
];

function renderPipeline(status: IPipelineStatus, totalTasks: number) {
  state.data = {
    id: 42,
    status,
    totalTasks,
    finishedTasks: 0,
    stage: 0,
    beganAt: null,
    finishedAt: null,
    message: 'Failed to generate pipeline plan; fix the blueprint and trigger a new run',
  };
  return render(
    <ThemeProvider theme={getTheme('light')}>
      <PipelineInfo id={42} />
    </ThemeProvider>,
  );
}

describe('pipeline task reruns', () => {
  beforeEach(() => {
    state.rerun.mockReset().mockResolvedValue(undefined);
  });

  it.each(terminalStatuses)('does not offer a task rerun for %s with no tasks', (status) => {
    renderPipeline(status, 0);
    expect(screen.queryByRole('button')).toBeNull();
    expect(state.rerun).not.toHaveBeenCalled();
  });

  it('keeps the planning error and recovery instructions visible', () => {
    renderPipeline(IPipelineStatus.FAILED, 0);
    expect(screen.getByText('Failed to generate pipeline plan; fix the blueprint and trigger a new run')).toBeTruthy();
    expect(screen.getByText('0/0')).toBeTruthy();
  });

  it.each(terminalStatuses)('allows existing tasks to be rerun for %s', async (status) => {
    renderPipeline(status, 1);
    fireEvent.click(screen.getByRole('button', { name: 'Rerun pipeline' }));
    await waitFor(() => expect(state.rerun).toHaveBeenCalledExactlyOnceWith(42));
    await waitFor(() =>
      expect(screen.getByRole('button', { name: 'Rerun pipeline' }).getAttribute('disabled')).toBeNull(),
    );
  });

  it('does not offer reruns while the pipeline is running', () => {
    renderPipeline(IPipelineStatus.RUNNING, 1);
    expect(screen.queryByRole('button', { name: 'Rerun pipeline' })).toBeNull();
  });
});
