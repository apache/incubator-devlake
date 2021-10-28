import { useState, useEffect } from 'react'
import {
  useHistory
} from 'react-router-dom'
import { ToastNotification } from '@/components/Toast'
import { DEVLAKE_ENDPOINT } from '@/utils/config'
import request from '@/utils/request'

function useConnectionManager ({
  activeProvider,
  activeConnection,
  connectionId,
  setActiveConnection,
  name = null,
  endpointUrl = null,
  token = null,
  username = null,
  password = null,
  // isTesting, setIsTesting,
  // isSaving, setIsSaving,
  // testStatus, setTestStatus,
  // errors, setErrors,
  // showError, setShowError
}, updateMode = false) {
  const history = useHistory()

  const [isSaving, setIsSaving] = useState(false)
  const [isFetching, setIsFetching] = useState(false)
  const [isRunning, setIsRunning] = useState(false)
  const [isTesting, setIsTesting] = useState(false)
  const [isDeleting, setIsDeleting] = useState(false)
  const [errors, setErrors] = useState([])
  const [showError, setShowError] = useState(false)
  const [testStatus, setTestStatus] = useState(0) //  0=Pending, 1=Success, 2=Failed

  const testConnection = () => {
    setIsTesting(true)
    setShowError(false)
    ToastNotification.clear()
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
        setShowError(false)
        ToastNotification.clear()
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

    const modifyConfiguration = async (configPayload) => {
      try {
        setShowError(false)
        ToastNotification.clear()
        const s = await request.put(`${DEVLAKE_ENDPOINT}/plugins/${activeProvider.id}/sources/${activeConnection.id}`, configPayload)
        console.log('>> CONFIGURATION MODIFIED SUCCESSFULLY', configPayload, s)
        saveResponse = {
          ...saveResponse,
          success: s.data.success,
          connection: { ...s.data },
          errors: s.isAxiosError ? [s.message] : []
        }
      } catch (e) {
        saveResponse.errors.push(e.message)
        setErrors(saveResponse.errors)
        console.log('>> CONFIGURATION FAILED TO UPDATE', configPayload, e)
      }
    }

    if (updateMode && activeConnection) {
      modifyConfiguration(connectionPayload)
    } else {
      saveConfiguration(connectionPayload)
    }

    setTimeout(() => {
      if (saveResponse.success && errors.length === 0) {
        ToastNotification.show({ message: 'Connection saved successfully.', intent: 'success', icon: 'small-tick' })
        setShowError(false)
        setIsSaving(false)
        if (!updateMode) {
          history.push(`/integrations/${activeProvider.id}`)
        }
      } else {
        ToastNotification.show({ message: 'Connection failed to save, please try again.', intent: 'danger', icon: 'error' })
        setShowError(true)
        setIsSaving(false)
      }
    }, 2000)
  }

  const runCollection = (options = {}) => {
    setIsRunning(true)
    ToastNotification.show({ message: 'Triggered Collection Process', intent: 'info', icon: 'info' })
    console.log('>> RUNNING COLLECTION PROCESS', isRunning)
    // Run Collection Tasks...
  }

  const fetchConnection = async () => {
    try {
      setIsFetching(true)
      console.log('>> FETCHING CONNECTION SOURCE', isFetching)
      const f = await request.get(`${DEVLAKE_ENDPOINT}/plugins/${activeProvider.id}/sources/${connectionId}`)
      const connectionData = f.data.data
      setActiveConnection(connectionData)
      setIsFetching(false)
    } catch (e) {
      setIsFetching(false)
      setActiveConnection({
        id: null,
        name: null,
        endpoint: null,
        token: null,
        username: null,
        password: null,
      })
      ToastNotification.show({ message: `${e}`, intent: 'danger', icon: 'error' })
      console.log('>> FAILED TO FETCH CONNECTION', e)
    }
  }

  const deleteConnection = () => {
    // @todo Implement DELETE
    try {
      setIsDeleting(true)
      console.log('>> TRYING TO DELETE CONNECTION...', isDeleting)
      // const d = await request.delete(`${DEVLAKE_ENDPOINT}/plugins/${activeProvider.id}/sources/${connectionId}`)
      // setIsDeleting(false)
    } catch (e) {
      setIsDeleting(false)
      console.log('>> FAILED TO DELETE CONNECTION', e)
    }
  }

  useEffect(() => {
    if (activeConnection && activeConnection.id !== null) {
      ToastNotification.clear()
      ToastNotification.show({ message: `Fetched settings for ${activeConnection.name}.`, intent: 'success', icon: 'small-tick' })
      console.log('>> FETCHED CONNECTION FOR MODIFY', activeConnection)
    }
  }, [activeConnection])

  return {
    fetchConnection,
    testConnection,
    saveConnection,
    deleteConnection,
    runCollection,
    isSaving,
    isTesting,
    isFetching,
    errors,
    showError,
    testStatus,
    setIsSaving,
    setIsTesting,
    setIsFetching,
    setErrors,
    setShowError,
    setTestStatus
  }
}

export default useConnectionManager
