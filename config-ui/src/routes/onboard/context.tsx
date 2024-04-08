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

import { createContext } from 'react';

export type Record = {
  plugin: string;
  connectionId: ID;
  blueprintId: ID;
  pipelineId: ID;
  scopeName: string;
};

const initialValue: {
  step: number;
  records: Record[];
  done: boolean;
  projectName?: string;
  plugin?: string;
  setStep: (value: number) => void;
  setRecords: (value: Record[]) => void;
  setProjectName: (value: string) => void;
  setPlugin: (value: string) => void;
} = {
  step: 0,
  records: [],
  done: false,
  setStep: () => {},
  setRecords: () => {},
  setProjectName: () => {},
  setPlugin: () => {},
};

export const Context = createContext(initialValue);
