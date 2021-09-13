import { useRouter } from 'next/router'
import { Button, Card, Elevation, Icon } from '@blueprintjs/core'
import styles from '../styles/Sidebar.module.css'

const Sidebar = () => {
  const { asPath } = useRouter()

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
  </Card>
}

export default Sidebar
