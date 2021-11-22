import React, { useEffect } from 'react'
import {
  Card,
  Elevation,
  Intent,
  Spinner
} from '@blueprintjs/core'
const ContentLoader = (props) => {
  const {
    title = 'Loading ...',
    message = 'Please wait while data is loaded.',
    spinnerSize = 24,
    spinnerIntent = Intent.PRIMARY
  } = props

  useEffect(() => {

  }, [title, message, spinnerSize])

  return (
    <Card interactive={false} elevation={Elevation.TWO} style={{ width: '100%', marginBottom: '20px' }}>
      <div style={{}}>
        <div style={{ display: 'flex' }}>
          <Spinner intent={spinnerIntent} size={spinnerSize} />
          <h4 className='bp3-heading' style={{ marginLeft: '10px' }}>
            {title}
          </h4>
        </div>

        <p className='bp3-ui-text bp3-text-large' style={{ margin: 0 }}>
          {message}
        </p>

      </div>
    </Card>
  )
}

export default ContentLoader
