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

import { useContext } from 'react';

import { ConnectionContext } from '@/store';

type UseConnectionsProps = {
  unique?: string;
  plugin?: string;
  filter?: string[];
  filterBeta?: boolean;
  filterPlugin?: string[];
};

export const useConnections = (props?: UseConnectionsProps) => {
  const { unique, plugin, filter, filterBeta, filterPlugin } = props || {};

  const { connections, onGet, onTest, onRefresh } = useContext(ConnectionContext);

  return {
    connection: unique ? connections.find((cs) => cs.unique === unique) : null,
    connections: connections
      .filter((cs) => (plugin ? cs.plugin === plugin : true))
      .filter((cs) => (filter ? !filter.includes(cs.unique) : true))
      .filter((cs) => (filterBeta ? !cs.isBeta : true))
      .filter((cs) => (filterPlugin ? !filterPlugin.includes(cs.plugin) : true)),
    onGet,
    onTest,
    onRefresh,
  };
};
