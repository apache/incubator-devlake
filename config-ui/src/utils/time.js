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
import * as dayjs from 'dayjs'
import * as relativeTime from 'dayjs/plugin/relativeTime'
import * as updateLocale from 'dayjs/plugin/updateLocale'
import * as LocalizedFormat from 'dayjs/plugin/localizedFormat'
import * as utc from 'dayjs/plugin/utc'

const localeConfiguration = {
  relativeTime: {
    future: 'in %s',
    past: '%s ago',
    s: '< 1min',
    m: 'a minute',
    mm: '%d minutes',
    h: 'an hour',
    hh: '%d hours',
    d: 'a day',
    dd: '%d days',
    M: 'a month',
    MM: '%d months',
    y: 'a year',
    yy: '%d years'
  }
}

dayjs.extend(relativeTime)
dayjs.extend(updateLocale)
dayjs.extend(LocalizedFormat)
dayjs.extend(utc)
dayjs.updateLocale('en', localeConfiguration)

export default dayjs
