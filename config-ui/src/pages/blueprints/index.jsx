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
    getCronPreset,
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
            {(!isFetchingBlueprints) && blueprints.length > 0 && (
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
                  cronPresets={cronPresets}
                />
              </>)}

            {!isFetchingBlueprints && blueprints.length === 0 && (
              <div style={{ marginTop: '36px' }}>
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
