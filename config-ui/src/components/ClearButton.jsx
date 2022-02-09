import React from 'react'
import { Button, Intent } from '@blueprintjs/core'

const ClearButton = ({ onClick, minimal = true, intent = Intent.NONE, disabled = false }) => {
  return (
    <Button
      disabled={disabled}
      intent={intent}
      icon='cross'
      minimal={minimal}
      onClick={onClick}
    />
  )
}

export default ClearButton
