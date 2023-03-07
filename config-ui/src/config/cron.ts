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
    description: '(at 00:00 AM) ',
  },
  {
    label: 'Weekly',
    config: '0 0 * * 1',
    description: '(on Monday at 00:00 AM) ',
  },
  {
    label: 'Monthly',
    config: '0 0 1 * *',
    description: '(on first day of the month at 00:00 AM) ',
  },
];

const getNextTime = (config: string) => {
  try {
    return parser.parseExpression(config, { tz: 'utc' }).next().toString();
  } catch {
    return null;
  }
};

export const getCron = (isManual: boolean, config: string) => {
  if (isManual) {
    return {
      label: 'Manual',
      value: 'manual',
      description: 'Manual',
      config: '',
      nextTime: '',
    };
  }

  const preset = cronPresets.find((preset) => preset.config === config);

  return preset
    ? {
        ...preset,
        value: preset.config,
        nextTime: getNextTime(preset.config),
      }
    : {
        label: 'Custom',
        value: 'custom',
        description: 'Custom',
        config,
        nextTime: getNextTime(config),
      };
};

export const getCronOptions = () => {
  return [
    {
      value: 'manual',
      label: 'Manual',
      subLabel: '',
    },
  ]
    .concat(
      cronPresets.map((cp) => ({
        value: cp.config,
        label: cp.label,
        subLabel: cp.description,
      })),
    )
    .concat([
      {
        value: 'custom',
        label: 'Custom',
        subLabel: '',
      },
    ]);
};
