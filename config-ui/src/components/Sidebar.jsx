import React, { useEffect, useState } from 'react'
import {
  // BrowserRouter as Router,
  useRouteMatch,
} from 'react-router-dom'
import { Button, Card, Elevation } from '@blueprintjs/core'
import request from '@/utils/request'
import SidebarMenu from '@/components/Sidebar/SidebarMenu'
import MenuConfiguration from '@/components/Sidebar/MenuConfiguration'
import { DEVLAKE_ENDPOINT, GRAFANA_URL } from '@/utils/config'

import '@/styles/sidebar.scss'

const Sidebar = () => {
  const activeRoute = useRouteMatch()

  const [menu, setMenu] = useState(MenuConfiguration(activeRoute))
  const [versionTag, setVersionTag] = useState()

  useEffect(() => {
    setMenu(MenuConfiguration(activeRoute))
  }, [activeRoute])

  useEffect(() => {
    const fetchVersion = async () => {
      try {
        const versionUrl = `${DEVLAKE_ENDPOINT}/version`
        const res = await request.get(versionUrl).catch(e => {
          console.log('>>> API VERSION ERROR...', e)
          setVersionTag('dev+error')
        })
        setVersionTag(res?.data ? res.data?.version : 'dev+error')
      } catch (e) {
        setVersionTag('dev+error')
      }
    }
    fetchVersion()
  }, [])

  return (
    <Card interactive={false} elevation={Elevation.ZERO} className='card sidebar-card'>
      <img src='/logo.svg' className='logo' />
      <a href={GRAFANA_URL} rel='noreferrer' target='_blank' className='dashboardBtnLink'>
        <Button icon='grouped-bar-chart' outlined={true} className='dashboardBtn'>View Dashboards</Button>
      </a>

      <h3
        className='sidebar-app-heading'
        style={{
          marginTop: '30px',
          letterSpacing: '3px',
          marginBottom: 0,
          fontFamily: '"Montserrat", sans-serif',
          fontWeight: 900,
          color: '#444444',
          textAlign: 'center'
        }}
      >
        <sup style={{ fontSize: '9px', color: '#cccccc', marginLeft: '-30px' }}>DEV</sup>LAKE
      </h3>
      <SidebarMenu menu={menu} />
      <span className='copyright-tag'>
        <span className='version-tag'>{versionTag || 'dev+unknown'}</span><br />
        <strong>Apache 2.0 License</strong><br />&copy; 2021 Merico
      </span>
    </Card>
  )
}

export default Sidebar
