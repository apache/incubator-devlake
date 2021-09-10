import { Button, Card, Elevation, Icon } from "@blueprintjs/core"
import styles from '../styles/Sidebar.module.css'

const Sidebar = () => {
  return <Card interactive={false} elevation={Elevation.ZERO} className={styles.card}>

    <img src="/logo.svg" className={styles.logo} />
    <a href="http://localhost:3002" target="_blank" className={styles.dashboardBtnLink}>
      <Button icon="grouped-bar-chart" outlined={true} large={true} className={styles.dashboardBtn}>View Dashboards</Button>
    </a>

    <ul className={styles.sidebarMenu}>
      <a href="/">
        <li>
          <Icon icon="layout-grid" size={16} className={styles.sidebarMenuListIcon} />
          Configuration
        </li>
          <div className={styles.sidebarMenuDash}></div>
      </a>
    </ul>
  </Card>
}

export default Sidebar
