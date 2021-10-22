import React, { useEffect, useState } from 'react'
import { FormGroup, InputGroup, Button, Tooltip, Position, Label, Utils } from '@blueprintjs/core'
import Nav from '../../components/Nav'
import Sidebar from '../../components/Sidebar'
import Content from '../../components/Content'
import SaveAlert from '../../components/SaveAlert'
import { DEVLAKE_ENDPOINT } from '../../utils/config'
import request from '../../utils/request'

export default function Configure () {
  const [alertOpen, setAlertOpen] = useState(false)
  const [dbUrl, setDbUrl] = useState()
  const [port, setPort] = useState()
  const [mode, setMode] = useState()
  const [config, setConfig] = useState()

  async function saveAll (e) {
    e.preventDefault()

    config.DB_URL = dbUrl
    config.PORT = port
    config.MODE = mode

    await request.post(`${DEVLAKE_ENDPOINT}/env`, config)

    setAlertOpen(true)
  }

  useEffect(async () => {
    let env = await request.get(`${DEVLAKE_ENDPOINT}/env`)
    setConfig(env.data)
    setDbUrl(env.data.DB_URL)
    setPort(env.data.PORT)
    setMode(env.data.MODE)
  }, [])

  return (
    <>
      <div className='container'>

        <Nav />
        <Sidebar />
        <Content>
          <main className='main'>

            <div className='headlineContainer'>
              <h1>Configuration</h1>
              <p className='description'>Configure your <code className='code'>.env</code> file values</p>
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
                >
                  <Tooltip content='The URL Connection string to the database' position={Position.TOP}>
                    <Label>
                      Database&nbsp;URL <span className='requiredStar'>*</span>
                      <InputGroup
                        id='db-url'
                        placeholder='Enter DB Connection String'
                        defaultValue={dbUrl}
                        onChange={(e) => setDbUrl(e.target.value)}
                        className='input'
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
                >
                  <Tooltip content='The main port for the REST server' position={Position.TOP}>
                    <Label>
                      Port <span className='requiredStar'>*</span>
                      <InputGroup
                        id='port'
                        placeholder='Enter Port eg. :8080'
                        defaultValue={port}
                        onChange={(e) => setPort(e.target.value)}
                        className='input'
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
                >
                  <Tooltip content='The development mode for the server' position={Position.TOP}>
                    <Label>
                      Mode
                      <InputGroup
                        id='mode'
                        placeholder='Enter Mode eg. debug'
                        defaultValue={mode}
                        onChange={(e) => setMode(e.target.value)}
                        className='input'
                      />
                    </Label>
                  </Tooltip>
                </FormGroup>
              </div>

              <Button
                type='submit'
                outlined
                large
                className='saveBtn'
                onClick={(e) => saveAll(e)}
              >
                Save Config
              </Button>
            </form>
          </main>
        </Content>
      </div>

      <SaveAlert alertOpen={alertOpen} onClose={() => setAlertOpen(false)} />
    </>
  )
}
