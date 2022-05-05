import React from 'react'
import { Providers, ProviderLabels } from '@/data/Providers'
import { ReactComponent as NullProviderIcon } from '@/images/integrations/null.svg'

const NullProvider = {
  id: Providers.NULL, // Unique ID, for a Provider (alphanumeric, lowercase)
  enabled: false, // Enabled Flag
  multiConnection: false, // If Provider is Multi-connection
  name: ProviderLabels.NULL, // Display Name of Data Provider
  // eslint-disable-next-line max-len
  icon: <NullProviderIcon className='providerIconSvg' width='30' height='30' style={{ float: 'left', marginTop: '5px' }} />, // Provider Icon
  iconDashboard: <NullProviderIcon className='providerIconSvg' width='48' height='48' />, // Provider Icon on INTEGRATIONS Dashboard
  settings: ({ activeProvider, activeConnection, isSaving, setSettings }) => (<></>) // REACT Settings Component for Render
}

export {
  NullProvider
}
