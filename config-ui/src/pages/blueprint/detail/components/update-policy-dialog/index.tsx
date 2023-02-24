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

import React, { useState, useEffect } from 'react';

import { Dialog } from '@/components';

import type { BlueprintType } from '../../../types';
import { ModeEnum } from '../../../types';
import { SyncPolicy } from '../../../components';

interface Props {
  blueprint: BlueprintType;
  isManual: boolean;
  cronConfig: string;
  skipOnFail: boolean;
  timeAfter: string | null;
  operating: boolean;
  onCancel: () => void;
  onSubmit: (params: any) => Promise<void>;
}

export const UpdatePolicyDialog = ({ blueprint, operating, onCancel, onSubmit, ...props }: Props) => {
  const [isManual, setIsManual] = useState(false);
  const [cronConfig, setCronConfig] = useState('');
  const [skipOnFail, setSkipOnFail] = useState(false);
  const [timeAfter, setTimeAfter] = useState<string | null>(null);

  useEffect(() => {
    setIsManual(props.isManual);
    setCronConfig(props.cronConfig);
    setSkipOnFail(props.skipOnFail);
    setTimeAfter(props.timeAfter);
  }, []);

  const handleSubmit = () => {
    onSubmit({
      isManual,
      cronConfig,
      skipOnFail,
      settings:
        blueprint.mode === ModeEnum.normal
          ? {
              ...blueprint.settings,
              timeAfter,
            }
          : undefined,
    });
  };

  return (
    <Dialog
      isOpen
      title="Change Sync Policy"
      style={{
        width: 720,
      }}
      okText="Save"
      okLoading={operating}
      onCancel={onCancel}
      onOk={handleSubmit}
    >
      <SyncPolicy
        isManual={isManual}
        cronConfig={cronConfig}
        skipOnFail={skipOnFail}
        showTimeFilter={blueprint.mode === ModeEnum.normal}
        timeAfter={timeAfter}
        onChangeIsManual={setIsManual}
        onChangeCronConfig={setCronConfig}
        onChangeSkipOnFail={setSkipOnFail}
        onChangeTimeAfter={setTimeAfter}
      />
    </Dialog>
  );
};
