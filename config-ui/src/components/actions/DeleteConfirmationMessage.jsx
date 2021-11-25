import React, { Fragment } from 'react'

const DeleteConfirmationMessage = (props) => {
  const { title = 'DELETE CONFIRMATION' } = props
  return (
    <>
      <h3>{title}</h3>
      <p className='confirmation-text'>
        <strong>
          Are you sure you want to continue?
        </strong>
        &nbsp;This instance will be permanently deleted and cannot be restored.
      </p>
    </>
  )
}

export default DeleteConfirmationMessage
