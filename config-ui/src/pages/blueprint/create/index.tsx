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

import React, { useMemo } from 'react'
import { useParams } from 'react-router-dom'

import { PageHeader } from '@/components'
import { ConnectionContextProvider } from '@/store'

import { BPContext, BPContextProvider } from './bp-context'

import { FromEnum, ModeEnum } from './types'
import { WorkFlow, Action, Inspector } from './components'
import { StepOne } from './step-one'
import { StepTwo } from './step-two'
import { StepThree } from './step-three'
import { StepFour } from './step-four'
import * as S from './styled'

interface Props {
  from: FromEnum
}

export const CreateBlueprintPage = ({ from }: Props) => {
  const { pname } = useParams<{ pname: string }>()

  const breadcrumbs = useMemo(
    () =>
      from === FromEnum.project
        ? [
            { name: 'Projects', path: '/projects' },
            { name: pname, path: `/projects/${pname}` },
            {
              name: 'Create a Blueprint',
              path: `/projects/${pname}/create-blueprint`
            }
          ]
        : [
            { name: 'Blueprints', path: '/blueprints' },
            { name: 'Create a Blueprint', path: '/blueprints/create' }
          ],
    [from, pname]
  )

  return (
    <ConnectionContextProvider>
      <BPContextProvider from={from} projectName={pname}>
        <BPContext.Consumer>
          {({ step, mode }) => (
            <PageHeader breadcrumbs={breadcrumbs}>
              <S.Container>
                <WorkFlow />
                <S.Content>
                  {step === 1 && <StepOne />}
                  {mode === ModeEnum.normal && step === 2 && <StepTwo />}
                  {step === 3 && <StepThree />}
                  {((mode === ModeEnum.normal && step === 4) ||
                    (mode === ModeEnum.advanced && step === 2)) && <StepFour />}
                </S.Content>
                <Action />
                <Inspector />
              </S.Container>
            </PageHeader>
          )}
        </BPContext.Consumer>
      </BPContextProvider>
    </ConnectionContextProvider>
  )
}
