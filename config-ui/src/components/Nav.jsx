import React from 'react'
import {
  Alignment,
  Button,
  Navbar,
  Icon,
} from '@blueprintjs/core'
// import styles from '../styles/Nav.module.css'

const Nav = () => {

  return <Navbar className='navbar'>
    <Navbar.Group align={Alignment.RIGHT}>
        <a href='https://github.com/merico-dev/lake' target='_blank' className='navIconLink'>
          <Icon icon='git-branch' size={16} />
        </a>
        <Navbar.Divider />
        <a href='mailto:hello@merico.dev' target='_blank' className='navIconLink'>
        <Icon icon='envelope' size={16} />
        </a>
    </Navbar.Group>
  </Navbar>
}

export default Nav
