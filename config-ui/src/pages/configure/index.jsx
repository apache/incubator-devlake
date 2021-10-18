import React, { useEffect, useState } from 'react'
import { FormGroup, InputGroup, Button, Tooltip, Position, Label } from '@blueprintjs/core'
import Nav from '../../components/Nav'
import Sidebar from '../../components/Sidebar'
import Content from '../../components/Content'
import SaveAlert from '../../components/SaveAlert'

export default function Home (props) {
  const [env, setEnv] = useState()
  const [alertOpen, setAlertOpen] = useState(false)
  const [dbUrl, setDbUrl] = useState()
  const [port, setPort] = useState()
  const [mode, setMode] = useState()
  const SERVER_HOST = 'http://localhost:5000'

  function updateEnv (key, value) {
    fetch(`/api/setenv/${key}/${encodeURIComponent(value)}`)
  }

  function saveAll (e) {
    e.preventDefault()
    updateEnv('DB_URL', dbUrl)
    updateEnv('PORT', port)
    updateEnv('MODE', mode)
    setAlertOpen(true)
  }

  useEffect(() => {
    fetch(`${SERVER_HOST}/api/getenv`)
      .then(response => response.json())
      .then(env => {
        setDbUrl(env.DB_URL)
        setPort(env.PORT)
        setMode(env.MODE)
        setEnv(env)
      })
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

              <Button type='submit' outlined large className='saveBtn' onClick={() => saveAll}>Save Config</Button>
            </form>
          </main>
        </Content>
      </div>

      <SaveAlert alertOpen={alertOpen} onClose={() => setAlertOpen(false)} />
    </>
  )
}
