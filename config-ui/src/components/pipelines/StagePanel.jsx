import React from 'react'
import { CSSTransition } from 'react-transition-group'
import {
  Card,
  Button, Icon,
  ButtonGroup,
  Elevation,
  Colors,
  // Alignment, Classes, Spinner
} from '@blueprintjs/core'
import { ReactComponent as PipelineRunningIcon } from '@/images/synchronize.svg'
import { ReactComponent as PipelineFailedIcon } from '@/images/no-synchronize.svg'
import { ReactComponent as PipelineCompleteIcon } from '@/images/check-circle.svg'

const StagePanel = (props) => {
  const { activePipeline, pipelineReady = false, stages, activeStageId = 1 } = props

  return (
    <>
      <CSSTransition
        in={pipelineReady}
        timeout={350}
        classNames='activity-panel'
      >
        <Card
          elevation={Elevation.TWO}
          style={{
            display: 'flex',
            width: '100%',
            justifySelf: 'flex-start',
            marginBottom: '8px',
            padding: 0,
            backgroundColor: activePipeline.status === 'TASK_COMPLETED' ? 'rgba(245, 255, 250, 0.99)' : 'inherit',
            overflow: 'hidden',
            whiteSpace: 'nowrap',
            textOverflow: 'ellipsis'
          }}
        >

          <ButtonGroup style={{ backgroundColor: 'transparent' }}>
            <Button minimal active style={{ backgroundColor: '#eeeeee' }}>
              <h3 style={{ margin: 0, fontSize: '20px', display: 'flex' }}>
                {/* Stage 1 */}
                {(() => {
                  let statusIcon = null
                  switch (activePipeline.status) {
                    case 'TASK_COMPLETED':
                      statusIcon = (
                        <Icon
                          icon={<PipelineCompleteIcon
                            width={24} height={24} style={{
                              margin: '0 0 0 10px',
                              display: 'flex',
                              alignSelf: 'center'
                            }}
                                />}
                          size={24}
                        />
                      )
                      break
                    case 'TASK_FAILED':
                      statusIcon = (
                        <Icon
                          icon={<PipelineFailedIcon
                            width={24}
                            height={24} style={{
                              margin: '0 0 0 10px',
                              display: 'flex',
                              alignSelf: 'center'
                            }}
                                />}
                          size={24}
                        />
                      )
                      break
                    case 'TASK_RUNNING':
                    default:
                      statusIcon = (
                        <Icon
                          icon={<PipelineRunningIcon
                            width={24}
                            height={24} style={{
                              margin: '0 0 0 10px',
                              display: 'flex',
                              alignSelf: 'center'
                            }}
                                />}
                          size={24}
                        />
                      )
                      break
                  }
                  return statusIcon
                })()}
              </h3>
            </Button>
            {/* @todo: re-active "stage" ux in a future release */}
            {Object.keys(stages).length > 0 && (
              <>
                {Object.keys(stages).map((s, sIdx) => (
                  // <Button
                  //   minimal style={{
                  //     backgroundColor: '#eeeeee',
                  //     color: '#cccccc',
                  //     fontSize: '35px',
                  //     lineHeight: '20px',
                  //     padding: 0,
                  //     fontWeight: 100,
                  //   }}
                  // >/
                  // </Button>
                  <Button
                    key={`stage-btn-key${sIdx}`}
                    disabled={activeStageId !== (sIdx + 1)}
                    minimal
                    style={{
                      position: 'relative',
                      backgroundColor: '#eeeeee',
                      paddingRight: '50px',
                    }}
                    rightIcon={
                      sIdx !== (Object.keys(stages).length - 1)
                        ? (
                          <Button
                            minimal style={{
                              backgroundColor: '#eeeeee',
                              color: '#cccccc',
                              fontSize: '35px',
                              lineHeight: '20px',
                              padding: 0,
                              margin: 0,
                              position: 'absolute',
                              right: 0,
                              fontWeight: 100,
                            }}
                          >/
                          </Button>
                          )
                        : null
                    }
                  >
                    <h3 style={{ margin: 0, fontSize: '20px', color: activeStageId === (sIdx + 1) ? Colors.BLACK : Colors.GRAY3 }}>
                      Stage {sIdx + 1}
                    </h3>
                  </Button>
                ))}

                <Button
                  className='btn-stage-endcap'
                  minimal
                  style={{
                    marginLeft: '1px',
                    background: '#ffffff!!important',
                    width: 0,
                    height: 0,
                    borderTop: '16px solid transparent',
                    borderBottom: '16px solid transparent',
                    borderLeft: '16px solid #eeeeee',
                    pointerEvents: 'none'
                  }}
                />
              </>
            )}
          </ButtonGroup>
          <h3 style={{
            textTransform: 'uppercase',
            lineHeight: '33px',
            margin: 0,
            fontFamily: 'Montserrat',
            fontWeight: 800,
            letterSpacing: '2px'
          }}
          >Finished Tasks &middot; <span style={{ color: Colors.GREEN5 }}>{activePipeline.finishedTasks}</span>
            <em style={{ color: '#dddddd', padding: '0 4px', textTransform: 'lowercase' }}>/</em>{activePipeline.totalTasks}
          </h3>
          <div style={{ display: 'flex', fontSize: '16px', fontWeight: 700, marginLeft: 'auto', lineHeight: '33px', padding: '0 10px' }}>

            {Number((activePipeline.finishedTasks / activePipeline.totalTasks) * 100).toFixed(1)}%
          </div>
        </Card>
      </CSSTransition>
    </>
  )
}

export default StagePanel
