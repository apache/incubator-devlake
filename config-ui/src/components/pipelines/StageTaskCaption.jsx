import React, { useState, useEffect } from 'react'
import { Providers } from '@/data/Providers'
import {
  Icon,
  Spinner,
  Colors,
  Tooltip,
  Position,
  Intent,
} from '@blueprintjs/core'

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
      {task.plugin !== Providers.GITHUB && (<>ID {options.projectId || options.boardId}</>)}
      {task.plugin === Providers.GITHUB && (<>@{options.owner}/{options.repositoryName}</>)}
    </span>
  )
}

export default StageTaskCaption
