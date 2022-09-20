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
import React, { useEffect, useState } from 'react'
import {
  Button,
  Colors,
  FormGroup,
  InputGroup,
  Label,
  Position,
  Tooltip
} from '@blueprintjs/core'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import AppCrumbs from '@/components/Breadcrumbs'
import Content from '@/components/Content'
import SaveAlert from '@/components/SaveAlert'
import { DEVLAKE_ENDPOINT } from '@/utils/config'
import request from '@/utils/request'

export default function Configure() {
  const [alertOpen, setAlertOpen] = useState(false)
  const [dbUrl, setDbUrl] = useState()
  const [port, setPort] = useState()
  const [mode, setMode] = useState()
  const [config, setConfig] = useState({
    DB_URL: null,
    PORT: null,
    MODE: null
  })

  async function saveAll(e) {
    e.preventDefault()

    // config.DB_URL = dbUrl
    // config.PORT = port
    // config.MODE = mode

    await request.post(`${DEVLAKE_ENDPOINT}/env`, { ...config })
    setAlertOpen(true)
  }

  const isValidForm = (dbUrl, port) => {
    return (
      port &&
      port.toString().length > 0 &&
      dbUrl &&
      dbUrl.toString().length >= 2
    )
  }

  useEffect(() => {
    const fetchEnv = async () => {
      const env = await request.get(`${DEVLAKE_ENDPOINT}/env`)
      setConfig(env.data)
      setDbUrl(env.data.DB_URL)
      setPort(env.data.PORT)
      setMode(env.data.MODE)
    }
    try {
      fetchEnv()
    } catch (e) {
      console.log('>> API CONFIGURATION / UNABLE TO FETCH ENV!', e)
    }
  }, [])

  return (
    <>
      <div className='container'>
        <Nav />
        <Sidebar />
        <Content>
          <main className='main'>
            <AppCrumbs
              items={[
                { href: '/', icon: false, text: 'Dashboard' },
                { href: '/', icon: false, text: 'DEV LAKE' },
                {
                  href: '/lake/api/configuration',
                  icon: false,
                  text: 'Lake API Setup'
                }
              ]}
            />
            <div className='headlineContainer'>
              <h1>
                Configuration{' '}
                {alertOpen && (
                  <span style={{ color: Colors.GRAY3 }}>(Please Restart)</span>
                )}
              </h1>
              <p className='description'>
                Configure your <code className='code'>.env</code> file values
              </p>
            </div>

            <form className='form'>
              <div className='headlineContainer'>
                <h2 className='headline'>Dev Lake Database Connection</h2>
                <p className='description'>Settings for the MySQL database</p>
              </div>

              <div className='formContainer'>
                <FormGroup
                  inline
                  labelFor='db-url'
                  className='formGroup'
                  helperText='DB_URL'
                  contentClassName='formGroup'
                  readOnly={alertOpen}
                >
                  <Tooltip
                    content='The URL Connection string to the database'
                    position={Position.TOP}
                  >
                    <Label>
                      Database&nbsp;URL <span className='requiredStar'>*</span>
                      <InputGroup
                        id='db-url'
                        placeholder='Enter DB Connection String'
                        defaultValue={dbUrl}
                        onChange={(e) => setDbUrl(e.target.value)}
                        className='input'
                        readOnly={alertOpen}
                      />
                    </Label>
                  </Tooltip>
                </FormGroup>
              </div>

              <div className='headlineContainer'>
                <h2 className='headline'>Dev Lake API Server</h2>
                <p className='description'>Configure main REST Settings</p>
              </div>

              <div className='formContainer'>
                <FormGroup
                  inline
                  labelFor='port'
                  className='formGroup'
                  helperText='PORT'
                  contentClassName='formGroup'
                  readOnly={alertOpen}
                >
                  <Tooltip
                    content='The main port for the REST server'
                    position={Position.TOP}
                  >
                    <Label>
                      Port <span className='requiredStar'>*</span>
                      <InputGroup
                        id='port'
                        placeholder='Enter Port eg. :8080'
                        defaultValue={port}
                        onChange={(e) => setPort(e.target.value)}
                        className='input'
                        readOnly={alertOpen}
                      />
                    </Label>
                  </Tooltip>
                </FormGroup>
              </div>

              <div className='formContainer'>
                <FormGroup
                  inline
                  labelFor='mode'
                  className='formGroup'
                  helperText='MODE'
                  contentClassName='formGroup'
                  readOnly={alertOpen}
                >
                  <Tooltip
                    content='The development mode for the server'
                    position={Position.TOP}
                  >
                    <Label>
                      Mode
                      <InputGroup
                        id='mode'
                        placeholder='Enter Mode eg. debug'
                        defaultValue={mode}
                        onChange={(e) => setMode(e.target.value)}
                        className='input'
                        readOnly={alertOpen}
                      />
                    </Label>
                  </Tooltip>
                </FormGroup>
              </div>

              <Button
                type='submit'
                icon='cloud-upload'
                intent='primary'
                loading={alertOpen}
                onClick={(e) => saveAll(e)}
                disabled={!isValidForm(dbUrl, port)}
              >
                Save Configuration
              </Button>
            </form>
          </main>
        </Content>
      </div>

      <SaveAlert alertOpen={alertOpen} onClose={() => setAlertOpen(false)} />
    </>
  )
}
