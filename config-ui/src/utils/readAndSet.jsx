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
import { findStrBetween } from './findStrBetween'

export const readAndSet = (tagName, tagLen, isStatus, str, fn1, fn2) => {
  if (isStatus) {
    const strValuesReq = findStrBetween(str, 'Requirement:', ';')
    const strValuesRes = findStrBetween(str, 'Resolved:', ';')

    if (strValuesReq) fn1(strValuesReq[0].slice(12, -1).split(','))
    if (strValuesRes) fn2(strValuesRes[0].slice(9, -1).split(','))
  } else {
    const strValues = findStrBetween(str, tagName, ';')

    if (strValues) fn1(strValues[0].slice(tagLen, -1).split(','))
  }
}
