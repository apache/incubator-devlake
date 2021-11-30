import React, { useEffect, Fragment } from 'react'
import { Menu } from '@blueprintjs/core'
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
        {/* <Menu.Item
          text='API Configuration'
          icon='cog' href='/lake/api/configuration'
          active={top.location.href.endsWith('/lake/api/configuration')}
        /> */}
        <Menu.Item text='Documentation' icon='help' href='https://github.com/merico-dev/lake/blob/main/README.md' target='_blank' />
        <Menu.Item text='Merico Network' icon='globe-network'>
          <Menu.Item text='Merico GitHub' href='https://github.com/merico-dev' target='_blank' />
          <Menu.Item text='DevLake Github' href='https://github.com/merico-dev/lake' target='_blank' />
          <Menu.Item text='Merico Enterprise' href='https://meri.co/' target='_blank' />
        </Menu.Item>
      </Menu>
    </>
  )
}

export default SidebarMenu
