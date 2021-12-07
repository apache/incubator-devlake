import React from 'react'
import { Alert, Intent, TextArea } from '@blueprintjs/core'

const SaveAlert = ({ alertOpen, onClose }) => {
  return (
    <Alert
      canEscapeKeyCancel={true}
      canOutsideClickCancel={true}
      confirmButtonText='Continue'
      isOpen={alertOpen}
      onClose={onClose}
      intent={Intent.PRIMARY}
    >
      <h2 style={{ fontWeight: 'bold' }}><span style={{ fontWeight: 800 }}>API</span> Configuration Updated</h2>
      <p style={{ fontSize: '16px', fontFamily: 'Montserrat', color: '#E8471C' }}>
        To apply new configuration, <strong>restart</strong> devlake by running <code>docker-compose up -d</code>.
      </p>
      <TextArea
        readOnly
        fill
        rows={1}
        style={{
          fontSize: '13px',
          resize: 'none',
          boxShadow: '0 0 0 1px #e8471c, 0 0 0 3px rgba(232, 71, 28, 0.3), inset 0 1px 1px rgba(16, 22, 26, 0.2)'
        }}
        growVertically={false}
        autoFocus
      >
        docker-compose up -d
      </TextArea>
      <p style={{ marginTop: '10px' }}>
        Click <strong>ESC</strong> or <strong>Continue</strong>&nbsp;
        when ready after running the command shown above.
      </p>
    </Alert>
  )
}

export default SaveAlert
