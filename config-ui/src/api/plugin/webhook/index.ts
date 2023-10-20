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

import * as connection from '../../connection';

export const list = () => connection.list('webhook');

export const get = (
  id: ID,
): Promise<{
  name: string;
  postIssuesEndpoint: string;
  closeIssuesEndpoint: string;
  postPipelineDeployTaskEndpoint: string;
  apiKey: {
    id: string;
    apiKey: string;
  };
}> => connection.get('webhook', id) as any;

export const create = (payload: any): Promise<{ id: string; apiKey: { apiKey: string } }> =>
  connection.create('webhook', payload) as any;

export const remove = (id: ID) => connection.remove('webhook', id);

export const update = (id: ID, payload: any) => connection.update('webhook', id, payload);
