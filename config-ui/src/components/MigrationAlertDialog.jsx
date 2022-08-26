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
  Elevation,
} from '@blueprintjs/core'
import ContentLoader from '@/components/loaders/ContentLoader'

const MigrationAlertDialog = (props) => {
  const {
    isOpen = false,
    icon = 'outdated',
    title = 'New Migration Scripts Detected',
    onClose = () => {},
    onClosed = () => {},
    onCancel = () => {},
    onConfirm = () => {},
    canEscapeKeyClose = false,
    canOutsideClickClose = false,
    isCloseButtonShown = false,
    isMigrating = false,
    wasSuccessful = false,
    hasFailed = false,
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
          {!isMigrating && wasSuccessful ? (
            <>
              <p style={{ margin: 0, padding: 0, color: Colors.GREEN4 }}>
                <strong>Database Migration Successful!</strong>
              </p>
              <p style={{ margin: 0, padding: 0 }}>
                There is no further action, You may continue using DevLake by
                clicking below or Reloading your browser.
              </p>
            </>
          ) : (
            <>
              {isMigrating ? (
                <>
                  <ContentLoader
                    title='Running Migrations...'
                    elevation={Elevation.ZERO}
                    cardStyleOverrides={{ backgroundColor: 'transparent', marginBottom: 0, fontSize: '12px' }}
                    messageClasses={['bp3-ui-text']}
                    message={
                      <>
                        Please wait for database migrations to complete, do{' '}
                        <strong>NOT</strong> close your browser at this time.
                      </>
                    }
                  />
                </>
              ) : (
                <>
                  <p style={{ margin: 0, padding: 0, color: Colors.RED4 }}>
                    WARNING: Performing migration may wipe collected data for
                    consistency and re-collecting data may be required.
                  </p>
                  <p style={{ margin: 0, padding: 0 }}>
                    A Database migration is required to launch{' '}
                    <strong>DevLake</strong>, to proceed, please send a request
                    to{' '}
                    <code style={{ backgroundColor: '#eeeeee' }}>
                      &lt;config-ui-endpoint&gt;/api/proceed-db-migration
                    </code>{' '}
                    ( or{' '}
                    <code style={{ backgroundColor: '#eeeeee' }}>
                      &lt;devlake-endpoint&gt;/proceed-db-migration
                    </code>
                    ) Alternatively, you may downgrade back to the previous
                    DevLake version.
                  </p>
                </>
              )}
            </>
          )}
        </div>
        <div className={Classes.DIALOG_FOOTER}>
          <div className={Classes.DIALOG_FOOTER_ACTIONS}>
            {wasSuccessful ? (
              <>
                <Button
                  icon='small-tick'
                  text='Continue'
                  intent={Intent.SUCCESS}
                  onClick={onClose}
                />
              </>
            ) : (
              <>
                <Button
                  disabled={isMigrating}
                  text='Downgrade'
                  intent={Intent.PRIMARY}
                  onClick={onCancel}
                  outlined
                />
                <Button
                  disabled={isMigrating}
                  loading={isMigrating}
                  text='Proceed to Database Migration'
                  intent={Intent.PRIMARY}
                  onClick={onConfirm}
                />
              </>
            )}
          </div>
        </div>
      </Dialog>
    </>
  )
}

export default MigrationAlertDialog
