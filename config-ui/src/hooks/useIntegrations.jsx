import React, { useState, useEffect } from 'react'
import { integrationsData } from '@/data/integrations'

function useIntegrations (data = []) {
  const [integrations, setIntegrations] = useState(integrationsData)

  useEffect(() => {
    // setIntegrations(integrationsData)
  }, [])

  return [
    integrations,
    setIntegrations
  ]
}

export default useIntegrations
