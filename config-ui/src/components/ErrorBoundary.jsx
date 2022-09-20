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
import React, { Fragment } from 'react'
import { withRouter } from 'react-router-dom'
import {
  Button,
  Intent,
  Icon,
  Colors,
  Elevation,
  Card
} from '@blueprintjs/core'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import AppCrumbs from '@/components/Breadcrumbs'
import Content from '@/components/Content'
import ContentLoader from '@/components/loaders/ContentLoader'
import { ReactComponent as Logo } from '@/images/devlake-logo.svg'
import { ReactComponent as LogoText } from '@/images/devlake-textmark.svg'

class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      hasError: false,
      error: null,
      errorInfo: null
    }
  }

  static getDerivedStateFromError(error) {
    console.log('>>> DEVLAKE APP ERROR:', error)
    return { hasError: true, error: error }
  }

  componentDidCatch(error, errorInfo) {
    console.log('>>> DEVLAKE ERROR STACKTRACE:', errorInfo, error)
  }

  render() {
    if (this.state.hasError) {
      return (
        <>
          <div className='container'>
            <Nav />
            <Content>
              <main className='main' style={{ marginLeft: 0 }}>
                <div className='headlineContainer'>
                  <div style={{ display: 'flex', justifyContent: 'center' }}>
                    <div>
                      <div className='devlake-logo' style={{ margin: 0 }}>
                        <Logo width={48} height={48} className='logo' />
                        <LogoText
                          width={100}
                          height={13}
                          className='logo-textmark'
                        />
                      </div>
                      <h1 style={{ margin: 0, textAlign: 'center' }}>
                        Application Error
                      </h1>
                      <Card
                        elevation={Elevation.TWO}
                        style={{ margin: '18px 0', maxWidth: '700px' }}
                      >
                        <h2 style={{ margin: 0 }}>
                          <span
                            style={{
                              display: 'inline-block',
                              marginRight: '10px'
                            }}
                          >
                            <Icon icon='error' color={Colors.RED5} size={16} />
                          </span>
                          {this.state.error?.toString() || 'Unknown Error'}
                        </h2>
                        <p style={{ margin: 0 }}>
                          Please try again, if the problem persists include the
                          above error message when filing a bug report on{' '}
                          <strong>GitHub</strong>. You can also message us on{' '}
                          <strong>Slack</strong> to engage with community
                          members for solutions to common issues.
                        </p>
                        <p
                          style={{ margin: '18px 0 0 0', textAlign: 'center' }}
                        >
                          <Button
                            onClick={() => this.props.history.push('/')}
                            text='Continue'
                            intent={Intent.PRIMARY}
                          />
                          <a
                            href='https://github.com/apache/incubator-devlake'
                            className='bp3-button bp3-outlined'
                            target='_blank'
                            style={{ marginLeft: '10px' }}
                            rel='noreferrer'
                          >
                            Visit GitHub
                          </a>
                        </p>
                      </Card>
                    </div>
                  </div>
                </div>
              </main>
            </Content>
          </div>
        </>
      )
    }

    return this.props.children
  }
}
export default withRouter(ErrorBoundary)
