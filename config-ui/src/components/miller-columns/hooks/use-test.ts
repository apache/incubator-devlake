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

import { useMemo, useState } from 'react'

const items = [
  {
    id: 1,
    title: 'merico-dev',
    total: 13,
    items: [
      {
        id: 11,
        title: 'devlake'
      },
      {
        id: 12,
        title: 'devstream'
      },
      {
        id: 13,
        title: 'another-repo'
      },
      {
        id: 14,
        title: 'repo2'
      },
      {
        id: 15,
        title: 'repo2'
      },
      {
        id: 16,
        title: 'repo2'
      },
      {
        id: 17,
        title: 'repo2'
      },
      {
        id: 18,
        title: 'repo2'
      },
      {
        id: 19,
        title: 'ae-repo',
        total: 4,
        items: [
          {
            id: 191,
            title: 'ae-repo-1'
          },
          {
            id: 192,
            title: 'ae-repo-2'
          },
          {
            id: 193,
            title: 'ae-repo-child',
            total: 1,
            items: [
              {
                id: 1931,
                title: 'ae-repo-child-1'
              }
            ]
          }
        ]
      }
    ]
  },
  {
    id: 2,
    title: 'mintsweet',
    total: 2,
    items: [
      {
        id: 21,
        title: 'reate'
      },
      {
        id: 22,
        title: 'mst-advanced'
      }
    ]
  },
  {
    id: 3,
    title: 'test'
  }
]

export const useTest = () => {
  const [ids, setIds] = useState([])

  console.log(ids)

  return useMemo(() => {
    return { items, ids, setIds }
  }, [items, ids, setIds])
}
