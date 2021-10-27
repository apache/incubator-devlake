import React from 'react'
import {
  Alignment,
  Breadcrumbs,
  Breadcrumb,
  Icon,
  Colors,
} from '@blueprintjs/core'
import '../styles/breadcrumbs.scss'

const AppCrumbs = (props) => {
  const { items } = props

  const renderBreadcrumb = ({ text, ...restProps }) => {
    return <Breadcrumb {...restProps}>{text}</Breadcrumb>
  }

  const renderCurrentBreadcrumb = ({ text, ...restProps }) => {
    return <Breadcrumb {...restProps}>{text} <Icon icon='symbol-circle' size={4} color={Colors.GREEN3} /></Breadcrumb>
  }

  return (
    <Breadcrumbs
      breadcrumbRenderer={renderBreadcrumb}
      currentBreadcrumbRenderer={renderCurrentBreadcrumb}
      items={items}
      icon='icon-slash'
    />
  )
}

export default AppCrumbs
