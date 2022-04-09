import React, { useState, useEffect, useCallback } from 'react'
import { DEVLAKE_ENDPOINT } from '@/utils/config'
import request from '@/utils/request'
import { NullPipelineRun } from '@/data/NullPipelineRun'
import { ToastNotification } from '@/components/Toast'
import { parseCronExpression } from 'cron-schedule'
import cron from 'cron-validate'

function useBlueprintManager (blueprintName = `BLUEPRINT WEEKLY ${Date.now()}`, initialConfiguration = {}) {
  const [isFetching, setIsFetching] = useState(false)
  const [blueprints, setBlueprints] = useState([])
  const [blueprintCount, setBlueprintCount] = useState(0)
  const [errors, setErrors] = useState([])

  const [name, setName] = useState('DAILY BLUEPRINT')
  const [cronConfig, setCronConfig] = useState('0 0 * * *')
  const [tasks, setTasks] = useState([])

  const [cronPresets, setCronPresets] = useState([
    // eslint-disable-next-line max-len
    { id: 0, name: 'hourly', label: 'Hourly', cronConfig: '59 * * * 1-5', description: 'At minute 59 on every day-of-week from Monday through Friday.' },
    { id: 1, name: 'daily', label: 'Daily', cronConfig: '0 0 * * *', description: 'At 00:00 (Midnight) on every day-of-week from Monday through Friday.' },
    { id: 2, name: 'weekly', label: 'Weekly', cronConfig: '0 0 * * 1', description: 'At 00:00 (Midnight) on Monday.' },
    { id: 3, name: 'monthly', label: 'Monthly', cronConfig: '0 0 1 * *', description: 'At 00:00 (Midnight) on day-of-month 1.' },
  ])

  const createCronExpression = (cronExpression = '0 0 * * *') => {
    let newCron = parseCronExpression('0 0 * * *')
    try {
      if (cron(cronExpression).isValid()) {
        newCron = parseCronExpression(cronExpression)
      }
    } catch (e) {
      console.log('>> INVALID CRON EXPRESSION INPUT!', e)
    }
    return newCron
  }

  const getCronSchedule = useCallback((cronExpression, events = 5) => {
    let schedule = []
    try {
      const cS = cron(cronExpression).isValid() ? parseCronExpression(cronExpression) : parseCronExpression('0 0 * * *')
      schedule = cS ? cS.getNextDates(events) : []
    } catch (e) {
      console.log('>> INVALID CRON SCHEDULE!', e)
    }
    return schedule
  }, [])

  const activateBlueprint = useCallback((blueprintId) => {

  }, [])

  const deactivateBlueprint = useCallback((blueprintId) => {

  }, [])

  return {
    blueprints,
    blueprintCount,
    name,
    cronConfig,
    cronPresets,
    tasks,
    createCronExpression,
    getCronSchedule,
    activateBlueprint,
    deactivateBlueprint,
    setName,
    setCronConfig,
    setTasks,
    isFetching,
    errors
  }
}

export default useBlueprintManager
