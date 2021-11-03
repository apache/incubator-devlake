import React, { useEffect, useState } from 'react'
import {
  BrowserRouter as Router,
  // Switch,
  // Route,
  // Link,
  useRouteMatch,
  // useParams
} from 'react-router-dom'
import { Button, Card, Colors, Elevation, Icon, Menu } from '@blueprintjs/core'
import '@/styles/sidebar-menu.scss'

const SidebarMenu = (props) => {
  const { menu } = props
  // const activeRoute = useRouteMatch()

  useEffect(() => {

  }, [menu])

  return (
    <>
      <Menu className='sidebarMenu' style={{ marginTop: '10px' }}>
        {menu.map((m, mIdx) => (
          m.children.length === 0
            ? (
              <Menu.Item
                active={m.active}
                key={`menu-item-key${mIdx}`}
                icon={m.icon}
                text={m.label}
                href={m.route}
                target={m.target}
                disabled={m.disabled}
              />
              )
            : (
              <Menu.Item
                className='is-submenu has-children'
                active={m.active}
                key={`menu-item-key${mIdx}`}
                text={m.label}
                icon={m.icon}
                href={m.route}
                disabled={m.disabled}
              >
                {m.children.map((mS, mSidx) => (
                  <Menu.Item
                    active={mS.active}
                    key={`submenu-${mIdx}-item-key${mSidx}`}
                    href={mS.route}
                    icon={mS.icon}
                    text={mS.label}
                    disabled={mS.disabled}
                    // className={mS.classNames.join(' ')}
                  />))}
              </Menu.Item>
              )
        ))}
        <Menu.Divider />
        <Menu.Item text='API Configuration' icon='cog' href='/lake/api/configuration' active={top.location.href.endsWith('/lake/api/configuration')}/>
        <Menu.Item text='Merico Network' icon='globe-network'>
          <Menu.Item text='Merico GitHub' />
          <Menu.Item text='Triforce Project' />
          <Menu.Item text='Merico Enterprise' />
        </Menu.Item>
      </Menu>
    </>
  )
}

export default SidebarMenu
