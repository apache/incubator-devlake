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

import React from 'react'

import { Divider } from '@/components'

import {
  IssueTracking,
  CiCd,
  CodeReview,
  AdditionalSettings
} from './components'

import * as S from './styled'

interface Props {
  transformation: any
  setTransformation: React.Dispatch<React.SetStateAction<any>>
}

export const GitHubTransformation = ({ ...props }: Props) => {
  return (
    <S.TransformationWrapper>
      <IssueTracking {...props} />
      <Divider />
      <CiCd {...props} />
      <Divider />
      <CodeReview {...props} />
      <AdditionalSettings {...props} />
    </S.TransformationWrapper>
  )
}
