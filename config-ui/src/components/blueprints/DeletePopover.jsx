import React from 'react'
import {
  Button,
  Intent,
  Colors,
  Classes
} from '@blueprintjs/core'
const DeletePopover = (props) => {
  const {
    activeBlueprint,
    onCancel = () => {},
    onConfirm = () => {},
    isRunning = false
  } = props
  return (
    <>
      <div style={{ padding: '10px', fontSize: '10px', maxWidth: '220px' }}>
        <h3 style={{ margin: '0 0 5px 0', color: Colors.RED3 }}>Delete {activeBlueprint?.name}?</h3>
        <p><strong>Are you sure? This Blueprint will be removed, all pipelines will be stopped.</strong></p>
        <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
          <Button
            className={Classes.POPOVER_DISMISS}
            intent={Intent.NONE}
            text='CANCEL'
            small style={{ marginRight: '5px' }}
            onClick={() => onCancel(activeBlueprint)}
            disabled={isRunning}
          />
          <Button disabled={isRunning} intent={Intent.DANGER} text='YES' small onClick={() => onConfirm(activeBlueprint)} />
        </div>
      </div>
    </>
  )
}

export default DeletePopover
