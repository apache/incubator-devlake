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
import { useCallback, useEffect, useState } from 'react'
import { DEVLAKE_ENDPOINT } from '@/utils/config'
import request from '@/utils/request'
import { ToastNotification } from '@/components/Toast'
import { NullBlueprint, BlueprintMode } from '@/data/NullBlueprint'
import cron from 'cron-validate'
import parser from 'cron-parser'
import { Intent } from '@blueprintjs/core'

function useBlueprintManager(
  blueprintName = `BLUEPRINT WEEKLY ${Date.now()}`,
  initialConfiguration = {}
) {
  const [isFetching, setIsFetching] = useState(false)
  const [isSaving, setIsSaving] = useState(false)
  const [isDeleting, setIsDeleting] = useState(false)
  const [blueprints, setBlueprints] = useState([])
  const [blueprintCount, setBlueprintCount] = useState(0)
  const [blueprint, setBlueprint] = useState(NullBlueprint)
  const [errors, setErrors] = useState([])

  const [name, setName] = useState('MY BLUEPRINT')
  const [cronConfig, setCronConfig] = useState('0 0 * * *')
  const [customCronConfig, setCustomCronConfig] = useState('0 0 * * *')
  const [interval, setInterval] = useState('daily')
  const [tasks, setTasks] = useState([])
  const [settings, setSettings] = useState({
    version: '1.0.0',
    connections: []
  })
  const [mode, setMode] = useState(BlueprintMode.NORMAL)
  const [enable, setEnable] = useState(true)
  const [detectedProviderTasks, setDetectedProviderTasks] = useState([])
  const [isManual, setIsManual] = useState(false)
  const [rawConfiguration, setRawConfiguration] = useState(
    JSON.stringify([tasks], null, '  ')
  )

  const [cronPresets, setCronPresets] = useState([
    // eslint-disable-next-line max-len
    {
      id: 0,
      name: 'hourly',
      label: 'Hourly',
      cronConfig: '59 * * * *',
      description:
        'At minute 59 on every day-of-week from Monday through Sunday.'
    },
    // eslint-disable-next-line max-len
    {
      id: 1,
      name: 'daily',
      label: 'Daily',
      cronConfig: '0 0 * * *',
      description:
        'At 00:00 (Midnight) on every day-of-week from Monday through Sunday.'
    },
    {
      id: 2,
      name: 'weekly',
      label: 'Weekly',
      cronConfig: '0 0 * * 1',
      description: 'At 00:00 (Midnight) on Monday.'
    },
    {
      id: 3,
      name: 'monthly',
      label: 'Monthly',
      cronConfig: '0 0 1 * *',
      description: 'At 00:00 (Midnight) on day-of-month 1.'
    }
  ])

  const [saveComplete, setSaveComplete] = useState(false)
  const [deleteComplete, setDeleteComplete] = useState(false)

  const parseCronExpression = useCallback(
    (expression, utc = true, additionalOptions = {}) => {
      if (expression.toLowerCase() === 'manual') {
        return {
          next: () => 'manual',
          prev: () => 'manual'
        }
      }
      return parser.parseExpression(expression, { utc, ...additionalOptions })
    },
    []
  )

  const detectCronInterval = useCallback(
    (cronConfig) => {
      if (cronConfig === 'manual') {
        return 'Manual'
      }
      return cronPresets.find((p) => p.cronConfig === cronConfig)
        ? cronPresets.find((p) => p.cronConfig === cronConfig).label
        : 'Custom'
    },
    [cronPresets]
  )

  const fetchAllBlueprints = useCallback(
    async (notify = false) => {
      try {
        setIsFetching(true)
        setErrors([])
        ToastNotification.clear()
        console.log('>> FETCHING ALL BLUEPRINTS')
        const b = await request.get(`${DEVLAKE_ENDPOINT}/blueprints`)
        console.log('>> RAW ALL BLUEPRINTS DATA FROM API...', b.data)
        const blueprints = []
          .concat(Array.isArray(b.data.blueprints) ? b.data.blueprints : [])
          .map((blueprint, idx) => {
            return {
              ...blueprint,
              id: blueprint.id,
              enable: blueprint.enable,
              status: 0,
              nextRunAt: null, // @todo: calculate next run date
              interval: detectCronInterval(blueprint.cronConfig)
            }
          })
        if (notify) {
          ToastNotification.show({
            message: 'Loaded all blueprints.',
            intent: 'success',
            icon: 'small-tick'
          })
        }
        setBlueprints(blueprints)
        setBlueprintCount(b.data.count || blueprints.length)
        setIsFetching(false)
      } catch (e) {
        console.log('>> FAILED TO FETCH ALL CONNECTIONS', e)
        ToastNotification.show({
          message: `Failed to Load Blueprints - ${e.message}`,
          intent: 'danger',
          icon: 'error'
        })
        setIsFetching(false)
        setBlueprints([])
        setBlueprintCount(0)
        setErrors([e.message])
      }
    },
    [detectCronInterval]
  )

  const fetchBlueprint = useCallback(
    (blueprintId = null) => {
      console.log('>> FETCHING BLUEPRINT....')
      try {
        setIsFetching(true)
        setErrors([])
        ToastNotification.clear()
        console.log('>> FETCHING BLUEPRINT #', blueprintId)
        const fetch = async () => {
          const b = await request.get(
            `${DEVLAKE_ENDPOINT}/blueprints/${blueprintId}`
          )
          const blueprintData = b.data
          console.log('>> RAW BLUEPRINT DATA FROM API...', b)
          setBlueprint((B) =>
            b?.status === 200
              ? {
                  ...B,
                  ...blueprintData,
                  id: blueprintData.id,
                  enable: blueprintData.enable,
                  status: 0,
                  nextRunAt: null, // @todo: calculate next run date
                  interval: detectCronInterval(blueprintData.cronConfig)
                }
              : NullBlueprint
          )
          setErrors(b.status !== 200 ? [b.data] : [])
          setTimeout(() => {
            setIsFetching(false)
          }, 500)
        }
        fetch()
      } catch (e) {
        setIsFetching(false)
        setBlueprint(null)
        setErrors([e.message])
        ToastNotification.show({
          message: `${e}`,
          intent: 'danger',
          icon: 'error'
        })
        console.log('>> FAILED TO FETCH BLUEPRINT', e)
      }
    },
    [detectCronInterval]
  )

  const saveBlueprint = useCallback(
    (blueprintId = null) => {
      console.log('>> SAVING BLUEPRINT....')
      try {
        setIsSaving(true)
        setErrors([])
        ToastNotification.clear()
        const detectCronConfig = () => {
          if (cronConfig === 'custom') {
            return customCronConfig
            // For "Manual" frequency, we'll save cronConfig as daily to comply with BE API expectation of valid cron
            // Once user re-enables an automated frequency, this will get overwritten
          } else if (cronConfig === 'manual') {
            return '0 0 * * *'
          } else {
            return cronConfig
          }
        }
        const blueprintPayload = {
          name,
          cronConfig: detectCronConfig(),
          // @todo: refactor tasks ===> plan at higher levels
          plan: tasks,
          settings,
          enable: enable,
          mode,
          isManual
        }
        console.log('>> DISPATCHING BLUEPRINT SAVE REQUEST', blueprintPayload)
        const run = async () => {
          // eslint-disable-next-line max-len
          // const b = await request.post(`${DEVLAKE_ENDPOINT}/blueprints`, blueprintPayload)
          const b = blueprintId
            ? await request.patch(
                `${DEVLAKE_ENDPOINT}/blueprints/${blueprintId}`,
                blueprintPayload
              )
            : await request.post(
                `${DEVLAKE_ENDPOINT}/blueprints`,
                blueprintPayload
              )
          console.log('>> RAW BLUEPRINT DATA FROM API...', b.data)
          const blueprintObject = {
            ...b.data,
            id: b.data?.id,
            enable: b.data?.enable,
            status: 0,
            nextRunAt: null,
            interval: detectCronInterval(b.data?.cronConfig)
          }
          setBlueprint(blueprintObject)
          setSaveComplete(blueprintObject)
          if ([200, 201].includes(b.status)) {
            ToastNotification.show({
              message: `${
                blueprintId ? 'Updated' : 'Created'
              } Blueprint - ${name}.`,
              intent: Intent.SUCCESS,
              icon: 'small-tick'
            })
          } else {
            ToastNotification.show({
              message: `Blueprint Failure - ${b.message}`,
              intent: Intent.NONE,
              icon: 'error'
            })
          }
          setTimeout(() => {
            setIsSaving(false)
          }, 500)
        }
        run()
      } catch (e) {
        setIsSaving(false)
        setErrors([e.message])
        setSaveComplete(false)
        console.log('>> FAILED TO SAVE BLUEPRINT!!', e)
      }
    },
    [
      name,
      mode,
      settings,
      cronConfig,
      customCronConfig,
      tasks,
      enable,
      isManual,
      detectCronInterval
    ]
  )

  const deleteBlueprint = useCallback(async (blueprint) => {
    try {
      setIsDeleting(true)
      setErrors([])
      console.log('>> TRYING TO DELETE BLUEPRINT...', blueprint)
      const d = await request.delete(
        `${DEVLAKE_ENDPOINT}/blueprints/${blueprint.id}`
      )
      console.log('>> BLUEPRINT DELETED...', d)
      setIsDeleting(false)
      setDeleteComplete({ status: d.status, data: d.data || null })
      setSaveComplete(null)
    } catch (e) {
      setIsDeleting(false)
      setDeleteComplete(false)
      setErrors([e.message])
      console.log('>> FAILED TO DELETE BLUEPRINT', e)
    }
  }, [])

  const createCronExpression = (cronExpression = '0 0 * * *') => {
    let newCron = parseCronExpression('0 0 * * *', false)
    try {
      newCron = parseCronExpression(cronExpression, false)
    } catch (e) {
      console.log('>> INVALID CRON EXPRESSION INPUT!', e)
    }
    return newCron
  }

  const patchBlueprint = useCallback(
    async (blueprint, settings = {}, callback = () => {}) => {
      try {
        setIsSaving(true)
        setErrors([])
        console.log('>> TRYING TO PATCH BLUEPRINT...', blueprint)
        const p = await request.patch(
          `${DEVLAKE_ENDPOINT}/blueprints/${blueprint.id}`,
          settings
        )
        console.log('>> BLUEPRINT PATCHED...', p)
        setBlueprint((b) => ({
          ...b,
          ...p?.data,
          interval: detectCronInterval(p?.data?.cronConfig)
        }))
        setIsSaving(false)
        setSaveComplete({ status: p.status, data: p.data || null })
        callback(p)
      } catch (e) {
        setIsSaving(false)
        setSaveComplete(null)
        setErrors([e.message])
        callback(e)
        console.log('>> FAILED TO PATCH BLUEPRINT', e)
      }
    },
    [detectCronInterval]
  )

  const getCronSchedule = useCallback(
    (cronExpression, events = 5) => {
      let schedule = []
      try {
        const cS = cron(cronExpression).isValid()
          ? parseCronExpression(cronExpression, false)
          : parseCronExpression('0 0 * * *', false)
        schedule = cS
          ? new Array(events)
              .fill(cS, 0, events)
              .map((interval) => interval.next().toString())
          : []
        console.log('>>> NEW CRON SCHEDULE= ', schedule)
      } catch (e) {
        console.log('>> INVALID CRON SCHEDULE!', e)
      }
      return schedule
    },
    [parseCronExpression]
  )

  const getNextRunDate = useCallback(
    (cronExpression) => {
      return (
        cronExpression && parseCronExpression(cronExpression, false).next().toString()
      )
    },
    [parseCronExpression]
  )

  const getCronPreset = useCallback(
    (presetName) => {
      return cronPresets.find((p) => p.name === presetName)
    },
    [cronPresets]
  )

  const getCronPresetByConfig = useCallback(
    (cronConfig) => {
      return cronPresets.find((p) => p.cronConfig === cronConfig)
    },
    [cronPresets]
  )

  const activateBlueprint = useCallback(
    (blueprint) => {
      console.log('>> ACTIVATING BLUEPRINT....')
      try {
        setIsSaving(true)
        setErrors([])
        ToastNotification.clear()
        const blueprintPayload = {
          name: blueprint.name,
          enable: true
        }
        console.log(
          '>> DISPATCHING BLUEPRINT ACTIVATION REQUEST',
          blueprintPayload
        )
        const run = async () => {
          // eslint-disable-next-line max-len
          const activateB = await request.patch(
            `${DEVLAKE_ENDPOINT}/blueprints/${blueprint.id}`,
            blueprintPayload
          )
          console.log('>> RAW BLUEPRINT DATA FROM API...', activateB.data)
          // eslint-disable-next-line no-unused-vars
          const updatedBlueprint = activateB.data
          if (activateB.status === 200) {
            setBlueprint((b) => ({ ...b, ...updatedBlueprint }))
            ToastNotification.show({
              message: `Activated Blueprint - ${blueprint.name}.`,
              intent: Intent.SUCCESS,
              icon: 'small-tick'
            })
          } else {
            ToastNotification.show({
              message: `Activation Failed ${activateB?.message || ''}`,
              intent: 'danger',
              icon: 'error'
            })
          }
          setTimeout(() => {
            setIsSaving(false)
            fetchAllBlueprints()
          }, 500)
        }
        run()
      } catch (e) {
        setIsSaving(false)
        setErrors([e.message])
        console.log('>> FAILED TO ACTIVATE BLUEPRINT!!', e)
      }
    },
    [fetchAllBlueprints]
  )

  const deactivateBlueprint = useCallback(
    (blueprint) => {
      console.log('>> DEACTIVATING BLUEPRINT....')
      try {
        setIsSaving(true)
        setErrors([])
        ToastNotification.clear()
        const blueprintPayload = {
          enable: false
        }
        console.log(
          '>> DISPATCHING BLUEPRINT ACTIVATION REQUEST',
          blueprintPayload
        )
        const run = async () => {
          // eslint-disable-next-line max-len
          const deactivateB = await request.patch(
            `${DEVLAKE_ENDPOINT}/blueprints/${blueprint.id}`,
            blueprintPayload
          )
          console.log('>> RAW BLUEPRINT DATA FROM API...', deactivateB.data)
          // eslint-disable-next-line no-unused-vars
          const updatedBlueprint = deactivateB.data
          if (deactivateB.status === 200) {
            setBlueprint((b) => ({ ...b, ...updatedBlueprint }))
            ToastNotification.show({
              message: `Deactivated Blueprint - ${blueprint.name}.`,
              intent: Intent.SUCCESS,
              icon: 'small-tick'
            })
          } else {
            ToastNotification.show({
              message: `Deactivation Failed ${deactivateB?.message || ''}`,
              intent: 'danger',
              icon: 'error'
            })
          }
          setTimeout(() => {
            setIsSaving(false)
            fetchAllBlueprints()
          }, 500)
        }
        run()
      } catch (e) {
        setIsSaving(false)
        setErrors([e.message])
        console.log('>> FAILED TO DEACTIVATE BLUEPRINT!!', e)
      }
    },
    [fetchAllBlueprints]
  )

  return {
    blueprint,
    blueprints,
    blueprintCount,
    name,
    cronConfig,
    customCronConfig,
    cronPresets,
    tasks,
    settings,
    detectedProviderTasks,
    enable,
    mode,
    interval,
    rawConfiguration,
    saveBlueprint,
    deleteBlueprint,
    fetchAllBlueprints,
    fetchBlueprint,
    patchBlueprint,
    saveComplete,
    deleteComplete,
    createCronExpression,
    getCronSchedule,
    getNextRunDate,
    getCronPreset,
    getCronPresetByConfig,
    detectCronInterval,
    activateBlueprint,
    deactivateBlueprint,
    setName,
    setCronConfig,
    setCustomCronConfig,
    setCronPresets,
    setTasks,
    setSettings,
    setEnable,
    setMode,
    setDetectedProviderTasks,
    setIsManual,
    setRawConfiguration,
    setInterval,
    isFetching,
    isSaving,
    isDeleting,
    isManual,
    errors
  }
}

export default useBlueprintManager
