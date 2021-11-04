import React, { useEffect, useState } from 'react'
import Nav from '../../components/Nav'
import Sidebar from '../../components/Sidebar'
import Content from '../../components/Content'
import { DEVLAKE_ENDPOINT } from '../../utils/config'
import request from '../../utils/request'

export default function Tasks () {

  let [tasks, setTasks] = useState([])
  useEffect(async () => {
    if (tasks.length === 0) {
      let res = await request.get(`${DEVLAKE_ENDPOINT}/task`)
      setTasks(res?.data?.tasks)
    } 
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
                <table className='bp3-html-table bp3-html-table-bordered connections-table' style={{ width: '100%' }}>
                  <thead>
                    <tr>
                      <th>ID</th>
                      <th>CreatedAt</th>
                      <th>Plugin</th>
                      <th>Progress</th>
                      <th>Status</th>
                    </tr>
                  </thead>
                  <tbody>
                    {
                      tasks.length > 0 ?
                      tasks.map((task, i) => 
                        <tr key={i}>
                          <td>{task.ID}</td>
                          <td>{task.CreatedAt}</td>
                          <td>{task.plugin}</td>
                          <td>{task.progress}</td>
                          <td>{task.status}</td>
                        </tr>
                        )
                      : null
                    }
                  </tbody>
                </table>
              </div>
            </>
          </main>
        </Content>
      </div>
    </>
  )
}
