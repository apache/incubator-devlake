import React from 'react'
import {
  Alignment,
  Navbar,
  Icon,
} from '@blueprintjs/core'
import '../styles/nav.scss'
import { ReactComponent as DiscordIcon } from '@/images/discord.svg'

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
          <DiscordIcon className='discordIcon' width={16} height={16} />
        </a>
      </Navbar.Group>
    </Navbar>
  )
}

export default Nav
