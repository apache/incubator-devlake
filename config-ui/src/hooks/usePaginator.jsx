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
import React, { useState, useEffect, useRef, useCallback, useMemo } from 'react'
import { Popover, Menu, MenuItem, Button } from '@blueprintjs/core'

import { integrationsData } from '@/data/integrations'

const PagingOptionsMenu = (props) => {
  const {
    pageOptions = [10, 25, 50, 75, 100],
    perPage = 10,
    setPerPage = (page) => undefined
  } = props
  return (
    <Menu>
      {pageOptions &&
        pageOptions.map((pageOption) => (
          <MenuItem
            key={pageOption}
            active={perPage === pageOption}
            icon='key-option'
            text={`${pageOption} Records`}
            onClick={() => setPerPage(pageOption)}
          />
        ))}
    </Menu>
  )
}

const Controls = (props) => {
  const {
    enabled = true,
    pagingOptionsMenu,
    currentPage,
    perPage = 10,
    maxPage,
    onPrevPage = () => {},
    onNextPage = () => {},
    isLoading = false
  } = props

  return (
    <div
      className='pagination-controls'
      style={{ display: 'flex', whiteSpace: 'nowrap' }}
    >
      <Popover placement='bottom'>
        <Button
          className='btn-select-page-size'
          style={{ whiteSpace: 'nowrap' }}
          icon='numbered-list'
          text={`Rows: ${perPage}`}
          disabled={isLoading}
          outlined
          minimal
        />
        <>{pagingOptionsMenu}</>
      </Popover>
      <Button
        onClick={onPrevPage}
        className='pagination-btn btn-prev-page'
        icon='step-backward'
        small
        text='PREV'
        style={{ marginLeft: '5px', marginRight: '5px', whiteSpace: 'nowrap' }}
        disabled={currentPage === 1 || isLoading}
      />
      <Button
        style={{ whiteSpace: 'nowrap' }}
        disabled={currentPage === maxPage || isLoading}
        onClick={onNextPage}
        className='pagination-btn btn-next-page'
        rightIcon='step-forward'
        text='NEXT'
        small
      />
    </div>
  )
}

function usePaginator(initialLoadingState = false) {
  // const [integrations, setIntegrations] = useState(integrationsData)

  const [data, setData] = useState([])

  // filter related
  const [filterParams, setFilterParams] = useState({})
  const [filterFunc, setFilterFunc] = useState(() => (params, item) => true)
  const filteredData = useMemo(() => {
    console.log('>> SET FILTER DATA BY', filterParams, data)
    const filteredData = []
    for (const item of data) {
      if (filterFunc(filterParams, item)) {
        filteredData.push(item)
      }
    }
    return filteredData
  }, [data, filterParams, filterFunc])

  // page related
  const [pagedData, setPagedData] = useState([])
  const [pageOptions, setPageOptions] = useState([5, 25, 50, 75, 100])
  const [currentPage, setCurrentPage] = useState(1)
  const [perPage, setPerPage] = useState(pageOptions[1])
  const maxPage = useMemo(
    () => Math.max(1, Math.ceil(filteredData.length / perPage)),
    [filteredData, perPage]
  )

  // others
  const [isLoading, setIsLoading] = useState(initialLoadingState || false)
  const [isProcessing, setIsProcessing] = useState(false)
  const [refresh, setRefresh] = useState(false)

  const setDataWithDefault = useCallback(
    (data) => {
      console.log('>> SET ALL DATA...', data)
      setData(data || [])
    },
    [setData]
  )

  const goNextPage = useCallback(() => {
    console.log('>>> PAGINATOR: GO NEXT PAGE ...')
    setCurrentPage((currentPage) => Math.min(maxPage, currentPage + 1))
    // setRefresh((r) => !r)
    console.log('>>>> NEXT PAGE', currentPage)
  }, [maxPage, currentPage])

  const goPrevPage = useCallback(() => {
    console.log('>>> PAGINATOR: GO PREV PAGE ...')
    setCurrentPage((currentPage) => Math.max(1, currentPage - 1))
    // setRefresh((r) => !r)
    console.log('>>>> PREV PAGE', currentPage)
  }, [currentPage])

  const changePerPage = useCallback(
    (perPage) => {
      setPerPage(perPage)
      setCurrentPage(1)
    },
    [setPerPage]
  )

  const resetPage = useCallback(() => {
    setCurrentPage(1)
  }, [setCurrentPage])

  const paginateData = useCallback(() => {
    console.log('>> FILTERED DATA...', filteredData)
    const sliceBegin = (currentPage - 1) * perPage
    const sliceEnd = currentPage * perPage
    setPagedData(filteredData.slice(sliceBegin, sliceEnd))
    console.log('>> PAGED DATA = ', filteredData.slice(sliceBegin, sliceEnd))
  }, [filteredData, perPage, currentPage, setPagedData])

  useEffect(() => {
    paginateData()
  }, [/* refresh, */ perPage, filteredData, currentPage, paginateData])

  useEffect(() => {
    setCurrentPage((currentPage) => Math.min(maxPage, currentPage))
  }, [maxPage, setCurrentPage])

  const renderControlsComponent = useCallback(() => {
    return (
      <Controls
        currentPage={currentPage}
        onNextPage={goNextPage}
        onPrevPage={goPrevPage}
        maxPage={maxPage}
        perPage={perPage}
        isLoading={isLoading}
        pagingOptionsMenu={
          <PagingOptionsMenu
            pageOptions={pageOptions}
            perPage={perPage}
            setPerPage={changePerPage}
          />
        }
      />
    )
  }, [
    currentPage,
    maxPage,
    perPage,
    pageOptions,
    isLoading,
    goNextPage,
    goPrevPage,
    changePerPage
  ])

  return {
    goNextPage,
    goPrevPage,
    resetPage,
    setPageOptions,
    renderControlsComponent,
    isLoading,
    isProcessing,
    data,
    filteredData,
    pagedData,
    perPage,
    maxPage,
    setData: setDataWithDefault,
    setPagedData,
    setFilterParams,
    setFilterFunc,
    // setMaxPage,
    setIsLoading,
    setIsProcessing
  }
}

export default usePaginator
