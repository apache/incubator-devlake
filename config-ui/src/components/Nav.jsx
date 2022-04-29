import React from 'react'
import {
  Alignment,
  Position,
  Popover,
  Navbar,
  Icon,
} from '@blueprintjs/core'
import '@/styles/nav.scss'
// import { ReactComponent as DiscordIcon } from '@/images/discord.svg'
import { ReactComponent as SlackIcon } from '@/images/slack-mark-monochrome-black.svg'
import { ReactComponent as SlackLogo } from '@/images/slack-rgb.svg'

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
        {/* DISCORD: !DISABLED! */}
        {/* <a href='https://discord.com/invite/83rDG6ydVZ' rel='noreferrer' target='_blank' className='navIconLink'>
          <DiscordIcon className='discordIcon' width={16} height={16} />
        </a> */}
        {/* SLACK: ENABLED (Primary) */}
        <Popover position={Position.LEFT}>
          <SlackIcon className='slackIcon' width={16} height={16} style={{ cursor: 'pointer' }} />
          <>
            <div style={{ maxWidth: '200px', padding: '10px', fontSize: '11px' }}>
              <SlackLogo width={131} height={49} style={{ display: 'block', margin: '0 auto' }} />
              <p style={{ textAlign: 'center' }}>
                Want to interact with the <strong>Merico Community</strong>? Join us on our Slack Channel.<br />
                <a
                  href='https://join.slack.com/t/devlake-io/shared_invite/zt-17b6vuvps-x98pqseoUagM7EAmKC82xQ'
                  rel='noreferrer'
                  target='_blank'
                  className='bp3-button bp3-intent-warning bp3-elevation-1 bp3-small'
                  style={{ marginTop: '10px' }}
                >
                  Message us on&nbsp;<strong>Slack</strong>
                </a>
              </p>
            </div>
          </>
        </Popover>
      </Navbar.Group>
    </Navbar>
  )
}

export default Nav
