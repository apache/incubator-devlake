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
  Tag,
  Intent,
  RadioGroup,
  Radio,
  FormGroup,
  InputGroup
} from '@blueprintjs/core'

interface Props {
  transformation: any
  setTransformation: React.Dispatch<React.SetStateAction<any>>
}

export const CiCd = ({ transformation, setTransformation }: Props) => {
  const [enable, setEnable] = useState(1)

  const handleChangeEnable = (e: number) => {
    if (e === 0) {
      setTransformation({
        ...transformation,
        deploymentPattern: undefined,
        productionPattern: undefined
      })
    } else {
      setTransformation({
        ...transformation,
        deploymentPattern: '',
        productionPattern: ''
      })
    }
    setEnable(e)
  }

  return (
    <>
      <h3>CI/CD</h3>
      <p>
        <strong>What is a deployment?</strong>{' '}
        <Tag minimal intent={Intent.PRIMARY} style={{ fontSize: '10px' }}>
          DORA
        </Tag>
      </p>
      <p>Define Deployment using one of the following options.</p>
      <RadioGroup
        selectedValue={enable}
        onChange={(e) =>
          handleChangeEnable(+(e.target as HTMLInputElement).value)
        }
      >
        <Radio label='Detect Deployment from Builds in Jenkins' value={1} />
        {enable === 1 && (
          <div style={{ paddingLeft: 20 }}>
            <p>
              A Jenkins build with a name that matches the given regEx will be
              considered as a Deployment.
            </p>
            <FormGroup inline label='Deployment'>
              <InputGroup
                placeholder='(?i)deploy'
                value={transformation.deploymentPattern}
                onChange={(e) =>
                  setTransformation({
                    ...transformation,
                    deploymentPattern: e.target.value
                  })
                }
              />
            </FormGroup>
            <p>
              A Jenkins build with a name that matches the given regEx will be
              considered as a build in the Production environment. If you leave
              this field empty, all data will be tagged as in the Production
              environment.
            </p>
            <FormGroup inline label='Production'>
              <InputGroup
                placeholder='(?i)production'
                value={transformation.productionPattern}
                onChange={(e) =>
                  setTransformation({
                    ...transformation,
                    productionPattern: e.target.value
                  })
                }
              />
            </FormGroup>
          </div>
        )}
        <Radio label='Not using Jenkins Builds as Deployments' value={0} />
      </RadioGroup>
    </>
  )
}
