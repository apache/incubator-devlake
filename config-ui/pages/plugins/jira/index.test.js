/**
 * @jest-environment jsdom
 */

import React from 'react'
import { render, screen } from '@testing-library/react'
import Home from './index'

describe('Check the correct heading', () => {
  it('renders a heading with proper text', () => {
    render(<Home env={{
      JIRA_ENDPOINT: "test-endpoint"
    }} />)

    const heading = screen.getByRole('heading', {
      name: 'Jira Plugin'
    })

    expect(heading).toBeInTheDocument()
  })
})
