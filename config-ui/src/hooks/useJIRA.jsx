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
import { useEffect, useState, useCallback } from 'react'
import request from '@/utils/request'
import { ToastNotification } from '@/components/Toast'
import { Providers } from '@/data/Providers'

const useJIRA = ({ apiProxyPath, issuesEndpoint, fieldsEndpoint, boardsEndpoint }, activeConnection = null) => {
  const [isFetching, setIsFetching] = useState(false)
  const [issueTypes, setIssueTypes] = useState([])
  const [fields, setFields] = useState([])
  const [boards, setBoards] = useState([])
  const [allResources, setAllResources] = useState({ boards, fields, issueTypes })
  const [issueTypesResponse, setIssueTypesResponse] = useState([])
  const [fieldsResponse, setFieldsResponse] = useState([])
  const [boardsResponse, setBoardsResponse] = useState([])
  const [error, setError] = useState()

  const fetchIssueTypes = useCallback(() => {
    // if (activeConnection?.plugin !== Providers.JIRA) {
    //   return
    // }
    try {
      if (apiProxyPath.includes('null')) {
        throw new Error('Connection ID is Null')
      }
      setError(null)
      setIsFetching(true)
      const fetchIssueTypes = async () => {
        const issues = await
        request
          .get(activeConnection?.connectionId ? issuesEndpoint.replace('[:connectionId:]', activeConnection?.connectionId) : issuesEndpoint)
          .catch(e => setError(e))
        console.log('>>> JIRA API PROXY: Issues Response...', issues)
        setIssueTypesResponse(issues && Array.isArray(issues.data) ? issues.data : [])
        // setTimeout(() => {
        setIsFetching(false)
        // }, 1000)
      }
      fetchIssueTypes()
    } catch (e) {
      setIsFetching(false)
      setError(e)
      ToastNotification.show({ message: e.message, intent: 'danger', icon: 'error' })
    }
  }, [issuesEndpoint, activeConnection, apiProxyPath])

  const fetchFields = useCallback(() => {
    // if (activeConnection?.plugin !== Providers.JIRA) {
    //   return
    // }
    try {
      if (apiProxyPath.includes('null')) {
        throw new Error('Connection ID is Null')
      }
      setError(null)
      setIsFetching(true)
      const fetchIssueFields = async () => {
        const fields = await
        request
          .get(activeConnection?.connectionId ? fieldsEndpoint.replace('[:connectionId:]', activeConnection?.connectionId) : fieldsEndpoint)
          .catch(e => setError(e))
        console.log('>>> JIRA API PROXY: Fields Response...', fields)
        setFieldsResponse(fields && Array.isArray(fields.data) ? fields.data : [])
        // setTimeout(() => {
        setIsFetching(false)
        // }, 1000)
      }
      fetchIssueFields()
    } catch (e) {
      setIsFetching(false)
      setError(e)
      ToastNotification.show({ message: e.message, intent: 'danger', icon: 'error' })
    }
  }, [fieldsEndpoint, activeConnection, apiProxyPath])

  const fetchBoards = useCallback(() => {
    // if (activeConnection?.plugin !== Providers.JIRA) {
    //   return
    // }
    try {
      if (apiProxyPath.includes('null')) {
        throw new Error('Connection ID is Null')
      }
      setError(null)
      setIsFetching(true)
      const fetchApiBoards = async () => {
        const boards = await
        request
          .get(activeConnection?.connectionId ? boardsEndpoint.replace('[:connectionId:]', activeConnection?.connectionId) : boardsEndpoint)
          .catch(e => setError(e))
        console.log('>>> JIRA API PROXY: Boards Response...', boards)
        setBoardsResponse(boards && Array.isArray(boards.data?.values) ? boards.data?.values : [])
        // setTimeout(() => {
        setIsFetching(false)
        // }, 1000)
      }
      fetchApiBoards()
    } catch (e) {
      setIsFetching(false)
      setError(e)
      ToastNotification.show({ message: e.message, intent: 'danger', icon: 'error' })
    }
  }, [boardsEndpoint, activeConnection, apiProxyPath])

  const fetchAllResources = useCallback(async (connectionId, callback = () => {}) => {
    try {
      if (apiProxyPath.includes('null')) {
        throw new Error('Connection ID is Null')
      }
      setError(null)
      setIsFetching(true)
      const aR = await Promise.all([
        request
          .get(activeConnection?.connectionId ? boardsEndpoint.replace('[:connectionId:]', connectionId || activeConnection?.connectionId) : boardsEndpoint),
        request
          .get(activeConnection?.connectionId ? fieldsEndpoint.replace('[:connectionId:]', connectionId || activeConnection?.connectionId) : fieldsEndpoint),
        request
          .get(activeConnection?.connectionId ? issuesEndpoint.replace('[:connectionId:]', connectionId || activeConnection?.connectionId) : issuesEndpoint),
      ])
      console.log('>>> JIRA API PROXY: ALL API RESOURCES...', aR)
      const apiResources = { boards: aR[0]?.data?.values || [], fields: aR[1]?.data || [], issues: aR[2].data || [] }
      setAllResources(apiResources)
      setBoardsResponse(Array.isArray(aR[0]?.data.values) ? aR[0]?.data.values : [])
      setFieldsResponse(Array.isArray(aR[1]?.data) ? aR[1]?.data : [])
      setIssueTypesResponse(Array.isArray(aR[2]?.data) ? aR[2]?.data : [])
      setIsFetching(false)
      callback(apiResources)
    } catch (e) {
      setIsFetching(false)
      setError(e)
      callback(e)
    }
  }, [
    boardsEndpoint,
    fieldsEndpoint,
    issuesEndpoint,
    activeConnection,
    apiProxyPath])

  const createListData = (data = [], titleProperty = 'name', valueProperty = 'id') => {
    return data.map((d, dIdx) => ({
      id: d[valueProperty],
      key: d[valueProperty],
      title: d[titleProperty],
      value: d[valueProperty],
      type: d.schema?.type || 'string'
    }))
  }

  useEffect(() => {
    setIssueTypes(issueTypesResponse
      ? createListData(issueTypesResponse).reduce((pV, cV) => !pV.some(i => i.value === cV.value) ? [...pV, cV] : [...pV], [])
      : [])
  }, [issueTypesResponse])

  useEffect(() => {
    setFields(fieldsResponse ? createListData(fieldsResponse, 'name', 'id') : [])
  }, [fieldsResponse])

  useEffect(() => {
    setBoards(boardsResponse ? createListData(boardsResponse, 'name', 'id') : [])
  }, [boardsResponse])

  useEffect(() => {
    console.log('>>> JIRA API PROXY: FIELD SELECTOR LIST DATA', fields)
  }, [fields])

  useEffect(() => {
    console.log('>>> JIRA API PROXY: FIELD SELECTOR BOARDS DATA', boards)
  }, [boards])

  useEffect(() => {
    if (error) {
      console.log('>>> JIRA PROXY API ERROR!', error)
    }
  }, [error])

  // useEffect(() => {
  //   console.log('>>> JIRA PROXY ACTIVE CONNECTION...', activeConnection)
  // }, [activeConnection])

  return {
    isFetching,
    fetchFields,
    fetchIssueTypes,
    fetchBoards,
    fetchAllResources,
    createListData,
    issueTypesResponse,
    fieldsResponse,
    boardsResponse,
    allResources,
    issueTypes,
    fields,
    boards,
    error
  }
}

export default useJIRA
