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
import { flatten } from 'lodash';

import API from '@/api';
import { RootState } from '@/app/store';
import { PluginConfig } from '@/plugins';
import { IConnection, IConnectionStatus, IWebhook, IStatus } from '@/types';

import { transformConnection, transformWebhook } from './utils';

const initialState: {
  connections: IConnection[];
  webhooks: IWebhook[];
  status: IStatus;
} = {
  connections: [],
  webhooks: [],
  status: 'idle',
};

export const init = createAsyncThunk('connections/init', async () => {
  const getConnections = async (plugin: string) => {
    try {
      return await API.connection.list(plugin);
    } catch {
      return [];
    }
  };

  const connections = await Promise.all(
    PluginConfig.filter((p) => p.plugin !== 'webhook').map(async ({ plugin }) => {
      const connections = await getConnections(plugin);
      return connections.map((connection) => transformConnection(plugin, connection));
    }),
  );

  const webhooks = await Promise.all(
    PluginConfig.filter((p) => p.plugin === 'webhook').map(async () => {
      const webhooks = await API.plugin.webhook.list();
      return webhooks.map((webhook) => transformWebhook(webhook));
    }),
  );

  return {
    connections: flatten(connections),
    webhooks: flatten(webhooks),
  };
});

export const addConnection = createAsyncThunk('connections/addConnection', async ({ plugin, ...payload }: any) => {
  const connection = await API.connection.create(plugin, payload);
  return transformConnection(plugin, connection);
});

export const updateConnection = createAsyncThunk(
  'connections/updateConnection',
  async ({ plugin, connectionId, ...payload }: any) => {
    const connection = await API.connection.update(plugin, connectionId, payload);
    return transformConnection(plugin, connection);
  },
);

export const removeConnection = createAsyncThunk(
  'connections/removeConnection',
  async ({ plugin, connectionId }: any) => {
    await API.connection.remove(plugin, connectionId);
    return `${plugin}-${connectionId}`;
  },
);

export const testConnection = createAsyncThunk(
  'connections/testConnection',
  async ({ unique, plugin, endpoint, proxy, token, username, password, authMethod, secretKey, appId }: IConnection) => {
    const res = await API.connection.test(plugin, {
      endpoint,
      proxy,
      token,
      username,
      password,
      authMethod,
      secretKey,
      appId,
    });

    return {
      unique,
      status: res.success ? IConnectionStatus.ONLINE : IConnectionStatus.OFFLINE,
    };
  },
);

export const addWebhook = createAsyncThunk('connections/addWebhook', async (payload: any) => {
  const webhook = await API.plugin.webhook.create(payload);
  return transformWebhook(webhook);
});

export const removeWebhook = createAsyncThunk('connections/removeWebhook', async (id: ID) => {
  await API.plugin.webhook.remove(id);
  return id;
});

export const updateWebhook = createAsyncThunk('connections/updateWebhook', async ({ id, ...payload }: any) => {
  const webhook = await API.plugin.webhook.update(id, payload);
  return webhook;
});

export const renewWebhookApiKey = createAsyncThunk('connections/renewWebhookApiKey', async (id: ID, { getState }) => {
  const webhook = (getState() as RootState).connections.webhooks.find((wh) => wh.id === id) as IWebhook;
  const apiKey = await API.apiKey.renew(webhook.apiKeyId);
  return {
    id: webhook.id,
    apiKey: apiKey.apiKey,
  };
});

export const ConnectionsSlice = createSlice({
  name: 'connections',
  initialState,
  reducers: {},
  extraReducers(builder) {
    builder
      .addCase(init.pending, (state) => {
        state.status = 'loading';
      })
      .addCase(init.fulfilled, (state, action) => {
        state.connections = action.payload.connections;
        state.webhooks = action.payload.webhooks;
        state.status = 'success';
      })
      .addCase(addConnection.fulfilled, (state, action) => {
        state.connections.push(action.payload);
      })
      .addCase(updateConnection.fulfilled, (state, action) => {
        state.connections = state.connections.map((cs) => {
          if (cs.unique === action.payload.unique) {
            return action.payload;
          }
          return cs;
        });
      })
      .addCase(removeConnection.fulfilled, (state, action) => {
        state.connections = state.connections.filter((cs) => cs.unique !== action.payload);
      })
      .addCase(testConnection.pending, (state, action) => {
        const existingConnection = state.connections.find((cs) => cs.unique === action.meta.arg.unique);
        if (existingConnection) {
          existingConnection.status = IConnectionStatus.TESTING;
        }
      })
      .addCase(testConnection.fulfilled, (state, action) => {
        const existingConnection = state.connections.find((cs) => cs.unique === action.payload.unique);
        if (existingConnection) {
          existingConnection.status = action.payload.status;
        }
      })
      .addCase(addWebhook.fulfilled, (state, action) => {
        state.webhooks.push(action.payload);
      })
      .addCase(removeWebhook.fulfilled, (state, action) => {
        state.webhooks = state.webhooks.filter((wh) => wh.id !== action.payload);
      })
      .addCase(updateWebhook.fulfilled, (state, action) => {
        state.webhooks = state.webhooks.map((wh) =>
          wh.id === action.payload.id ? { ...wh, name: action.payload.name } : wh,
        );
      })
      .addCase(renewWebhookApiKey.fulfilled, (state, action) => {
        state.webhooks = state.webhooks.map((wh) =>
          wh.id === action.payload.id ? { ...wh, apiKey: action.payload.apiKey } : wh,
        );
      });
  },
});

export const {} = ConnectionsSlice.actions;

export default ConnectionsSlice.reducer;

export const selectStatus = (state: RootState) => state.connections.status;

export const selectAllConnections = (state: RootState) => state.connections.connections;

export const selectConnections = (state: RootState, plugin: string) =>
  state.connections.connections.filter((connection) => connection.plugin === plugin);

export const selectConnection = (state: RootState, unique: string) =>
  state.connections.connections.find((cs) => cs.unique === unique);

export const selectWebhooks = (state: RootState) => state.connections.webhooks;

export const selectWebhook = (state: RootState, id: ID) => state.connections.webhooks.find((wh) => wh.id === id);
