import React from 'react'
import 'normalize.css'
// import '@blueprintjs/core/lib/css/blueprint.css'
// Theme variables (@styles/theme.scss) injected via Webpack w/ @sass-loader additionalData option!
import '@blueprintjs/core/src/blueprint.scss'
import '@/styles/libraries/blueprint.scss'
import '@/styles/globals.scss'
import '@/styles/common.scss'

import {
  BrowserRouter as Router,
  Route,
} from 'react-router-dom'
import Configure from './pages/configure/index'
import Integration from '@/pages/configure/integration/index'
import ManageIntegration from '@/pages/configure/integration/manage'
import AddConnection from '@/pages/configure/connections/AddConnection'
import EditConnection from '@/pages/configure/connections/EditConnection'
import ConfigureConnection from '@/pages/configure/connections/ConfigureConnection'
import Triggers from '@/pages/triggers/index'
import Documentation from '@/pages/documentation/index'
import Jira from '@/pages/plugins/jira/index'
import Gitlab from '@/pages/plugins/gitlab/index'
import Jenkins from '@/pages/plugins/jenkins/index'

function App () {
  return (
    <Router>
      {/* Admin */}
      <Route exact path='/'>
        <Integration />
      </Route>
      <Route path='/integrations/:providerId'>
        <ManageIntegration />
      </Route>
      <Route path='/connections/add/:providerId'>
        <AddConnection />
      </Route>
      <Route path='/connections/edit/:providerId/:connectionId'>
        <EditConnection />
      </Route>
      <Route path='/connections/configure/:providerId/:connectionId'>
        <ConfigureConnection />
      </Route>
      <Route exact path='/integrations'>
        <Integration />
      </Route>
      <Route exact path='/triggers'>
        <Triggers />
      </Route>
      <Route exact path='/lake/api/configuration'>
        <Configure />
      </Route>
      <Route exact path='/documentation'>
        <Documentation />
      </Route>
      {/* Plugins */}
      <Route exact path='/plugins/jira'>
        <Jira />
      </Route>
      <Route exact path='/plugins/gitlab'>
        <Gitlab />
      </Route>
      <Route exact path='/plugins/jenkins'>
        <Jenkins />
      </Route>
    </Router>
  )
}

export default App
