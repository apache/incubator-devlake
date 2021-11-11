// import React from 'react'
// import { Tooltip, Position, FormGroup, Label, Tag, TagInput } from '@blueprintjs/core'

// const MappingTag = ({ labelIntent, labelName, onChange, rightElement, helperText, typeOrStatus, values, placeholderText }) => {
//   return (
//     <>
//       <p>Issue {typeOrStatus === 'type' ? 'types' : 'statuses'} mapped to&nbsp;&nbsp;<Tag intent={labelIntent}>{labelName}</Tag></p>

//       <div className='formContainer'>
//         <FormGroup
//           inline={true}
//           labelFor='jira-issue-type-mapping'
//           helperText={helperText}
//           className='formGroup'
//           contentClassName='formGroup'
//         >
//           <Tooltip content={`Map custom Jira types to main ${labelName} status`} position={Position.TOP}>
//             <Label>
//               <TagInput
//                 placeholder={placeholderText}
//                 values={values || []}
//                 fill={true}
//                 onChange={value => setTimeout(() => onChange([...new Set(value)]), 0)}
//                 addOnPaste={true}
//                 rightElement={rightElement}
//                 onKeyDown={e => e.key === 'Enter' && e.preventDefault()}
//                 className='tagInput'
//               />
//             </Label>
//           </Tooltip>
//         </FormGroup>
//       </div>
//     </>
//   )
// }

// export default MappingTag
