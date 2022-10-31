/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
import React, {
  useEffect,
  useState,
  useContext,
  useMemo,
  useCallback
} from 'react'
import {
  // BrowserRouter as Router,
  useRouteMatch
} from 'react-router-dom'
import { Card, Elevation } from '@blueprintjs/core'
import request from '@/utils/request'
import SidebarMenu from '@/components/Sidebar/SidebarMenu'
import MenuConfiguration from '@/components/Sidebar/MenuConfiguration'
import { DEVLAKE_ENDPOINT } from '@/utils/config'
import { ReactComponent as Logo } from '@/images/devlake-logo.svg'
import { ReactComponent as LogoText } from '@/images/devlake-textmark.svg'

import '@/styles/sidebar.scss'
import UIContext from '@/store/UIContext'

const Sidebar = (props) => {
  const { integrations = [] } = props
  const activeRoute = useRouteMatch()
  const uiContext = useContext(UIContext)

  const getMenu = useCallback(
    () => MenuConfiguration(activeRoute, integrations),
    [activeRoute, integrations]
  )

  const ActiveMenu = useMemo(() => getMenu(), [getMenu])

  const [menu, setMenu] = useState(ActiveMenu)
  const [versionTag, setVersionTag] = useState('')

  // useEffect(() => {
  //   setMenu(ActiveMenu)
  // }, [ActiveMenu])

  useEffect(() => {
    // @todo: re-enable version fetch
    // const fetchVersion = async () => {
    //   try {
    //     const versionUrl = `${DEVLAKE_ENDPOINT}/version`
    //     const res = await request.get(versionUrl).catch((e) => {
    //       console.log('>>> API VERSION ERROR...', e)
    //       setVersionTag('')
    //     })
    //     setVersionTag(res?.data ? res.data?.version : '')
    //   } catch (e) {
    //     setVersionTag('')
    //   }
    // }
    // fetchVersion()
  }, [])

  // useEffect(() => {}, [integrations])

  return uiContext.sidebarVisible ? (
    <Card
      interactive={false}
      elevation={Elevation.ZERO}
      className='card sidebar-card'
    >
      <div className='devlake-logo'>
        <Logo width={48} height={48} className='logo' />
        <LogoText width={100} height={13} className='logo-textmark' />
      </div>
      {/* <a href={GRAFANA_URL} rel='noreferrer' target='_blank' className='dashboardBtnLink'>
        <Button icon='grouped-bar-chart' outlined={true} className='dashboardBtn'>View Dashboards</Button>
      </a> */}

      {/* <h3
        className='sidebar-app-heading'
        style={{
          marginTop: '30px',
          letterSpacing: '3px',
          marginBottom: 0,
          fontWeight: 900,
          color: '#444444',
          textAlign: 'center'
        }}
      >
        <sup style={{ fontSize: '9px', color: '#cccccc', marginLeft: '-30px' }}>DEV</sup>LAKE
      </h3> */}
      <SidebarMenu menu={ActiveMenu} />
      <span className='copyright-tag'>
        {/* <span className='version-tag'>{versionTag || ''}</span><br /> */}
        <strong>Apache 2.0 License</strong>
        <br />
      </span>
    </Card>
  ) : (
    <></>
  )
}

export default Sidebar
