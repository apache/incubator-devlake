import React, { useState, useEffect, useCallback } from 'react'
// import { ToastNotification } from '@/components/Toast'
import {
  Providers,
} from '@/data/Providers'

function usePipelineValidation ({
  enabledProviders = [],
  pipelineName,
  projectId,
  boardId,
  owner,
  repositoryName,
  gitExtractorUrl,
  gitExtractorRepoId,
  refDiffRepoId,
  refDiffTasks = [],
  refDiffPairs = [],
  sourceId,
  tasks = [],
  tasksAdvanced = [],
  advancedMode
}) {
  const [errors, setErrors] = useState([])
  const [isValid, setIsValid] = useState(false)
  const [detectedProviders, setDetectedProviders] = useState([])
  const [allowedProviders, setAllowedProviders] = useState([
    Providers.JIRA,
    Providers.GITLAB,
    Providers.JENKINS,
    Providers.GITHUB,
    Providers.REFDIFF,
    Providers.GITEXTRACTOR,
    Providers.FEISHU
  ])

  const clear = () => {
    setErrors([])
  }

  const validateNumericSet = (set = []) => {
    return Array.isArray(set) ? set.every(i => !isNaN(i)) : false
  }

  const validate = useCallback(() => {
    const errs = []
    console.log('>> VALIDATING PIPELINE RUN ', pipelineName)

    if (!pipelineName || pipelineName.length <= 2) {
      errs.push('Name: Enter a valid Pipeline Name')
    }

    if (enabledProviders.includes(Providers.GITLAB) && (!projectId || projectId.length === 0 || projectId.toString() === '')) {
      errs.push('GitLab: Enter one or more valid Project IDs (Numeric)')
    }

    if (enabledProviders.includes(Providers.GITLAB) && !validateNumericSet(projectId)) {
      errs.push('GitLab: One of the entered Project IDs is NOT numeric!')
    }

    if (enabledProviders.includes(Providers.JIRA) && (!sourceId || isNaN(sourceId))) {
      errs.push('JIRA: Select a valid Connection Source ID (Numeric)')
    }

    if (enabledProviders.includes(Providers.JIRA) && (!boardId || boardId.length === 0 || boardId.toString() === '')) {
      errs.push('JIRA: Enter one or more valid Board IDs (Numeric)')
    }

    if (enabledProviders.includes(Providers.JIRA) && !validateNumericSet(boardId)) {
      errs.push('JIRA: One of the entered Board IDs is NOT numeric!')
    }

    if (enabledProviders.includes(Providers.GITHUB) && (!owner || owner <= 2)) {
      errs.push('GitHub: Owner/Developer is required')
    }

    if (enabledProviders.includes(Providers.GITHUB) && (owner.match(/^[a-zA-Z0-9_-]+$/g) === null)) {
      errs.push('GitHub: Owner invalid format')
    }

    if (enabledProviders.includes(Providers.GITHUB) && !repositoryName) {
      errs.push('GitHub: Repository Name is required')
    }

    if (enabledProviders.includes(Providers.GITHUB) && repositoryName.match(/^[a-zA-Z0-9._-]+$/g) === null) {
      errs.push('GitHub: Repository name invalid format')
    }

    if (enabledProviders.includes(Providers.GITEXTRACTOR) && !gitExtractorUrl) {
      errs.push('GitExtractor: Repository Git URL is required')
    }

    if (enabledProviders.includes(Providers.GITEXTRACTOR) &&
        gitExtractorUrl.toLowerCase().match(/^(http:\/\/|https:\/\/|ssh:\/\/|git@)+/g) === null) {
      errs.push('GitExtractor: Repository Git URL must be valid HTTP/S, SSH or Git@ protocol')
    }

    if (enabledProviders.includes(Providers.GITEXTRACTOR) && !gitExtractorUrl.toLowerCase().endsWith('.git')) {
      errs.push('GitExtractor: Invalid Git URL Extension')
    }

    if (enabledProviders.includes(Providers.GITEXTRACTOR) && !gitExtractorRepoId) {
      errs.push('GitExtractor: Repository Column ID Code is required')
    }

    if (enabledProviders.includes(Providers.REFDIFF) && !refDiffRepoId) {
      errs.push('RefDiff: Repository Column ID Code is required')
    }

    if (enabledProviders.includes(Providers.REFDIFF) && refDiffTasks.length === 0) {
      errs.push('RefDiff: Please select at least ONE (1) Plugin Task')
    }

    if (enabledProviders.includes(Providers.REFDIFF) && refDiffPairs.length === 0) {
      errs.push('RefDiff: Please enter at least ONE (1) Tag Ref Pair')
    }

    if (enabledProviders.length === 0) {
      errs.push('Pipeline: Invalid/Empty Configuration')
    }
    setErrors(errs)
  }, [
    enabledProviders,
    pipelineName,
    projectId,
    boardId,
    owner,
    repositoryName,
    gitExtractorUrl,
    gitExtractorRepoId,
    refDiffRepoId,
    refDiffTasks,
    refDiffPairs,
    sourceId
  ])

  const validateAdvanced = useCallback(() => {
    const errs = []
    let parsed = []
    if (advancedMode) {
      console.log('>> VALIDATING ADVANCED PIPELINE RUN ', tasksAdvanced, pipelineName)

      if (Array.isArray(tasksAdvanced)) {
        setDetectedProviders([...new Set(tasksAdvanced?.flat().filter(aT => allowedProviders.includes(aT.Plugin || aT.plugin)).map(p => p.Plugin || p.plugin))])
      }

      if (!pipelineName || pipelineName.length <= 2) {
        errs.push('Name: Enter a valid Pipeline Name')
      }

      try {
        parsed = JSON.parse(JSON.stringify(tasksAdvanced))
      } catch (e) {
        errs.push('Advanced Pipeline: Invalid JSON Configuration')
      }

      if (Array.isArray(tasksAdvanced) && tasksAdvanced?.flat().length === 0) {
        errs.push('Advanced Pipeline: Invalid/Empty Configuration')
      }

      if (!Array.isArray(tasksAdvanced) || !Array.isArray(tasksAdvanced[0])) {
        errs.push('Advanced Pipeline: Invalid Tasks Array Structure!')
      }

      if (Array.isArray(tasksAdvanced) && !tasksAdvanced?.flat().every(aT => allowedProviders.includes(aT.Plugin || aT.plugin))) {
        errs.push('Advanced Pipeline: Unsupported Data Provider Plugin Detected!')
      }

      console.log('>>> Advanced Pipeline Validation Errors? ...', errs)
    }
    setErrors(errs)
  }, [
    advancedMode,
    tasksAdvanced,
    pipelineName,
    allowedProviders
  ])

  useEffect(() => {
    console.log('>>> PIPELINE RUN FORM ERRORS...', errors)
    setIsValid(errors.length === 0)
    if (errors.length > 0) {
      // ToastNotification.clear()
    }
  }, [errors])

  useEffect(() => {
    console.log('>>> DETECTED PLUGIN PROVIDERS...', detectedProviders)
  }, [detectedProviders])

  return {
    errors,
    setErrors,
    isValid,
    validate,
    validateAdvanced,
    clear,
    setAllowedProviders,
    allowedProviders,
    detectedProviders
  }
}

export default usePipelineValidation
