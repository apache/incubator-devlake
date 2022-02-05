import React, { useEffect, useState, useCallback } from 'react'
import request from '@/utils/request'
import { ToastNotification } from '@/components/Toast'

const useJIRA = ({ apiProxyPath, issuesEndpoint, fieldsEndpoint }) => {
  const [isFetching, setIsFetching] = useState(false)
  const [issueTypes, setIssueTypes] = useState([])
  const [fields, setFields] = useState([])
  const [issueTypesResponse, setIssueTypesResponse] = useState([])
  const [fieldsResponse, setFieldsResponse] = useState([])
  const [error, setError] = useState()

  const fetchIssueTypes = useCallback(() => {
    try {
      if (apiProxyPath.includes('null')) {
        throw new Error('Connection Source ID is Null')
      }
      const fetchIssueTypes = async () => {
        const issues = await
        request
          .get(issuesEndpoint)
          .catch(e => setError(e))
        console.log('>>> JIRA API PROXY: Issues Response...', issues)
        setIssueTypesResponse(issues ? issues.data : [])
        setTimeout(() => {
          setIsFetching(false)
        }, 1000)
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
      if (apiProxyPath.includes('null')) {
        throw new Error('Connection Source ID is Null')
      }
      const fetchIssueFields = async () => {
        const fields = await
        request
          .get(fieldsEndpoint)
          .catch(e => setError(e))
        console.log('>>> JIRA API PROXY: Fields Response...', fields)
        setFieldsResponse(fields ? fields.data : [])
        setTimeout(() => {
          setIsFetching(false)
        }, 1000)
        setError(null)
      }
      fetchIssueFields()
    } catch (e) {
      setIsFetching(false)
      setError(e)
      ToastNotification.show({ message: e.message, intent: 'danger', icon: 'error' })
    }
  }, [fieldsEndpoint])

  const createListData = (data = [], titleProperty = 'name', valueProperty = 'name') => {
    return data.map((d, dIdx) => {
      return {
        ...d,
        id: dIdx,
        key: d.key ? d.key : dIdx,
        title: d[titleProperty],
        value: d[valueProperty],
        type: d.schema?.type || 'string'
      }
    })
  }

  useEffect(() => {
    setIssueTypes(createListData(issueTypesResponse))
  }, [issueTypesResponse])

  useEffect(() => {
    setFields(createListData(fieldsResponse, 'name', 'key'))
  }, [fieldsResponse])

  useEffect(() => {
    console.log('>>> JIRA API PROXY: FIELD SELECTOR LIST DATA', fields)
  }, [fields])

  useEffect(() => {
    if (error) {
      console.log('>>> JIRA PROXY API ERROR!', error)
    }
  }, [error])

  return {
    isFetching,
    fetchFields,
    fetchIssueTypes,
    createListData,
    issueTypesResponse,
    fieldsResponse,
    issueTypes,
    fields,
    error
  }
}

export default useJIRA
