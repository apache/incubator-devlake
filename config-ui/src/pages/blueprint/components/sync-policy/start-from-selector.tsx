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
import React, { useCallback, useMemo } from 'react';
import { Tag } from '@blueprintjs/core';
import { TimePrecision } from '@blueprintjs/datetime';
import { DateInput2 } from '@blueprintjs/datetime2';

import { formatTime } from '@/utils';

import * as S from './styled';

interface Props {
  value: string | null;
  onChange: (val: string | null) => void;
}

const StartFromSelector = ({ value, onChange }: Props) => {
  const formatDate = useCallback((date: Date) => formatTime(date), []);
  const parseDate = useCallback((str: string) => new Date(str), []);

  const quickDates = useMemo(() => {
    const now = new Date();
    now.setHours(0, 0, 0, 0);
    const ago6m = new Date(now);
    ago6m.setMonth(now.getMonth() - 6);
    const ago90d = new Date(now);
    ago90d.setDate(ago90d.getDate() - 90);
    const ago30d = new Date(now);
    ago30d.setDate(ago30d.getDate() - 30);
    const ago1y = new Date(now);
    ago1y.setUTCFullYear(ago1y.getFullYear() - 1);

    return [
      { label: 'Last 6 months', date: ago6m },
      { label: 'Last 90 days', date: ago90d },
      { label: 'Last 30 days', date: ago30d },
      { label: 'Last Year', date: ago1y },
    ];
  }, []);

  const handleChangeDate = (val: string | Date | null) => {
    onChange(val ? formatTime(val, 'YYYY-MM-DD[T]HH:mm:ssZ') : null);
  };

  return (
    <S.FromTimeWrapper>
      <div className="quick-selection">
        {quickDates.map((quickDate) => (
          <Tag
            key={quickDate.date.toISOString()}
            minimal={formatDate(quickDate.date) !== formatTime(value)}
            intent="primary"
            onClick={() => handleChangeDate(quickDate.date)}
            style={{ marginRight: 5 }}
          >
            {quickDate.label}
          </Tag>
        ))}
      </div>

      <div className="time-selection">
        <DateInput2
          timePrecision={TimePrecision.MINUTE}
          formatDate={formatDate}
          parseDate={parseDate}
          placeholder="Select start from"
          popoverProps={{ placement: 'bottom' }}
          value={value}
          onChange={handleChangeDate}
        />{' '}
        <strong>to Now</strong>
      </div>
    </S.FromTimeWrapper>
  );
};

export default StartFromSelector;
