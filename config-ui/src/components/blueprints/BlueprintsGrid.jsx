import React, { useEffect } from 'react'
import dayjs from '@/utils/time'
import cron from 'cron-validate'
import {
  Classes, FormGroup, InputGroup, ButtonGroup,
  Button, Icon, Intent,
  Dialog, DialogProps,
  RadioGroup, Radio,
  Menu, MenuItem,
  Card, Elevation,
  Popover,
  Tooltip,
  Position,
  Spinner,
  Colors,
  Label,
  Collapse,
  NonIdealState,
  Divider,
  H5,
  Switch,
  Pre,
  Tag
} from '@blueprintjs/core'
import DeletePopover from '@/components/blueprints/DeletePopover'
import EventIcon from '@/images/calendar-3.png'
import EventOffIcon from '@/images/calendar-4.png'

const BlueprintsGrid = (props) => {
  const {
    blueprints = [],
    activeBlueprint,
    blueprintSchedule,
    isActiveBlueprint = (b) => {},
    expandBlueprint = (b) => {},
    deleteBlueprint = (b) => {},
    createCron = () => {},
    handleBlueprintActivation = (b) => {},
    configureBlueprint = (b) => {},
    isDeleting = false,
    expandDetails = false,
    cronPresets
  } = props

  return (
    <>
      <div style={{ display: 'flex', marginTop: '30px', minHeight: '36px', width: '100%', justifyContent: 'flex-start' }}>
        <div
          className='blueprints-list-grid' style={{
            display: 'flex',
            flexDirection: 'column',
            width: '100%',
            minWidth: '830px'
          }}
        >
          {blueprints.map((b, bIdx) => (
            <div key={`blueprint-row-key-${bIdx}`}>
              <div
                style={{
                  display: 'flex',
                  width: '100%',
                  minHeight: '48px',
                  borderBottom: isActiveBlueprint(b.id) && expandDetails ? 'none' : '1px solid #eee',
                  backgroundColor: !b.enable ? '#f8f8f8' : 'inherit',
                  color: !b.enable ? '#555555' : 'inherit',
                }}
              >
                <div
                  className='blueprint-row' style={{
                    display: 'flex',
                    width: '100%',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    padding: '8px 5px',
                    paddingBottom: '16px',
                    position: 'relative'
                  }}
                >
                  <div className='blueprint-id' style={{ flex: 1, maxWidth: '100px' }}>
                    <div style={{ height: '24px', lineHeight: '24px' }}>
                      <label style={{
                        marginLeft: '25px',
                        fontSize: '9px',
                        fontWeight: '400',
                        fontFamily: 'Montserrat, sans-serif',
                        color: '#777777'
                      }}
                      >
                        ID
                      </label>
                    </div>
                    <Button
                      className='bp-row-expand-trigger'
                      onClick={() => expandBlueprint(b)}
                      small minimal style={{
                        minHeight: '20px',
                        minWidth: '20px',
                        marginTop: '-3px',
                        padding: 0,
                        marginRight: '5px',
                        float: 'left'
                      }}
                    >
                      <Icon
                        size={12} color={isActiveBlueprint(b.id) && expandDetails ? Colors.BLUE3 : Colors.GRAY2}
                        icon={isActiveBlueprint(b.id) && expandDetails ? 'collapse-all' : 'expand-all'}
                        style={{ margin: '0' }}
                      />
                    </Button>
                    {b.id}
                  </div>
                  <div
                    className='blueprint-name'
                    style={{ flex: 2, minWidth: '176px', fontWeight: 800 }}
                  >
                    <div style={{ height: '24px', lineHeight: '24px' }}>
                      <label style={{
                        fontSize: '9px',
                        fontWeight: '400',
                        fontFamily: 'Montserrat, sans-serif',
                        color: '#777777'
                      }}
                      >
                        Blueprint Name
                      </label>
                    </div>
                    <Icon
                      size={16}
                      icon={(
                        <img
                          src={b.enable ? EventIcon : EventOffIcon} width={16} height={16}
                          style={{ float: 'left', marginRight: '5px' }}
                        />)}
                      style={{

                      }}
                    />
                    {b.name}
                  </div>
                  <div className='blueprint-interval' style={{ flex: 1, minWidth: '60px' }}>
                    <div style={{ height: '24px', lineHeight: '24px' }}>
                      <label style={{
                        fontSize: '9px',
                        fontWeight: '400',
                        fontFamily: 'Montserrat, sans-serif',
                        color: '#777777'
                      }}
                      >
                        Frequency
                      </label>
                    </div>
                    {b.interval}
                  </div>
                  <div className='blueprint-next-rundate' style={{ flex: 1, whiteSpace: 'nowrap' }}>
                    <div style={{ height: '24px', lineHeight: '24px' }}>
                      <label style={{
                        fontSize: '9px',
                        fontWeight: '400',
                        fontFamily: 'Montserrat, sans-serif',
                        color: '#777777'
                      }}
                      >
                        Next Run Date
                      </label>
                    </div>
                    <div>
                      {dayjs(createCron(b.cronConfig).getNextDate().toString()).format('L LTS')}
                    </div>
                    <div>
                      <span style={{ color: b.enable ? Colors.GREEN5 : Colors.GRAY3, position: 'absolute', bottom: '4px' }}>{b.cronConfig}</span>
                    </div>
                  </div>
                  <div className='blueprint-actions' style={{ flex: 1, textAlign: 'right' }}>
                    <div style={{ height: '24px', lineHeight: '24px' }}>
                      <label style={{
                        fontSize: '9px',
                        fontWeight: '400',
                        fontFamily: 'Montserrat, sans-serif',
                        color: '#777777'
                      }}
                      >
                   &nbsp;
                      </label>
                    </div>
                    <div style={{ display: 'flex', alignItems: 'center', justifySelf: 'flex-end' }}>
                      <Button small minimal style={{ marginLeft: 'auto', marginRight: '5px' }} onClick={() => configureBlueprint(b)}>
                        <Tooltip content='Blueprint Settings'>
                          <Icon icon='cog' size={16} color={Colors.GRAY3} />
                        </Tooltip>
                      </Button>
                      <Popover position={Position.LEFT}>
                        <Button small minimal style={{ marginRight: '10px' }}>
                          <Icon icon='trash' color={Colors.GRAY3} size={15} />
                        </Button>
                        <DeletePopover
                          activeBlueprint={b}
                          onCancel={() => {}}
                          onConfirm={deleteBlueprint}
                          isRunning={isDeleting}
                        />
                      </Popover>

                      <Switch
                        checked={b.enable}
                        label={false}
                        onChange={() => handleBlueprintActivation(b)}
                        style={{ marginBottom: '0' }}
                      />
                    </div>
                  </div>
                </div>
              </div>
              <Collapse isOpen={expandDetails && activeBlueprint.id === b.id}>
                <Card elevation={Elevation.TWO} style={{ padding: '0', margin: '30px 30px', backgroundColor: !b.enable ? '#f8f8f8' : 'initial' }}>
                  <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', margin: '0', padding: '10px' }}>
                    <div>
                      <span style={{ float: 'left', display: 'block', marginRight: '10px' }}>
                        <Spinner size={14} />
                      </span>
                      LOADING ASSOCIATED PIPELINES...
                    </div>
                    <div>
                      <Tag style={{ backgroundColor: b.enable ? Colors.GREEN3 : Colors.GRAY3 }} round='true'>{b.enable ? 'ACTIVE' : 'INACTIVE'}</Tag>
                    </div>
                  </div>
                  <Divider style={{ marginRight: 0, marginLeft: 0 }} />
                  <div style={{ padding: '20px', display: 'flex' }}>
                    <div style={{ flex: 2, paddingRight: '20px' }}>
                      <h3 style={{ margin: 0, textTransform: 'uppercase' }}>Pipeline Run Schedule</h3>
                      <p style={{ margin: 0 }}>Based on the current CRON settings, here are next <strong>5</strong> expected pipeline collection dates.</p>
                      <div style={{ margin: '10px 0' }}>
                        {activeBlueprint?.id && blueprintSchedule.map((s, sIdx) => (
                          <div key={`run-schedule-event-key${sIdx}`} style={{ padding: '6px 4px', opacity: b.enable ? 1 : 0.5 }}>
                            <Icon icon='calendar' size={14} color={b.enable ? Colors.BLUE4 : Colors.GRAY4} style={{ marginRight: '10px' }} />
                            {dayjs(s).format('L LTS')}
                          </div>
                        ))}
                      </div>

                      {!b.enable && (
                        <p style={{ margin: 0, fontSize: '9px', fontFamily: 'Montserrat, sans-serif' }}>
                          <Icon icon='warning-sign' size={11} color={Colors.ORANGE5} style={{ float: 'left', marginRight: '5px' }} />
                          Blueprint is NOT Enabled / Active this schedule will not run.
                        </p>
                      )}
                    </div>
                    <div style={{ flex: 1 }}>
                      <label style={{ color: Colors.GRAY1, fontFamily: 'Montserrat,sans-serif' }}>Blueprint</label>
                      <h3 style={{ marginTop: 0, fontSize: '18px', fontWeight: 800 }}>
                        {b.name}
                      </h3>
                      <label style={{ color: Colors.GRAY1, fontFamily: 'Montserrat,sans-serif' }}>Crontab Configuration</label>
                      <h3 style={{ margin: '0 0 20px 0', fontSize: '18px' }}>{b.cronConfig}</h3>

                      <label style={{ color: Colors.GRAY1, fontFamily: 'Montserrat,sans-serif' }}>Next Run</label>
                      <h3 style={{ margin: '0 0 20px 0', fontSize: '18px' }}>
                        {dayjs(createCron(b.cronConfig).getNextDate().toString()).fromNow()}
                      </h3>

                      <label style={{ color: Colors.GRAY3, fontFamily: 'Montserrat,sans-serif' }}>Operations</label>
                      <div style={{ marginTop: '5px', display: 'flex', justifySelf: 'flex-start', alignItems: 'center', justifyContent: 'left', fontSize: '10px' }}>
                        <Button
                          intent={Intent.PRIMARY}
                          icon='cog'
                          text='Settings'
                          small
                          style={{ marginRight: '8px' }}
                          onClick={() => configureBlueprint(b)}
                        />
                        <Popover>
                          <Button icon='trash' text='Delete' small minimal style={{ marginRight: '8px' }} />
                          <DeletePopover activeBlueprint={activeBlueprint} onCancel={() => {}} onConfirm={deleteBlueprint} isRunning={isDeleting} />
                        </Popover>
                        <Switch
                          checked={b.enable}
                          label={b.enable ? 'Disable' : 'Enable'}
                          onChange={() => handleBlueprintActivation(b)}
                          style={{ marginBottom: '0', fontSize: '11px' }}
                        />
                      </div>
                    </div>

                  </div>
                </Card>

              </Collapse>
            </div>
          ))}
        </div>
      </div>
      <div style={{
        display: 'flex',
        margin: '20px 10px',
        alignSelf: 'flex-start',
        width: '50%',
        fontSize: '11px',
        color: '#555555'
      }}
      >
        <Icon icon='user' size={14} style={{ marginRight: '8px' }} />
        <div>
          <span>by {' '} <strong>Administrator</strong></span><br />
          Displaying {blueprints.length} Blueprints from API.
        </div>
      </div>
    </>
  )
}

export default BlueprintsGrid
