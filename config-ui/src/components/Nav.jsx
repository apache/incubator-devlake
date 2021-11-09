import React from 'react'
import {
  Alignment,
  Navbar,
  Icon,
} from '@blueprintjs/core'
import '../styles/nav.scss'

const DISCORD_IMG_URL = './src/images/icon_clyde_black_RGB.svg'
const Nav = () => {
  return (
    <Navbar className='navbar'>
      <Navbar.Group align={Alignment.RIGHT}>
        <a href='https://github.com/merico-dev/lake' rel='noreferrer' target='_blank' className='navIconLink'>
          <Icon icon='git-branch' size={16} />
        </a>
        <Navbar.Divider />
        <a href='mailto:hello@merico.dev' rel='noreferrer' target='_blank' className='navIconLink'>
          <Icon icon='envelope' size={16} />
        </a>
        <Navbar.Divider />
        <a href='https://discord.com/invite/83rDG6ydVZ' rel='noreferrer' target='_blank' className='navIconLink'>
          <img className='discordIcon' src={DISCORD_IMG_URL}/>
        </a>
      </Navbar.Group>
    </Navbar>
  )
}

export default Nav
