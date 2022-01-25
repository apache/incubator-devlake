import React, { useState, useEffect } from 'react'
// import { CSSTransition } from 'react-transition-group'
import {
  Providers,
  ProviderLabels
} from '@/data/Providers'
import {
  Card,
  Elevation,
} from '@blueprintjs/core'
import dayjs from '@/utils/time'
import StageTaskName from '@/components/pipelines/StageTaskName'
import StageTaskIndicator from '@/components/pipelines/StageTaskIndicator'
import StageTaskCaption from '@/components/pipelines/StageTaskCaption'

const StageTask = (props) => {
  const {
    stages = [],
    task,
    sK,
    sIdx,
  } = props

  const [taskModuleOpened, setTaskModuleOpened] = useState(null)

  const generateStageTaskCssClasses = () => {
    return `pipeline-task-module task-${task.status.split('_')[1].toLowerCase()} ${task.ID === taskModuleOpened?.ID ? 'active' : ''}`
  }

  const determineCardElevation = (status, isElevated = false) => {
    let elevation = Elevation.ONE
    if (status === 'TASK_RUNNING' && isElevated) {
      elevation = Elevation.THREE
    } else if (status === 'TASK_RUNNING' && !isElevated) {
      elevation = Elevation.TWO
    } else if (isElevated) {
      elevation = Elevation.THREE
    } else {
      elevation = Elevation.ONE
    }
    return elevation
  }

  return (
    <>
      <Card
        elevation={determineCardElevation(task.status, taskModuleOpened !== null)}
        className={generateStageTaskCssClasses()}
        onClick={() => setTaskModuleOpened(task)}
        style={{

        }}
      >
        <StageTaskIndicator task={task} />
        <div
          className='task-module-name'
          style={{
            flex: 1
          }}
        >
          <div style={{ padding: '4px 2px 4px 0' }}>
            <StageTaskName task={task} showDetails={taskModuleOpened} onClose={() => setTaskModuleOpened(null)} />
            <StageTaskCaption task={task} options={task.options} />
          </div>
        </div>
        <div
          className='task-module-progress'
          style={{
            fontWeight: task.status === 'TASK_COMPLETED' ? 700 : 300
          }}
        >
          {Number(task.status === 'TASK_COMPLETED' ? 100 : (task.progress / 1) * 100).toFixed(0)}%
        </div>
      </Card>
    </>
  )
}

export default StageTask
