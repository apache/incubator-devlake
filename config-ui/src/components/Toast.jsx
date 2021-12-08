import { Position, Toaster } from '@blueprintjs/core'

export const ToastNotification = Toaster.create({
  className: 'recipe-toaster',
  position: Position.BOTTOM_RIGHT,
})
