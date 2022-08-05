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
import React, { Fragment, useEffect, useState, useCallback } from 'react'
import { useParams, useHistory, useRouteMatch } from 'react-router-dom'

const BlueprintNavigationLinks = (props) => {
  const history = useHistory()
  const activeRoute = useRouteMatch()

  const {
    blueprint,
    links = [
      {
        id: 0,
        name: 'status',
        label: 'Status',
        route: `/blueprints/detail/${blueprint?.id}`,
        active: activeRoute?.url.endsWith(`/blueprints/detail/${blueprint?.id}`)
      },
      {
        id: 1,
        name: 'status',
        label: 'Settings',
        route: `/blueprints/settings/${blueprint?.id}`,
        active: activeRoute?.url.endsWith(`/blueprints/settings/${blueprint?.id}`)
      }
    ]
  } = props

  const routeToLocation = useCallback((route) => {
    history.push(route)
  }, [history])

  return (
    <div
      className='blueprint-navigation'
      style={{
        alignSelf: 'center',
        display: 'flex',
        margin: '20px auto',
      }}
    >
      {links.map((link) => (
        <div key={`blueprint-nav-link-key-${link?.id}`} style={{ marginRight: '10px' }}>
          <a
            href='#'
            className={`blueprint-navigation-link ${link?.active ? 'active' : ''}`}
            onClick={() => routeToLocation(link?.route)}
          >
            {link.label}
          </a>
        </div>
      ))}
    </div>
  )
}

export default BlueprintNavigationLinks
