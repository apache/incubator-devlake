import { Button, Card, Elevation } from "@blueprintjs/core"
import styles from '../styles/Sidebar.module.css'

const Sidebar = () => {
  return <Card interactive={false} elevation={Elevation.ZERO} className={styles.card}>

    <img src="/logo.svg" className={styles.logo} />
    <Button icon="grouped-bar-chart" outlined={true} large={true} className={styles.dashboardBtn}>View Dashboards</Button>
  </Card>
}

export default Sidebar
