
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
          label: 'JIRA',
          route: '/integrations/jira',
          active: activeRoute.url.endsWith('/integrations/jira'),
          icon: 'layers',
          classNames: [],
        },
        {
          id: 1,
          label: 'GitLab',
          route: '/integrations/gitlab',
          active: activeRoute.url.endsWith('/integrations/gitlab'),
          icon: 'layers',
          classNames: [],
        },
        {
          id: 2,
          label: 'Jenkins',
          route: '/integrations/jenkins',
          active: activeRoute.url.endsWith('/integrations/jenkins'),
          icon: 'layers',
          classNames: [],
        }
      ]
    },
    {
      id: 1,
      label: 'Jobs & Tasks',
      icon: 'automatic-updates',
      route: '/tasks',
      disabled: true,
      active: activeRoute.url === '/tasks',
      children: [
      ]
    },
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
