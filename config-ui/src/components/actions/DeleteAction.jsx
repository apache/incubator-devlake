import React, { Fragment } from 'react'
import {
  Button, Card, Elevation, Colors,
  Spinner,
  Tooltip,
  Position,
  Icon,
  Intent,
  Popover,
  Popover2,
  Classes
} from '@blueprintjs/core'

const DeleteAction = (props) => {
  const {
    id, connection, showConfirmation = () => {}, onConfirm = () => {}, onCancel = () => {},
    isDisabled = false,
    isLoading = false,
    text = 'Delete',
    children
  } = props
  return (
    <Popover
      key={`delete-popover-key-${connection.ID}`}
      className='trigger-delete-connection'
      popoverClassName='popover-delete-connection'
      position={Position.RIGHT}
      autoFocus={false}
      enforceFocus={false}
      isOpen={id === connection.ID}
      usePortal={false}
    >
      <a
        href='#'
        intent={Intent.DANGER}
        data-provider={connection.id}
        className='table-action-link actions-link'
        onClick={showConfirmation}
        style={{ color: '#DB3737' }}
      >
        <Icon icon='trash' color={Colors.RED3} size={12} />
        Delete
      </a>
      <>
        <div style={{ padding: '15px 20px 15px 15px' }}>
          {children}
          {/* <h3 style={{ color: 'rgb(219, 55, 55)' }}>DELETE CONFIRMATION</h3>
          <p className='confirmation-text'>
            <strong style={{
              fontFamily: 'Montserrat, sans-serif',
              fontSize: '14px',
              fontWeight: '800',
              color: 'rgb(219, 55, 55)'
            }}
            >
              Are you sure you want to continue?
            </strong>
            &nbsp;This instance will be permanently deleted and cannot be restored.
          </p> */}
          <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: 15 }}>
            <Button
              className={Classes.POPOVER2_DISMISS}
              style={{ marginRight: 10 }}
              disabled={isDisabled || isLoading}
              onClick={onCancel}
            >
              Cancel
            </Button>
            <Button
              disabled={isDisabled}
              loading={isLoading}
              onClick={(e) => onConfirm(connection, e)}
              intent={Intent.DANGER}
              icon='remove'
              className={Classes.POPOVER2_DISMISS}
              style={{ fontWeight: 'bold' }}
            >
              {text}
            </Button>
          </div>
        </div>
      </>
    </Popover>
  )
}

export default DeleteAction
