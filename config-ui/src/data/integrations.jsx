import React from 'react'

import { Providers, ProviderLabels, ProviderTypes } from '@/data/Providers'

import JiraSettings from '@/pages/configure/settings/jira'
import GitlabSettings from '@/pages/configure/settings/gitlab'
import JenkinsSettings from '@/pages/configure/settings/jenkins'
import GithubSettings from '@/pages/configure/settings/github'

import { ReactComponent as GitlabProvider } from '@/images/integrations/gitlab.svg'
import { ReactComponent as JenkinsProvider } from '@/images/integrations/jenkins.svg'
import { ReactComponent as JiraProvider } from '@/images/integrations/jira.svg'
import { ReactComponent as GitHubProvider } from '@/images/integrations/github.svg'
import GitExtractorProvider from '@/images/git.png'
import RefDiffProvider from '@/images/git-diff.png'
// import { ReactComponent as NullProvider } from '@/images/integrations/null.svg'

const integrationsData = [
  {
    id: Providers.GITLAB,
    type: ProviderTypes.INTEGRATION,
    enabled: true,
    multiSource: false,
    name: ProviderLabels.GITLAB,
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
    type: ProviderTypes.INTEGRATION,
    enabled: true,
    multiSource: false,
    name: ProviderLabels.JENKINS,
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
    type: ProviderTypes.INTEGRATION,
    enabled: true,
    multiSource: true,
    name: ProviderLabels.JIRA,
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
    type: ProviderTypes.INTEGRATION,
    enabled: true,
    multiSource: false,
    name: ProviderLabels.GITHUB,
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

const pluginsData = [
  {
    id: Providers.GITEXTRACTOR,
    type: ProviderTypes.PIPELINE,
    enabled: true,
    multiSource: false,
    name: ProviderLabels.GITEXTRACTOR,
    icon: <img src={GitExtractorProvider} className='providerIconPng' width='30' height='30' style={{ float: 'left', marginTop: '5px' }} />,
    iconDashboard: <img src={GitExtractorProvider} className='providerIconPng' width='48' height='48' />,
    settings: ({ activeProvider, activeConnection, isSaving, setSettings }) => (
      null
    )
  },
  {
    id: Providers.REFDIFF,
    type: ProviderTypes.PIPELINE,
    enabled: true,
    multiSource: false,
    name: ProviderLabels.REFDIFF,
    icon: <img src={RefDiffProvider} className='providerIconPng' width='30' height='30' style={{ float: 'left', marginTop: '5px' }} />,
    iconDashboard: <img src={RefDiffProvider} className='providerIconPng' width='48' height='48' />,
    settings: ({ activeProvider, activeConnection, isSaving, setSettings }) => (
      null
    )
  },
]

export {
  integrationsData,
  pluginsData
}
