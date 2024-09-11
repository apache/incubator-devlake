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

import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';

import API from '@/api';
import type { RootState } from '@/app/store';
import type { IStatus } from '@/types';

type DataType = {
  initial: boolean;
  step: number;
  records: Array<{
    plugin: string;
    connectionId: ID;
    blueprintId: ID;
    pipelineId: ID;
    scopeName: string;
  }>;
  projectName: string;
  plugin: string;
  done: boolean;
};

export const request = createAsyncThunk('onboard/request', async () => {
  const res = await API.store.get('onboard');
  return {
    ...res,
    initial: res ? true : res?.initial ?? false,
    step: res?.step ?? 0,
    records: res?.records ?? [],
    done: res?.done ?? false,
  };
});

export const update = createAsyncThunk('onboard/update', async (payload: Partial<DataType>, { getState }) => {
  const { data } = (getState() as RootState).onboard;
  const res = await API.store.set('onboard', {
    ...data,
    ...payload,
    initial: true,
    step: payload.step ?? data.step + 1,
  });
  return res;
});

export const done = createAsyncThunk('onboard/done', async (_, { getState }) => {
  const { data } = (getState() as RootState).onboard;
  await API.store.set('onboard', {
    ...data,
    done: true,
  });
  return {};
});

const initialState: { status: IStatus; error?: unknown; data: DataType } = {
  status: 'loading',
  data: {
    initial: true,
    step: 0,
    records: [],
    projectName: '',
    plugin: '',
    done: false,
  },
};

export const onboardSlice = createSlice({
  name: 'onboard',
  initialState,
  reducers: {
    previous: (state) => {
      state.data.step -= 1;
    },
    changeProjectName: (state, action) => {
      state.data.projectName = action.payload;
    },
    changePlugin: (state, action) => {
      state.data.plugin = action.payload;
    },
    changeRecords: (state, action) => {
      state.data.records = action.payload;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(request.pending, (state) => {
        state.status = 'loading';
      })
      .addCase(request.fulfilled, (state, action) => {
        state.status = 'success';
        state.data = action.payload;
      })
      .addCase(request.rejected, (state, action) => {
        state.status = 'failed';
        state.error = action.error;
      })
      .addCase(update.fulfilled, (state, action) => {
        state.data = {
          ...state.data,
          ...action.payload,
        };
      })
      .addCase(done.fulfilled, (state) => {
        state.data.done = true;
      });
  },
});

export default onboardSlice.reducer;

export const { previous, changeProjectName, changePlugin, changeRecords } = onboardSlice.actions;

export const selectOnboardStatus = (state: RootState) => state.onboard.status;

export const selectOnboardError = (state: RootState) => state.onboard.error;

export const selectOnboard = (state: RootState) => state.onboard.data;

export const selectRecord = (state: RootState) => {
  const { plugin, records } = state.onboard.data;
  return records.find((it) => it.plugin === plugin);
};
