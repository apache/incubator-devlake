// import React, { useState, useEffect } from 'react'
// import { FormGroup, InputGroup, Button, Label, Tooltip, Position } from '@blueprintjs/core'
// import Nav from '../../../components/Nav'
// import Sidebar from '../../../components/Sidebar'
// import Content from '../../../components/Content'
// import SaveAlert from '../../../components/SaveAlert'
// import { DEVLAKE_ENDPOINT } from '../../../utils/config'

// export default function Gitlab () {
//   const [alertOpen, setAlertOpen] = useState(false)
//   const [gitlabEndpoint, setGitlabEndpoint] = useState()
//   const [gitlabAuth, setGitlabAuth] = useState()
//   const [jiraBoardGitlabeProjects, setJiraBoardGitlabeProjects] = useState()

//   function updateEnv (key, value) {
//     fetch(`${DEVLAKE_ENDPOINT}/api/setenv/${key}/${encodeURIComponent(value)}`)
//   }

//   function saveAll (e) {
//     e.preventDefault()
//     updateEnv('GITLAB_ENDPOINT', gitlabEndpoint)
//     updateEnv('GITLAB_AUTH', gitlabAuth)
//     updateEnv('JIRA_BOARD_GITLAB_PROJECTS', jiraBoardGitlabeProjects)
//     setAlertOpen(true)
//   }

//   useEffect(() => {
//     fetch(`${DEVLAKE_ENDPOINT}/api/getenv`)
//       .then(response => response.json())
//       .then(env => {
//         setGitlabEndpoint(env.GITLAB_ENDPOINT)
//         setGitlabAuth(env.GITLAB_AUTH)
//         setJiraBoardGitlabeProjects(env.JIRA_BOARD_GITLAB_PROJECTS)
//       })
//   }, [])

//   return (
//     <div className='container'>

//       <Nav />
//       <Sidebar />
//       <Content>
//         <main className='main'>

//           <form className='form'>
//             <div className='headlineContainer'>
//               <h2 className='headline'>Gitlab Configuration</h2>
//               <p className='description'>Gitlab account and config settings</p>
//             </div>

//             <div className='formContainer'>
//               <FormGroup
//                 inline={true}
//                 labelFor='gitlab-endpoint'
//                 helperText='GITLAB_ENDPOINT'
//                 className='formGroup'
//                 contentClassName='formGroup'
//               >
//                 <Label>
//                   API&nbsp;Endpoint <span className='requiredStar'>*</span>
//                   <InputGroup
//                     id='gitlab-endpoint'
//                     placeholder='Enter Gitlab API endpoint'
//                     defaultValue={gitlabEndpoint}
//                     onChange={(e) => setGitlabEndpoint(e.target.value)}
//                     className='input'
//                   />
//                 </Label>
//               </FormGroup>
//             </div>

//             <div className='formContainer'>
//               <FormGroup
//                 inline={true}
//                 labelFor='gitlab-auth'
//                 helperText='GITLAB_AUTH'
//                 className='formGroup'
//                 contentClassName='formGroup'
//               >
//                 <Label>
//                   Auth&nbsp;Token <span className='requiredStar'>*</span>
//                   <InputGroup
//                     id='gitlab-auth'
//                     placeholder='Enter Gitlab Auth Token eg. uJVEDxabogHbfFyu2riz'
//                     defaultValue={gitlabAuth}
//                     onChange={(e) => setGitlabAuth(e.target.value)}
//                     className='input'
//                   />
//                 </Label>
//               </FormGroup>
//             </div>

//             <div className='headlineContainer'>
//               <h3 className='headline'>Jira / Gitlab Connection</h3>
//               <p className='description'>Connect jira board to gitlab projects</p>
//             </div>

//             <div className='formContainer'>
//               <FormGroup
//                 inline={true}
//                 labelFor='jira-board-projects'
//                 helperText='JIRA_BOARD_GITLAB_PROJECTS'
//                 className='formGroup'
//                 contentClassName='formGroup'
//               >
//                 <Tooltip content='Jira board and Gitlab projects relationship' position={Position.TOP}>
//                   <Label>
//                     Jira&nbsp;Board&nbsp;Gitlab&nbsp;Projects
//                     <InputGroup
//                       id='jira-storypoint-field'
//                       placeholder='<JIRA_BOARD>:<GITLAB_PROJECT_ID>,...; eg. 8:8967944,8967945;9:8967946,8967947'
//                       defaultValue={jiraBoardGitlabeProjects}
//                       onChange={(e) => setJiraBoardGitlabeProjects(e.target.value)}
//                       className='input'
//                     />
//                   </Label>
//                 </Tooltip>
//               </FormGroup>
//             </div>

//             <Button
//               type='submit'
//               outlined={true}
//               large={true}
//               className='saveBtn'
//               onClick={(e) => saveAll(e)}
//             >
//               Save Config
//             </Button>

//             <SaveAlert alertOpen={alertOpen} onClose={() => setAlertOpen(false)} />
//           </form>
//         </main>
//       </Content>
//     </div>
//   )
// }
