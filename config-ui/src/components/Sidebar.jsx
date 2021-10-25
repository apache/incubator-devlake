import React, { useEffect, useState } from 'react'
import { Button, Card, Colors, Elevation, Icon, Tree } from '@blueprintjs/core'
import '../styles/sidebar.scss'

const Sidebar = () => {
  const [isOpen, setIsOpen] = useState(true)
  const [pluginData, setPluginData] = useState()

  useEffect(() => {
    setPluginData([
      {
        id: 5,
        label: 'Collection Plugins',
        isExpanded: isOpen,
        childNodes: [
          {
            id: 0,
            label: <a href='/plugins/jira' className='pluginListItemLink'>Jira</a>,
          },
          {
            id: 1,
            label: <a href='/plugins/gitlab' className='pluginListItemLink'>Gitlab</a>,
          },
          {
            id: 2,
            label: <a href='/plugins/jenkins' className='pluginListItemLink'>Jenkins</a>,
          }
        ]
      },
    ])
  }, [isOpen])

  return (
    <Card interactive={false} elevation={Elevation.ZERO} className='card'>
      <img src='/logo.svg' className='logo' />
      <a href='http://localhost:3002' rel='noreferrer' target='_blank' className='dashboardBtnLink'>
        <Button icon='grouped-bar-chart' outlined={true} large={true} className='dashboardBtn'>View Dashboards</Button>
      </a>

      <ul className='sidebarMenu'>
        <a href='/integrations'>
          {/* IN DEVELOPMENT */}
          <li style={{ color: Colors.RED4, fontWeight: 'bold' }}>
            <Icon icon='data-connection' size={16} className='sidebarMenuListIcon' />
            Data Integrations
          </li>
        </a>
        <a href='/'>
          <li>
            <Icon icon='layout-grid' size={16} className='sidebarMenuListIcon' />
            Configuration
          </li>
          {/* {pagePath === '/' && <div className='sidebarMenuDash'></div>} */}
        </a>
        <a href='/triggers'>
          <li>
            <Icon icon='repeat' size={16} className='sidebarMenuListIcon' />
            Triggers
          </li>
          {/* {pagePath === '/triggers' && <div className='sidebarMenuDash'></div>} */}
        </a>
      </ul>

      <Tree
        contents={pluginData}
        className='pluginMenu'
        // onNodeClick={()=>alert('clicked item')}
        onNodeExpand={() => setIsOpen(true)}
        onNodeCollapse={() => setIsOpen(false)}
      />
    </Card>
  )
}

export default Sidebar
