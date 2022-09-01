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

const useGitlab = ({ apiProxyPath, projectsEndpoint }, activeConnection = null) => {
  const [isFetching, setIsFetching] = useState(false)
  const [projects, setProjects] = useState([])
  const [error, setError] = useState()

  const fetchProjects = useCallback(async (search = '') => {
    try {
      if (apiProxyPath.includes('null')) {
        throw new Error('Connection ID is Null')
      }
      setError(null)
      setIsFetching(true)
      const endpoint = projectsEndpoint
        .replace('[:connectionId:]', activeConnection?.connectionId)
        .replace('[:search:]', search)
      const projectsResponse = await request.get(endpoint)
      if (projectsResponse && projectsResponse.status === 200 && projectsResponse.data) {
        setProjects(createListData(projectsResponse.data))
      } else {
        throw new Error('request projects fail')
      }
    } catch (e) {
      setError(e)
      ToastNotification.show({ message: e.message, intent: 'danger', icon: 'error' })
    } finally {
      setIsFetching(false)
    }
  }, [projectsEndpoint, activeConnection, apiProxyPath])

  const createListData = (data = [], titleProperty = 'name_with_namespace', valueProperty = 'name') => {
    return data.map((d, dIdx) => ({
      id: d.id || dIdx,
      key: d.key ? d.key : dIdx,
      title: d[titleProperty],
      value: d[valueProperty],
      type: 'string'
    }))
  }

  useEffect(() => {
    console.log('>>> GITLAB API PROXY: FIELD SELECTOR PROJECTS DATA', projects)
  }, [projects])

  return {
    isFetching,
    fetchProjects,
    projects,
    error
  }
}

export default useGitlab
