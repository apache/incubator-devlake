
import React, { useState, useEffect } from 'react'
import { CSSTransition } from 'react-transition-group'
import {
  Icon,
  Colors,
  Tooltip,
  Position,
  Intent,
  Card,
  Elevation,
  ProgressBar,
  H4,
  // Alignment
} from '@blueprintjs/core'
import dayjs from '@/utils/time'
import StageTask from '@/components/pipelines/StageTask'
import StageLaneStatus from '@/components/pipelines/StageLaneStatus'

const StageLane = (props) => {
  const { stages = [], sK = 1, sIdx } = props

  const [activeStage, setActiveStage] = useState(stages[sK])

  const isStageActive = (stageId) => {
    return stages[stageId].some(s => s.status === 'TASK_RUNNING')
  }

  const isStageCompleted = (stageId) => {
    return stages[stageId].every(s => s.status === 'TASK_COMPLETED')
  }

  const isStageFailed = (stageId) => {
    return !isStageActive(stageId) && stages[stageId].some(s => s.status === 'TASK_FAILED')
  }

  const isStagePending = (stageId) => {
    return stages[stageId].every(s => s.status === 'TASK_CREATED')
  }

  const generateStageCssClasses = (stageId) => {
    const classes = []
    if (isStageCompleted(stageId)) {
      classes.push('stage-completed')
    }
    if (isStageFailed(stageId)) {
      classes.push('stage-failed')
    }
    if (isStageActive(stageId)) {
      classes.push('stage-running')
    }
    if (isStagePending(stageId)) {
      classes.push('stage-created')
    }
    return classes.join(' ')
  }

  const calculateStageLaneProgress = (stageTasks) => {
    const completed = stageTasks.filter(s => s.status === 'TASK_COMPLETED')
    const remaining = stageTasks.filter(s => s.status !== 'TASK_COMPLETED')
    console.log('>>> STAGE LANE PROGRESS  = ', completed, remaining, completed.length / remaining.length)
    return completed.length / remaining.length
  }

  const getRunningTaskCount = (stageTasks) => {
    return stageTasks.filter(s => s.status === 'TASK_RUNNING').length
  }

  const getCompletedTaskCount = (stageTasks) => {
    return stageTasks.filter(s => s.status === 'TASK_COMPLETED').length
  }

  const getTotalTasksCount = (stageTasks) => {
    return stageTasks.length
  }

  useEffect(() => {
    setActiveStage(stages[sK])
    console.log('>> ACTIVE STAGE LANE', stages[sK])
  }, [stages, sK])

  return (
    <>
      <div
        // key={`stage-lane-key-${sIdx}`}
        className={`stage-lane ${generateStageCssClasses(sK)} ${isStageActive(sK) ? 'bp3-elevation-2' : ''}`}
        style={{
          position: 'relative',
          display: 'flex',
          flexDirection: 'column',
          flex: 1,
          justifyContent: 'center',
          alignContent: 'flex-start',
          alignItems: 'center',
          padding: '0'
        }}
      >
        {isStageActive(sK) && (
          <Icon
            icon='dot'
            color={Colors.GREEN5}
            size={14}
            style={{ position: 'absolute', display: 'inline-block', right: '5px', top: '5px' }}
          />
        )}
        {isStageFailed(sK) && (
          <Icon
            icon='warning-sign'
            color={Colors.RED5}
            size={10}
            style={{ position: 'absolute', display: 'inline-block', right: '5px', top: '5px' }}
          />
        )}
        {isStageCompleted(sK) && (
          <Icon
            icon='tick'
            color={Colors.GREEN5}
            size={12}
            style={{ position: 'absolute', display: 'inline-block', right: '5px', top: '5px' }}
          />
        )}
        <H4
          className='stage-title'
          style={{

          }}
        >
          Stage {sIdx + 1}
        </H4>
        {/* {sIdx} */}
        {stages[sK].map((t, tIdx) => (
          <CSSTransition
            key={`fx-key-stage-task-${tIdx}`}
            in={true}
            timeout={350}
            classNames='provider-datarow'
            unmountOnExit
          >
            <StageTask task={t} key={`stage-task-key-${tIdx}`} />
          </CSSTransition>
        ))}
        <StageLaneStatus
          sK={sK}
          stage={activeStage}
          stages={stages}
          isStageCompleted={isStageCompleted}
          isStagePending={isStagePending}
          isStageActive={isStageActive}
          isStageFailed={isStageFailed}
          calculateStageLaneProgress={calculateStageLaneProgress}
          getTotalTasksCount={getTotalTasksCount}
          getCompletedTaskCount={getCompletedTaskCount}
          getRunningTaskCount={getRunningTaskCount}
        />
      </div>
    </>
  )
}

export default StageLane
