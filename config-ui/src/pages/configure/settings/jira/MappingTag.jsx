import React from 'react'
import { FormGroup, Label, Tag, TagInput } from '@blueprintjs/core'

const MappingTag = ({ labelIntent, labelName, onChange, rightElement, helperText, typeOrStatus, values, placeholderText }) => {
  return (
    <>
      <div className='formContainer'>
        <FormGroup
            // disabled={isTesting || isSaving}
          label=''
          inline={true}
          labelFor='jira-issue-type-mapping'
          helperText={helperText}
          className='formGroup'
          contentClassName='formGroupContent'
        >
          <Label style={{ display: 'inline' }}>
            <span style={{ marginRight: '10px' }}><Tag intent={labelIntent}>{labelName}</Tag></span>
          </Label>
          <TagInput
            placeholder={placeholderText}
            values={values || []}
            fill={true}
            onChange={value => setTimeout(() => onChange([...new Set(value)]), 0)}
            addOnPaste={true}
            rightElement={rightElement}
            onKeyDown={e => e.key === 'Enter' && e.preventDefault()}
            className='tagInput'
          />
        </FormGroup>
      </div>
    </>
  )
}

export default MappingTag
