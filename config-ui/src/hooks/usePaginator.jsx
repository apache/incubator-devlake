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
import React, { useState, useEffect, useRef, useCallback } from 'react'
import {
  Popover,
  Menu,
  MenuItem,
  Button
} from '@blueprintjs/core'

import { integrationsData } from '@/data/integrations'

const PagingOptionsMenu = (props) => {
  const { perPage = 10, setPerPage = () => {} } = props
  return (
    <Menu>
      <MenuItem
        active={perPage === 10}
        icon='key-option'
        text='10 Records'
        onClick={() => setPerPage(10)}
      />
      <MenuItem
        active={perPage === 25}
        icon='key-option'
        text='25 Records'
        onClick={() => setPerPage(25)}
      />
      <MenuItem
        active={perPage === 50}
        icon='key-option'
        text='50 Records'
        onClick={() => setPerPage(50)}
      />
      <MenuItem
        active={perPage === 75}
        icon='key-option'
        text='75 Records'
        onClick={() => setPerPage(75)}
      />
      <MenuItem
        active={perPage === 100}
        icon='key-option'
        text='100 Records'
        onClick={() => setPerPage(100)}
      />
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
    isLoading = false,
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
        disabled={currentPage?.current === 1 || isLoading}
      />
      <Button
        style={{ whiteSpace: 'nowrap' }}
        disabled={currentPage?.current === maxPage || isLoading}
        onClick={onNextPage}
        className='pagination-btn btn-next-page'
        rightIcon='step-forward'
        text='NEXT'
        small
      />
    </div>
  )
}

function usePaginator (initialLoadingState = false) {
  // const [integrations, setIntegrations] = useState(integrationsData)

  const [data, setData] = useState([])

  const [filteredData, setFilteredData] = useState([])
  const [pagedData, setPagedData] = useState([])
  const [pageOptions, setPageOptions] = useState([10, 25, 50, 75, 100])
  const currentPage = useRef(1)
  const [perPage, setPerPage] = useState(pageOptions[0])
  const [maxPage, setMaxPage] = useState(
    Math.max(1, Math.ceil(pagedData.length / perPage))
  )

  const [isLoading, setIsLoading] = useState(initialLoadingState || false)
  const [isProcessing, setIsProcessing] = useState(false)
  const [refresh, setRefresh] = useState(false)

  const nextPage = useCallback(() => {
    console.log('>>> PAGINATOR: GO NEXT PAGE ...')
    currentPage.current = Math.min(maxPage, currentPage.current + 1)
    // setRefresh((r) => !r)
    console.log('>>>> NEXT PAGE', currentPage.current)
  }, [maxPage, currentPage.current, setRefresh])

  const prevPage = useCallback(() => {
    console.log('>>> PAGINATOR: GO PREV PAGE ...')
    currentPage.current = Math.max(1, currentPage.current - 1)
    // setRefresh((r) => !r)
    console.log('>>>> PREV PAGE', currentPage.current)
  }, [maxPage, currentPage.current, setRefresh])

  const resetPage = useCallback(() => {
    currentPage.current = 1
  }, [currentPage.current])

  const paginateData = useCallback(() => {
    console.log('>> PAGINATING DATA...', data)
    const sliceOffset = currentPage.current >= 2 ? -1 : 0
    const sliceBegin = currentPage.current === 1
      ? 0
      : (currentPage.current + sliceOffset) * perPage
    const sliceEnd = currentPage.current === 1
      ? perPage
      : ((currentPage.current + sliceOffset) * perPage) + perPage
    console.log('>> CURRENT PAGE = ', currentPage.current)
    console.log('>> START RECORD INDEX ====', sliceBegin)
    console.log('>> END RECORD INDEX ====', sliceEnd)
    setPagedData(data.slice(sliceBegin, sliceEnd))
  }, [data, perPage, currentPage?.current, setPagedData])

  useEffect(() => {
    paginateData()
  }, [refresh, perPage, data, currentPage?.current, paginateData])

  useEffect(() => {
    console.log('>> PAGINATOR: DATA ON PAGE...', pagedData)
    setMaxPage(Math.max(1, Math.ceil(pagedData.length / perPage)))
  }, [pagedData, perPage])

  const controls = useCallback(() => {
    return (
      <Controls
        currentPage={currentPage}
        onNextPage={nextPage}
        onPrevPage={prevPage}
        maxPage={maxPage}
        perPage={perPage}
        isLoading={isLoading}
        pagingOptionsMenu={
          <PagingOptionsMenu perPage={perPage} setPerPage={setPerPage} />
        }
      />
    )
  }, [currentPage, maxPage, perPage, isLoading, nextPage, prevPage, setPerPage])

  useEffect(() => {
    console.log('>>> PAGINATOR: DATA ...', data)
  }, [data])

  return {
    nextPage,
    prevPage,
    resetPage,
    controls,
    isLoading,
    isProcessing,
    data,
    pagedData,
    filteredData,
    perPage,
    maxPage,
    setIsLoading,
    setIsProcessing,
    setData,
    setPagedData,
    setFilteredData,
    setMaxPage
  }
}

export default usePaginator
