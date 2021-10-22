import React, { useState, useEffect } from 'react'
import axios from 'axios'
import { AnchorButton, Spinner, Button, TextArea, Intent } from '@blueprintjs/core'
import defaultTriggerValue from '../../data/defaultTriggerValue.js'
import Nav from '../../components/Nav'
import Sidebar from '../../components/Sidebar'
import Content from '../../components/Content'
import Config from '../../../config'
import request from '../../utils/request'

export default function Triggers () {
  const [textAreaBody, setTextAreaBody] = useState(JSON.stringify(defaultTriggerValue, null, 2))

  const sendTrigger = async (e) => {
    e.preventDefault()
    console.log('JON >>> Config.DEVLAKE_ENDPOINT', Config.DEVLAKE_ENDPOINT);
    try {
      await request.post(
        `${Config.DEVLAKE_ENDPOINT}/task`,
        textAreaBody
      )
    } catch (e) {
      console.error(e)
    }
  }

  const [pendingTasks, setPendingTasks] = useState([])
  const [stage, setStage] = useState(0)
  const [grafanaUrl, setGrafanaUrl] = useState(3002)
  useEffect(() => {
    let s = 0
    const interval = setInterval(async () => {
      try {
        const res = await request.get('/api/triggers/pendings')
        console.log(await res.data)
        if (res.data.tasks.length > 0) {
          s = 1
        } else if (s === 1) {
          s = 2
        }
        setStage(s)
        setPendingTasks(res.data.tasks)
        setGrafanaUrl(`${location.protocol}//${location.hostname}:${res.data.grafanaPort}`)
      } catch (e) {
        console.log(e)
      }
    }, 3000)
    return () => clearInterval(interval)
  }, [])

  return (
    <div className='container'>
      <Nav />
      <Sidebar />
      <Content>
        <main className='main'>
          {
          stage === 2 &&
            <div className='headlineContainer'>
              <h1>Done</h1>
              <p className='description'>Navigate to Grafana to view updated metrics</p>
              <AnchorButton
                href={grafanaUrl}
                icon='grouped-bar-chart'
                target='_blank'
                text='View Dashboards'
              />
            </div>
          }
          {stage === 1 &&
            <div className='headlineContainer'>
              <h1>Collecting Data</h1>
              <p className='description'>Please wait... </p>

            {pendingTasks.map(task => <div className='pluginSpinnerWrap' key={`key-${task.ID}`}>
                <Spinner
                  size={12}
                  value={task.progress ? task.progress : null}
                  className='pluginSpinner'
                />
              <div key={`progress-${task.ID}`}>
                  {task.plugin}: <strong>{task.progress * 100}%</strong>
                </div>
              </div>
              )}
            </div>
          }
          {stage === 0 && <>
            <div className='headlineContainer'>
              <h1>Triggers</h1>
              <p className='description'>Trigger data collection on your plugins</p>
            </div>

            <form className='form'>
              <div className='headlineContainer'>
                <p className='description'>Create a http request to trigger data collect tasks, please replace your&nbsp;
                  <code>gitlab projectId</code> and <code>jira boardId</code> in the request body. This can take&nbsp;
                  up to 20 minutes for large projects. (gitlab 10k+ commits or jira 5k+ issues)
                </p>
              </div>

              <div className='formContainer'>
                <TextArea
                  growVertically={true}
                  large={true}
                  intent={Intent.PRIMARY}
                  fill={true}
                  className='codeArea'
                  defaultValue={textAreaBody}
                  onChange={(e) => setTextAreaBody(e.target.value)}
                />
              </div>

              <Button outlined={true} large={true} className='saveBtn' onClick={(e) => sendTrigger(e)}>Trigger Collection</Button>
            </form>
            </>
          }
        </main>
      </Content>
    </div>
  )
}
