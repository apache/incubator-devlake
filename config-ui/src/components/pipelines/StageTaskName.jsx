import React, { useState, useEffect, useRef } from 'react'
import {
  Providers,
  ProviderLabels,
  ProviderIcons
} from '@/data/Providers'
import {
  Icon,
  Colors,
  Position,
  Popover,
  TextArea,
  H3
} from '@blueprintjs/core'
import dayjs from '@/utils/time'

const StageTaskName = (props) => {
  const {
    task,
    showDetails = null,
    onClose = () => {}
  } = props

  const popoverTriggerRef = useRef()

  useEffect(() => {
    if (showDetails !== null && popoverTriggerRef.current) {
      popoverTriggerRef.current.click()
    }
  }, [showDetails])

  return (
    <>
      <Popover
        className='trigger-pipeline-activity-help'
        popoverClassName='popover-help-pipeline-activity'
        // isOpen={showDetails && showDetails.ID === task.ID}
        onClosed={onClose}
        position={Position.RIGHT}
        autoFocus={false}
        enforceFocus={false}
        usePortal={true}
      >
        <span className='task-plugin-text' ref={popoverTriggerRef}>{task.plugin}</span>
        <>
          <div style={{ textShadow: 'none', fontSize: '12px', padding: '12px', maxWidth: '360px' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between' }}>
              <div style={{
                marginBottom: '10px',
                color: Colors.GRAY2,
                fontWeight: 700,
                fontSize: '14px',
                fontFamily: '"Montserrat", sans-serif',
                maxWidth: '80%'
              }}
              >
                <H3 style={{
                  margin: 0,
                  fontFamily: '"JetBrains Mono", monospace',
                  color: Colors.BLACK,
                  textOverflow: 'ellipsis',
                  overflow: 'hidden',
                  whiteSpace: 'nowrap',
                }}
                >
                  {task.plugin === 'jenkins' && (<>{ProviderLabels.JENKINS}</>)}
                  {task.plugin !== 'github' && task.plugin !== 'jenkins' && (<>ID {task.options.projectId || task.options.boardId}</>)}
                  {task.plugin === 'github' && task.plugin !== 'jenkins' && (<>@{task.options.owner}/{task.options.repositoryName}</>)}
                </H3>
                {ProviderLabels[task.plugin.toUpperCase()] || 'System Task'}<br />
              </div>
              <div style={{ padding: '0 0 10px 20px' }}>
                {ProviderIcons[task.plugin.toLowerCase()](32, 32)}
              </div>
            </div>

            <div style={{ fontSize: '10px' }}>
              <div style={{ display: 'flex', flexDirection: 'row', justifyContent: 'space-between' }}>
                <div>
                  <label style={{ color: Colors.GRAY2 }}>ID</label><br />
                  <span>{task.ID}</span>
                </div>
                <div style={{ marginLeft: '20px' }}>
                  <label style={{ color: Colors.GRAY2 }}>Created</label><br />
                  <span>{dayjs(task.CreatedAt).format('L LT')}</span>
                </div>
                {task.finishedAt && (
                  <div style={{ marginLeft: '20px' }}>
                    <label style={{ color: Colors.GRAY2 }}>Finished</label><br />
                    <span>{dayjs(task.finishedAt).format('L LT')}</span>
                  </div>
                )}
              </div>
              <div style={{ marginTop: '6px' }}>
                <label style={{ color: Colors.GRAY2 }}>Status</label><br />
                <strong>{task.status}</strong>{' '}
                <strong
                  className='bp3-tag bp3-intent-primary'
                  style={{ minHeight: '16px', fontSize: '10px', padding: '2px 6px', borderRadius: '6px' }}
                >PROGRESS {task.progress * 100}%
                </strong>
              </div>
              <div style={{ marginTop: '6px' }}>
                <label style={{ color: Colors.GRAY2 }}>Options</label><br />
                <TextArea
                  readOnly
                  fill value={JSON.stringify(task.options)} style={{ fontSize: '10px', backgroundColor: '#f8f8f8', resize: 'none' }}
                />
                {/* <span>
                        <pre style={{ margin: 0 }}>
                          <code>
                            {JSON.stringify(task.options)}
                          </code>
                        </pre>
                      </span> */}
              </div>
              {task.message !== '' && (
                <div style={{ marginTop: '6px' }}>
                  <label style={{ color: Colors.DARK_GRAY1 }}>Message</label><br />
                  <span style={{ color: task.status === 'TASK_FAILED' ? Colors.RED3 : Colors.BLACK }}>
                    {task.status === 'TASK_FAILED' && (
                      <Icon
                        icon='warning-sign'
                        color={Colors.RED5}
                        size={10}
                        style={{ marginRight: '3px' }}
                      />
                    )}
                    {task.message}
                  </span>
                </div>
              )}
            </div>
          </div>
        </>
      </Popover>
    </>
  )
}

export default StageTaskName
