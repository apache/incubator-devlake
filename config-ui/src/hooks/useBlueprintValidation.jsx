import React, { useState, useEffect, useCallback } from 'react'
// import { ToastNotification } from '@/components/Toast'
import {
  Providers,
} from '@/data/Providers'
import { parseCronExpression } from 'cron-schedule'
import cron from 'cron-validate'

function useBlueprintValidation ({
  name,
  cronConfig,
  customCronConfig,
  enable,
  tasks = [],
}) {
  const [errors, setErrors] = useState([])
  const [isValid, setIsValid] = useState(false)

  const clear = () => {
    setErrors([])
  }

  const validate = useCallback(() => {
    const errs = []
    // console.log('>> VALIDATING BLUEPRINT ', name)

    if (!name) {
      errs.push('Blueprint Name: Enter a valid Name')
    }

    if (name && name.length <= 2) {
      errs.push('Blueprint Name: Name too short, 3 chars minimum.')
    }

    if (!cronConfig) {
      errs.push('Blueprint Cron: No Crontab schedule defined.')
    }

    if (cronConfig && cronConfig !== 'custom' && !cron(cronConfig).isValid()) {
      errs.push('Blueprint Cron: Invalid Crontab Expression, unable to parse.')
    }

    if (cronConfig === 'custom' && !cron(customCronConfig).isValid()) {
      errs.push('Blueprint Cron: Invalid Custom Crontab Expression, unable to parse.')
    }

    if (enable && tasks?.length === 0) {
      errs.push('Blueprint Tasks: Invalid/Empty Configuration')
    }

    setErrors(errs)
  }, [
    name,
    cronConfig,
    customCronConfig,
    tasks,
    enable
  ])

  const fieldHasError = useCallback((fieldId) => {
    return errors.some(e => e.includes(fieldId))
  }, [errors])

  const getFieldError = useCallback((fieldId) => {
    return errors.find(e => e.includes(fieldId))
  }, [errors])

  useEffect(() => {
    // console.log('>>> BLUEPRINT FORM ERRORS...', errors)
    setIsValid(errors.length === 0)
    if (errors.length > 0) {
      // ToastNotification.clear()
    }
  }, [errors])

  return {
    errors,
    setErrors,
    isValid,
    validate,
    clear,
    fieldHasError,
    getFieldError
  }
}

export default useBlueprintValidation
