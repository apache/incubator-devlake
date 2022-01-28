import React from 'react'
import { Providers } from '@/data/Providers'

const StageTaskCaption = (props) => {
  const { task, options } = props

  return (
    <span
      className='task-module-caption'
      style={{
        opacity: 0.4,
        display: 'block',
        width: '90%',
        fontSize: '9px',
        overflow: 'hidden',
        whiteSpace: 'nowrap',
        textOverflow: 'ellipsis'
      }}
    >
      {(task.plugin === Providers.GITLAB || task.plugin === Providers.JIRA) && (<>ID {options.projectId || options.boardId}</>)}
      {task.plugin === Providers.GITHUB && (<>@{options.owner}/{options.repositoryName}</>)}
      {task.plugin === Providers.JENKINS && (<>Task #{task.ID}</>)}
      {(task.plugin === Providers.GITEXTRACTOR || task.plugin === Providers.REFDIFF) && (<>{options.repoId || `ID ${task.ID}`}</>)}
    </span>
  )
}

export default StageTaskCaption
