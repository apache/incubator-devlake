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
import React, {useCallback, useEffect, useMemo} from 'react'
import {Tag} from '@blueprintjs/core'
import {TimePrecision} from "@blueprintjs/datetime"
import {DateInput2} from "@blueprintjs/datetime2"
import dayjs from '@/utils/time'

const StartFromSelector = (
  {
    placeholder = 'Select start from',
    disabled = false,
    autoFillDefault = false,
    date,
    onSave,
  }: {
    placeholder?: string,
    disabled?: boolean,
    autoFillDefault?: boolean,
    date: string | null,
    onSave: (newDate: string | null, isUserChange?: boolean) => void,
  }) => {

  const formatDate = useCallback((date: Date | string) => dayjs(date).format("YYYY-MM-DD[T]HH:mm:ssZ"), []);
  const displayDate = useCallback((date: Date) => dayjs(date).format('L LTS'), []);
  const parseDate = useCallback((str: string) => dayjs(str).toDate(), []);
  const chooseDate = (date: Date | string | null) => {
    onSave(date ? formatDate(date) : null)
  }

  const quickDates = useMemo(() => {
    const now = new Date()
    now.setHours(0, 0, 0, 0)
    const ago6m = new Date(now)
    ago6m.setMonth(now.getMonth() - 6)
    const ago90d = new Date(now)
    ago90d.setDate(ago90d.getDate() - 90)
    const ago30d = new Date(now)
    ago30d.setDate(ago30d.getDate() - 30)
    const ago1y = new Date(now)
    ago1y.setUTCFullYear(ago1y.getFullYear() - 1)
    return [
      {label: 'Last 6 months', date: ago6m},
      {label: 'Last 90 days', date: ago90d},
      {label: 'Last 30 days', date: ago30d},
      {label: 'Last Year', date: ago1y}
    ]
  }, [])

  useEffect(() => {
    if (!date && autoFillDefault) {
      onSave(formatDate(quickDates[0].date))
    }
  }, [date])

  return (
    <>
      <div
        className='start-from'
      >
        <div
          className="123"
          style={{display: 'flex', marginBottom: '10px'}}>
          {quickDates.map(quickDate => <Tag
            key={quickDate.date.toISOString()} minimal={formatDate(quickDate.date) != date}
            intent="primary"
            interactive={!disabled}
            onClick={() => chooseDate(quickDate.date)}
            style={{marginRight: 5}}>{quickDate.label}
          </Tag>)}
        </div>

        <div style={{display: "flex", alignItems: "baseline"}}>
          <DateInput2
            disabled={disabled}
            timePrecision={TimePrecision.MINUTE}
            formatDate={displayDate}
            parseDate={parseDate}
            fill={false}
            placeholder={placeholder}
            onChange={chooseDate}
            popoverProps={{placement: "bottom"}}
            value={date}
          />
          <span style={{fontWeight: 'bold'}}> to Now</span>
        </div>
      </div>
    </>
  )
}

export default StartFromSelector
