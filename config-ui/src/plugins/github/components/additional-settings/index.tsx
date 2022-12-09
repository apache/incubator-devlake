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

import React, { useState } from 'react'
import {
  Checkbox,
  FormGroup,
  InputGroup,
  NumericInput
} from '@blueprintjs/core'

interface Props {
  transformation: any
  setTransformation: React.Dispatch<React.SetStateAction<any>>
}

export const AdditionalSettings = ({
  transformation,
  setTransformation
}: Props) => {
  const [enable, setEnable] = useState(false)

  const handleToggleEnable = () => {
    setEnable(!enable)
  }

  return (
    <>
      <h3>Additional Settings</h3>
      <Checkbox
        checked={enable}
        label='Enable calculation of commit and issue difference'
        onChange={handleToggleEnable}
      />
      {enable && (
        <div className='additional-settings-refdiff'>
          <FormGroup inline label='Tags Limit'>
            <NumericInput
              placeholder='10'
              allowNumericCharactersOnly={true}
              value={transformation.refdiff?.tagsLimit}
              onValueChange={(tagsLimit) =>
                setTransformation({
                  ...transformation,
                  refdiff: {
                    ...transformation?.refdiff,
                    tagsLimit
                  }
                })
              }
            />
          </FormGroup>
          <FormGroup inline label='Tags Pattern'>
            <InputGroup
              placeholder='(regex)$'
              value={transformation.refdiff?.tagsPattern}
              onChange={(e) =>
                setTransformation({
                  ...transformation,
                  refdiff: {
                    ...transformation?.refdiff,
                    tagsPattern: e.target.value
                  }
                })
              }
            />
          </FormGroup>
          <FormGroup inline label='Tags Order'>
            <InputGroup
              placeholder='reverse semver'
              value={transformation.refdiff?.tagsOrder}
              onChange={(e) =>
                setTransformation({
                  ...transformation,
                  refdiff: {
                    ...transformation?.refdiff,
                    tagsOrder: e.target.value
                  }
                })
              }
            />
          </FormGroup>
        </div>
      )}
    </>
  )
}
