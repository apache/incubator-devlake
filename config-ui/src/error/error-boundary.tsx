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
 */

import React from 'react'

import { Logo } from '@/components'

import { Error } from './types'
import { DBMigrate, Offline, Default } from './components'

import * as S from './styled'

type Props = {
  children: React.ReactNode
}

type State = {
  hasError: boolean
  error?: any
}

export class ErrorBoundary extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props)
    this.state = { hasError: false }
  }

  static getDerivedStateFromError(error: any) {
    return {
      hasError: true,
      error
    }
  }

  handleResetError = () => {
    this.setState({
      hasError: false,
      error: undefined
    })
  }

  render() {
    const { hasError, error } = this.state

    if (!hasError) {
      return this.props.children
    }

    return (
      <S.Wrapper>
        <Logo />
        <S.Inner>
          {error === Error.DB_NEED_MIGRATE && (
            <DBMigrate onResetError={this.handleResetError} />
          )}
          {error === Error.API_OFFLINE && (
            <Offline onResetError={this.handleResetError} />
          )}
          {!Object.keys(Error).includes(error) && (
            <Default error={error} onResetError={this.handleResetError} />
          )}
        </S.Inner>
      </S.Wrapper>
    )
  }
}
