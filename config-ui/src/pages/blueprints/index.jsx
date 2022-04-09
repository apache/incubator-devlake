import React, { Fragment, useEffect, useState, useRef, useCallback } from 'react'
import { useHistory } from 'react-router-dom'
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
  Divider,
  H5,
  Switch,
  Pre,
  Tag
} from '@blueprintjs/core'
import { parseCronExpression } from 'cron-schedule'
import usePipelineManager from '@/hooks/usePipelineManager'
import useBlueprintManager from '@/hooks/useBlueprintManager'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import AppCrumbs from '@/components/Breadcrumbs'
import Content from '@/components/Content'
import ContentLoader from '@/components/loaders/ContentLoader'
import { ReactComponent as HelpIcon } from '@/images/help.svg'
import ManageBlueprintsIcon from '@/images/blueprints.png'
import EventIcon from '@/images/calendar-3.png'
import EventOffIcon from '@/images/calendar-4.png'
import { NullBlueprint } from '@/data/NullBlueprint'

const DeletePopover = (props) => {
  const {
    activeBlueprint,
    onCancel = () => {},
    onConfirm = () => {},
    isRunning = false
  } = props
  return (
    <>
      <div style={{ padding: '10px', fontSize: '10px', maxWidth: '220px' }}>
        <h3 style={{ margin: '0 0 5px 0', color: Colors.RED3 }}>Delete {activeBlueprint?.name}?</h3>
        <p><strong>Are you sure? This Blueprint will be removed, all pipelines will be stopped.</strong></p>
        <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
          <Button
            className={Classes.POPOVER_DISMISS}
            intent={Intent.NONE}
            text='CANCEL'
            small style={{ marginRight: '5px' }}
            onClick={() => onCancel(activeBlueprint)}
            disabled={isRunning}
          />
          <Button disabled={isRunning} intent={Intent.DANGER} text='YES' small onClick={() => onConfirm(activeBlueprint)} />
        </div>
      </div>
    </>
  )
}

const Blueprints = (props) => {
  const history = useHistory()
  // const { providerId } = useParams()
  // const [activeProvider, setActiveProvider] = useState(integrationsData[0])

  const {
    blueprint,
    blueprints,
    name,
    cronConfig,
    customCronConfig,
    cronPresets,
    tasks,
    enable,
    setName: setBlueprintName,
    setCronConfig,
    setCustomCronConfig,
    setTasks: setBlueprintTasks,
    setEnable: setEnableBlueprint,
    isFetching: isFetchingBlueprints,
    isSaving,
    isDeleting,
    createCronExpression: createCron,
    getCronSchedule: getSchedule,
    activateBlueprint,
    deactivateBlueprint,
    fetchBlueprint,
    fetchAllBlueprints,
    saveBlueprint,
    deleteBlueprint,
    saveComplete,
    deleteComplete
  } = useBlueprintManager()

  // BLUEPRINTS MOCK DATA
  // const [blueprints, setBlueprints] = useState([
  //   { id: 5, name: 'GITHUB DAILY', cronConfig: '0 0 * * *', nextRunAt: null, enable: true, interval: 'Daily', tasks: [[]], createdAt: Date.now(), updatedAt: null },
  //   { id: 6, name: 'GITLAB WEEKLY', cronConfig: '0 0 * * 1', nextRunAt: null, enable: true, interval: 'Weekly', tasks: [[]], createdAt: Date.now(), updatedAt: null },
  //   { id: 7, name: 'GITHUB MONTHLY', cronConfig: '0 0 30 * 1', nextRunAt: null, enable: true, interval: 'Monthly', tasks: [[]], createdAt: Date.now(), updatedAt: null },
  //   { id: 8, name: 'JIRA DAILY', cronConfig: '0 23 * * 1-5', nextRunAt: null, enable: false, interval: 'Daily', tasks: [[]], createdAt: Date.now(), updatedAt: null },
  //   { id: 9, name: 'JENKINS DAILY 8AM @hezyin', cronConfig: '0 0 * * 1-5', nextRunAt: null, enable: false, interval: 'Daily', tasks: [[]], createdAt: Date.now(), updatedAt: null },
  //   { id: 10, name: 'GITLAB CUSTOM @klesh', cronConfig: '0 4 8-14 * *', nextRunAt: null, enable: false, interval: 'Custom', tasks: [[]], createdAt: Date.now(), updatedAt: null },
  //   { id: 11, name: 'JIRA CUSTOM @e2corporation', cronConfig: '0 22 * * 1-5', nextRunAt: null, enable: true, interval: 'Custom', tasks: [[]], createdAt: Date.now(), updatedAt: null },
  // ])

  const [expandDetails, setExpandDetails] = useState(false)
  const [activeBlueprint, setActiveBlueprint] = useState(null)
  const [draftBlueprint, setDraftBlueprint] = useState(null)
  const [blueprintSchedule, setBlueprintSchedule] = useState([])
  // const [customCron, setCustomCron] = useState('0 0 * * *')

  const [blueprintDialogIsOpen, setBlueprintDialogIsOpen] = useState(false)
  const [blueprintErrors, setBlueprintErrors] = useState([])

  const handleBlueprintActivation = (blueprint) => {
    if (blueprint.enable) {
      deactivateBlueprint(blueprint)
    } else {
      activateBlueprint(blueprint)
    }
  }

  const expandBlueprint = (blueprint) => {
    setExpandDetails(opened => blueprint.id === activeBlueprint?.id && opened ? false : !opened)
    setActiveBlueprint(blueprint)
  }

  const configureBlueprint = (blueprint) => {
    setDraftBlueprint(b => ({ ...b, ...blueprint }))
  }

  const createNewBlueprint = () => {
    setDraftBlueprint(null)
    setExpandDetails(false)
    setBlueprintName('DAILY BLUEPRINT')
    setCronConfig('0 0 * * *')
    setCustomCronConfig('0 0 * * *')
    setBlueprintDialogIsOpen(true)
  }

  const isActiveBlueprint = (bId) => {
    return activeBlueprint?.id === bId
  }

  const isStandardCronPreset = (cronConfig) => {
    return cronPresets.some(p => p.cronConfig === cronConfig)
  }

  const fieldHasError = (fieldId) => {
    return blueprintErrors.some(e => e.includes(fieldId))
  }

  const getFieldError = (fieldId) => {
    return blueprintErrors.find(e => e.includes(fieldId))
  }

  useEffect(() => {
    if (activeBlueprint) {
      console.log(getSchedule(activeBlueprint?.cronConfig))
    }
    setBlueprintSchedule(activeBlueprint?.id ? getSchedule(activeBlueprint.cronConfig) : [])
    console.log('>>> ACTIVE/EXPANDED BLUEPRINT', activeBlueprint)
  }, [activeBlueprint, getSchedule])

  useEffect(() => {
    if (draftBlueprint && draftBlueprint.id) {
      setBlueprintName(draftBlueprint.name)
      setCronConfig(!isStandardCronPreset(draftBlueprint.cronConfig) ? 'custom' : draftBlueprint.cronConfig)
      setCustomCronConfig(draftBlueprint.cronConfig)
      setBlueprintDialogIsOpen(true)
    }
  }, [draftBlueprint, setBlueprintName, setCronConfig])

  useEffect(() => {
    if (saveComplete?.id) {
      setBlueprintDialogIsOpen(false)
      fetchAllBlueprints()
    }
  }, [saveComplete])

  useEffect(() => {
    if (deleteComplete.status === 200) {
      fetchAllBlueprints()
    }
  }, [deleteComplete])

  useEffect(() => {
    fetchAllBlueprints()
  }, [fetchAllBlueprints])

  return (
    <>
      <div className='container'>
        <Nav />
        <Sidebar />
        <Content>
          <main className='main'>
            <AppCrumbs
              items={[
                { href: '/', icon: false, text: 'Dashboard' },
                { href: '/pipelines', icon: false, text: 'Pipelines' },
                { href: '/blueprints', icon: false, text: 'Blueprints', current: true },
              ]}
            />
            <div className='headlineContainer'>
              <div style={{ display: 'flex' }}>
                <div>
                  <span style={{ marginRight: '10px' }}>
                    <Icon icon={<img src={ManageBlueprintsIcon} width='38' height='38' />} size={38} />
                  </span>
                </div>
                <div>
                  <h1 style={{ margin: 0 }}>
                    Pipeline Blueprints
                    <Popover
                      className='trigger-manage-blueprints-help'
                      popoverClassName='popover-help-manage-blueprints'
                      position={Position.RIGHT}
                      autoFocus={false}
                      enforceFocus={false}
                      usePortal={false}
                    >
                      <a href='#' rel='noreferrer'><HelpIcon width={19} height={19} style={{ marginLeft: '10px' }} /></a>
                      <>
                        <div style={{ textShadow: 'none', fontSize: '12px', padding: '12px', maxWidth: '300px' }}>
                          <div style={{ marginBottom: '10px', fontWeight: 700, fontSize: '14px', fontFamily: '"Montserrat", sans-serif' }}>
                            <Icon icon='help' size={16} /> Schedule Recurring Pipelines
                          </div>
                          <p>Need Help? &mdash; Manage, Stop running and Restart failed pipelines.
                            Access <strong>Task Progress</strong> and Activity for all your pipelines.
                          </p>
                        </div>
                      </>
                    </Popover>
                  </h1>
                  <p className='page-description mb-0'>Create scheduled plans for automating pipelines with CRON.</p>
                  <p className=''>Choose a preset schedule or use your custom crontab configuration.</p>
                </div>
                <div style={{ marginLeft: 'auto' }}>
                  <Button icon='add' intent={Intent.PRIMARY} text='Create Blueprint' onClick={() => createNewBlueprint()} />
                </div>
              </div>
            </div>
            {(!isFetchingBlueprints) && (
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
                    {/* <Card
                      elevation={Elevation.ZERO}
                      style={{ boxShadow: 'none', padding: '8px', marginBottom: '10px', borderBottom: '1px solid #bbbbbb' }}
                    >
                      <div
                        className='blueprint-header-row'
                        style={{
                          margin: 'auto auto',
                          display: 'flex',
                          width: '100%',
                          justifyContent: 'space-between',
                          color: '#777777',
                          fontFamily: 'Montserrat, sans-serif',
                          fontSize: '11px',
                        }}
                      >
                        <div className='blueprint-header-id' style={{ flex: 1, maxWidth: '100px' }}>ID</div>
                        <div className='blueprint-header-name' style={{ flex: 2 }}>Blueprint Name</div>
                        <div className='blueprint-header-interval' style={{ flex: 2 }}>Frequency</div>
                        <div className='blueprint-header-next-rundate' style={{ flex: 1 }}>Next Run Date</div>
                        <div className='blueprint-header-actions' style={{ flex: 1, textAlign: 'right' }}>&nbsp;</div>
                      </div>
                    </Card> */}
                    {blueprints.map((b, bIdx) => (
                      <div key={`blueprint-row-key-${bIdx}`}>
                        {/* <div
                          className='blueprint-header-row'
                          style={{
                            margin: 'auto auto',
                            marginBottom: '-10px',
                            display: 'flex',
                            width: '100%',
                            justifyContent: 'space-between',
                            color: '#777777',
                            fontFamily: 'Montserrat, sans-serif',
                            fontSize: '9px',
                            backgroundColor: !b.enable ? '#f8f8f8' : 'inherit',
                            paddingTop: '10px'
                          }}
                        >
                          <div className='blueprint-header-id' style={{ flex: 1, maxWidth: '100px', paddingLeft: '30px' }}>ID</div>
                          <div className='blueprint-header-name' style={{ flex: 2 }}>Blueprint Name</div>
                          <div className='blueprint-header-interval' style={{ flex: 2 }}>Frequency</div>
                          <div className='blueprint-header-next-rundate' style={{ flex: 1 }}>Next Run Date</div>
                          <div className='blueprint-header-actions' style={{ flex: 1, textAlign: 'right' }}>&nbsp;</div>
                        </div> */}
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
                              <div>{dayjs(createCron(b.cronConfig).getNextDate().toString()).format('L LTS')}</div>
                              <div>
                                <span style={{ color: b.enable ? Colors.GREEN5 : Colors.GRAY3 }}>{b.cronConfig}</span>
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
                                    {/* <>
                                      <div style={{ padding: '10px', fontSize: '10px', maxWidth: '220px' }}>
                                        <h3 style={{ margin: '0 0 5px 0', color: Colors.RED3 }}>Delete {activeBlueprint?.name}?</h3>
                                        <p><strong>Are you sure? This Blueprint will be removed, all pipelines will be stopped.</strong></p>
                                        <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
                                          <Button className={Classes.POPOVER_DISMISS} intent={Intent.NONE} text='CANCEL' small style={{ marginRight: '5px' }} />
                                          <Button intent={Intent.DANGER} text='YES' small />
                                        </div>
                                      </div>
                                    </> */}
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
              </>)}
          </main>
        </Content>
      </div>

      <Dialog
        className='dialog-manage-blueprint'
        icon={draftBlueprint ? 'edit' : 'add'}
        title={draftBlueprint ? `Edit ${draftBlueprint.name}` : 'Create Pipeline Blueprint'}
        isOpen={blueprintDialogIsOpen}
        onClose={() => setBlueprintDialogIsOpen(false)}
        onClosed={() => setDraftBlueprint(null)}
        style={{ backgroundColor: '#ffffff' }}
      >
        <div className={Classes.DIALOG_BODY}>

          <div className='pipeline-form-container'>
            <div className='formContainer'>
              <FormGroup
                label=''
                inline={true}
                labelFor='blueprint-name'
                className='formGroup-inline'
                contentClassName='formGroupContent'
              >
                <Label style={{ display: 'inline', marginRight: 0 }}>

                  Blueprint Name

                  <span className='requiredStar'>*</span>
                </Label>
                <InputGroup
                  id='blueprint-name'
                  placeholder='Enter Blueprint Name'
                  value={name}
                  onChange={(e) => setBlueprintName(e.target.value)}
                  className={`blueprint-name-input ${fieldHasError('Blueprint Name') ? 'invalid-field' : ''}`}
                  inline={true}
                  style={{ marginBottom: '10px' }}
                />
                <Label style={{ display: 'inline', marginRight: 0, marginBottom: 0 }}>
                  Frequency
                  <span className='requiredStar'>*</span>
                </Label>
                <RadioGroup
                  inline={true}
                  label={false}
                  name='blueprint-frequency'
                  onChange={(e) => setCronConfig(e.target.value)}
                  selectedValue={cronConfig}
                  required
                >
                  <Radio label='Hourly' value='59 * * * 1-5' style={{ fontWeight: cronConfig === '59 * * * 1-5' ? 'bold' : 'normal' }} />
                  <Radio label='Daily' value='0 0 * * *' style={{ fontWeight: cronConfig === '0 0 * * *' ? 'bold' : 'normal' }} />
                  <Radio label='Weekly' value='0 0 * * 1' style={{ fontWeight: cronConfig === '0 0 * * 1' ? 'bold' : 'normal' }} />
                  <Radio label='Monthly' value='0 0 1 * *' style={{ fontWeight: cronConfig === '0 0 1 * *' ? 'bold' : 'normal' }} />
                  <Radio label='Custom' value='custom' style={{ fontWeight: cronConfig === 'custom' ? 'bold' : 'normal' }} />
                </RadioGroup>

                {/* {cronConfig === 'custom' && ( */}
                <>
                  <div className='formContainer'>
                    <FormGroup
                      disabled={cronConfig !== 'custom'}
                      label=''
                      inline={true}
                      labelFor='connection-name'
                      className='formGroup-inline'
                      contentClassName='formGroupContent'
                    >
                      <Label style={{ display: 'inline', marginRight: 0 }}>
                        Custom Shedule
                        <span className='requiredStar'>*</span>
                      </Label>
                      <InputGroup
                        id='cron-custom'
                        // disabled={cronConfig !== 'custom'}
                        readOnly={cronConfig !== 'custom'}
                        rightElement={cronConfig !== 'custom' ? <Icon icon='lock' size={11} style={{ alignSelf: 'center', margin: '4px 10px -2px 2px' }} /> : null}
                        placeholder='Enter Crontab Syntax'
                        // defaultValue='0 0 * * *'
                        value={cronConfig !== 'custom' ? cronConfig : customCronConfig}
                        onChange={(e) => setCustomCronConfig(e.target.value)}
                        className={`cron-custom-input ${fieldHasError('Cron Custom') ? 'invalid-field' : ''}`}
                        inline={true}
                        style={{ backgroundColor: cronConfig !== 'custom' ? '#ffffdd' : 'inherit' }}
                      />
                    </FormGroup>
                  </div>

                </>
                {/* )} */}

              </FormGroup>
            </div>
            <div>
              <div>
                <Label style={{ display: 'inline', marginRight: 0, marginBottom: 0 }}>
                  Next Run Date
                </Label>
              </div>
              <div style={{ fontSize: '14px', fontWeight: 800 }}>
                {!cron(cronConfig === 'custom' ? customCronConfig : cronConfig).isValid() && <Icon icon='warning-sign' size={14} color={Colors.RED4} style={{ marginRight: '5px' }} />}
                {dayjs(createCron(cronConfig === 'custom' ? customCronConfig : cronConfig).getNextDate().toString()).format('L LTS')} &middot;{' '}
                <span style={{ color: Colors.GRAY3 }}>({dayjs(createCron(cronConfig === 'custom' ? customCronConfig : cronConfig).getNextDate().toString()).fromNow()})</span>
              </div>
            </div>
          </div>

        </div>
        <div className={Classes.DIALOG_FOOTER}>
          <div className={Classes.DIALOG_FOOTER_ACTIONS}>
            <Button disabled={isSaving} onClick={() => setBlueprintDialogIsOpen(false)}>Cancel</Button>
            <Button
              disabled={isSaving}
              icon='cloud-upload'
              intent={Intent.PRIMARY}
              onClick={() => saveBlueprint(draftBlueprint ? draftBlueprint.id : null)}
            >
              {draftBlueprint ? 'Modify Blueprint' : 'Save Blueprint'}
            </Button>
          </div>
        </div>
      </Dialog>

    </>
  )
}

export default Blueprints
