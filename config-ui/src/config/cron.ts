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

export const cronPresets = [
  {
    label: 'Hourly',
    config: '59 * * * *',
    description: 'At minute 59 on every day-of-week from Monday through Sunday.'
  },
  {
    label: 'Daily',
    config: '0 0 * * *',
    description:
      'At 00:00 (Midnight) on every day-of-week from Monday through Sunday.'
  },
  {
    label: 'Weekly',
    config: '0 0 * * 1',
    description: 'At 00:00 (Midnight) on Monday.'
  },
  {
    label: 'Monthly',
    config: '0 0 1 * *',
    description: 'At 00:00 (Midnight) on day-of-month 1.'
  }
]

export const getCron = (isManual: boolean, config: string) => {
  if (isManual) {
    return {
      label: 'Manual',
      value: 'manual',
      description: 'Manual',
      config: ''
    }
  }

  const preset = cronPresets.find((preset) => preset.config === config)

  return preset
    ? {
        ...preset,
        value: preset.config
      }
    : {
        label: 'Custom',
        value: 'custom',
        description: 'Custom',
        config
      }
}

export const getCronOptions = () => {
  return [
    {
      label: 'Manual',
      value: 'manual'
    }
  ]
    .concat(
      cronPresets.map((cp) => ({
        label: cp.label,
        value: cp.config
      }))
    )
    .concat([
      {
        label: 'Custom',
        value: 'custom'
      }
    ])
}
