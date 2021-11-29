import React, { useState, useEffect } from 'react'
import {
  AnchorButton,
  Spinner,
  Button,
  TextArea,
  Card,
  Elevation,
  Colors,
  Intent,
  Icon
} from '@blueprintjs/core'
import Nav from '../../components/Nav'
import Sidebar from '../../components/Sidebar'
import AppCrumbs from '@/components/Breadcrumbs'
import Content from '../../components/Content'
import Config from '../../../config'
import request from '../../utils/request'

export default function Documentation () {

  const [pendingTasks, setPendingTasks] = useState([])
  const [stage, setStage] = useState(0)
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
          <AppCrumbs
            items={[
              { href: '/', icon: false, text: 'Dashboard' },
              { href: '/triggers', icon: false, text: 'Data Triggers' },
            ]}
          />

          <>
            <div className='headlineContainer'>
              <h1 style={{ margin: 0 }}>Documentation</h1>
              <h2 style={{ margin: '0 0 20px 0' }}>DOWNLOAD READMEs ON GITHUB</h2>
              <p style={{ fontSize: '20px', color: '#444444' }}>
                Dev Lake is the one-stop solution for engineering teams. that <strong>integrates</strong>
                , <strong>analyzes</strong>, and <strong>visualizes</strong>
                &nbsp; data throughout the software development life cycle (SDLC).
              </p>
            </div>

            <div style={{ justifyContent: 'flex-start', width: '100%' }}>
              <h3>DATA SOURCES <Icon icon='backlink' size={14} color={Colors.BLUE4} /></h3>

              <h3>USER SETUP <Icon icon='backlink' size={14} color={Colors.BLUE4} /></h3>

              <h3>DEVELOPER SETUP <Icon icon='backlink' size={14} color={Colors.BLUE4} /></h3>

            </div>
          </>
        </main>
      </Content>
    </div>
  )
}
