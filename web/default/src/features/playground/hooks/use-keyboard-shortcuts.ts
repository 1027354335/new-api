import { useEffect } from 'react'

interface ShortcutActions {
  onNewSession: () => void
  onToggleSessions?: () => void
  onToggleSettings?: () => void
  onClearChat: () => void
  onStopGeneration: () => void
  onFocusInput: () => void
  isGenerating: boolean
}

export function useKeyboardShortcuts(actions: ShortcutActions) {
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      const isMod = e.ctrlKey || e.metaKey

      // Ctrl/Cmd + Shift + N -> New Session
      if (isMod && e.shiftKey && e.key.toUpperCase() === 'N') {
        e.preventDefault()
        actions.onNewSession()
        return
      }

      // Ctrl/Cmd + Shift + O -> Toggle Sessions Panel
      if (isMod && e.shiftKey && e.key.toUpperCase() === 'O') {
        e.preventDefault()
        actions.onToggleSessions?.()
        return
      }

      // Ctrl/Cmd + Shift + S -> Toggle Settings Panel
      if (isMod && e.shiftKey && e.key.toUpperCase() === 'S') {
        e.preventDefault()
        actions.onToggleSettings?.()
        return
      }

      // Ctrl/Cmd + Shift + L -> Clear Chat
      if (isMod && e.shiftKey && e.key.toUpperCase() === 'L') {
        e.preventDefault()
        actions.onClearChat()
        return
      }

      // Escape -> Stop Generation
      if (e.key === 'Escape' && actions.isGenerating) {
        e.preventDefault()
        actions.onStopGeneration()
        return
      }

      // Ctrl/Cmd + / -> Focus Input
      if (isMod && e.key === '/') {
        e.preventDefault()
        actions.onFocusInput()
        return
      }
    }

    window.addEventListener('keydown', handler)
    return () => window.removeEventListener('keydown', handler)
  }, [actions])
}
