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

import { ModeEnum } from '../types';

export type ContextType = {
  step: number;

  mode: ModeEnum;
  name: string;
  connections: MixConnection[];
  rawPlan: string;
  cronConfig: string;
  isManual: boolean;
  skipOnFail: boolean;
  timeAfter: string | null;

  onPrev: () => void;
  onNext: () => void;

  onSave: () => void;
  onSaveAndRun: () => void;

  onChangeMode: (mode: ModeEnum) => void;
  onChangeName: React.Dispatch<React.SetStateAction<string>>;
  onChangeConnections: React.Dispatch<React.SetStateAction<MixConnection[]>>;
  onChangeRawPlan: React.Dispatch<React.SetStateAction<string>>;
  onChangeCronConfig: React.Dispatch<React.SetStateAction<string>>;
  onChangeIsManual: React.Dispatch<React.SetStateAction<boolean>>;
  onChangeSkipOnFail: React.Dispatch<React.SetStateAction<boolean>>;
  onChangeTimeAfter: React.Dispatch<React.SetStateAction<string | null>>;
};
