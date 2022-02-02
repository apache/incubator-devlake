import React, { useEffect, useState, useCallback } from 'react'
import request from '@/utils/request'
import { ToastNotification } from '@/components/Toast'

const useJIRA = ({ apiProxyPath, issuesEndpoint, fieldsEndpoint }) => {
  const [isFetching, setIsFetching] = useState(false)
  const [issueTypes, setIssueTypes] = useState([])
  const [fields, setFields] = useState([])
  const [error, setError] = useState()

  const fetchIssueTypes = useCallback(() => {
    try {
      const fetchIssueTypes = async () => {
        const issues = await request.get(issuesEndpoint)
        console.log('>>> JIRA API PROXY: Issues Response...', issues)
        setIssueTypes(issues.data)
        setIsFetching(false)
        setError(null)
      }
      fetchIssueTypes()
    } catch (e) {
      setIsFetching(false)
      setError(e)
      ToastNotification.show({ message: e.message, intent: 'danger', icon: 'error' })
    }
  }, [issuesEndpoint])

  const fetchFields = useCallback(() => {
    try {
      const fetchIssueFields = async () => {
        const fields = await request.get(fieldsEndpoint)
        console.log('>>> JIRA API PROXY: Fields Response...', fields)
        setFields(fields.data)
        setIsFetching(false)
        setError(null)
      }
      fetchIssueFields()
    } catch (e) {
      setIsFetching(false)
      setError(e)
      ToastNotification.show({ message: e.message, intent: 'danger', icon: 'error' })
    }
  }, [fieldsEndpoint])

  return {
    isFetching,
    fetchFields,
    fetchIssueTypes,
    issueTypes,
    fields,
    error
  }
}

export default useJIRA
