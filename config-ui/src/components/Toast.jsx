import { Button, Position, Toast, Toaster } from '@blueprintjs/core'
import React, { useEffect, useState } from 'react'

export const ToastNotification = Toaster.create({
  className: 'recipe-toaster',
  position: Position.BOTTOM_RIGHT,
})
