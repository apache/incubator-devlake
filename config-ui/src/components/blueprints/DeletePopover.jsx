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
import { Button, Intent, Colors, Classes } from '@blueprintjs/core'
const DeletePopover = (props) => {
  const {
    activeBlueprint,
    onCancel = () => {},
    onConfirm = () => {},
    isRunning = false
  } = props
  return (
    <>
      <div style={{ padding: '10px', fontSize: '10px', maxWidth: '220px' }}>
        <h3 style={{ margin: '0 0 5px 0', color: Colors.RED3 }}>
          Delete {activeBlueprint?.name}?
        </h3>
        <p>
          <strong>
            Are you sure? This Blueprint will be removed, all pipelines will be
            stopped.
          </strong>
        </p>
        <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
          <Button
            className={Classes.POPOVER_DISMISS}
            intent={Intent.NONE}
            text='CANCEL'
            small
            style={{ marginRight: '5px' }}
            onClick={() => onCancel(activeBlueprint)}
            disabled={isRunning}
          />
          <Button
            disabled={isRunning}
            intent={Intent.DANGER}
            text='YES'
            small
            onClick={() => onConfirm(activeBlueprint)}
          />
        </div>
      </div>
    </>
  )
}

export default DeletePopover
