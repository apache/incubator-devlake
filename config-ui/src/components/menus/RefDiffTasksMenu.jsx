import React from 'react'
import {
  Menu,
  MenuItem,
  Checkbox,
  Intent
} from '@blueprintjs/core'

const RefDiffTasksMenu = (props) => {
  const {
    tasks = [
      { task: 'calculateCommitsDiff', label: 'Calculate Commits Diff' },
      { task: 'calculateIssuesDiff', label: 'Calculate Issues Diff' }],
    selected = [],
    onSelect = () => {}
  } = props
  return (
    <Menu className='tasks-menu refdiff-tasks-menu' minimal='true'>
      <label style={{
        fontSize: '10px',
        fontWeight: 800,
        fontFamily: '"Montserrat", sans-serif',
        textTransform: 'uppercase',
        padding: '6px 8px',
        display: 'block'
      }}
      >AVAILABLE PLUGIN TASKS
      </label>
      {tasks.map((t, tIdx) => (
        <MenuItem
          key={`refdiff-item-task-key-${tIdx}`}
          intent={Intent.DANGER}
          minimal='true'
          active={Boolean(selected.includes(t.task))}
          // onClick={(e) => e.preventDefault()}
          data-task={t}
          text={(
            <>
              <Checkbox
                intent={Intent.WARNING}
                checked={Boolean(selected.includes(t.task))}
                label={t.label}
                onChange={(e) => onSelect(e, t)}
              />
            </>
          )}
        />
      ))}
    </Menu>
  )
}

export default RefDiffTasksMenu
