import { useState, useEffect } from 'react'
import {
  useHistory
} from 'react-router-dom'
import { ToastNotification } from '@/components/Toast'
import { DEVLAKE_ENDPOINT } from '@/utils/config'
import request from '@/utils/request'

function useConnectionManager ({
  activeProvider,
  name, endpointUrl, token, username, password,
  isTesting, setIsTesting,
  isSaving, setIsSaving,
  testStatus, setTestStatus,
  errors, setErrors,
  showError, setShowError
}) {
  const history = useHistory()

  const testConnection = () => {
    setIsTesting(true)
    const connectionTestPayload = {
      name,
      endpointUrl,
      token,
      username,
      password
    }
    const testResponse = {
      success: false,
      connection: {
        ...connectionTestPayload
      },
      errors: []
    }
    console.log(testResponse)
    setTimeout(() => {
      if (testResponse.success) {
        setIsTesting(false)
        setTestStatus(1)
        ToastNotification.show({ message: 'Connection test OK.', intent: 'success', icon: 'small-tick' })
      } else {
        setIsTesting(false)
        setTestStatus(2)
        ToastNotification.show({ message: 'Connection test FAILED.', intent: 'danger', icon: 'error' })
      }
    }, 2000)
  }

  const saveConnection = () => {
    setIsSaving(true)
    let connectionPayload
    switch (activeProvider.id) {
      case 'jira':
        connectionPayload = { name: name, JIRA_ENDPOINT: endpointUrl, JIRA_BASIC_AUTH_ENCODED: token }
        break
      case 'jenkins':
        connectionPayload = { name: name, JENKINS_ENDPOINT: endpointUrl, JENKINS_USERNAME: username, JENKINS_PASSWORD: password }
        break
      case 'gitlab':
        connectionPayload = { name: name, GITLAB_ENDPOINT: endpointUrl, GITLAB_AUTH: token }
        break
    }

    let saveResponse = {
      success: false,
      connection: {
        ...connectionPayload
      },
      errors: []
    }

    const saveConfiguration = async (configPayload) => {
      try {
        const s = await request.post(`${DEVLAKE_ENDPOINT}/plugins/${activeProvider.id}/source`, configPayload)
        console.log('>> CONFIGURATION SAVED SUCCESSFULLY', configPayload, s)
        saveResponse = {
          ...saveResponse,
          success: s.data.success,
          connection: { ...s.data },
          errors: s.isAxiosError ? [s.message] : []
        }
      } catch (e) {
        saveResponse.errors.push(e.message)
        setErrors(saveResponse.errors)
        console.log('>> CONFIGURATION FAILED TO SAVE', configPayload, e)
      }
    }

    saveConfiguration(connectionPayload)

    setTimeout(() => {
      if (saveResponse.success && errors.length === 0) {
        ToastNotification.show({ message: 'Connection added successfully.', intent: 'success', icon: 'small-tick' })
        setShowError(false)
        setIsSaving(false)
        history.push(`/integrations/${activeProvider.id}`)
      } else {
        ToastNotification.show({ message: 'Connection failed to add, please try again.', intent: 'danger', icon: 'error' })
        setShowError(true)
        setIsSaving(false)
      }
    }, 2000)
  }

  return {
    testConnection,
    saveConnection,
  }
}

export default useConnectionManager
