import React from 'react'
import {
  BrowserRouter as Router,
  Route,
} from 'react-router-dom'

import 'normalize.css'
import '@/styles/app.scss'
// Theme variables (@styles/theme.scss) injected via Webpack w/ @sass-loader additionalData option!
// import '@/styles/theme.scss'

import Configure from './pages/configure/index'
import Integration from '@/pages/configure/integration/index'
import ManageIntegration from '@/pages/configure/integration/manage'
import AddConnection from '@/pages/configure/connections/AddConnection'
import EditConnection from '@/pages/configure/connections/EditConnection'
import ConfigureConnection from '@/pages/configure/connections/ConfigureConnection'
import Triggers from '@/pages/triggers/index'
import Offline from '@/pages/offline/index'
import Pipelines from '@/pages/pipelines/index'
import CreatePipeline from '@/pages/pipelines/create'
import PipelineActivity from '@/pages/pipelines/activity'
import Blueprints from '@/pages/blueprints/index'

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
      <Route exact path='/pipelines/create'>
        <CreatePipeline />
      </Route>
      <Route exact path='/pipelines'>
        <Pipelines />
      </Route>
      <Route exact path='/pipelines/activity'>
        <PipelineActivity />
      </Route>
      <Route exact path='/pipelines/activity/:pId'>
        <PipelineActivity />
      </Route>
      <Route exact path='/blueprints'>
        <Blueprints />
      </Route>
      <Route exact path='/lake/api/configuration'>
        <Configure />
      </Route>
      <Route exact path='/offline'>
        <Offline />
      </Route>
    </Router>
  )
}

export default App
