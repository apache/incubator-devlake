import React, { Fragment } from 'react'
import {
  Icon,
  Colors
} from '@blueprintjs/core'
const FormValidationErrors = (props) => {
  const { errors = [], textAlign = 'right', styles = {} } = props

  return (
    <>
      {errors.length > 0 && (
        <div className='validation-errors'>
          <p style={{ margin: '5px 0 5px 0', textAlign: textAlign, ...styles }}>
            <Icon icon='warning-sign' size={13} color={Colors.ORANGE4} style={{ marginRight: '6px', marginBottom: '2px' }} />
            <span>
              {errors[0]}
            </span>
          </p>
        </div>
      )}
    </>
  )
}

export default FormValidationErrors
