import React from 'react'
import {
  Colors,
  Icon,
  Popover,
  PopoverInteractionKind,
  Intent,
  Position
} from '@blueprintjs/core'

const InputValidationError = (props) => {
  const { error, position = Position.TOP } = props
  return error
    ? (
      <div className='inline-input-error' style={{ outline: 'none', cursor: 'pointer', margin: '5px 5px 3px 5px' }}>
        <Popover
          position={position}
          usePortal={true}
          openOnTargetFocus={true}
          intent={Intent.WARNING}
          interactionKind={PopoverInteractionKind.HOVER_TARGET_ONLY}
        >
          <Icon icon='warning-sign' size={12} color={Colors.RED5} style={{ outline: 'none' }} />
          <div style={{ outline: 'none', padding: '5px', borderTop: `2px solid ${Colors.RED5}` }}>{error}</div>
        </Popover>
      </div>
      )
    : null
}

export default InputValidationError
