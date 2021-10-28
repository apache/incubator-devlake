import { useState, useEffect } from 'react'
import {
  useHistory
} from 'react-router-dom'
import { ToastNotification } from '@/components/Toast'
import { DEVLAKE_ENDPOINT } from '@/utils/config'
import request from '@/utils/request'

function useSettingsManager ({
  activeProvider,
  activeConnection,
  settings,
}) {
  const history = useHistory()

  const [isSaving, setIsSaving] = useState(false)
  const [isTesting, setIsTesting] = useState(false)
  const [errors, setErrors] = useState([])
  const [showError, setShowError] = useState(false)

  const saveSettings = () => {
    setIsSaving(true)

    const settingsPayload = {
      ...settings,
      // DEV: true
    }

    let saveResponse = {
      success: false,
      settings: {
        ...settingsPayload
      },
      errors: []
    }

    const saveConfiguration = async (settingsPayload) => {
      try {
        setShowError(false)
        ToastNotification.clear()
        const s = await request.put(`${DEVLAKE_ENDPOINT}/plugins/${activeProvider.id}/sources/${activeConnection.id}`, settingsPayload)
        console.log('>> SETTINGS SAVED SUCCESSFULLY', settingsPayload, s)
        saveResponse = {
          ...saveResponse,
          success: s.data.success,
          settings: { ...s.data },
          errors: s.isAxiosError ? [s.message] : []
        }
      } catch (e) {
        saveResponse.errors.push(e.message)
        setErrors(saveResponse.errors)
        console.log('>> SETTINGS FAILED TO SAVE', settingsPayload, e)
      }
    }

    saveConfiguration(settingsPayload)

    setTimeout(() => {
      if (saveResponse.success && errors.length === 0) {
        ToastNotification.show({ message: 'Instance Settings saved successfully.', intent: 'success', icon: 'small-tick' })
        setShowError(false)
        setIsSaving(false)
      } else {
        ToastNotification.show({ message: 'Instance Settings failed to save, please try again.', intent: 'danger', icon: 'error' })
        setShowError(true)
        setIsSaving(false)
      }
    }, 2000)
  }

  const clear = () => {

  }

  const restore = () => {

  }

  return {
    saveSettings,
    clear,
    restore,
    isSaving,
    isTesting,
    errors,
    showError,
    setIsSaving,
    setIsTesting,
    setErrors,
    setShowError,
  }
}

export default useSettingsManager
