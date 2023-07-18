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

import { useMemo } from 'react';
import dayjs from 'dayjs';
import { Tag, Checkbox, FormGroup, InputGroup, Radio, RadioGroup } from '@blueprintjs/core';
import { TimePrecision } from '@blueprintjs/datetime';
import { DateInput2 } from '@blueprintjs/datetime2';

import { FormItem, ExternalLink } from '@/components';
import { getCron, getCronOptions } from '@/config';
import { formatTime } from '@/utils';

import * as S from './styled';

interface Props {
  isManual: boolean;
  cronConfig: string;
  skipOnFail: boolean;
  showTimeFilter: boolean;
  timeAfter: string | null;
  onChangeIsManual: (val: boolean) => void;
  onChangeCronConfig: (val: string) => void;
  onChangeSkipOnFail: (val: boolean) => void;
  onChangeTimeAfter: (val: string | null) => void;
}

export const SyncPolicy = ({
  isManual,
  cronConfig,
  skipOnFail,
  showTimeFilter,
  timeAfter,
  onChangeIsManual,
  onChangeCronConfig,
  onChangeSkipOnFail,
  onChangeTimeAfter,
}: Props) => {
  const [timezone, quickTimeOpts, cronOpts] = useMemo(() => {
    const timezone = dayjs().format('ZZ').replace('00', '');
    const quickTimeOpts = [
      { label: 'Last 6 months', date: dayjs().subtract(6, 'month').toDate() },
      { label: 'Last 90 days', date: dayjs().subtract(90, 'day').toDate() },
      { label: 'Last 30 days', date: dayjs().subtract(30, 'day').toDate() },
      { label: 'Last Year', date: dayjs().subtract(1, 'year').toDate() },
    ];

    const cronOpts = getCronOptions();

    return [timezone, quickTimeOpts, cronOpts];
  }, []);

  const cron = useMemo(() => getCron(isManual, cronConfig), [isManual, cronConfig]);

  const [mintue, hour, day, month, week] = useMemo(() => cronConfig.split(' '), [cronConfig]);

  const handleChangeFrequency = (e: React.FormEvent<HTMLInputElement>) => {
    const value = (e.target as HTMLInputElement).value;
    if (value === 'manual') {
      onChangeIsManual(true);
    } else if (!value) {
      onChangeIsManual(false);
      onChangeCronConfig('* * * * *');
    } else {
      onChangeIsManual(false);
      onChangeCronConfig(value);
    }
  };

  return (
    <S.Wrapper>
      <div className="timezone">
        Your local time zone is <strong>UTC {timezone}</strong>. All time listed below is shown in your local time.
      </div>
      {showTimeFilter && (
        <FormItem
          label="Time Range"
          subLabel="Select the time range for the data you wish to collect. DevLake will collect the last six months of data by default."
        >
          <div className="quick-selection">
            {quickTimeOpts.map((opt, i) => (
              <Tag
                key={i}
                style={{ marginRight: 5, cursor: 'pointer' }}
                minimal={formatTime(opt.date) !== formatTime(timeAfter)}
                intent="primary"
                onClick={() => onChangeTimeAfter(dayjs(opt.date).utc().format('YYYY-MM-DD[T]HH:mm:ssZ'))}
              >
                {opt.label}
              </Tag>
            ))}
          </div>

          <div className="time-selection">
            <DateInput2
              timePrecision={TimePrecision.MINUTE}
              showTimezoneSelect={false}
              formatDate={formatTime}
              parseDate={(str: string) => new Date(str)}
              placeholder="Select start from"
              popoverProps={{ placement: 'bottom' }}
              value={timeAfter}
              onChange={(date) => onChangeTimeAfter(date ? dayjs(date).utc().format('YYYY-MM-DD[T]HH:mm:ssZ') : null)}
            />
            <strong>to Now</strong>
          </div>
        </FormItem>
      )}
      <div className="cron">
        <FormItem
          style={{ flex: '0 0 400px', marginRight: 50 }}
          label="Sync Frequency"
          subLabel="Blueprints will run on creation and recurringly based on the schedule."
        >
          <RadioGroup selectedValue={cron.value} onChange={handleChangeFrequency}>
            {cronOpts.map(({ value, label, subLabel }) => (
              <Radio key={value} label={`${label} ${subLabel}`} value={value} />
            ))}
          </RadioGroup>
          {!cron.value && (
            <>
              <S.Input>
                <FormGroup label="Minute">
                  <InputGroup
                    value={mintue}
                    onChange={(e) => onChangeCronConfig([e.target.value, hour, day, month, week].join(' '))}
                  />
                </FormGroup>
                <FormGroup label="Hour">
                  <InputGroup
                    value={hour}
                    onChange={(e) => onChangeCronConfig([mintue, e.target.value, day, month, week].join(' '))}
                  />
                </FormGroup>
                <FormGroup label="Day">
                  <InputGroup
                    value={day}
                    onChange={(e) => onChangeCronConfig([mintue, hour, e.target.value, month, week].join(' '))}
                  />
                </FormGroup>
                <FormGroup label="Month">
                  <InputGroup
                    value={month}
                    onChange={(e) => onChangeCronConfig([mintue, hour, day, e.target.value, week].join(' '))}
                  />
                </FormGroup>
                <FormGroup label="Week">
                  <InputGroup
                    value={week}
                    onChange={(e) => onChangeCronConfig([mintue, hour, day, month, e.target.value].join(' '))}
                  />
                </FormGroup>
              </S.Input>
              {!cron.nextTime && <S.Error>Invalid Cron code, please enter again.</S.Error>}
            </>
          )}
          <ExternalLink style={{ display: 'block', marginTop: 16 }} link="https://crontab.cronhub.io/">
            Learn about how to use cron code
          </ExternalLink>
        </FormItem>
        <FormItem label="Next Three Runs:">
          {cron.nextTimes.length ? (
            <ul>
              {cron.nextTimes.map((it, i) => (
                <li key={i}>
                  {dayjs(it).format('YYYY-MM-DD HH:mm A')}({dayjs(it).fromNow()})
                </li>
              ))}
            </ul>
          ) : (
            'N/A'
          )}
        </FormItem>
      </div>
      <FormItem label="Running Policy">
        <Checkbox
          label="Skip failed tasks (Recommended when collecting a large volume of data, eg. 10+ GitHub repos, Jira boards, etc.)"
          checked={skipOnFail}
          onChange={(e) => onChangeSkipOnFail((e.target as HTMLInputElement).checked)}
        />
        <p style={{ paddingLeft: 28 }}>
          A task is a unit of a pipeline, an execution of a blueprint. By default, when a task is failed, the whole
          pipeline will fail and all the data that has been collected will be discarded. By skipping failed tasks, the
          pipeline will continue to run, and the data collected by successful tasks will not be affected. After the
          pipeline is finished, you can rerun these failed tasks.
        </p>
      </FormItem>
    </S.Wrapper>
  );
};
