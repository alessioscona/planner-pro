import { render, screen } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import App from './App'

// Mock per bypassare Keycloak/Auth
vi.mock('./RequireAuth', () => ({
  RequireAuth: ({ children }: { children: React.ReactNode }) => <>{children}</>,
}))

describe('App component', () => {
  it('renders welcome message', () => {
    render(<App />)
    expect(screen.getByText(/Welcome to Scheduler/i)).toBeInTheDocument()
  })
})
