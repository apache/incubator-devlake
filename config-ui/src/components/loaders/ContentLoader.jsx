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
    spinnerIntent = Intent.PRIMARY,
    elevation = Elevation.TWO,
    cardStyle = { width: '100%', marginBottom: '20px', boxShadow: elevation === Elevation.ZERO ? 'none' : 'initial' }
  } = props

  useEffect(() => {

  }, [title, message, spinnerSize])

  return (
    <Card interactive={false} elevation={elevation} style={cardStyle}>
      <div style={{}}>
        <div style={{ display: 'flex' }}>
          <Spinner intent={spinnerIntent} size={spinnerSize} />
          <div style={{ marginLeft: '10px' }}>
            <h4 className='bp3-heading' style={{ margin: '0 0 2px 0' }}>
              {title}
            </h4>
            <p className='bp3-ui-text bp3-text-large' style={{ margin: 0 }}>
              {message}
            </p>
          </div>
        </div>
      </div>
    </Card>
  )
}

export default ContentLoader
