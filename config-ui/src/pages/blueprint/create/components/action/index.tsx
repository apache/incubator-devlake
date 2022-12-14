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
import {
  ButtonGroup,
  Button,
  Icon,
  Intent,
  Position,
  Colors
} from '@blueprintjs/core'
import { Popover2 } from '@blueprintjs/popover2'

import { ModeEnum } from '../../types'
import { useCreateBP } from '../../bp-context'

import * as S from './styled'

export const Action = () => {
  const {
    step,
    mode,
    error,
    showDetail,
    onChangeStep,
    onChangeShowInspector,
    onSave,
    onSaveAndRun
  } = useCreateBP()

  const [isFirst, isLast] = useMemo(() => {
    return [
      step === 1,
      (mode === ModeEnum.normal && step === 4) ||
        (mode === ModeEnum.advanced && step === 2)
    ]
  }, [step, mode])

  if (showDetail) {
    return null
  }

  return (
    <S.Container>
      <ButtonGroup>
        {!isFirst && (
          <Button
            outlined
            intent={Intent.PRIMARY}
            text='Previous Step'
            onClick={() => onChangeStep(step - 1)}
          />
        )}
      </ButtonGroup>
      <ButtonGroup>
        <Button
          minimal
          intent={Intent.PRIMARY}
          icon='code'
          text='Inspect'
          onClick={() => onChangeShowInspector(true)}
        />
        {isLast ? (
          <>
            <Button
              intent={Intent.PRIMARY}
              text='Save Blueprint'
              onClick={onSave}
            />
            <Button
              intent={Intent.DANGER}
              text='Save and Run Now'
              onClick={onSaveAndRun}
            />
          </>
        ) : (
          <Button
            intent={Intent.PRIMARY}
            disabled={!!error}
            icon={
              error ? (
                <Popover2
                  defaultIsOpen
                  placement={Position.TOP}
                  content={
                    <S.Error>
                      <Icon icon='warning-sign' color={Colors.ORANGE5} />
                      <span>{error}</span>
                    </S.Error>
                  }
                >
                  <Icon
                    icon='warning-sign'
                    color={Colors.ORANGE5}
                    style={{ margin: 0 }}
                  />
                </Popover2>
              ) : null
            }
            text='Next Step'
            onClick={() => onChangeStep(step + 1)}
          />
        )}
      </ButtonGroup>
    </S.Container>
  )
}
