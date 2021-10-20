import React from 'react'
import 'normalize.css'
import '@blueprintjs/core/lib/css/blueprint.css'
import './styles/libraries/blueprint.scss'
import './styles/globals.scss'
import './styles/common.scss'

import {
  BrowserRouter as Router,
  Route,
} from 'react-router-dom'
import Configure from './pages/configure/index'
import Triggers from './pages/triggers/index'
import Jira from './pages/plugins/jira/index'
import Gitlab from './pages/plugins/gitlab/index'
import Jenkins from './pages/plugins/jenkins/index'

function App () {
  return (
    <Router>
      {/* Admin */}
      <Route exact path='/'>
        <Configure />
      </Route>
      <Route exact path='/triggers'>
        <Triggers />
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
