// import React, { useState, useEffect } from 'react'
// import {
//   Tooltip, Position, FormGroup, InputGroup, Button, Label, Icon, Classes, Dialog
// } from '@blueprintjs/core'
// import Nav from '../../../components/Nav'
// import Sidebar from '../../../components/Sidebar'
// import Content from '../../../components/Content'
// import SaveAlert from '../../../components/SaveAlert'
// import MappingTag from './MappingTag'
// import MappingTagStatus from './MappingTagStatus'
// import ClearButton from './ClearButton'
// import { findStrBetween } from '../../../utils/findStrBetween'
// import { DEVLAKE_ENDPOINT } from '../../../utils/config'

// export default function Jira () {
//   const [alertOpen, setAlertOpen] = useState(false)
//   const [jiraEndpoint, setJiraEndpoint] = useState()
//   const [jiraBasicAuthEncoded, setJiraBasicAuthEncoded] = useState()
//   const [jiraIssueEpicKeyField, setJiraIssueEpicKeyField] = useState()
//   const [jiraIssueStoryCoefficient, setJiraIssueStoryCoefficient] = useState()
//   const [jiraIssueStoryPointField, setJiraIssueStoryPointField] = useState()

//   // Type mappings state
//   const [typeMappingBug, setTypeMappingBug] = useState()
//   const [typeMappingIncident, setTypeMappingIncident] = useState()
//   const [typeMappingRequirement, setTypeMappingRequirement] = useState()
//   const [typeMappingAll, setTypeMappingAll] = useState()

//   const [statusMappings, setStatusMappings] = useState()

//   function setStatusMapping (key, values, status) {
//     setStatusMappings(statusMappings.map(mapping => {
//       if (mapping.key === key) {
//         mapping.mapping[status] = values
//       }
//       return mapping
//     }))
//   }

//   const [customStatusOverlay, setCustomStatusOverlay] = useState(false)
//   const [customStatusName, setCustomStatusName] = useState('')

//   function addStatusMapping (e) {
//     const type = customStatusName.trim().toUpperCase()
//     if (statusMappings.find(e => e.type === type)) {
//       return
//     }
//     const result = [
//       ...statusMappings,
//       {
//         type,
//         key: `JIRA_ISSUE_${type}_STATUS_MAPPING`,
//         mapping: {
//           Resolved: [],
//           Rejected: [],
//         }
//       }
//     ]
//     setStatusMappings(result)
//     setCustomStatusOverlay(false)
//     e.preventDefault()
//   }

//   function updateEnv (key, value) {
//     fetch(`${DEVLAKE_ENDPOINT}/api/setenv/${key}/${encodeURIComponent(value)}`)
//   }

//   function saveAll (e) {
//     e.preventDefault()
//     updateEnv('JIRA_ENDPOINT', jiraEndpoint)
//     updateEnv('JIRA_BASIC_AUTH_ENCODED', jiraBasicAuthEncoded)
//     updateEnv('JIRA_ISSUE_EPIC_KEY_FIELD', jiraIssueEpicKeyField)
//     updateEnv('JIRA_ISSUE_TYPE_MAPPING', typeMappingAll)
//     updateEnv('JIRA_ISSUE_STORYPOINT_COEFFICIENT', jiraIssueStoryCoefficient)
//     updateEnv('JIRA_ISSUE_STORYPOINT_FIELD', jiraIssueStoryPointField)

//     // Save all custom status data
//     statusMappings.forEach(mapping => {
//       const { Resolved, Rejected } = mapping.mapping
//       updateEnv(mapping.key,
//         `Rejected:${Rejected ? Rejected.join(',') : ''};Resolved:${Resolved ? Resolved.join(',') : ''};`)
//     })

//     setAlertOpen(true)
//   }

//   useEffect(() => {
//     if (typeMappingBug && typeMappingIncident && typeMappingRequirement) {
//       const typeBug = 'Bug:' + typeMappingBug.toString() + ';'
//       const typeIncident = 'Incident:' + typeMappingIncident.toString() + ';'
//       const typeRequirement = 'Requirement:' + typeMappingRequirement.toString() + ';'
//       const all = typeBug + typeIncident + typeRequirement
//       setTypeMappingAll(all)
//     }
//   }, [typeMappingBug, typeMappingIncident, typeMappingRequirement])

//   useEffect(() => {
//     fetch(`${DEVLAKE_ENDPOINT}/api/getenv`)
//       .then(response => response.json())
//       .then(env => {
//         setJiraEndpoint(env.JIRA_ENDPOINT)
//         setJiraBasicAuthEncoded(env.JIRA_BASIC_AUTH_ENCODED)
//         setJiraIssueEpicKeyField(env.JIRA_ISSUE_EPIC_KEY_FIELD)
//         setJiraIssueStoryCoefficient(env.JIRA_ISSUE_STORYPOINT_COEFFICIENT)
//         setJiraIssueStoryPointField(env.JIRA_ISSUE_STORYPOINT_FIELD)
//         setTypeMappingBug(findStrBetween(env.JIRA_ISSUE_TYPE_MAPPING, 'Bug:', 4))
//         setTypeMappingIncident(findStrBetween(env.JIRA_ISSUE_TYPE_MAPPING, 'Incident:', 9))
//         setTypeMappingRequirement(findStrBetween(env.JIRA_ISSUE_TYPE_MAPPING, 'Requirement:', 12))

//         // status mapping
//         const defaultStatusMappings = []
//         for (const [key, value] of Object.entries(env)) {
//           const m = /^JIRA_ISSUE_([A-Z]+)_STATUS_MAPPING$/.exec(key)
//           if (!m) {
//             continue
//           }
//           const type = m[1]
//           const rejected = findStrBetween(value, 'Rejected:', 9)
//           const resolved = findStrBetween(value, 'Resolved:', 9)

//           defaultStatusMappings.push({
//             type,
//             key,
//             mapping: {
//               Resolved: resolved || [],
//               Rejected: rejected || [],
//             }
//           })
//         }
//         setStatusMappings(defaultStatusMappings)
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
//               <h2 className='headline'>Jira Plugin</h2>
//               <p className='description'>Jira Account and config settings</p>
//             </div>

//             <div className='formContainer'>
//               <FormGroup
//                 label=''
//                 inline={true}
//                 labelFor='jira-endpoint'
//                 helperText='JIRA_ENDPOINT'
//                 className='formGroup'
//                 contentClassName='formGroup'
//               >
//                 <Label>
//                   Endpoint&nbsp;URL <span className='requiredStar'>*</span>
//                   <InputGroup
//                     id='jira-endpoint'
//                     placeholder='Enter Jira endpoint eg. https://merico.atlassian.net/rest'
//                     defaultValue={jiraEndpoint}
//                     onChange={(e) => setJiraEndpoint(e.target.value)}
//                     className='input'
//                   />
//                 </Label>
//               </FormGroup>
//             </div>

//             <div className='formContainer'>
//               <FormGroup
//                 inline={true}
//                 labelFor='jira-basic-auth'
//                 helperText='JIRA_BASIC_AUTH_ENCODED'
//                 className='formGroup'
//                 contentClassName='formGroup'
//               >
//                 <Label>
//                   Basic&nbsp;Auth&nbsp;Token <span className='requiredStar'>*</span>
//                   <InputGroup
//                     id='jira-basic-auth'
//                     placeholder='Enter Jira Auth eg. EJrLG8DNeXADQcGOaaaX4B47'
//                     defaultValue={jiraBasicAuthEncoded}
//                     onChange={(e) => setJiraBasicAuthEncoded(e.target.value)}
//                     className='input'
//                   />
//                 </Label>
//               </FormGroup>
//             </div>

//             <div className='headlineContainer'>
//               <h3 className='headline'>Issue Type Mappings</h3>
//               <p className='description'>Map your own issue types to Dev Lake's standard types</p>
//             </div>

//             <MappingTag
//               labelName='Bug'
//               labelIntent='danger'
//               typeOrStatus='type'
//               placeholderText='Add Issue Types...'
//               values={typeMappingBug}
//               helperText='JIRA_ISSUE_TYPE_MAPPING'
//               rightElement={<ClearButton onClick={() => setTypeMappingBug([])} />}
//               onChange={(values) => setTypeMappingBug(values)}
//             />

//             <MappingTag
//               labelName='Incident'
//               labelIntent='warning'
//               typeOrStatus='type'
//               placeholderText='Add Issue Types...'
//               values={typeMappingIncident}
//               helperText='JIRA_ISSUE_TYPE_MAPPING'
//               rightElement={<ClearButton onClick={() => setTypeMappingIncident([])} />}
//               onChange={(values) => setTypeMappingIncident(values)}
//             />

//             <MappingTag
//               labelName='Requirement'
//               labelIntent='primary'
//               typeOrStatus='type'
//               placeholderText='Add Issue Types...'
//               values={typeMappingRequirement}
//               helperText='JIRA_ISSUE_TYPE_MAPPING'
//               rightElement={<ClearButton onClick={() => setTypeMappingRequirement([])} />}
//               onChange={(values) => setTypeMappingRequirement(values)}
//             />

//             <div className='headlineContainer'>
//               <h3 className='headline'>Issue Status Mappings</h3>
//               <p className='description'>Map your own issue statuses to Dev Lake's standard statuses for every issue type</p>
//             </div>

//             <div className='jiraFormContainer'>

//               {(statusMappings && statusMappings.length > 0) && statusMappings.map((statusMapping, i) =>
//                 <div key={statusMapping.key} className='jiraFormContainerItem'>
//                   <p>Mapping {statusMapping.type} </p>
//                   <div>
//                     <MappingTagStatus
//                       reqValue={statusMapping.mapping.Rejected || []}
//                       resValue={statusMapping.mapping.Resolved || []}
//                       envName={statusMapping.key}
//                       clearBtnReq={<ClearButton onClick={() => setStatusMapping(statusMapping.key, [], 'Rejected')} />}
//                       clearBtnRes={<ClearButton onClick={() => setStatusMapping(statusMapping.key, [], 'Resolved')} />}
//                       onChangeReq={values => setStatusMapping(statusMapping.key, values, 'Rejected')}
//                       onChangeRes={values => setStatusMapping(statusMapping.key, values, 'Resolved')}
//                       className='mappingTagStatus'
//                     />
//                   </div>
//                 </div>
//               )}
//               <Button icon='add' onClick={() => setCustomStatusOverlay(true)} className='addNewStatusBtn'>Add New</Button>

//               <Dialog
//                 icon='diagram-tree'
//                 onClose={() => setCustomStatusOverlay(false)}
//                 title='Add a New Status Mapping'
//                 isOpen={customStatusOverlay}
//                 onOpened={() => setCustomStatusName('')}
//                 autoFocus={false}
//                 className='customStatusDialog'
//               >
//                 <div className={Classes.DIALOG_BODY}>
//                   <form onSubmit={() => addStatusMapping}>
//                     <FormGroup
//                       className='formGroup customStatusFormGroup'
//                     >
//                       <InputGroup
//                         id='custom-status'
//                         placeholder='Enter custom status name'
//                         onChange={(e) => setCustomStatusName(e.target.value)}
//                         className='customStatusInput'
//                         autoFocus={true}
//                       />
//                       <Button
//                         icon='add'
//                         onClick={() => addStatusMapping}
//                         className='addNewStatusBtnDialog'
//                         onSubmit={(e) => e.preventDefault()}
//                       >
//                         Add New
//                       </Button>
//                     </FormGroup>
//                   </form>
//                 </div>
//               </Dialog>

//             </div>

//             <div className='headlineContainer'>
//               <h3 className='headline'>Additional Customization Settings</h3>
//               <p className='description'>Additional Jira settings</p>
//             </div>

//             <div className='formContainer'>
//               <FormGroup
//                 inline={true}
//                 labelFor='jira-epic-key'
//                 helperText='JIRA_ISSUE_EPIC_KEY_FIELD'
//                 className='formGroup'
//                 contentClassName='formGroup'
//               >
//                 <Label>
//                   Issue&nbsp;Epic&nbsp;Key&nbsp;Field

//                   <div>
//                     <Tooltip content='Get help with Issue Epic Key Field' position={Position.TOP}>
//                       <a
//                         href='https://github.com/merico-dev/lake/tree/main/plugins/jira#set-jira-custom-fields'
//                         rel='noreferrer'
//                         target='_blank'
//                         className='helpIcon'
//                       >
//                         <Icon icon='help' size={15} />
//                       </a>
//                     </Tooltip>
//                   </div>

//                   <Tooltip content='Your custom epic key field' position={Position.TOP}>
//                     <InputGroup
//                       id='jira-epic-key'
//                       placeholder='Enter Jira epic key field'
//                       defaultValue={jiraIssueEpicKeyField}
//                       onChange={(e) => setJiraIssueEpicKeyField(e.target.value)}
//                       className='helperInput'
//                     />
//                   </Tooltip>
//                 </Label>
//               </FormGroup>
//             </div>

//             <div className='formContainer'>
//               <FormGroup
//                 inline={true}
//                 labelFor='jira-storypoint-field'
//                 helperText='JIRA_ISSUE_STORYPOINT_FIELD'
//                 className='formGroup'
//                 contentClassName='formGroup'
//               >
//                 <Tooltip content='Your custom story point key field' position={Position.TOP}>
//                   <Label>
//                     Issue&nbsp;Storypoint&nbsp;Field

//                     <div>
//                       <Tooltip content='Get help with Issue Story Point Field' position={Position.TOP}>
//                         <a
//                           href='https://github.com/merico-dev/lake/tree/main/plugins/jira#set-jira-custom-fields'
//                           target='_blank'
//                           className='helpIcon'
//                           rel='noreferrer'
//                         >
//                           <Icon icon='help' size={15} />
//                         </a>
//                       </Tooltip>
//                     </div>
//                     <InputGroup
//                       id='jira-storypoint-field'
//                       placeholder='Enter Jira Story Point Field'
//                       defaultValue={jiraIssueStoryPointField}
//                       onChange={(e) => setJiraIssueStoryPointField(e.target.value)}
//                       className='input'
//                     />
//                   </Label>
//                 </Tooltip>
//               </FormGroup>
//             </div>

//             <div className='formContainer'>
//               <FormGroup
//                 inline={true}
//                 labelFor='jira-storypoint-coef'
//                 helperText='JIRA_ISSUE_STORYPOINT_COEFFICIENT'
//                 className='formGroup'
//                 contentClassName='formGroup'
//               >
//                 <Tooltip content='Your custom story point coefficent (optional)' position={Position.TOP}>
//                   <Label>
//                     Issue&nbsp;Storypoint&nbsp;Coefficient <span className='requiredStar'>*</span>
//                     <InputGroup
//                       id='jira-storypoint-coef'
//                       placeholder='Enter Jira Story Point Coefficient'
//                       defaultValue={jiraIssueStoryCoefficient}
//                       onChange={(e) => setJiraIssueStoryCoefficient(e.target.value)}
//                       className='input'
//                     />
//                   </Label>
//                 </Tooltip>
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
