import { Button, Card, Elevation } from "@blueprintjs/core"
import styles from '../styles/Sidebar.module.css'

const Sidebar = () => {
  return <Card interactive={false} elevation={Elevation.ZERO} className={styles.card}>

    <img src="/logo.svg" className={styles.logo} />
    <a href="http://localhost:3002" target="_blank" className={styles.dashboardBtnLink}>
      <Button icon="grouped-bar-chart" outlined={true} large={true} className={styles.dashboardBtn}>View Dashboards</Button>
    </a>
  </Card>
}

export default Sidebar
