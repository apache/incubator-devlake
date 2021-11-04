import React, { useEffect, useState } from 'react'
import Nav from '../../components/Nav'
import Sidebar from '../../components/Sidebar'
import AppCrumbs from '@/components/Breadcrumbs'
import Content from '../../components/Content'
import {
  Tooltip, Position, FormGroup, InputGroup, Button, Label, Icon, Classes, Dialog
} from '@blueprintjs/core'
import axios from 'axios'
import { DEVLAKE_ENDPOINT } from '../../utils/config'
import { LABEL } from '@blueprintjs/core/lib/esm/common/classes'

export default function Tasks () {


  useEffect(async () => {
    let res = await axios.get(`${DEVLAKE_ENDPOINT}/task`)
    console.log('res', res)
  }, [])

  return (
    <>
      <div className='container'>
      <Nav />
      <Sidebar />
        <Content>
          <main className='main'>
            <>
              <div className='headlineContainer'>
                <h1>Tasks</h1>
                <table>

                </table>
                {/* <table className='bp3-html-table bp3-html-table-bordered connections-table' style={{ width: '100%' }}>
                      <thead>
                        <tr>
                          <th>Task Name</th>
                          <th>Endpoint</th>
                          <th>Status</th>
                          <th />
                        </tr>
                      </thead>
                      <tbody>
                        {}
                      </tbody>
                    </table> */}
              </div>
            </>
          </main>
        </Content>
      </div>
    </>
  )
}
