const endpoints = {
  projects: [
    {
      requiredFields: {
        id: true
      },
      path: 'access_requests'
    }
  ],
  groups: [
    {
      requiredFields: {
        id: true
      },
      path: 'epics',
      name: 'Epics'
    }
  ],
}