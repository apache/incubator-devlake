import React, { useEffect, useState } from 'react'
import {
  // BrowserRouter as Router,
  useRouteMatch,
} from 'react-router-dom'
import { Button, Card, Elevation } from '@blueprintjs/core'
import SidebarMenu from '@/components/Sidebar/SidebarMenu'
import MenuConfiguration from '@/components/Sidebar/MenuConfiguration'
import { GRAFANA_ENDPOINT } from '@/utils/config'

import '@/styles/sidebar.scss'

const Sidebar = () => {
  const activeRoute = useRouteMatch()

  const [menu, setMenu] = useState(MenuConfiguration(activeRoute))

  useEffect(() => {
    setMenu(MenuConfiguration(activeRoute))
  }, [activeRoute])

  return (
    <Card interactive={false} elevation={Elevation.ZERO} className='card sidebar-card'>
      <img src='/logo.svg' className='logo' />
      <a href={GRAFANA_ENDPOINT} rel='noreferrer' target='_blank' className='dashboardBtnLink'>
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

      <span className='copyright-tag'><strong>Apache 2.0 License</strong><br />&copy; 2021 Merico</span>
    </Card>
  )
}

export default Sidebar
