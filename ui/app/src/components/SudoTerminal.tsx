import {useRef, useState, type FormEvent, type KeyboardEvent} from 'react'
import {type DemoUser} from '../models'

type Props = {
  users: DemoUser[]
  onLogin: (username: string, password: string) => Promise<void>
  onClose: () => void
}

export function SudoTerminal({users, onLogin, onClose}: Props) {
  const [user, setUser] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const form = useRef<HTMLFormElement>(null)
  const passwordInput = useRef<HTMLInputElement>(null)

  async function submit(event: FormEvent) {
    event.preventDefault()
    setError('')
    try {
      await onLogin(user, password)
      setPassword('')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'sudo failed')
    }

  }

  function handleUserKeyDown(event: KeyboardEvent<HTMLInputElement>) {
    if (event.key === 'Enter') {
      event.preventDefault()
      passwordInput.current?.focus()
    }
  }

  function handlePasswordKeyDown(event: KeyboardEvent<HTMLInputElement>) {
    if (event.key === 'Enter') {
      event.preventDefault()
      form.current?.requestSubmit()
    }
  }

  return <div className="terminal-backdrop" role="presentation" onMouseDown={event => {
    if (event.target === event.currentTarget) onClose()
  }}>
    <form ref={form} className="terminal-window" role="dialog" aria-modal="true" aria-labelledby="sudo-title" onSubmit={submit}>
      <div className="terminal-titlebar">
        <span id="sudo-title">urpi@backlog: ~</span>
        <button type="button" className="terminal-close" aria-label="Close sudo terminal" onClick={onClose}>×</button>
      </div>
      <div className="terminal-body">
        <label><span className="terminal-prompt">$</span> sudo <input autoFocus value={user}
          onChange={event => setUser(event.target.value)} onKeyDown={handleUserKeyDown}
          placeholder="user" aria-label="Username"/></label>
        <label><span className="terminal-prompt">$</span> password: <input ref={passwordInput} type="password" value={password}
          onChange={event => setPassword(event.target.value)} onKeyDown={handlePasswordKeyDown} aria-label="Password"/></label>
        {error && <p className="terminal-error">{error}</p>}
      </div>
    </form>
  </div>
}
