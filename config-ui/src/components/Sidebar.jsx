import { useEffect, useState } from 'react'
import { useRouter } from 'next/router'
import { Button, Card, Elevation, Icon, Tree, Classes } from '@blueprintjs/core'
import styles from '../styles/Sidebar.module.css'

const Sidebar = () => {
  const { asPath } = useRouter() || {asPath: { text: '' }}

  const [isOpen, setIsOpen] = useState(true)
  const [pluginData, setPluginData] = useState()

  useEffect(() => {
    setPluginData([
      {
        id: 5, label: "Collection Plugins",
        isExpanded: isOpen,
        childNodes: [
          {
            id: 0,
            label: <a href="/plugins/jira" className={styles.pluginListItemLink}>Jira</a>,
          },
          {
            id: 1,
            label: <a href="/plugins/gitlab" className={styles.pluginListItemLink}>Gitlab</a>,
          },
          {
            id: 2,
            label: <a href="/plugins/jenkins" className={styles.pluginListItemLink}>Jenkins</a>,
          }
        ]
      },
    ])
  }, [isOpen])

  return <Card interactive={false} elevation={Elevation.ZERO} className={styles.card}>

    <img src="/logo.svg" className={styles.logo} />
    <a href="http://localhost:3002" target="_blank" className={styles.dashboardBtnLink}>
      <Button icon="grouped-bar-chart" outlined={true} large={true} className={styles.dashboardBtn}>View Dashboards</Button>
    </a>

    <ul className={styles.sidebarMenu}>
      <a href="/" className={asPath === "/" ? styles.sidebarMenuActive : ''}>
        <li>
          <Icon icon="layout-grid" size={16} className={styles.sidebarMenuListIcon} />
          Configuration
        </li>
          {asPath === "/" && <div className={styles.sidebarMenuDash}></div>}
      </a>
      <a href="/triggers" className={asPath === "/triggers" ? styles.sidebarMenuActive: ''}>
        <li>
          <Icon icon="repeat" size={16} className={styles.sidebarMenuListIcon} />
          Triggers
        </li>
          {asPath === "/triggers" && <div className={styles.sidebarMenuDash}></div>}
      </a>
    </ul>

    <Tree
      contents={pluginData}
      className={styles.pluginMenu}
      // onNodeClick={()=>alert('clicked item')}
      onNodeExpand={()=>setIsOpen(true)}
      onNodeCollapse={()=>setIsOpen(false)}
    />
  </Card>
}

export default Sidebar
