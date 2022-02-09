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
  sourceId,
  tasks
}) {
  const [errors, setErrors] = useState([])
  const [isValid, setIsValid] = useState(false)

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

    if (enabledProviders.length === 0) {
      errs.push('Pipeline: Invalid/Empty Configuration')
    }

    // try {
    //   JSON.parse(tasks)
    // } catch (e) {
    //   errs.push(e.message)
    // }

    setErrors(errs)
  }, [
    enabledProviders,
    pipelineName,
    projectId,
    boardId,
    owner,
    repositoryName,
    sourceId,
    tasks
  ])

  useEffect(() => {
    console.log('>>> PIPELINE RUN FORM ERRORS...', errors)
    setIsValid(errors.length === 0)
    if (errors.length > 0) {
      // ToastNotification.clear()
    }
  }, [errors])

  return {
    errors,
    isValid,
    validate,
    clear
  }
}

export default usePipelineValidation
