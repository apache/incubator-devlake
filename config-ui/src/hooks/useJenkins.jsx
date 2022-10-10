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

const useJenkins = (
  { apiProxyPath, jobsEndpoint },
  activeConnection = null
) => {
  const [isFetching, setIsFetching] = useState(false)
  const [jobs, setJobs] = useState([])
  const [error, setError] = useState()

  const fetchJobs = useCallback(
    async () => {
      try {
        if (apiProxyPath.includes('null')) {
          throw new Error('Connection ID is Null')
        }
        setError(null)
        setIsFetching(true)
        // only search when type more than 2 chars
        const endpoint = jobsEndpoint
          .replace('[:connectionId:]', activeConnection?.connectionId)
        const jobsResponse = await request.get(endpoint)
        if (
          jobsResponse &&
          jobsResponse.status === 200 &&
          jobsResponse.data &&
          jobsResponse.data.jobs
        ) {
          setJobs(createListData(jobsResponse.data?.jobs))
        } else {
          throw new Error('request jobs fail')
        }
      } catch (e) {
        setError(e)
        ToastNotification.show({
          message: e.message,
          intent: 'danger',
          icon: 'error'
        })
      } finally {
        setIsFetching(false)
      }
    },
    [jobsEndpoint, activeConnection, apiProxyPath]
  )

  const createListData = (
    data = [],
    titleProperty = 'name',
    valueProperty = 'name',
  ) => {
    return data.map((d, dIdx) => ({
      id: d[valueProperty],
      key: d[valueProperty],
      title: d[titleProperty],
      value: d[valueProperty],
      type: 'string'
    }))
  }

  useEffect(() => {
    console.log('>>> Jenkins API PROXY: FIELD SELECTOR JOBS DATA', jobs)
  }, [jobs])

  return {
    isFetching,
    fetchJobs,
    jobs,
    error
  }
}

export default useJenkins
