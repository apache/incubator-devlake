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
import React, { useContext, useEffect, useState } from 'react'
import {
  Alignment,
  Position,
  Popover,
  Navbar,
  Icon,
} from '@blueprintjs/core'
import '@/styles/nav.scss'
import { ReactComponent as SlackIcon } from '@/images/slack-mark-monochrome-black.svg'
import { ReactComponent as SlackLogo } from '@/images/slack-rgb.svg'
import UIContext from '../store/UIContext'
import useWindowSize from '../hooks/useWIndowSize'


const Nav = () => {
  const uiContext = useContext(UIContext)
  const [menuClass, setMenuClass] = useState('navbarMenuButton')
  const size = useWindowSize()

  const toggleSidebarOpen = (open) => {
    uiContext.changeSidebarVisibility(open)
    setMenuClass((currentMenuClass) => open ? 'navbarMenuButtonSidebarOpened' : currentMenuClass)
  }

  useEffect(() => {
    toggleSidebarOpen((size.width >= uiContext.desktopBreakpointWidth && uiContext.sidebarVisible != true) ? true : false)
  }, [size])

  return (
    <Navbar className='navbar'>
      <Navbar.Group className={menuClass}>
        <Icon icon={uiContext.sidebarVisible ? 'menu-closed' : 'menu-open'} onClick={toggleSidebarOpen.bind(null, !uiContext.sidebarVisible)} size={16} />
      </Navbar.Group>
      {!uiContext.sidebarVisible && <Navbar.Group className='navbarItems'>
        <a href='https://github.com/apache/incubator-devlake' rel='noreferrer' target='_blank' className='navIconLink'>
          <Icon icon='git-branch' size={16} />
        </a>
        <Navbar.Divider />
        <a href='mailto:hello@merico.dev' rel='noreferrer' target='_blank' className='navIconLink'>
          <Icon icon='envelope' size={16} />
        </a>
        <Navbar.Divider />
        {/* DISCORD: !DISABLED! */}
        {/* <a href='https://discord.com/invite/83rDG6ydVZ' rel='noreferrer' target='_blank' className='navIconLink'>
          <DiscordIcon className='discordIcon' width={16} height={16} />
        </a> */}
        {/* SLACK: ENABLED (Primary) */}
        <Popover position={Position.LEFT}>
          <SlackIcon className='slackIcon' width={16} height={16} style={{ cursor: 'pointer' }} />
          <>
            <div style={{ maxWidth: '200px', padding: '10px', fontSize: '11px' }}>
              <SlackLogo width={131} height={49} style={{ display: 'block', margin: '0 auto' }} />
              <p style={{ textAlign: 'center' }}>
                Want to interact with the <strong>Merico Community</strong>? Join us on our Slack Channel.<br />
                <a
                  href='https://join.slack.com/t/devlake-io/shared_invite/zt-17b6vuvps-x98pqseoUagM7EAmKC82xQ'
                  rel='noreferrer'
                  target='_blank'
                  className='bp3-button bp3-intent-warning bp3-elevation-1 bp3-small'
                  style={{ marginTop: '10px' }}
                >
                  Message us on&nbsp
                  <strong>Slack</strong>
                </a>
              </p>
            </div>
          </>
        </Popover>
      </Navbar.Group>}
    </Navbar>
  )
}

export default Nav
