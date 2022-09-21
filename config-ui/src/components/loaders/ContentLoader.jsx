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
import React, { useEffect } from 'react'
import { Card, Elevation, Intent, Spinner } from '@blueprintjs/core'
const ContentLoader = (props) => {
  const {
    title = 'Loading ...',
    message = 'Please wait while data is loaded.',
    spinnerSize = 24,
    spinnerIntent = Intent.PRIMARY,
    elevation = Elevation.TWO,
    cardStyle = {
      width: '100%',
      marginBottom: '20px',
      boxShadow: elevation === Elevation.ZERO ? 'none' : 'initial'
    },
    cardStyleOverrides = {},
    messageClasses = ['bp3-ui-text', 'bp3-text-large']
  } = props

  useEffect(() => {}, [title, message, spinnerSize])

  return (
    <Card
      interactive={false}
      elevation={elevation}
      style={{ ...cardStyle, ...cardStyleOverrides }}
    >
      <div style={{}}>
        <div style={{ display: 'flex' }}>
          <Spinner intent={spinnerIntent} size={spinnerSize} />
          <div style={{ marginLeft: '10px' }}>
            <h4 className='bp3-heading' style={{ margin: '0 0 2px 0' }}>
              {title}
            </h4>
            <p className={messageClasses.join(' ')} style={{ margin: 0 }}>
              {message}
            </p>
          </div>
        </div>
      </div>
    </Card>
  )
}

export default ContentLoader
