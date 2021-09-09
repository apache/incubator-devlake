import {
  Alignment,
  Button,
  Navbar,
  Icon,
} from '@blueprintjs/core'
import styles from '../styles/Nav.module.css'

const Nav = () => {

  return <Navbar className={styles.navbar}>
    <Navbar.Group align={Alignment.RIGHT}>
        <a href="https://github.com/merico-dev/lake" target="_blank">
          <Icon icon="git-branch" size={14} />
        </a>
        <Navbar.Divider />
        <a href="mailto:hello@merico.dev" target="_blank">
        <Icon icon="envelope" size={14} />
        </a>
    </Navbar.Group>
  </Navbar>
}

export default Nav
