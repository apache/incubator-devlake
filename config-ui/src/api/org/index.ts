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

import { request } from '@/utils';

// The org plugin stores teams, users and their account mapping as whole CSV files:
// every PUT replaces the entire table (deleteAll + save), so the UI always sends
// back the complete set of rows rather than a delta.
const upload = (path: string, csv: string, filename: string) => {
  const data = new FormData();
  data.append('file', new Blob([csv], { type: 'text/csv' }), filename);
  return request(path, { method: 'put', data });
};

export const getTeams = (): Promise<string> => request('/plugins/org/teams.csv');

export const updateTeams = (csv: string) => upload('/plugins/org/teams.csv', csv, 'teams.csv');

export const getUsers = (): Promise<string> => request('/plugins/org/users.csv');

export const updateUsers = (csv: string) => upload('/plugins/org/users.csv', csv, 'users.csv');

// Returns every known account (one per connected tool identity) with the UserId it
// is currently mapped to, not only the accounts that already have one.
export const getUserAccountMapping = (): Promise<string> => request('/plugins/org/user_account_mapping.csv');

export const updateUserAccountMapping = (csv: string) =>
  upload('/plugins/org/user_account_mapping.csv', csv, 'user_account_mapping.csv');
