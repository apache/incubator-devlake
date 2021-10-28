import React from 'react'
import { ReactComponent as GitlabProvider } from '@/images/integrations/gitlab.svg'
import { ReactComponent as JenkinsProvider } from '@/images/integrations/jenkins.svg'
import { ReactComponent as JiraProvider } from '@/images/integrations/jira.svg'

const integrationsData = [
  {
    id: 'gitlab',
    name: 'GitLab',
    icon: <GitlabProvider className='providerIconSvg' width='30' height='30' style={{ float: 'left', marginTop: '5px' }} />
  },
  {
    id: 'jenkins',
    name: 'Jenkins',
    icon: <JenkinsProvider className='providerIconSvg' width='30' height='30' style={{ float: 'left', marginTop: '5px' }} />
  },
  {
    id: 'jira',
    name: 'JIRA',
    icon: <JiraProvider className='providerIconSvg' width='30' height='30' style={{ float: 'left', marginTop: '5px' }} />
  },
]

export {
  integrationsData
}
