import React from 'react'
// import { CSSTransition } from 'react-transition-group'
import {
  Classes,
  Drawer,
  DrawerSize,
  Card,
  Elevation,
  Position,
} from '@blueprintjs/core'

const CodeInspector = (props) => {
  const { activePipeline, isOpen, onClose, hasBackdrop = true } = props

  return (
    <Drawer
      className='drawer-json-inspector'
      icon='code'
      onClose={() => onClose(false)}
      title={`Inspect RUN #${activePipeline.ID}`}
      position={Position.RIGHT}
      size={DrawerSize.SMALL}
      autoFocus
      canEscapeKeyClose
      canOutsideClickClose
      enforceFocus
      hasBackdrop={hasBackdrop}
      isOpen={isOpen}
      usePortal
    >
      <div className={Classes.DRAWER_BODY}>
        <div className={Classes.DIALOG_BODY}>
          <h3 style={{ margin: 0, padding: '8px 0' }}>
            <span style={{ float: 'right', fontSize: '9px', color: '#aaaaaa' }}>application/json</span> JSON RESPONSE
          </h3>
          <p>
            If you are submitting a
            <strong> Bug-Report</strong> regarding a Pipeline Run, include the output below for better debugging.
          </p>
          <div className='formContainer'>
            <Card
              className='code-inspector-card'
              interactive={false}
              elevation={Elevation.ZERO}
              style={{ padding: '6px 12px', minWidth: '320px', width: '100%', maxWidth: '601px', marginBottom: '20px', overflow: 'auto' }}
            >

              <code>
                <pre style={{ fontSize: '10px' }}>
                  {JSON.stringify(activePipeline, null, '  ')}
                </pre>
              </code>
            </Card>
          </div>
        </div>
      </div>
    </Drawer>

  )
}

export default CodeInspector
