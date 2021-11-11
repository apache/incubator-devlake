// import React, { useState, useEffect } from 'react'
// import { FormGroup, InputGroup, Button, Label } from '@blueprintjs/core'
// import Nav from '../../../components/Nav'
// import Sidebar from '../../../components/Sidebar'
// import Content from '../../../components/Content'
// import SaveAlert from '../../../components/SaveAlert'
// import { DEVLAKE_ENDPOINT } from '../../../utils/config'

// export default function Jenkins () {
//   const [alertOpen, setAlertOpen] = useState(false)
//   const [jenkinsEndpoint, setJenkinsEndpoint] = useState()
//   const [jenkinsUsername, setJenkinsUsername] = useState()
//   const [jenkinsPassword, setJenkinsPassword] = useState()

//   function updateEnv (key, value) {
//     fetch(`${DEVLAKE_ENDPOINT}/api/setenv/${key}/${encodeURIComponent(value)}`)
//   }

//   function saveAll (e) {
//     e.preventDefault()
//     updateEnv('JENKINS_ENDPOINT', jenkinsEndpoint)
//     updateEnv('JENKINS_USERNAME', jenkinsUsername)
//     updateEnv('JENKINS_PASSWORD', jenkinsPassword)
//     setAlertOpen(true)
//   }

//   useEffect(() => {
//     fetch(`${DEVLAKE_ENDPOINT}/api/getenv`)
//       .then(response => response.json())
//       .then(env => {
//         setJenkinsEndpoint(env.JENKINS_ENDPOINT)
//         setJenkinsUsername(env.JENKINS_USERNAME)
//         setJenkinsPassword(env.JENKINS_PASSWORD)
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
//               <h2 className='headline'>Jenkins Configuration</h2>
//               <p className='description'>Jenkins account and config settings</p>
//             </div>

//             <div className='formContainer'>
//               <FormGroup
//                 inline={true}
//                 labelFor='jenkins-endpoint'
//                 helperText='JENKINS_ENDPOINT'
//                 className='formGroup'
//                 contentClassName='formGroup'
//               >
//                 <Label>
//                   API&nbsp;Endpoint <span className='requiredStar'>*</span>
//                   <InputGroup
//                     id='jenkins-endpoint'
//                     placeholder='Enter Jenkins API endpoint'
//                     defaultValue={jenkinsEndpoint}
//                     onChange={(e) => setJenkinsEndpoint(e.target.value)}
//                     className='input'
//                   />
//                 </Label>
//               </FormGroup>
//             </div>

//             <div className='formContainer'>
//               <FormGroup
//                 inline={true}
//                 labelFor='jenkins-username'
//                 helperText='JENKINS_USERNAME'
//                 className='formGroup'
//                 contentClassName='formGroup'
//               >
//                 <Label>
//                   Username <span className='requiredStar'>*</span>
//                   <InputGroup
//                     id='jenkins-username'
//                     placeholder='Enter Jenkins Username'
//                     defaultValue={jenkinsUsername}
//                     onChange={(e) => setJenkinsUsername(e.target.value)}
//                     className='input'
//                   />
//                 </Label>
//               </FormGroup>
//             </div>

//             <div className='formContainer'>
//               <FormGroup
//                 inline={true}
//                 labelFor='jenkins-password'
//                 helperText='JENKINS_PASSWORD'
//                 className='formGroup'
//                 contentClassName='formGroup'
//               >
//                 <Label>
//                   Password <span className='requiredStar'>*</span>
//                   <InputGroup
//                     id='jenkins-password'
//                     placeholder='Enter Jenkins Password'
//                     defaultValue={jenkinsPassword}
//                     onChange={(e) => setJenkinsPassword(e.target.value)}
//                     className='input'
//                   />
//                 </Label>
//               </FormGroup>
//             </div>

//             <SaveAlert alertOpen={alertOpen} onClose={() => setAlertOpen(false)} />
//             <Button
//               type='submit'
//               outlined={true}
//               large={true}
//               className='saveBtn'
//               onClick={(e) => saveAll(e)}
//             >
//               Save Config
//             </Button>
//           </form>
//         </main>
//       </Content>
//     </div>
//   )
// }
