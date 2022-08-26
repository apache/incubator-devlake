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
import React, { useEffect, useState } from 'react'
import {
  Button,
  Classes,
  Colors,
  Dialog,
  Intent,
} from '@blueprintjs/core'

const MigrationAlertDialog = (props) => {
  const {
    isOpen = false,
    icon = 'outdated',
    title = 'New Migration Scripts Detected',
    onClose = () => {},
    onClosed = () => {},
    onConfirm = () => {},
    canEscapeKeyClose = false,
    canOutsideClickClose = false,
    isCloseButtonShown = false
  } = props

  return (
    <>
      <Dialog
        className='dialog-db-migration'
        icon={icon}
        title={title}
        isOpen={isOpen}
        onClose={onClose}
        onClosed={onClosed}
        canEscapeKeyClose={canEscapeKeyClose}
        canOutsideClickClose={canOutsideClickClose}
        isCloseButtonShown={isCloseButtonShown}
      >
        <div className={Classes.DIALOG_BODY}>
          <p style={{ margin: 0, padding: 0, color: Colors.RED4 }}>
            WARNING: Performing migration may wipe collected data for consistency and re-collecting data may be required.
          </p>
          <p style={{ margin: 0, padding: 0 }}>
            A Database migration is required to launch <strong>DevLake</strong>, to proceed, please send a request to <code style={{ backgroundColor: '#eeeeee' }}>&lt;config-ui-endpoint&gt;/api/proceed-db-migration</code>{' '}
            ( or <code style={{ backgroundColor: '#eeeeee' }}>&lt;devlake-endpoint&gt;/proceed-db-migration</code>){' '}
            Alternatively, you may downgrade back to the previous DevLake version.
          </p>
        </div>
        <div className={Classes.DIALOG_FOOTER}>
          <div className={Classes.DIALOG_FOOTER_ACTIONS}>
            <Button text='Downgrade' intent={Intent.PRIMARY} outlined onClick={onClose} />
            <Button text='Proceed to Database Migration' intent={Intent.PRIMARY} onClick={onConfirm} />
          </div>
        </div>
      </Dialog>
    </>
  )
}

export default MigrationAlertDialog
