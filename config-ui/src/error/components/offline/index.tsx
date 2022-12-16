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
  Icon,
  Tag,
  ButtonGroup,
  Button,
  Intent,
  Colors,
  IconName
} from '@blueprintjs/core'

import { Card } from '@/components'
import { DEVLAKE_ENDPOINT } from '@/config'

import type { UseOfflineProps } from './use-offline'
import { useOffline } from './use-offline'

interface Props extends UseOfflineProps {}

export const Offline = ({ ...props }: Props) => {
  const { processing, offline, onRefresh, onContinue } = useOffline({
    ...props
  })

  const [icon, color, text] = useMemo(
    () => [
      offline ? 'offline' : 'endorsed',
      offline ? Colors.RED3 : Colors.GREEN3,
      offline ? 'Offline' : 'Online'
    ],
    [offline]
  )

  return (
    <Card>
      <h2>
        <Icon icon={icon as IconName} color={color} size={30} />
        <span>DevLake API</span>
        <strong style={{ marginLeft: 4, color }}>{text}</strong>
      </h2>
      <p>
        <Tag>DEVLAKE_ENDPOINT: {DEVLAKE_ENDPOINT}</Tag>
      </p>
      {offline ? (
        <>
          <p>
            Please wait for the&nbsp;
            <strong>Lake API</strong> to start before accessing the{' '}
            <strong>Configuration Interface</strong>.
          </p>
          <ButtonGroup>
            <Button
              loading={processing}
              icon='refresh'
              intent={Intent.PRIMARY}
              text='Refresh'
              onClick={onRefresh}
            />
          </ButtonGroup>
        </>
      ) : (
        <>
          <p>Connectivity to the Lake API service was successful.</p>
          <ButtonGroup>
            <Button
              intent={Intent.PRIMARY}
              text='Continue'
              onClick={onContinue}
            />
            <Button
              icon='help'
              text='Read Documentation'
              onClick={() =>
                window.open(
                  'https://github.com/apache/incubator-devlake/blob/main/README.md',
                  '_blank',
                  'noopener,noreferrer'
                )
              }
            />
          </ButtonGroup>
        </>
      )}
    </Card>
  )
}
