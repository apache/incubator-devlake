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
import { useState, useEffect, useCallback } from 'react'
import request from '@/utils/request'
import { MigrationOptions } from '@/config/migration'
import { Intent } from '@blueprintjs/core'
import { ToastNotification } from '@/components/Toast'

function useDatabaseMigrations(Configuration = MigrationOptions) {
  const [isProcessing, setIsProcessing] = useState(false)

  const [migrationWarning, setMigrationWarning] = useState()
  const [migrationAlertOpened, setMigrationAlertOpened] = useState(false)
  const [wasMigrationSuccessful, setWasMigrationSuccessful] = useState(false)
  const [hasMigrationFailed, setHasMigrationFailed] = useState(false)

  const handleConfirmMigration = useCallback(() => {
    setIsProcessing(true)
    const migrate = async () => {
      const m = await request.get(Configuration.apiProceedEndpoint)
      setWasMigrationSuccessful(m?.status === 200 && m?.data?.success === true)
      setIsProcessing(false)
      setHasMigrationFailed(m?.status !== 200)
    }
    migrate()
  }, [Configuration.apiProceedEndpoint])

  const handleCancelMigration = useCallback(() => {
    setIsProcessing(true)
    setMigrationAlertOpened(false)
    setIsProcessing(false)
    ToastNotification.clear()
    ToastNotification.show({
      // eslint-disable-next-line max-len
      message: Configuration.cancelToastMessage,
      intent: Intent.NONE,
      icon: 'warning-sign'
    })
  }, [Configuration.cancelToastMessage])

  const handleMigrationDialogClose = useCallback(() => {
    setMigrationAlertOpened(false)
  }, [setMigrationAlertOpened])

  useEffect(() => {
    setMigrationAlertOpened(!!migrationWarning)
  }, [migrationWarning, setMigrationAlertOpened])

  useEffect(() => {
    if (hasMigrationFailed) {
      ToastNotification.clear()
      ToastNotification.show({
        // eslint-disable-next-line max-len
        message: MigrationOptions.failedToastMessage,
        intent: Intent.DANGER,
        icon: 'error'
      })
    }
  }, [hasMigrationFailed])

  const setNeedMigrateCallback = useCallback((errorObject) => {
    if (migrationWarning) {
      console.log(`>>> MIGRATION WARNING DETECTED !!:`, migrationWarning)
    }
    setMigrationWarning(errorObject)
    setWasMigrationSuccessful(false)
    setIsProcessing(false)
    setHasMigrationFailed(false)
  }, [setMigrationWarning])
  useEffect(() => {
    MigrationOptions.setNeedMigrateCallback(setNeedMigrateCallback)
  }, [setNeedMigrateCallback])


  return {
    migrationWarning,
    migrationAlertOpened,
    wasMigrationSuccessful,
    hasMigrationFailed,
    isProcessing,
    setIsProcessing,
    handleConfirmMigration,
    handleCancelMigration,
    handleMigrationDialogClose
  }
}

export default useDatabaseMigrations
