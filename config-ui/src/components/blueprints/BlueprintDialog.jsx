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
import React, { useEffect, useState, useRef, useCallback } from 'react'
import {
  Button,
  Classes,
  Colors,
  Dialog,
  Elevation,
  FormGroup,
  Icon,
  Intent,
  Label,
  MenuItem,
} from '@blueprintjs/core'

import { NullBlueprint } from '@/data/NullBlueprint'

const Modes = {
  CREATE: 'create',
  EDIT: 'edit',
}

const BlueprintDialog = (props) => {
  const {
    isOpen = false,
    title = 'Manage Blueprint',
    blueprint = NullBlueprint,
    mode = Modes.EDIT,
    canOutsideClickClose = false,
    onClose = () => {},
    onCancel = () => {},
    onSave = () => {},
    isSaving = false,
    isValid = true,
    isTesting = false,
    content = null
  } = props

  useEffect(() => {

  }, [content])

  return (
    <>
      <Dialog
        className='dialog-manage-blueprint'
        // icon={mode === Modes.EDIT ? 'edit' : 'add'}
        title={title}
        isOpen={isOpen}
        onClose={onClose}
        onClosed={() => {}}
        canOutsideClickClose={canOutsideClickClose}
        style={{ backgroundColor: '#ffffff' }}
      >
        <div className={Classes.DIALOG_BODY}>
          {content}
        </div>
        <div className={Classes.DIALOG_FOOTER}>
          <div className={Classes.DIALOG_FOOTER_ACTIONS}>
            <Button
              className='btn-cancel'
              disabled={isSaving}
              intent={Intent.NONE}
              onClick={onCancel}
              loading={isSaving}
              outlined
            >
              Cancel
            </Button>
            <Button
              className='btn-save'
              disabled={isSaving || !isValid || isTesting}
              intent={Intent.PRIMARY}
              onClick={onSave}
              loading={isSaving}
              // outlined
            >
              Save Changes
            </Button>
          </div>
        </div>
      </Dialog>
    </>
  )
}

export default BlueprintDialog
