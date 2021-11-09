import React from 'react'

import { Providers } from '@/data/Providers'

import JiraSettings from '@/pages/configure/settings/jira'
import GitlabSettings from '@/pages/configure/settings/gitlab'
import JenkinsSettings from '@/pages/configure/settings/jenkins'
import GithubSettings from '@/pages/configure/settings/github'

import { ReactComponent as GitlabProvider } from '@/images/integrations/gitlab.svg'
import { ReactComponent as JenkinsProvider } from '@/images/integrations/jenkins.svg'
import { ReactComponent as JiraProvider } from '@/images/integrations/jira.svg'
import { ReactComponent as GitHubProvider } from '@/images/integrations/github.svg'

const integrationsData = [
  {
    id: Providers.GITLAB,
    enabled: true,
    name: 'GitLab',
    icon: <GitlabProvider className='providerIconSvg' width='30' height='30' style={{ float: 'left', marginTop: '5px' }} />,
    iconDashboard: <GitlabProvider className='providerIconSvg' width='48' height='48' />,
    settings: ({ activeProvider, activeConnection, isSaving, setSettings }) => (
      <GitlabSettings
        provider={activeProvider}
        connection={activeConnection}
        isSaving={isSaving}
        onSettingsChange={setSettings}
      />
    )
  },
  {
    id: Providers.JENKINS,
    enabled: true,
    name: 'Jenkins',
    icon: <JenkinsProvider className='providerIconSvg' width='30' height='30' style={{ float: 'left', marginTop: '5px' }} />,
    iconDashboard: <JenkinsProvider className='providerIconSvg' width='48' height='48' />,
    settings: ({ activeProvider, activeConnection, isSaving, setSettings }) => (
      <JenkinsSettings
        provider={activeProvider}
        connection={activeConnection}
        isSaving={isSaving}
        onSettingsChange={setSettings}
      />
    )
  },
  {
    id: Providers.JIRA,
    enabled: true,
    name: 'JIRA',
    icon: <JiraProvider className='providerIconSvg' width='30' height='30' style={{ float: 'left', marginTop: '5px' }} />,
    iconDashboard: <JiraProvider className='providerIconSvg' width='48' height='48' />,
    settings: ({ activeProvider, activeConnection, isSaving, setSettings }) => (
      <JiraSettings
        provider={activeProvider}
        connection={activeConnection}
        isSaving={isSaving}
        onSettingsChange={setSettings}
      />
    )
  },
  {
    id: Providers.GITHUB,
    enabled: true,
    name: 'GitHub',
    icon: <GitHubProvider className='providerIconSvg' width='30' height='30' style={{ float: 'left', marginTop: '5px' }} />,
    iconDashboard: <GitHubProvider className='providerIconSvg' width='48' height='48' />,
    settings: ({ activeProvider, activeConnection, isSaving, setSettings }) => (
      <GithubSettings
        provider={activeProvider}
        connection={activeConnection}
        isSaving={isSaving}
        onSettingsChange={setSettings}
      />
    )
  },
]

export {
  integrationsData
}
