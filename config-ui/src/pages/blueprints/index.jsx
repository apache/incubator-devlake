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
  NonIdealState,
  Divider,
  H5,
  Switch,
  Pre,
  Tag
} from '@blueprintjs/core'
import { parseCronExpression } from 'cron-schedule'
import usePipelineManager from '@/hooks/usePipelineManager'
import useBlueprintManager from '@/hooks/useBlueprintManager'
import useBlueprintValidation from '@/hooks/useBlueprintValidation'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import AppCrumbs from '@/components/Breadcrumbs'
import Content from '@/components/Content'
// import ContentLoader from '@/components/loaders/ContentLoader'
import AddBlueprintDialog from '@/components/blueprints/AddBlueprintDialog'
import { ReactComponent as HelpIcon } from '@/images/help.svg'
import ManageBlueprintsIcon from '@/images/blueprints.png'
import EventIcon from '@/images/calendar-3.png'
import EventOffIcon from '@/images/calendar-4.png'
import { NullBlueprint } from '@/data/NullBlueprint'
import InputValidationError from '@/components/validation/InputValidationError'
import DeletePopover from '@/components/blueprints/DeletePopover'
import BlueprintsGrid from '../../components/blueprints/BlueprintsGrid'

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

  const {
    pipelines,
    isFetchingAll: isFetchingAllPipelines,
    fetchAllPipelines,
    allowedProviders,
    detectPipelineProviders
  } = usePipelineManager()

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
  const [pipelineTemplates, setPipelineTemplates] = useState([])
  const [selectedPipelineTemplate, setSelectedPipelineTemplate] = useState()
  // const [blueprintErrors, setBlueprintErrors] = useState([])

  const {
    validate,
    errors: blueprintValidationErrors,
    setErrors: setPipelineErrors,
    isValid: isValidBlueprint,
  } = useBlueprintValidation({
    name,
    cronConfig,
    customCronConfig,
    enable,
    tasks
  })

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
    setEnableBlueprint(true)
    setBlueprintTasks([])
    setSelectedPipelineTemplate(null)
    setBlueprintDialogIsOpen(true)
  }

  const isActiveBlueprint = (bId) => {
    return activeBlueprint?.id === bId
  }

  const isStandardCronPreset = (cronConfig) => {
    return cronPresets.some(p => p.cronConfig === cronConfig)
  }

  const fieldHasError = (fieldId) => {
    return blueprintValidationErrors.some(e => e.includes(fieldId))
  }

  const getFieldError = (fieldId) => {
    return blueprintValidationErrors.find(e => e.includes(fieldId))
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
      setBlueprintTasks(draftBlueprint.tasks)
      setEnableBlueprint(draftBlueprint.enable)
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

  useEffect(() => {
    // console.log('>> BLUEPRINT VALIDATION....')
    validate()
  }, [name, cronConfig, customCronConfig, tasks, enable, validate])

  useEffect(() => {
    if (blueprintDialogIsOpen) {
      fetchAllPipelines('TASK_COMPLETED', 100)
    }
  }, [blueprintDialogIsOpen, fetchAllPipelines])

  useEffect(() => {
    setPipelineTemplates(pipelines.slice(0, 100).map(p => ({ ...p, id: p.id, title: p.name, value: p.id })))
  }, [pipelines])

  useEffect(() => {
    if (selectedPipelineTemplate) {
      setBlueprintTasks(selectedPipelineTemplate.tasks)
    }
  }, [selectedPipelineTemplate, setBlueprintTasks])

  useEffect(() => {
    setSelectedPipelineTemplate(pipelineTemplates.find(pT => pT.tasks.flat().toString() === tasks.flat().toString()))
  }, [pipelineTemplates])

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
                { href: '/blueprints', icon: false, text: 'Pipeline Blueprints', current: true },
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
                          <p>Need Help? &mdash; Automate pipelines by creating a Blueprint.
                            Schedule data collection with Crontab and save hours of time.
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
                <BlueprintsGrid
                  blueprints={blueprints}
                  activeBlueprint={activeBlueprint}
                  blueprintSchedule={blueprintSchedule}
                  isActiveBlueprint={isActiveBlueprint}
                  expandBlueprint={expandBlueprint}
                  deleteBlueprint={deleteBlueprint}
                  createCron={createCron}
                  handleBlueprintActivation={handleBlueprintActivation}
                  configureBlueprint={configureBlueprint}
                  isDeleting={isDeleting}
                  expandDetails={expandDetails}
                />
                {/* <div style={{ display: 'flex', marginTop: '30px', minHeight: '36px', width: '100%', justifyContent: 'flex-start' }}>
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
                </div> */}
              </>)}

            {!isFetchingBlueprints && blueprints.length === 0 && (
              <div>
                <NonIdealState
                  icon='grid'
                  title='No Defined Blueprints'
                  description={(
                    <>
                      Please create a new blueprint to get started. Need Help? Visit the DevLake Wiki on <strong>GitHub</strong>.{' '}
                      <div style={{
                        display: 'flex',
                        alignSelf: 'center',
                        justifyContent: 'center',
                        marginTop: '5px'
                      }}
                      >
                        <Button
                          intent={Intent.PRIMARY} text='Create Blueprint' small
                          style={{ marginRight: '10px' }}
                          onClick={createNewBlueprint}
                        />
                      </div>
                    </>
                  )}
                  // action={createNewBlueprint}
                />
              </div>
            )}
          </main>
        </Content>
      </div>

      <AddBlueprintDialog
        isLoading={isFetchingAllPipelines}
        isOpen={blueprintDialogIsOpen}
        setIsOpen={setBlueprintDialogIsOpen}
        name={name}
        cronConfig={cronConfig}
        customCronConfig={customCronConfig}
        enable={enable}
        tasks={tasks}
        draftBlueprint={draftBlueprint}
        setDraftBlueprint={setDraftBlueprint}
        setBlueprintName={setBlueprintName}
        setCronConfig={setCronConfig}
        setCustomCronConfig={setCustomCronConfig}
        setEnableBlueprint={setEnableBlueprint}
        setBlueprintTasks={setBlueprintTasks}
        createCron={createCron}
        saveBlueprint={saveBlueprint}
        isSaving={isSaving}
        isValidBlueprint={isValidBlueprint}
        fieldHasError={fieldHasError}
        getFieldError={getFieldError}
        pipelines={pipelineTemplates}
        selectedPipelineTemplate={selectedPipelineTemplate}
        setSelectedPipelineTemplate={setSelectedPipelineTemplate}
        detectedProviders={detectPipelineProviders(tasks, allowedProviders)}
      />

    </>
  )
}

export default Blueprints
