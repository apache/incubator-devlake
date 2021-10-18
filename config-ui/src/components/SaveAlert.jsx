import React from 'react'
import { Alert } from '@blueprintjs/core'

const SaveAlert = ({alertOpen, onClose}) => {

  return <Alert
    canEscapeKeyCancel={true}
    canOutsideClickCancel={true}
    confirmButtonText="Ok"
    isOpen={alertOpen}
    onClose={onClose}
    >
    <h4>Config File Updated</h4>
    <p>To apply new configuration, restart devlake by running: <br/><br/><code>docker-compose up -d</code></p>
  </Alert>
}

export default SaveAlert
