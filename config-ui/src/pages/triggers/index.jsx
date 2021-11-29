import React, { useState, useEffect, useCallback } from 'react'
import {
  AnchorButton,
  Spinner,
  Button,
  TextArea,
  Card,
  Elevation
} from '@blueprintjs/core'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import AppCrumbs from '@/components/Breadcrumbs'
import Content from '@/components/Content'
import request from '@/utils/request'
import { DEVLAKE_ENDPOINT, GRAFANA_ENDPOINT } from '@/utils/config.js'
import TriggersUtil  from '@/utils/triggersUtil'
import SourcesUtil from '@/utils/sourcesUtil'

const STAGE_INIT = 0
const STAGE_PENDING = 1
const STAGE_FINISHED = 2

export default function Triggers () {
  const [triggerJson, setTriggerJson] = useState([[]])
  const [pipeline, setPipeline] = useState(null)
  const [tasks, setTasks] = useState(null)

  // update stage based on pipeline existence and its status
  const stage = useCallback(() => {
    if (!pipeline) {
      return STAGE_INIT;
    }
    if (pipeline.status !== 'TASK_COMPLETED') {
      return STAGE_PENDING;
    }
    return STAGE_FINISHED;
  }, [pipeline])()

  const triggerDisabled = useCallback(() => {
    return pipeline ? 'true' : 'false'
  }, [pipeline])()

  // try to reload pipeline/tasks from server every 3s after triggered
  useEffect(() => {
    const interval = setInterval(async () => {
      if (!pipeline || pipeline.status === 'TASK_COMPLETED') {
        return
      }
      try {
        const [pipelineRes, tasksRes] = await Promise.all([
          request.get(`${DEVLAKE_ENDPOINT}/pipelines/${pipeline.ID}`),
          request.get(`${DEVLAKE_ENDPOINT}/pipelines/${pipeline.ID}/tasks`),
        ])
        setPipeline(pipelineRes.data)

        // convert to 2d array
        const newTasks = []
        for (const newTask of tasksRes.data.tasks) {
          if (!newTasks[newTask.pipelineRow-1]) {
            newTasks[newTask.pipelineRow-1] = []
          }
          newTasks[newTask.pipelineRow-1][newTask.pipelineCol-1] = newTask
        }
        console.log(newTasks)
        setTasks(newTasks)
      } finally { }
    }, 3000)
    return () => clearInterval(interval)
  }, [pipeline])

  useEffect(() => {
    console.log('Setting JSON based on active plugins...');
    const setTriggerJsonBasedOnActivePlugins = async () => {
      let pluginsToSet = await SourcesUtil.getPluginSources()
      let collectionJson = TriggersUtil.getCollectionJson(pluginsToSet)
      setTriggerJson(collectionJson)
    }
    setTriggerJsonBasedOnActivePlugins()
  }, [])

  useEffect(() => {
    console.log('Setting text area based on updated triggers JSON...');
    setTextAreaBody(JSON.stringify(triggerJson, null, 2))
  }, [triggerJson])

  // user clicked on trigger button
  const sendTrigger = async (e) => {
    e.preventDefault()
    // @todo RE_ACTIVATE Trigger Process!
    try {
      const res = await request.post(
        `${DEVLAKE_ENDPOINT}/pipelines`,
        JSON.stringify({
          name: `config-ui trigger ${new Date()}`,
          tasks: JSON.parse(textAreaBody)
        })
      )
      setPipeline(res.data)
      console.log('waiting following pipeline to complete: ', pipeline)
    } catch (e) {
      console.error(e)
    }
  }

  const [textAreaBody, setTextAreaBody] = useState(JSON.stringify(triggerJson, null, 2))

  return (
    <div className='container'>
      <Nav />
      <Sidebar />
      <Content>
        <main className='main'>
          <AppCrumbs
            items={[
              { href: '/', icon: false, text: 'Dashboard' },
              { href: '/triggers', icon: false, text: 'Data Triggers' },
            ]}
          />
          {stage === STAGE_FINISHED &&
            <div className='headlineContainer'>
              <h1>Done</h1>
              <p className='description'>Navigate to Grafana to view updated metrics</p>
              <AnchorButton
                href={GRAFANA_ENDPOINT}
                icon='grouped-bar-chart'
                target='_blank'
                text='View Dashboards'
              />
            </div>}
          {stage === STAGE_PENDING &&
            <div className='headlineContainer'>
              <h1>Collecting Data</h1>
              <p className='description'>Please wait... </p>
              {tasks && tasks.map((step, index) =>
                <div key={index}>
                <h2> Step {index+1}</h2>
                {step.map(task =>
                  <div className='pluginSpinnerWrap' key={`key-${task.ID}`}>
                    <div key={`progress-${task.ID}`}>
                      <span style={{ display: 'inline-block', width: '200px' }}>{task.plugin}</span>
                      {task.status === 'TASK_RUNNING' &&
                        <>
                          <Spinner
                            size={12}
                            className='pluginSpinner'
                          />
                          <strong>{task.progress * 100}%</strong>
                        </>}
                      {task.status === 'TASK_COMPLETED' &&
                        <>
                          <span style={{ display: 'inline-block', color: 'green', fontWeight: 'bold', width: '80px' }}>Succeeded</span>
                        </>}
                      {task.status === 'TASK_FAILED' &&
                        <>
                          <span style={{ display: 'inline-block', color: 'red', fontWeight: 'bold', width: '80px' }}>Failed</span>
                          <span style={{ display: 'inline-block', color: 'red' }}>{task.message}</span>
                        </>}
                    </div>
                  </div>
                )}
                </div>
              )}
            </div>
          }
          {stage === STAGE_INIT && (
            <>
              <div className='headlineContainer'>
                <h1>Triggers</h1>
                <p className='description'>Trigger data collection on your plugins</p>
              </div>

              <form className='form'>
                <div className='headlineContainer'>
                  <p className='description'>Create a <strong>http</strong> request to trigger data collect tasks,&nbsp;
                    {/* eslint-disable-next-line max-len */}
                    please customize the following JSON by removing the plugins you don't need and replace with your own&nbsp;
                    {/* eslint-disable-next-line max-len */}
                    <strong>JIRA</strong> <code>boardId</code> / <strong>GitLab</strong> <code>projectId</code> / <strong>GitHub</strong> <code>repositoryName</code> and <code>owner</code> in the request body. &nbsp;
                    {/* eslint-disable-next-line max-len */}
                    For a project with 10k commits and 5k JIRA issues, this can take up to <em>20 minutes</em> for collecting JIRA, GitLab, and Jenkins data.&nbsp;
                    {/* eslint-disable-next-line max-len */}
                    The data collection will take longer for GitHub since they have a rate limit of 2k requests per hour. You can accelerate the process by configuring multiple personal access tokens.
                  </p>
                  <p className='description' style={{ fontSize: '13px' }}>
                    <span style={{ fontWeight: 'bold' }}>Detailed configuration guide:</span>&nbsp;
                    <a
                      href='https://github.com/merico-dev/lake/wiki/How-to-use-the-triggers-page' target='_blank'
                      rel='noreferrer' style={{ fontWeight: 'bold', color: '#E8471C', textDecoration: 'underline' }}
                    >
                      How to use the Triggers page
                    </a>
                  </p>
                </div>

                <div className='formContainer'>
                  <Card
                    interactive={false}
                    elevation={Elevation.TWO}
                    style={{ padding: '2px', minWidth: '320px', width: '100%', maxWidth: '601px', marginBottom: '20px' }}
                  >
                    <h3 style={{ borderBottom: '1px solid #eeeeee', margin: 0, padding: '8px 10px' }}>
                      <span style={{ float: 'right', fontSize: '9px', color: '#aaaaaa' }}>application/json</span> JSON
                    </h3>
                    <TextArea
                      growVertically={true}
                      fill={true}
                      className='codeArea'
                      value={textAreaBody}
                      onChange={(e) => setTextAreaBody(e.target.value)}
                    />
                  </Card>
                </div>

                <Button icon='rocket' intent='primary' onClick={(e) => sendTrigger(e)} disable={triggerDisabled}>Trigger Collection</Button>
              </form>
            </>
          )}
        </main>
      </Content>
    </div>
  )
}
