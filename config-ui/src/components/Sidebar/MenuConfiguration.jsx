
import { ProviderLabels } from '@/data/Providers'

const MenuConfiguration = (activeRoute) => {
  return [
    {
      id: 0,
      label: 'Data Integrations',
      route: '/integrations',
      active: activeRoute.url.startsWith('/integrations') || activeRoute.url === '/',
      icon: 'data-connection',
      classNames: [],
      children: [
        {
          id: 0,
          label: ProviderLabels.JIRA,
          route: '/integrations/jira',
          active: activeRoute.url.endsWith('/integrations/jira') || activeRoute.url.endsWith('/jira'),
          icon: 'layers',
          classNames: [],
        },
        {
          id: 1,
          label: ProviderLabels.GITHUB,
          route: '/integrations/github',
          active: activeRoute.url.endsWith('/integrations/github') || activeRoute.url.endsWith('/github'),
          icon: 'layers',
          classNames: [],
        },
        {
          id: 2,
          label: ProviderLabels.GITLAB,
          route: '/integrations/gitlab',
          active: activeRoute.url.endsWith('/integrations/gitlab') || activeRoute.url.endsWith('/gitlab'),
          icon: 'layers',
          classNames: [],
        },
        {
          id: 3,
          label: ProviderLabels.JENKINS,
          route: '/integrations/jenkins',
          active: activeRoute.url.endsWith('/integrations/jenkins') || activeRoute.url.endsWith('/jenkins'),
          icon: 'layers',
          classNames: [],
        }
      ]
    },
    // {
    //   id: 1,
    //   label: 'Tasks',
    //   icon: 'automatic-updates',
    //   route: '/tasks',
    //   disabled: true,
    //   active: activeRoute.url === '/tasks',
    //   children: [
    //   ]
    // },
    {
      id: 2,
      label: 'Triggers',
      icon: 'asterisk',
      classNames: [],
      route: '/triggers',
      active: activeRoute.url === '/triggers',
      children: [
      ]
    },
    {
      id: 3,
      label: 'Pipelines',
      icon: 'git-merge',
      classNames: [],
      route: '/pipelines',
      active: activeRoute.url.startsWith('/pipelines'),
      children: [
        {
          id: 0,
          label: 'Create Pipeline Run',
          route: '/pipelines/create',
          active: activeRoute.url.endsWith('/pipelines/create'),
          icon: 'git-pull',
          classNames: [],
        },
        {
          id: 1,
          label: 'All Pipeline Runs',
          route: '/pipelines',
          active: activeRoute.url.endsWith('/pipelines'),
          icon: 'layers',
          classNames: [],
          disabled: false
        },
        {
          id: 2,
          label: 'Pipeline Blueprints',
          route: '/blueprints',
          active: activeRoute.url.endsWith('/blueprints'),
          icon: 'grid',
          classNames: [],
        },
      ]
    },
    // {
    //   id: 3,
    //   label: 'Documentation',
    //   icon: 'help',
    //   classNames: [],
    //   route: 'https://github.com/merico-dev/lake/wiki',
    //   target: "_blank",
    //   external: true,
    //   active: activeRoute.url === '/documentation',
    //   children: [
    //   ]
    // },
  ]
}

export default MenuConfiguration
