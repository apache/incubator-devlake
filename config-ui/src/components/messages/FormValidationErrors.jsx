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
import { Icon, Colors } from '@blueprintjs/core'
const FormValidationErrors = (props) => {
  const { errors = [], textAlign = 'right', styles = {} } = props

  return (
    <>
      {errors.length > 0 && (
        <div className='validation-errors'>
          <p style={{ margin: '5px 0 5px 0', textAlign: textAlign, ...styles }}>
            <Icon
              icon='warning-sign'
              size={13}
              color={Colors.ORANGE4}
              style={{ marginRight: '6px', marginBottom: '2px' }}
            />
            <span>{errors[0]}</span>
          </p>
        </div>
      )}
    </>
  )
}

export default FormValidationErrors
