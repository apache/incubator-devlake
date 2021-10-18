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

function App () {
  return (
    <Router>
      <Route exact path='/'>
        <Configure />
      </Route>
      <Route exact path='/triggers'>
        <Triggers />
      </Route>
      <Route path='/plugins/jira'>
        <p>jira</p>
      </Route>
    </Router>
  )
}

export default App
