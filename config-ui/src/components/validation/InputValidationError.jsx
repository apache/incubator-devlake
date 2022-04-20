import React, { useEffect, useState, useCallback } from 'react'
import {
  Colors,
  Icon,
  Popover,
  PopoverInteractionKind,
  Intent,
  Position
} from '@blueprintjs/core'

const InputValidationError = (props) => {
  const {
    error,
    position = Position.TOP,
    validateOnFocus = false,
    elementRef, onError = () => {},
    onSuccess = () => {}
  } = props

  const [elementIsFocused, setElementIsFocused] = useState(false)
  const [inputElement, setInputElement] = useState(null)

  const handleElementFocus = useCallback((isFocused, ref) => {
    setElementIsFocused(isFocused)
    if (error) {
      elementRef?.current.parentElement.classList.remove('valid-field')
      elementRef?.current.parentElement.classList.add('invalid-field')
    } else {
      elementRef?.current.parentElement.classList.remove('invalid-field')
      elementRef?.current.parentElement.classList.add('valid-field')
    }
  }, [elementRef, error])

  const handleElementBlur = useCallback((isFocused, ref) => {
    setElementIsFocused(isFocused)
    if (!error) {
      elementRef?.current.parentElement.classList.remove('invalid-field')
    }
  }, [elementRef])

  useEffect(() => {
    console.log('cName Ref===', elementRef)
    const iRef = elementRef?.current
    if (iRef) {
      setInputElement(iRef)
      // iRef.addEventListener('focus', (e) => setElementIsFocused(true), true)
      // iRef.addEventListener('blur', (e) => setElementIsFocused(false), true)
      iRef.addEventListener('focus', (e) => handleElementFocus(true, iRef), true)
      iRef.addEventListener('keyup', (e) => handleElementFocus(true, iRef), true)
      iRef.addEventListener('blur', (e) => handleElementBlur(false, iRef), true)
    } else {
      setInputElement(null)
    }

    return () => {
      iRef?.removeEventListener('focus', setElementIsFocused, true)
      iRef?.removeEventListener('keyup', setElementIsFocused, true)
      iRef?.removeEventListener('blur', setElementIsFocused, true)
      setInputElement(null)
    }
  }, [elementRef, handleElementBlur, handleElementFocus])

  useEffect(() => {
    if (error && elementIsFocused) {
      onError(elementRef?.current?.id ? elementRef?.current?.id : null)
    } else {
      onSuccess()
    }
  }, [error, onError, onSuccess, elementIsFocused, elementRef])

  return error
    ? (
      <div className='inline-input-error' style={{ outline: 'none', cursor: 'pointer', margin: '5px 5px 3px 5px' }}>
        <Popover
          position={position}
          usePortal={true}
          openOnTargetFocus={true}
          intent={Intent.WARNING}
          interactionKind={PopoverInteractionKind.HOVER_TARGET_ONLY}
          enforceFocus={false}
          // autoFocus={false}
        >
          <Icon
            icon='warning-sign'
            size={12}
            color={elementIsFocused ? Colors.RED5 : Colors.GRAY5}
            style={{ outline: 'none' }}
            onClick={(e) => e.stopPropagation()}
          />
          <div style={{ outline: 'none', padding: '5px', borderTop: `2px solid ${Colors.RED5}` }}>{error}</div>
        </Popover>
      </div>
      )
    : null
}

export default InputValidationError
