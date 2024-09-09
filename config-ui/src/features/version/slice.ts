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
import { RootState } from '@/app/store';
import type { IStatus } from '@/types';

export const request = createAsyncThunk('version/request', async () => {
  const res = await API.version();
  return res;
});

const initialState: {
  status: IStatus;
  error?: unknown;
  version: string;
} = {
  status: 'idle',
  version: '',
};

export const versionSlice = createSlice({
  name: 'version',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(request.pending, (state) => {
        state.status = 'loading';
      })
      .addCase(request.fulfilled, (state, action) => {
        state.status = 'success';
        state.version = action.payload.version;
      })
      .addCase(request.rejected, (state, action) => {
        state.status = 'failed';
        state.error = action.error;
      });
  },
});

export const selectVersionStatus = (state: RootState) => state.version.status;

export const selectVersionError = (state: RootState) => state.version.error;

export const selectVersion = (state: RootState) => state.version.version;
