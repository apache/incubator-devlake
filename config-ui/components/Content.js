import { Button, Card, Elevation } from "@blueprintjs/core"
import styles from '../styles/Content.module.css'

const Content = () => {
  return <Card interactive={false} elevation={Elevation.THREE} className={styles.card}>
    <h5><a href="#">Card heading</a></h5>
    <p>Card content</p>
    <Button>Submit</Button>
  </Card>
}

export default Content
