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

import { Button, Intent } from '@blueprintjs/core';

import { Card, Divider } from '@/components';

import { ModeEnum } from '../../types';
import { SyncPolicy } from '../../components';

import { useCreate } from '../context';

import * as S from './styled';

export const Step4 = () => {
  const {
    mode,
    isManual,
    cronConfig,
    skipOnFail,
    timeAfter,
    onChangeIsManual,
    onChangeCronConfig,
    onChangeSkipOnFail,
    onChangeTimeAfter,
    onPrev,
    onSave,
    onSaveAndRun,
  } = useCreate();

  return (
    <S.Wrapper>
      <Card>
        <h2>Set Sync Policy</h2>
        <Divider />
        <SyncPolicy
          isManual={isManual}
          cronConfig={cronConfig}
          skipOnFail={skipOnFail}
          showTimeFilter={mode === ModeEnum.normal}
          timeAfter={timeAfter}
          onChangeIsManual={onChangeIsManual}
          onChangeCronConfig={onChangeCronConfig}
          onChangeSkipOnFail={onChangeSkipOnFail}
          onChangeTimeAfter={onChangeTimeAfter}
        />
      </Card>
      <S.Btns>
        <Button intent={Intent.PRIMARY} outlined text="Previous Step" onClick={onPrev} />
        <div>
          <Button intent={Intent.PRIMARY} outlined text="Save Blueprint" onClick={onSave} />
          <Button intent={Intent.PRIMARY} text="Save and Run Now" onClick={onSaveAndRun} />
        </div>
      </S.Btns>
    </S.Wrapper>
  );
};
