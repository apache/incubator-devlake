/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
import React from 'react'
import { CSSTransition } from 'react-transition-group'
import {
  Icon,
  Button,
  Colors,
  Intent,
  Classes,
  Popover,
  Spinner,
  Position
} from '@blueprintjs/core'

const PipelineIndicator = (props) => {
  const {
    pipeline,
    graphsUrl = '#',
    isVisible = true,
    onFetch = () => {},
    onCancel = () => {},
    onView = () => {},
    onRetry = () => {}
  } = props

  return (
    <>
      <CSSTransition
        in={pipeline && pipeline.ID !== null && isVisible}
        timeout={300}
        classNames='lastrun-module'
        unmountOnExit
      >
        <div
          className='trigger-module-lastrun'
          style={{
            position: 'fixed',
            borderRadius: '40px',
            backgroundColor: '#ffffff',
            width: '40px',
            height: '40px',
            right: '30px',
            bottom: '20px',
            zIndex: 500,
            boxShadow: '0px 0px 6px rgba(0, 0, 0, 0.25)',
            display: 'flex',
            alignItems: 'center',
            alignContent: 'center',
            justifyContent: 'center',
            cursor: 'pointer'
          }}
        >
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              width: '40px',
              height: '40px',
              position: 'relative'
            }}
          >
            <Popover
              key='popover-lastrun-info'
              className='popover-trigger-lastrun'
              popoverClassName='popover-lastrun'
              position={Position.LEFT_BOTTOM}
              autoFocus={false}
              enforceFocus={false}
              usePortal={false}
              onOpening={() => onFetch(pipeline.ID, true)}
              // onOpened=
            >
              <>
                <Spinner
                  value={pipeline.status === 'TASK_COMPLETED' ? 100 : null}
                  className='lastrun-spinner'
                  intent={
                    pipeline.status === 'TASK_COMPLETED'
                      ? Intent.WARNING
                      : Intent.PRIMARY
                  }
                  size={40}
                  style={{}}
                />
                {(() => {
                  switch (pipeline.status) {
                    case 'TASK_COMPLETED':
                      return <Icon icon='tick-circle' size={40} />
                    case 'TASK_FAILED':
                      return <Icon icon='error' size={40} />
                    case 'TASK_RUNNING':
                    default:
                      return <Icon icon='refresh' size={40} />
                  }
                })()}
                {/* <Icon
                  icon={pipeline.status === 'TASK_COMPLETED'
                    ? <PipelineCompleteIcon width={40} height={40} style={{ marginTop: '3px', display: 'flex', alignSelf: 'center' }} />
                    : <pipelineningIcon width={40} height={40} style={{ marginTop: '3px', display: 'flex', alignSelf: 'center' }} />}
                  size={40}
                /> */}
              </>
              <>
                <div
                  style={{
                    fontSize: '12px',
                    padding: '12px',
                    minWidth: '420px',
                    maxWidth: '420px',
                    maxHeight: '400px',
                    overflow: 'hidden',
                    overflowY: 'auto'
                  }}
                >
                  <h3
                    onClick={() => onView(pipeline.ID)}
                    className='group-header'
                    style={{
                      marginTop: '0',
                      marginBottom: '6px',
                      textOverflow: 'ellipsis',
                      overflow: 'hidden',
                      whiteSpace: 'nowrap'
                    }}
                  >
                    <Icon
                      icon={
                        pipeline.status === 'TASK_FAILED'
                          ? 'warning-sign'
                          : 'help'
                      }
                      size={16}
                      style={{ marginRight: '5px' }}
                    />{' '}
                    {pipeline.name || 'Last Pipeline Run'}
                  </h3>
                  <p
                    style={{
                      fontSize: '11px',
                      color:
                        pipeline.status === 'TASK_FAILED'
                          ? Colors.RED4
                          : Colors.DARK_GRAY4
                    }}
                  >
                    {pipeline.message.length > 150
                      ? `${pipeline.message.slice(0, 150)}...`
                      : pipeline.message}
                  </p>
                  <div
                    style={{
                      display: 'flex',
                      width: '100%',
                      justifyContent: 'space-between'
                    }}
                  >
                    <div>
                      <label>
                        <strong>Pipeline ID</strong>
                      </label>
                      <div style={{ fontSize: '13px', fontWeight: 800 }}>
                        {pipeline.ID}
                      </div>
                    </div>
                    <div style={{ padding: '0 12px' }}>
                      <label>
                        <strong>Tasks</strong>
                      </label>
                      <div style={{ fontSize: '13px' }}>
                        {pipeline.finishedTasks}/{pipeline.totalTasks}
                      </div>
                    </div>
                    <div>
                      <label>
                        <strong>Status</strong>
                      </label>
                      <div style={{ fontSize: '13px' }}>{pipeline.status}</div>
                    </div>
                    <div
                      style={{
                        paddingLeft: '10px',
                        justifyContent: 'flex-end',
                        alignSelf: 'flex-end'
                      }}
                    >
                      {pipeline.status === 'TASK_COMPLETED' && (
                        <a
                          className='bp3-button bp3-intent-primary bp3-small'
                          href={graphsUrl}
                          target='_blank'
                          rel='noreferrer'
                          style={{
                            backgroundColor: '#3bd477',
                            color: '#ffffff'
                          }}
                        >
                          <Icon icon='doughnut-chart' size={13} />{' '}
                          <span className='bp3-button-text'>Graphs</span>
                        </a>
                      )}
                      {pipeline.status === 'TASK_RUNNING' && (
                        <Button
                          className={`btn-cancel-pipeline ${Classes.POPOVER_DISMISS}`}
                          small
                          icon='stop'
                          text='CANCEL'
                          intent='primary'
                          onClick={() => onCancel(pipeline.ID)}
                        />
                      )}
                      {pipeline.status === 'TASK_FAILED' && (
                        <Button
                          className={`btn-retry-pipeline ${Classes.POPOVER_DISMISS}`}
                          intent='danger'
                          icon='reset'
                          text='RETRY'
                          style={{ color: '#ffffff' }}
                          onClick={() => onRetry(pipeline)}
                          small
                        />
                      )}
                      <Button
                        minimal
                        className={`btn-ok ${Classes.POPOVER_DISMISS}`}
                        small
                        text='OK'
                        style={{ marginLeft: '3px' }}
                      />
                    </div>
                  </div>
                  <div
                    style={{
                      paddingTop: '7px',
                      borderTop: '1px solid #f5f5f5',
                      marginTop: '14px'
                    }}
                  >
                    {pipeline?.tasks &&
                      pipeline.tasks.map((t, tIdx) => (
                        <div
                          className='pipeline-task-'
                          key={`pipeline-task-key-${tIdx}`}
                          style={{
                            display: 'flex',
                            padding: '4px 6px',
                            justifyContent: 'space-between'
                          }}
                        >
                          <div style={{ paddingRight: '8px' }}>
                            <Spinner
                              className='mini-task-spinner'
                              size={14}
                              intent={
                                t.status === 'TASK_COMPLETED'
                                  ? 'success'
                                  : 'warning'
                              }
                              value={
                                t.status === 'TASK_COMPLETED' ? 1 : t.progress
                              }
                            />
                          </div>
                          <div style={{ padding: '0 8px', width: '100%' }}>
                            <strong
                              style={{
                                overflow: 'hidden',
                                textOverflow: 'ellipsis',
                                whiteSpace: 'nowrap'
                              }}
                            >
                              {t.plugin}
                            </strong>
                            {t.status === 'TASK_COMPLETED' && (
                              <Icon
                                icon='small-tick'
                                size={14}
                                color={Colors.GREEN5}
                                style={{ marginLeft: '5px' }}
                              />
                            )}
                            {t.status === 'TASK_FAILED' && (
                              <Icon
                                icon='warning-sign'
                                size={11}
                                color={Colors.RED5}
                                style={{
                                  marginLeft: '5px',
                                  marginBottom: '3px'
                                }}
                              />
                            )}
                          </div>
                          <div
                            style={{
                              padding: '0',
                              minWidth: '80px',
                              textAlign: 'right'
                            }}
                          >
                            <strong>
                              {Number(t.spentSeconds / 60).toFixed(2)}mins
                            </strong>
                          </div>
                          <div
                            style={{
                              padding: '0 8px',
                              minWidth: '100px',
                              textAlign: 'right'
                            }}
                          >
                            <span color={Colors.GRAY5}>
                              {Number(
                                t.status === 'TASK_COMPLETED'
                                  ? 100
                                  : (t.progress / 1) * 100
                              ).toFixed(2)}
                              %
                            </span>
                          </div>
                          <div />
                        </div>
                      ))}
                  </div>
                </div>
              </>
            </Popover>
          </div>
        </div>
      </CSSTransition>
    </>
  )
}

export default PipelineIndicator
