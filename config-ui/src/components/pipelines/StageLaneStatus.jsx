import React from 'react'
// import { CSSTransition } from 'react-transition-group'
import {
  // Classes,
  Icon,
  Spinner,
  Colors,
  Tooltip,
  Position,
  Intent,
  Card,
  Elevation,
  H4,
  ProgressBar
  // Alignment
} from '@blueprintjs/core'
import dayjs from '@/utils/time'

const StageLaneStatus = (props) => {
  const {
    stage,
    stages = [],
    sK = 1,
    isStageActive = () => {},
    isStagePending = () => {},
    isStageCompleted = () => {},
    isStageFailed = () => {},
    calculateStageLaneProgress = () => {},
    getRunningTaskCount = () => {},
    getTotalTasksCount = () => {},
    getCompletedTaskCount = () => {}

  } = props

  return (
    <>
      <div
        className='stage-footer' style={{
          padding: '5px 10px',
          marginBottom: '4px',
          alignSelf: 'flex-end',
          marginTop: 'auto',
          fontSize: '10px',
          color: '#999999',
          alignContent: 'flex-end',
          textAlign: 'right'
        }}
      >
        <div
          className='stage-status'
          style={{
            fontFamily: 'Montserrat',
            fontWeight: 900,
            fontSize: '10px',
            letterSpacing: '1px',
            color: isStageCompleted(sK) ? Colors.GREEN5 : Colors.GRAY5
          }}
        >
          {isStageActive(sK) && <span style={{ color: Colors.BLACK }}> ACTIVE</span>}
          {isStageCompleted(sK) && <span style={{ color: Colors.GREEN5 }}>COMPLETED</span>}
          {isStageFailed(sK) && <span style={{ color: Colors.RED5 }}>FAILED</span>}
          {isStagePending(sK) && <>WAITING</>}
          <span style={{ fontWeight: 500, color: isStageActive(sK) ? '#000' : 'inherit', opacity: 0.6 }}>
            {' '}&middot;{' '}
            {isStageCompleted(sK) ? getCompletedTaskCount(stage) : getRunningTaskCount(stage)}/{getTotalTasksCount(stage)}
          </span>
        </div>
        <div className='stage-caption'>
          {isStageActive(sK) && <>Stage Running</>}
          {/* {isStageFailed(sK) && <>Stage Failed</>} */}
          {(isStageCompleted(sK) || isStageFailed(sK)) && <>{dayjs(stage.UpdatedAt).from(stage.CreatedAt, true)}</>}
          {isStagePending(sK) && <><Icon icon='more' color={Colors.GRAY5} size={12} /></>}
        </div>
      </div>
      {isStageActive(sK) && (
        <ProgressBar
          className='stage-lane-progressbar'
          stripes={true}
          intent={Intent.SUCCESS} value={calculateStageLaneProgress(stage)} style={{ borderRadius: 0 }}
        />
      )}
    </>
  )
}

export default StageLaneStatus
