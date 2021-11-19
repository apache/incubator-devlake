import React, { useState, useEffect, useCallback } from 'react'
import { ToastNotification } from '@/components/Toast'
import {
  Providers,
} from '@/data/Providers'

function useConnectionValidation ({
  activeProvider,
  name,
  endpointUrl,
  token,
  username,
  password
}) {
  const [errors, setErrors] = useState([])
  const [isValid, setIsValid] = useState(false)

  const clear = () => {
    setErrors([])
  }

  const validate = useCallback(() => {
    const errs = []
    console.log('>> VALIDATING PROVIDER ID: ', activeProvider.id)
    console.log('>> RUNNING FORM VALIDATIONS AGAINST FIELD VALUES...')
    console.log(
      'NAME', name,
      'ENDPOINT URL', endpointUrl,
      'TOKEN', token,
      'USERNAME', username,
      'PASSWORD', password
    )

    if (!name || name.length <= 2) {
      errs.push('Connection Source name is required')
    }

    if (!endpointUrl || endpointUrl.length <= 2) {
      errs.push('Endpoint URL is required')
    }

    if (!endpointUrl?.startsWith('http')) {
      errs.push('Endpoint URL must be valid HTTP/S protocol')
    }

    if (!endpointUrl?.endsWith('/')) {
      errs.push('Endpoint URL must end in trailing slash (/)')
    }

    switch (activeProvider.id) {
      case Providers.GITHUB:
      case Providers.JIRA:
      case Providers.GITLAB:
        if (!token || token.length <= 2) {
          errs.push('Authentication token(s) are required')
        }
        break
      case Providers.JENKINS:
        if (!username || username.length <= 2) {
          errs.push('Username is required')
        }
        if (!password || password.length <= 2) {
          errs.push('Password is required')
        }
        break
    }

    setErrors(errs)
  }, [name, endpointUrl, token, username, password, activeProvider])

  useEffect(() => {
    console.log('>>> CONNECTION FORM ERRORS...', errors)
    setIsValid(errors.length === 0)
    if (errors.length > 0) {
      // ToastNotification.clear()
      // ToastNotification.show({ message: errors[0], intent: 'danger', icon: 'warning-sign' })
    }
  }, [errors])

  return {
    errors,
    isValid,
    validate,
    clear
  }
}

export default useConnectionValidation
