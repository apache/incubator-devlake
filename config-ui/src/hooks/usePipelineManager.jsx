import React, { useState, useEffect, useCallback } from 'react'
import { DEVLAKE_ENDPOINT } from '@/utils/config'
import request from '@/utils/request'
import { NullPipelineRun } from '@/data/NullPipelineRun'
import { ToastNotification } from '@/components/Toast'
// import { integrationsData } from '@/data/integrations'

function usePipelineManager (pipelineName = `COLLECTION ${Date.now()}`, initialTasks = []) {
  // const [integrations, setIntegrations] = useState(integrationsData)
  const [isFetching, setIsFetching] = useState(false)
  const [isRunning, setIsRunning] = useState(false)
  const [isCancelling, setIsCancelling] = useState(false)
  const [errors, setErrors] = useState([])
  const [settings, setSettings] = useState({
    name: pipelineName,
    tasks: [
      [...initialTasks]
    ]
  })

  const [activePipeline, setActivePipeline] = useState(NullPipelineRun)
  const [lastRunId, setLastRunId] = useState(null)
  const [pipelineRun, setPipelineRun] = useState(NullPipelineRun)

  const runPipeline = useCallback(() => {
    console.log('>> RUNNING PIPELINE....')
    try {
      setIsRunning(true)
      setErrors([])
      ToastNotification.clear()
      console.log('>> DISPATCHING PIPELINE REQUEST', settings)
      const run = async () => {
        const p = await request.post(`${DEVLAKE_ENDPOINT}/pipelines`, settings)
        const t = await request.get(`${DEVLAKE_ENDPOINT}/pipelines/${p.data?.ID}/tasks`)
        console.log('>> RAW PIPELINE DATA FROM API...', p.data)
        setPipelineRun({ ...p.data, tasks: [...t.data.tasks] })
        setLastRunId(p.data?.ID)
        ToastNotification.show({ message: `Created New Pipeline - ${pipelineName}.`, intent: 'danger', icon: 'small-tick' })
        setTimeout(() => {
          setIsRunning(false)
        }, 500)
      }
      run()
    } catch (e) {
      setIsRunning(false)
      setErrors([e.message])
      console.log('>> FAILED TO RUN PIPELINE!!', e)
    }
  }, [pipelineName, settings])

  const cancelPipeline = useCallback((pipelineID) => {
    try {
      setIsCancelling(true)
      setErrors([])
      ToastNotification.clear()
      console.log('>> DISPATCHING CANCEL PIPELINE REQUEST, RUN ID =', pipelineID)
      const cancel = async () => {
        const c = await request.delete(`${DEVLAKE_ENDPOINT}/pipelines/${pipelineID}`)
        console.log('>> RAW PIPELINE CANCEL RUN RESPONSE FROM API...', c)
        setPipelineRun(NullPipelineRun)
        ToastNotification.show({ message: `Pipeline RUN ID - ${pipelineID} Cancelled`, intent: 'danger', icon: 'small-tick' })
        setTimeout(() => {
          setIsCancelling(false)
        }, 500)
      }
      cancel()
    } catch (e) {
      setIsCancelling(false)
      setErrors([e.message])
      console.log('>> FAILED TO FETCH CANCEL PIPELINE RUN!!', pipelineID, e)
    }
  }, [])

  const fetchPipeline = useCallback((pipelineID, refresh = false) => {
    if (!pipelineID) {
      console.log('>> !ABORT! Pipeline ID Missing! Aborting Fetch...')
      // return ToastNotification.show({ message: 'Pipeline ID Missing! Aborting Fetch...', intent: 'danger', icon: 'warning-sign' })
    }
    try {
      setIsFetching(true)
      setErrors([])
      ToastNotification.clear()
      console.log('>> FETCHING PIPELINE RUN DETAILS...')
      const fetch = async () => {
        const p = await request.get(`${DEVLAKE_ENDPOINT}/pipelines/${pipelineID}`)
        const t = await request.get(`${DEVLAKE_ENDPOINT}/pipelines/${pipelineID}/tasks`)
        console.log('>> RAW PIPELINE RUN DATA FROM API...', p.data)
        console.log('>> RAW PIPELINE TASKS DATA FROM API...', t.data)
        setActivePipeline({
          ...p.data,
          tasks: [...t.data.tasks]
        })
        setPipelineRun((pR) => refresh ? { ...p.data, tasks: [...t.data.tasks] } : pR)
        setLastRunId((lrId) => refresh ? p.data?.ID : lrId)
        // ToastNotification.show({ message: `Fetched Pipeline ID - ${p.data?.ID}.`, intent: 'danger', icon: 'small-tick' })
        setTimeout(() => {
          setIsFetching(false)
        }, 500)
      }
      fetch()
    } catch (e) {
      setIsFetching(false)
      setErrors([e.message])
      setActivePipeline(NullPipelineRun)
      console.log('>> FAILED TO FETCH PIPELINE RUN!!', e)
    }
  }, [])

  const fetchPipelineTasks = useCallback(() => {
    try {
      setIsFetching(true)
      setErrors([])
      ToastNotification.clear()
    } catch (e) {
      setIsFetching(false)
      setErrors([e.message])
      console.log('>> FAILED TO FETCH PIPELINE RUN TASKS!!', e)
    }
  }, [])

  // const fetchAllRuns = () => {

  // }

  // const fetchRun = () => {

  // }

  useEffect(() => {
    // setIntegrations(integrationsData)
  }, [])

  useEffect(() => {
    console.log('>> PIPELINE MANAGER - RECEIVED RUN/TASK SETTINGS', settings)
  }, [settings])

  useEffect(() => {

  }, [pipelineName, initialTasks])

  return {
    errors,
    isRunning,
    isFetching,
    settings,
    setSettings,
    pipelineRun,
    activePipeline,
    lastRunId,
    runPipeline,
    cancelPipeline,
    fetchPipeline,
    fetchPipelineTasks
  }
}

export default usePipelineManager
