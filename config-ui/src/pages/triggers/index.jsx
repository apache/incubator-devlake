import React, { useState, useEffect } from 'react'
import {
  AnchorButton,
  Spinner,
  Button,
  TextArea,
  Card,
  Elevation
} from '@blueprintjs/core'
import defaultTriggerValue from '@/data/defaultTriggerValue.js'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import AppCrumbs from '@/components/Breadcrumbs'
import Content from '@/components/Content'
import request from '@/utils/request'
import { DEVLAKE_ENDPOINT, GRAFANA_ENDPOINT } from '@/utils/config.js'

const STAGE_INIT = 0
const STAGE_PENDING = 1
const STAGE_COMPELTED = 2
let stage = STAGE_INIT
let targetTaskIds = []

export default function Triggers () {
  const [pendingTasks, setPendingTasks] = useState([])
  const [triggerDisabled, setTriggerDisabled] = useState([])

  // component mounted, run once
  // @todo FIXME: React exhaustive dep warning, this needs to be wrapped in a useCallback (or async function moved inside)
  useEffect(async () => {
    stage = STAGE_INIT
    targetTaskIds = []
    const interval = setInterval(async () => {
      if (stage !== STAGE_PENDING) {
        return
      }
      try {
        const res = await request.get(`${DEVLAKE_ENDPOINT}/task/pending`)
        const tasks = res.data.tasks.filter(t => targetTaskIds.includes(t.ID))
        if (tasks.length === 0) {
          stage = STAGE_COMPELTED
        }
        setPendingTasks(tasks)
      } finally { }
    }, 3000)
    return () => clearInterval(interval)
  }, [])

  // user clicked on trigger button
  const sendTrigger = async (e) => {
    e.preventDefault()
    // @todo RE_ACTIVATE Trigger Process!
    try {
      const res = await request.post(
        `${DEVLAKE_ENDPOINT}/task`,
        textAreaBody
      )
      stage = STAGE_PENDING
      setTriggerDisabled(true)
      targetTaskIds = res.data.flat().map(t => t.ID)
      console.log('waiting following tasks to complete: ', targetTaskIds)
    } catch (e) {
      console.error(e)
    }
  }

  const [textAreaBody, setTextAreaBody] = useState(JSON.stringify(defaultTriggerValue, null, 2))

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
          {stage === STAGE_COMPELTED &&
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

              {pendingTasks.map(task => (
                <div className='pluginSpinnerWrap' key={`key-${task.ID}`}>
                  <div key={`progress-${task.ID}`}>
                    <span style={{ display: 'inline-block', width: '100px' }}>{task.plugin}</span>
                    {task.status === 'TASK_CREATED' &&
                      <>
                        <Spinner
                          size={12}
                          className='pluginSpinner'
                        />
                        <strong>{task.progress * 100}%</strong>
                      </>}
                    {task.status === 'TASK_FAILED' &&
                      <>
                        <span style={{ color: 'red', fontWeight: 'bold' }}>{task.status} </span>
                        {task.message}
                      </>}
                  </div>
                </div>
              ))}
            </div>}
          {stage === STAGE_INIT && (
            <>
              <div className='headlineContainer'>
                <h1>Triggers</h1>
                <p className='description'>Trigger data collection on your plugins</p>
              </div>

              <form className='form'>
                <div className='headlineContainer'>
                  <p className='description'>Create a <strong>http</strong> request to trigger data collect tasks, please replace your&nbsp;
                    Gitlab <code>projectId</code> and JIRA <code>boardId</code> in the request body. This can take&nbsp;
                    up to 20 minutes for large projects. (<strong>Gitlab</strong> 10k+ commits or <strong>JIRA</strong> 5k+ issues)
                  </p>
                  <p className='description'>
                    {/* eslint-disable-next-line max-len */}
                    There are two types of plugins in our application, corresponding to the 2 lists in the following JSON.&nbsp;
                    {/* eslint-disable-next-line max-len */}
                    The regular plugins collect and enrich data (<em>the 1st list</em>) while the domain layer plugins (<em>the 2nd list</em>) prepare the data for the graphs in Grafana dashboards.&nbsp;
                    {/* eslint-disable-next-line max-len */}
                    You <strong>SHOULD ONLY</strong> have to edit the regular plugins. Editing domain layer plugins is for advanced usage only.&nbsp;
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
                      defaultValue={textAreaBody}
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
