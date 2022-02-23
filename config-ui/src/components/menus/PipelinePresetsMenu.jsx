import React from 'react'
import { Menu } from '@blueprintjs/core'
import PipelineConfigsMenu from '@/components/menus/PipelineConfigsMenu'

const PipelinePresetsMenu = (props) => {
  const {
    namePrefix,
    pipelineSuffixes,
    setNamePrefix = () => {},
    setNameSuffix = () => {},
    setRawConfiguration = () => {},
    advancedMode = false
  } = props
  return (
    <Menu className='pipeline-presets-menu'>
      <label style={{
        fontSize: '10px',
        fontWeight: 800,
        fontFamily: '"Montserrat", sans-serif',
        textTransform: 'uppercase',
        padding: '6px 8px',
        display: 'block'
      }}
      >Preset Naming Options
      </label>
      <Menu.Item text='COLLECTION ...' active={namePrefix === 'COLLECT'}>
        <Menu.Item
          icon='key-option'
          text='COLLECT [UNIXTIME]'
          onClick={() => setNamePrefix('COLLECT') | setNameSuffix(pipelineSuffixes[0])}
        />
        <Menu.Item
          icon='key-option'
          text='COLLECT [YYYYMMDDHHMMSS]' onClick={() => setNamePrefix('COLLECT') | setNameSuffix(pipelineSuffixes[3])}
        />
        <Menu.Item
          icon='key-option' text='COLLECT [ISO]'
          onClick={() => setNamePrefix('COLLECT') | setNameSuffix(pipelineSuffixes[2])}
        />
        <Menu.Item icon='key-option' text='COLLECT [UTC]' onClick={() => setNamePrefix('COLLECT') | setNameSuffix(pipelineSuffixes[4])} />
      </Menu.Item>
      <Menu.Item text='SYNCHRONIZE ...' active={namePrefix === 'SYNC'}>
        <Menu.Item
          icon='key-option' text='SYNC [UNIXTIME]'
          onClick={() => setNamePrefix('SYNC') | setNameSuffix(pipelineSuffixes[0])}
        />
        <Menu.Item
          icon='key-option' text='SYNC [YYYYMMDDHHMMSS]'
          onClick={() => setNamePrefix('SYNC') | setNameSuffix(pipelineSuffixes[3])}
        />
        <Menu.Item
          icon='key-option' text='SYNC [ISO]'
          onClick={() => setNamePrefix('SYNC') | setNameSuffix(pipelineSuffixes[2])}
        />
        <Menu.Item
          icon='key-option' text='SYNC [UTC]'
          onClick={() => setNamePrefix('SYNC') | setNameSuffix(pipelineSuffixes[4])}
        />
      </Menu.Item>
      <Menu.Item text='RUN ...' active={namePrefix === 'RUN'}>
        <Menu.Item
          icon='key-option'
          text='RUN [UNIXTIME]'
          onClick={() => setNamePrefix('RUN') | setNameSuffix(pipelineSuffixes[0])}
        />
        <Menu.Item
          icon='key-option' text='RUN [YYYYMMDDHHMMSS]'
          onClick={() => setNamePrefix('RUN') | setNameSuffix(pipelineSuffixes[3])}
        />
        <Menu.Item
          icon='key-option'
          text='RUN [ISO]'
          onClick={() => setNamePrefix('RUN') | setNameSuffix(pipelineSuffixes[2])}
        />
        <Menu.Item
          icon='key-option'
          text='RUN [UTC]'
          onClick={() => setNamePrefix('RUN') | setNameSuffix(pipelineSuffixes[4])}
        />
      </Menu.Item>
      <Menu.Divider />
      <Menu.Item text='Advanced Options' icon='cog'>
        <Menu.Item icon='new-object' text='Save Pipeline Blueprint' disabled />
        {advancedMode && (
          <PipelineConfigsMenu
            setRawConfiguration={setRawConfiguration}
            advancedMode={advancedMode}
          />
        )}
      </Menu.Item>
    </Menu>
  )
}

export default PipelinePresetsMenu
