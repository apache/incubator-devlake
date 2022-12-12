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
import React, {useCallback, useEffect} from 'react'
import {Tag} from '@blueprintjs/core'
import {TimePrecision} from "@blueprintjs/datetime"
import {DateInput2} from "@blueprintjs/datetime2"
import { format, parse } from "date-fns"

const StartFromSelector = (
  {
    placeholder = 'Select start from',
    disabled = false,
    date,
    onSave,
  }: {
    placeholder?: string,
    disabled?: boolean,
    date: string | null,
    onSave: (newDate: string | null, isUserChange?: boolean) => void,
  }) => {

  const formatDate = useCallback((date: Date) => format(date, "yyyy-MM-dd HH:mm"), []);
  const parseDate = useCallback((str: string) => new Date(str), []);
  const chooseDate = (date: Date) => {
    onSave(format(date, "yyyy-MM-dd'T'HH:mm:ssxxx"))
  }

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

  useEffect(() => {
    if (!date) {
      onSave(ago6m.toISOString())
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
          <Tag
            key="-6m" minimal={ago6m.toISOString()!=date}
            intent="primary"
            interactive={!disabled}
            onClick={() => chooseDate(ago6m)}
            style={{marginRight: 5}}>Last 6 months</Tag>
          <Tag
            key="-90d" minimal={ago90d.toISOString()!=date}
            intent="primary"
            interactive={!disabled}
            onClick={() => chooseDate(ago90d)}
            style={{marginRight: 5}}>Last 90 days</Tag>
          <Tag
            key="-30d" minimal={ago30d.toISOString()!=date}
            intent="primary"
            interactive={!disabled}
            onClick={() => chooseDate(ago30d)}
            style={{marginRight: 5}}>Last 30 days</Tag>
          <Tag
            key="-1y" minimal={ago1y.toISOString()!=date}
            intent="primary"
            interactive={!disabled}
            onClick={() => chooseDate(ago1y)}
            style={{marginRight: 5}}>Last Year</Tag>
        </div>

        <DateInput2
          disabled={disabled}
          timePrecision={TimePrecision.MINUTE}
          formatDate={formatDate}
          parseDate={parseDate}
          fill={false}
          placeholder={placeholder}
          onChange={onSave}
          popoverProps={{placement: "bottom"}}
          value={date}
        />
        <span style={{ fontWeight: 'bold'}}> to Now</span>
      </div>
    </>
  )
}

export default StartFromSelector
