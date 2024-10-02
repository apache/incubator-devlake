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

import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';

import API from '@/api';
import type { RootState } from '@/app/store';
import type { IStatus, IProject } from '@/types';

export const request = createAsyncThunk('project/request', async (name: string) => {
  const res = await API.project.get(name);
  return res;
});

const initialState: { status: IStatus; data?: IProject } = {
  status: 'loading',
};

export const projectSlice = createSlice({
  name: 'project',
  initialState,
  reducers: {
    updateBlueprint: (state, action) => {
      if (state.data) {
        state.data.blueprint = action.payload;
      }
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
      .addCase(request.rejected, (state) => {
        state.status = 'failed';
      });
  },
});

export const { updateBlueprint } = projectSlice.actions;

export default projectSlice.reducer;

export const selectProjectStatus = (state: RootState) => state.project.status;

export const selectProject = (state: RootState) => state.project.data;
