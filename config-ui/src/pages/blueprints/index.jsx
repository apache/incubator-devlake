/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
import React, { Fragment, useEffect, useState, useCallback } from 'react'
import { useHistory } from 'react-router-dom'
// import dayjs from '@/utils/time'
// import cron from 'cron-validate'
import {
  Button, Icon, Intent,
  Popover,
  Position,
  Tooltip,
  Colors,
  Card,
  NonIdealState,
  Elevation,
} from '@blueprintjs/core'
import usePipelineManager from '@/hooks/usePipelineManager'
import useBlueprintManager from '@/hooks/useBlueprintManager'
import useBlueprintValidation from '@/hooks/useBlueprintValidation'
import usePaginator from '@/hooks/usePaginator'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import AppCrumbs from '@/components/Breadcrumbs'
import Content from '@/components/Content'
import AddBlueprintDialog from '@/components/blueprints/AddBlueprintDialog'
import { ReactComponent as NoBlueprintsIcon } from '@/images/no-blueprints.svg'
import BlueprintsGrid from '@/components/blueprints/BlueprintsGrid'

const Blueprints = (props) => {
  const history = useHistory()

  const {
    // eslint-disable-next-line no-unused-vars
    blueprint,
    blueprints,
    name,
    cronConfig,
    customCronConfig,
    cronPresets,
    tasks,
    detectedProviderTasks,
    enable,
    setName: setBlueprintName,
    setCronConfig,
    setCustomCronConfig,
    setTasks: setBlueprintTasks,
    setDetectedProviderTasks,
    setEnable: setEnableBlueprint,
    isFetching: isFetchingBlueprints,
    isSaving,
    isDeleting,
    createCronExpression: createCron,
    getCronSchedule: getSchedule,
    getCronPreset,
    getCronPresetByConfig,
    getNextRunDate,
    activateBlueprint,
    deactivateBlueprint,
    // eslint-disable-next-line no-unused-vars
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
    // eslint-disable-next-line no-unused-vars
    allowedProviders,
    // eslint-disable-next-line no-unused-vars
    detectPipelineProviders
  } = usePipelineManager()

  const {
    pagedData,
    setFilterParams,
    setFilterFunc,
    setData: setPaginatorData,
    renderControlsComponent: renderPagnationControls
  } = usePaginator()

  const [expandDetails, setExpandDetails] = useState(false)
  const [activeBlueprint, setActiveBlueprint] = useState(null)
  const [draftBlueprint, setDraftBlueprint] = useState(null)
  const [blueprintSchedule, setBlueprintSchedule] = useState([])

  const [blueprintDialogIsOpen, setBlueprintDialogIsOpen] = useState(false)
  const [pipelineTemplates, setPipelineTemplates] = useState([])
  const [selectedPipelineTemplate, setSelectedPipelineTemplate] = useState()

  const [activeFilterStatus, setActiveFilterStatus] = useState()

  const [relatedPipelines, setRelatedPipelines] = useState([])

  const {
    validate,
    errors: blueprintValidationErrors,
    // setErrors: setBlueprintErrors,
    isValid: isValidBlueprint,
  } = useBlueprintValidation({
    name,
    cronConfig,
    customCronConfig,
    enable,
    tasks
  })

  const handleBlueprintActivation = useCallback((blueprint) => {
    if (blueprint.enable) {
      deactivateBlueprint(blueprint)
    } else {
      activateBlueprint(blueprint)
    }
  }, [activateBlueprint, deactivateBlueprint])

  const expandBlueprint = useCallback((blueprint) => {
    setExpandDetails(opened => blueprint.id === activeBlueprint?.id && opened ? false : !opened)
    fetchAllPipelines()
    setActiveBlueprint(blueprint)
  }, [fetchAllPipelines, setExpandDetails, setActiveBlueprint, activeBlueprint?.id])

  const configureBlueprint = useCallback((blueprint) => {
    history.push(`/blueprints/detail/${blueprint.id}`)
  }, [history])

  const configureBlueprintSettings = useCallback((blueprint) => {
    history.push(`/blueprints/settings/${blueprint.id}`)
  }, [history])

  const createNewBlueprint = useCallback(() => {
    history.push('/blueprints/create')
  }, [history])

  const isActiveBlueprint = useCallback((bId) => {
    return activeBlueprint?.id === bId
  }, [activeBlueprint?.id])

  const isStandardCronPreset = useCallback((cronConfig) => {
    return cronPresets.some(p => p.cronConfig === cronConfig)
  }, [cronPresets])

  const fieldHasError = useCallback((fieldId) => {
    return blueprintValidationErrors.some(e => e.includes(fieldId))
  }, [blueprintValidationErrors])

  const getFieldError = useCallback((fieldId) => {
    return blueprintValidationErrors.find(e => e.includes(fieldId))
  }, [blueprintValidationErrors])

  const viewPipeline = useCallback((blueprintId) => {
    // history.push(`/pipelines/activity/${runId}`)
    history.push(`/blueprints/detail/${blueprintId}`)
  }, [history])

  useEffect(() => {
    setBlueprintSchedule(activeBlueprint?.id ? getSchedule(activeBlueprint.cronConfig) : [])
    setRelatedPipelines(pipelines.filter(p => p.blueprintId === activeBlueprint?.id))
    console.log('>>> ACTIVE/EXPANDED BLUEPRINT', activeBlueprint)
  }, [activeBlueprint, getSchedule, pipelines])

  useEffect(() => {
    if (draftBlueprint && draftBlueprint.id) {
      console.log('>>> DRAFT = ', draftBlueprint)
      setBlueprintName(draftBlueprint.name)
      setCronConfig(!isStandardCronPreset(draftBlueprint.cronConfig) ? 'custom' : draftBlueprint.cronConfig)
      setCustomCronConfig(draftBlueprint.cronConfig)
      setBlueprintTasks(draftBlueprint.tasks)
      setEnableBlueprint(draftBlueprint.enable)
      setDetectedProviderTasks(draftBlueprint.tasks.flat())
      setBlueprintDialogIsOpen(true)
    }
  }, [
    draftBlueprint,
    setBlueprintName,
    setCronConfig,
    isStandardCronPreset,
    setBlueprintTasks,
    setEnableBlueprint,
    setCustomCronConfig,
    setDetectedProviderTasks
  ])

  useEffect(() => {
    if (saveComplete?.id) {
      setBlueprintDialogIsOpen(false)
      fetchAllBlueprints()
    }
  }, [saveComplete, fetchAllBlueprints])

  useEffect(() => {
    if (deleteComplete.status === 200) {
      fetchAllBlueprints()
    }
  }, [deleteComplete, fetchAllBlueprints])

  useEffect(() => {
    fetchAllBlueprints()
  }, [fetchAllBlueprints])

  // useEffect(() => {
  //   // console.log('>> BLUEPRINT VALIDATION....')
  //   validate()
  // }, [name, cronConfig, customCronConfig, tasks, enable, validate])

  useEffect(() => {
    if (blueprintDialogIsOpen) {
      fetchAllPipelines('TASK_COMPLETED', 100)
    }
  }, [blueprintDialogIsOpen, fetchAllPipelines])

  // useEffect(() => {
  //   setPipelineTemplates(pipelines.slice(0, 100).map(p => ({ ...p, id: p.id, title: p.name, value: p.id })))
  // }, [pipelines, activeBlueprint?.id])

  // useEffect(() => {
  //   if ((!draftBlueprint?.id && selectedPipelineTemplate) || (tasks.length === 0 && selectedPipelineTemplate)) {
  //     console.log('>>>> SELECTED TEMPLATE?', selectedPipelineTemplate.tasks)
  //     setBlueprintTasks(selectedPipelineTemplate.tasks)
  //   }
  // }, [selectedPipelineTemplate, setBlueprintTasks, tasks?.length, draftBlueprint?.id])

  // useEffect(() => {
  //   setSelectedPipelineTemplate(pipelineTemplates.find(pT => pT.tasks?.flat().toString() === tasks.flat().toString()))
  // }, [pipelineTemplates])

  useEffect(() => {
    fetchAllPipelines()
  }, [fetchAllPipelines])

  useEffect(() => {
    setFilterFunc(() => (activeFilterStatus, blueprint) => {
      switch (activeFilterStatus) {
        case 'hourly':
        case 'daily':
        case 'weekly':
        case 'monthly':
          console.log(blueprint.cronConfig === getCronPreset(activeFilterStatus).cronConfig)
          return blueprint.cronConfig === getCronPreset(activeFilterStatus).cronConfig
        case 'manual':
          return blueprint.isManual
        case 'custom':
          return blueprint.cronConfig !== getCronPreset('hourly').cronConfig &&
            blueprint.cronConfig !== getCronPreset('daily').cronConfig &&
            blueprint.cronConfig !== getCronPreset('weekly').cronConfig &&
            blueprint.cronConfig !== getCronPreset('monthly').cronConfig
        default:
          return true
      }
    })
    setFilterParams(activeFilterStatus)
  }, [activeFilterStatus, setFilterParams, getCronPreset, setFilterFunc])

  // useEffect(() => {
  //   if (Array.isArray(tasks)) {
  //     setDetectedProviderTasks([...tasks.flat()])
  //   }
  //   return () => setDetectedProviderTasks([])
  // }, [tasks, setDetectedProviderTasks])

  useEffect(() => {
    console.log('>>>> DETECTED PROVIDERS TASKS....', detectedProviderTasks)
  }, [detectedProviderTasks])

  useEffect(() => {
    setPaginatorData(blueprints)
  }, [blueprints, setPaginatorData])

  return (
    <>
      <div className='container'>
        <Nav />
        <Sidebar />
        <Content>
          <main className='main'>
            {/* <AppCrumbs
              items={[
                { href: '/', icon: false, text: 'Dashboard' },
                { href: '/pipelines', icon: false, text: 'Pipelines' },
                { href: '/blueprints', icon: false, text: 'Pipeline Blueprints', current: true },
              ]}
            /> */}
            <div className='headlineContainer'>
              <div style={{ display: 'flex' }}>
                <div>
                  <h1 style={{ margin: 0 }}>
                    Blueprints
                  </h1>
                </div>
                <div style={{ marginLeft: 'auto' }}>
                  {blueprints.length > 0 && (<Button
                    // disabled={pipelines.length === 0}
                    icon='plus' intent={Intent.PRIMARY}
                    text='New Blueprint'
                    onClick={() => createNewBlueprint()}
                                             />)}
                </div>
              </div>
            </div>
            {(!isFetchingBlueprints) && blueprints.length > 0 && (
              <>
                <BlueprintsGrid
                  blueprints={pagedData}
                  pipelines={relatedPipelines}
                  activeFilterStatus={activeFilterStatus}
                  onFilter={setActiveFilterStatus}
                  activeBlueprint={activeBlueprint}
                  blueprintSchedule={blueprintSchedule}
                  isActiveBlueprint={isActiveBlueprint}
                  expandBlueprint={expandBlueprint}
                  deleteBlueprint={deleteBlueprint}
                  createCron={createCron}
                  getNextRunDate={getNextRunDate}
                  handleBlueprintActivation={handleBlueprintActivation}
                  configureBlueprint={configureBlueprint}
                  configureBlueprintSettings={configureBlueprintSettings}
                  isDeleting={isDeleting}
                  isLoading={isFetchingAllPipelines}
                  expandDetails={expandDetails}
                  cronPresets={cronPresets}
                  onViewPipeline={viewPipeline}
                />
              </>)}

            {!isFetchingBlueprints && blueprints.length === 0 && (
              <Card style={{ marginTop: '36px' }} elevation={Elevation.TWO}>
                <NonIdealState
                  className='blueprints-non-ideal-state'
                  icon={<NoBlueprintsIcon width={120} height={120} style={{ marginBottom: '-30px' }} />}
                  title=''
                  description={(
                    <>
                      <p style={{ color: '#292B3F', fontSize: '15px', padding: '0 14%', textAlign: 'center' }}>A blueprint is a plan that covers all the work, such as selecting and transforming the data you wish to collect, to get your raw data ready for query and metric computation in the dashboards. Try adding your first blueprint!</p>
                      <div style={{
                        display: 'flex',
                        alignSelf: 'center',
                        justifyContent: 'center',
                        marginTop: '5px'
                      }}
                      >
                        {pipelines.length === 0
                          ? (
                            <Button
                                // disabled={pipelines.length === 0}
                              icon='plus'
                              intent={Intent.PRIMARY} text='New Blueprint'
                              style={{ marginRight: '10px' }}
                              onClick={createNewBlueprint}
                            />
                            )
                          : (
                            <Button
                              icon='plus'
                              intent={Intent.PRIMARY} text='New Blueprint'
                              style={{ marginRight: '10px' }}
                              onClick={createNewBlueprint}
                            />
                            )}

                      </div>
                    </>
                  )}
                  // action={createNewBlueprint}
                />
              </Card>
            )}
            <div style={{ alignSelf: 'flex-end', padding: '10px' }}>{renderPagnationControls()}</div>
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
        detectedProviders={detectedProviderTasks}
        getCronPreset={getCronPreset}
        getCronPresetByConfig={getCronPresetByConfig}
      />

    </>
  )
}

export default Blueprints
