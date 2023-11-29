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

import parser from 'cron-parser';

export const cronPresets = [
  {
    label: 'Daily',
    config: '0 0 * * *',
    description: '(at 00:00AM UTC)',
  },
  {
    label: 'Weekly',
    config: '0 0 * * 1',
    description: '(on Monday at 00:00AM UTC)',
  },
  {
    label: 'Monthly',
    config: '0 0 1 * *',
    description: '(on first day of the month at 00:00AM UTC)',
  },
];

const getNextTime = (config: string) => {
  try {
    return parser.parseExpression(config, { tz: 'utc' }).next().toString();
  } catch {
    return null;
  }
};

const getNextTimes = (config: string) => {
  try {
    return parser
      .parseExpression(config, { tz: 'utc' })
      .iterate(3)
      .map((date) => date.toString());
  } catch {
    return [];
  }
};

export const getCron = (isManual: boolean, config: string) => {
  if (isManual) {
    return {
      label: 'Manual',
      config: '',
      description: '',
      nextTime: '',
      nextTimes: [],
    };
  }

  const preset = cronPresets.find((preset) => preset.config === config);

  return preset
    ? {
        ...preset,
        nextTime: getNextTime(preset.config),
        nextTimes: getNextTimes(preset.config),
      }
    : {
        label: 'Custom',
        config,
        description: '',
        nextTime: getNextTime(config),
        nextTimes: getNextTimes(config),
      };
};

export const getCronOptions = () => {
  return [
    {
      label: 'Manual',
      value: 'manual',
      subLabel: '',
    },
  ]
    .concat(
      cronPresets.map((cp) => ({
        label: cp.label,
        value: cp.config,
        subLabel: cp.description,
      })),
    )
    .concat([
      {
        label: 'Custom',
        value: 'custom',
        subLabel: '',
      },
    ]);
};
