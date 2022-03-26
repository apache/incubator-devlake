import React from 'react'

import { Providers, ProviderLabels } from '@/data/Providers'

import JiraSettings from '@/pages/configure/settings/jira'
import GitlabSettings from '@/pages/configure/settings/gitlab'
import JenkinsSettings from '@/pages/configure/settings/jenkins'
import GithubSettings from '@/pages/configure/settings/github'

import { ReactComponent as GitlabProvider } from '@/images/integrations/gitlab.svg'
import { ReactComponent as JenkinsProvider } from '@/images/integrations/jenkins.svg'
import { ReactComponent as JiraProvider } from '@/images/integrations/jira.svg'
import { ReactComponent as GitHubProvider } from '@/images/integrations/github.svg'
// import { ReactComponent as NullProvider } from '@/images/integrations/null.svg'

const integrationsData = [
  {
    id: Providers.GITLAB,
    enabled: true,
    multiSource: false,
    name: ProviderLabels.GITLAB,
    icon: <GitlabProvider className='providerIconSvg' width='30' height='30' style={{ float: 'left', marginTop: '5px' }} />,
    iconDashboard: <GitlabProvider className='providerIconSvg' width='48' height='48' />,
    settings: ({ activeProvider, activeConnection, isSaving, isSavingConnection, setSettings }) => (
      <GitlabSettings
        provider={activeProvider}
        connection={activeConnection}
        isSaving={isSaving}
        isSavingConnection={isSavingConnection}
        onSettingsChange={setSettings}
      />
    )
  },
  {
    id: Providers.JENKINS,
    enabled: true,
    multiSource: false,
    name: ProviderLabels.JENKINS,
    icon: <JenkinsProvider className='providerIconSvg' width='30' height='30' style={{ float: 'left', marginTop: '5px' }} />,
    iconDashboard: <JenkinsProvider className='providerIconSvg' width='48' height='48' />,
    settings: ({ activeProvider, activeConnection, isSaving, isSavingConnection, setSettings }) => (
      <JenkinsSettings
        provider={activeProvider}
        connection={activeConnection}
        isSaving={isSaving}
        isSavingConnection={isSavingConnection}
        onSettingsChange={setSettings}
      />
    )
  },
  {
    id: Providers.JIRA,
    enabled: true,
    multiSource: true,
    name: ProviderLabels.JIRA,
    icon: <JiraProvider className='providerIconSvg' width='30' height='30' style={{ float: 'left', marginTop: '5px' }} />,
    iconDashboard: <JiraProvider className='providerIconSvg' width='48' height='48' />,
    settings: ({ activeProvider, activeConnection, isSaving, isSavingConnection, setSettings }) => (
      <JiraSettings
        provider={activeProvider}
        connection={activeConnection}
        isSaving={isSaving}
        isSavingConnection={isSavingConnection}
        onSettingsChange={setSettings}
      />
    )
  },
  {
    id: Providers.GITHUB,
    enabled: true,
    multiSource: false,
    name: ProviderLabels.GITHUB,
    icon: <GitHubProvider className='providerIconSvg' width='30' height='30' style={{ float: 'left', marginTop: '5px' }} />,
    iconDashboard: <GitHubProvider className='providerIconSvg' width='48' height='48' />,
    settings: ({ activeProvider, activeConnection, isSaving, isSavingConnection, setSettings }) => (
      <GithubSettings
        provider={activeProvider}
        connection={activeConnection}
        isSaving={isSaving}
        isSavingConnection={isSavingConnection}
        onSettingsChange={setSettings}
      />
    )
  },
]

export {
  integrationsData
}
