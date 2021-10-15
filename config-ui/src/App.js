import React from 'react'
// import 'normalize.css'
// import '@blueprintjs/core/lib/css/blueprint.css'
// import '../styles/libraries/blueprint.css'
// import '../styles/globals.css'
import {
  BrowserRouter as Router,
  Route,
} from 'react-router-dom'
// import Configure from './pages/configure/index'

function App () {
  return (
    <Router>
      <Route exact path='/'>
        {/* <Configure /> */}
        <p>configure</p>
      </Route>
      <Route path='/plugins/jira'>
        <p>jira</p>
      </Route>
    </Router>
  )
}

export default App
