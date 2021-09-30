import { Tooltip, Position, FormGroup, Label, Tag, TagInput } from '@blueprintjs/core'
import styles from '../../../styles/Home.module.css'

const MappingTag = ({labelIntent, labelName, onChange, rightElement, helperText, typeOrStatus, values}) => {
  return <>
    <p>Issue {typeOrStatus === 'type' ? 'types' : 'statuses' } mapped to&nbsp;&nbsp;<Tag intent={labelIntent}>{labelName}</Tag></p>

    <div className={styles.formContainer}>
      <FormGroup
        inline={true}
        labelFor="jira-issue-type-mapping"
        helperText={helperText}
        className={styles.formGroup}
        contentClassName={styles.formGroup}
      >
        <Tooltip content={`Map custom Jira types to main ${labelName} status`} position={Position.TOP}>
          <Label>
          <TagInput
            placeholder="Add Tags..."
            values={values || []}
            fill={true}
            onChange={value => onChange([...new Set(value)])}
            addOnPaste={true}
            rightElement={rightElement}
            onKeyDown={e => e.key === 'Enter' && e.preventDefault()}
            className={styles.tagInput}
          />
          </Label>
        </Tooltip>
      </FormGroup>
    </div>
  </>
}

export default MappingTag
