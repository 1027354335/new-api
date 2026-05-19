/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
import { useState, useCallback, useEffect, useMemo, useRef } from 'react'
import { useAuthStore } from '@/stores/auth-store'
import {
  createPlaygroundSession,
  deletePlaygroundSession as deletePlaygroundSessionApi,
  updatePlaygroundSession,
} from '../api'
import { DEFAULT_CONFIG, DEFAULT_PARAMETER_ENABLED } from '../constants'
import {
  loadConfig,
  saveConfig,
  loadParameterEnabled,
  saveParameterEnabled,
  createEmptySession,
  getSessionTitle,
  loadActiveSessionId,
  loadSessions,
  saveActiveSessionId,
  saveSessions,
} from '../lib'
import type {
  GeneratedImage,
  Message,
  PlaygroundConfig,
  ParameterEnabled,
  ModelOption,
  GroupOption,
  PlaygroundSession,
} from '../types'

/**
 * Main state management hook for playground
 */
export function usePlaygroundState() {
  const currentUserId = useAuthStore((state) => state.auth.user?.id)
  const initialSessions = useMemo(() => loadSessions(), [])
  const pendingCreatesRef = useRef(
    new Map<string, Promise<number | undefined>>()
  )
  const pendingSyncTimersRef = useRef(
    new Map<string, ReturnType<typeof setTimeout>>()
  )
  // Load initial state from localStorage
  const [config, setConfig] = useState<PlaygroundConfig>(() => {
    const savedConfig = loadConfig()
    return { ...DEFAULT_CONFIG, ...savedConfig }
  })

  const [parameterEnabled, setParameterEnabled] = useState<ParameterEnabled>(
    () => {
      const saved = loadParameterEnabled()
      return { ...DEFAULT_PARAMETER_ENABLED, ...saved }
    }
  )

  const [sessions, setSessions] = useState<PlaygroundSession[]>(() => {
    return initialSessions
  })
  const sessionsRef = useRef(initialSessions)
  const userIdRef = useRef(currentUserId)

  const [activeSessionId, setActiveSessionId] = useState<string>(() => {
    const savedId = loadActiveSessionId()
    return (
      initialSessions.find((session) => session.id === savedId)?.id ??
      initialSessions[0]?.id ??
      createEmptySession().id
    )
  })

  const [models, setModels] = useState<ModelOption[]>([])
  const [groups, setGroups] = useState<GroupOption[]>([])

  const activeSession =
    sessions.find((session) => session.id === activeSessionId) ?? sessions[0]
  const messages = activeSession?.messages ?? []
  const selectedImage = activeSession?.selectedImage ?? null

  useEffect(() => {
    if (userIdRef.current === currentUserId) return
    userIdRef.current = currentUserId

    queueMicrotask(() => {
      const nextSessions = loadSessions()
      sessionsRef.current = nextSessions
      setSessions(nextSessions)
      setActiveSessionId(() => {
        const savedId = loadActiveSessionId()
        return (
          nextSessions.find((session) => session.id === savedId)?.id ??
          nextSessions[0]?.id ??
          createEmptySession().id
        )
      })
      setConfig({ ...DEFAULT_CONFIG, ...loadConfig() })
      setParameterEnabled({
        ...DEFAULT_PARAMETER_ENABLED,
        ...loadParameterEnabled(),
      })
    })
  }, [currentUserId])

  const hasInFlightMessage = useCallback(
    (session: PlaygroundSession) =>
      session.messages.some(
        (message) =>
          message.status === 'loading' || message.status === 'streaming'
      ),
    []
  )

  const getRemoteKey = useCallback((session: PlaygroundSession) => {
    if (session.remoteId) return String(session.remoteId)
    return /^\d+$/.test(session.id) ? session.id : null
  }, [])

  const persistSessions = useCallback((nextSessions: PlaygroundSession[]) => {
    sessionsRef.current = nextSessions
    setSessions(nextSessions)
    saveSessions(nextSessions)
  }, [])

  const persistSessionRemote = useCallback(
    async (session: PlaygroundSession) => {
      try {
        if (session.remoteId || /^\d+$/.test(session.id)) {
          await updatePlaygroundSession(session)
          return
        }
        let createPromise = pendingCreatesRef.current.get(session.id)
        if (!createPromise) {
          createPromise = createPlaygroundSession(session)
            .then((created) => {
              setSessions((current) => {
                const next = current.map((item) =>
                  item.id === session.id
                    ? {
                        ...item,
                        remoteId: created.remoteId,
                        createdAt: created.createdAt,
                        updatedAt: Math.max(item.updatedAt, created.updatedAt),
                      }
                    : item
                )
                sessionsRef.current = next
                saveSessions(next)
                return next
              })
              return created.remoteId
            })
            .finally(() => {
              pendingCreatesRef.current.delete(session.id)
            })
          pendingCreatesRef.current.set(session.id, createPromise)
        }
        const remoteId = await createPromise
        if (remoteId) {
          const latestSession =
            sessionsRef.current.find((item) => item.id === session.id) ??
            session
          await updatePlaygroundSession({ ...latestSession, remoteId })
        }
      } catch (error) {
        // eslint-disable-next-line no-console
        console.error('Failed to sync playground session:', error)
      }
    },
    []
  )

  const queuePersistSessionRemote = useCallback(
    (session: PlaygroundSession) => {
      const sessionKey = session.remoteId ? String(session.remoteId) : session.id
      const existingTimer = pendingSyncTimersRef.current.get(sessionKey)
      if (existingTimer) clearTimeout(existingTimer)

      const delay = hasInFlightMessage(session) ? 2500 : 400
      const timer = setTimeout(() => {
        pendingSyncTimersRef.current.delete(sessionKey)
        const latestSession =
          sessionsRef.current.find((item) => item.id === session.id) ?? session
        void persistSessionRemote(latestSession)
      }, delay)

      pendingSyncTimersRef.current.set(sessionKey, timer)
    },
    [hasInFlightMessage, persistSessionRemote]
  )

  const updateActiveSession = useCallback(
    (
      updater:
        | Partial<PlaygroundSession>
        | ((session: PlaygroundSession) => Partial<PlaygroundSession>)
    ) => {
      let changedSession: PlaygroundSession | null = null
      const nextSessions = sessionsRef.current.map((session) => {
        if (session.id !== activeSessionId) return session
        const patch = typeof updater === 'function' ? updater(session) : updater
        changedSession = {
          ...session,
          ...patch,
          updatedAt: Date.now(),
        }
        return changedSession
      })
      persistSessions(nextSessions)
      if (changedSession) queuePersistSessionRemote(changedSession)
    },
    [activeSessionId, persistSessions, queuePersistSessionRemote]
  )

  // Update config with automatic save
  const updateConfig = useCallback(
    <K extends keyof PlaygroundConfig>(key: K, value: PlaygroundConfig[K]) => {
      setConfig((prev) => {
        const updated = { ...prev, [key]: value }
        saveConfig(updated)
        return updated
      })
    },
    []
  )

  // Update parameter enabled with automatic save
  const updateParameterEnabled = useCallback(
    (key: keyof ParameterEnabled, value: boolean) => {
      setParameterEnabled((prev) => {
        const updated = { ...prev, [key]: value }
        saveParameterEnabled(updated)
        return updated
      })
    },
    []
  )

  // Update messages with automatic save
  const updateMessages = useCallback(
    (updater: Message[] | ((prev: Message[]) => Message[])) => {
      updateActiveSession((session) => {
        const newMessages =
          typeof updater === 'function' ? updater(session.messages) : updater
        const shouldRetitle =
          !session.title || session.title === 'Untitled conversation'
        return {
          messages: newMessages,
          title: shouldRetitle ? getSessionTitle(newMessages) : session.title,
        }
      })
    },
    [updateActiveSession]
  )

  const updateSessionMessages = useCallback(
    (
      sessionId: string,
      updater: Message[] | ((prev: Message[]) => Message[])
    ) => {
      let changedSession: PlaygroundSession | null = null
      const nextSessions = sessionsRef.current.map((session) => {
        if (session.id !== sessionId) return session
        const newMessages =
          typeof updater === 'function' ? updater(session.messages) : updater
        const shouldRetitle =
          !session.title || session.title === 'Untitled conversation'
        changedSession = {
          ...session,
          messages: newMessages,
          title: shouldRetitle ? getSessionTitle(newMessages) : session.title,
          updatedAt: Date.now(),
        }
        return changedSession
      })
      persistSessions(nextSessions)
      if (changedSession) queuePersistSessionRemote(changedSession)
    },
    [persistSessions, queuePersistSessionRemote]
  )

  const updateSelectedImage = useCallback(
    (image: GeneratedImage | null) => {
      updateActiveSession({ selectedImage: image })
    },
    [updateActiveSession]
  )

  const updateSessionSelectedImage = useCallback(
    (sessionId: string, image: GeneratedImage | null) => {
      let changedSession: PlaygroundSession | null = null
      const nextSessions = sessionsRef.current.map((session) => {
        if (session.id !== sessionId) return session
        changedSession = {
          ...session,
          selectedImage: image,
          updatedAt: Date.now(),
        }
        return changedSession
      })
      persistSessions(nextSessions)
      if (changedSession) queuePersistSessionRemote(changedSession)
    },
    [persistSessions, queuePersistSessionRemote]
  )

  const replaceSessions = useCallback(
    (nextSessions: PlaygroundSession[]) => {
      if (nextSessions.length === 0) return
      setSessions((current) => {
        const usedLocalIds = new Set<string>()
        const merged = nextSessions.map((remoteSession) => {
          const remoteKey = getRemoteKey(remoteSession)
          const localSession = current.find((session) => {
            const localKey = getRemoteKey(session)
            return remoteKey && localKey === remoteKey
          })

          if (!localSession) return remoteSession

          usedLocalIds.add(localSession.id)
          if (
            hasInFlightMessage(localSession) ||
            localSession.updatedAt >= remoteSession.updatedAt
          ) {
            return {
              ...localSession,
              remoteId: remoteSession.remoteId ?? localSession.remoteId,
            }
          }

          return {
            ...remoteSession,
            id: localSession.id,
            remoteId: remoteSession.remoteId ?? localSession.remoteId,
          }
        })

        const localOnly = current.filter((session) => {
          if (usedLocalIds.has(session.id)) return false
          return !getRemoteKey(session) || hasInFlightMessage(session)
        })

        const next = [...localOnly, ...merged].sort(
          (a, b) => b.updatedAt - a.updatedAt
        )

        if (JSON.stringify(current) === JSON.stringify(next)) {
          return current
        }

        saveSessions(next)
        sessionsRef.current = next
        setActiveSessionId((currentActiveId) => {
          if (next.some((session) => session.id === currentActiveId)) {
            return currentActiveId
          }
          saveActiveSessionId(next[0].id)
          return next[0].id
        })
        return next
      })
    },
    [getRemoteKey, hasInFlightMessage]
  )

  const switchSession = useCallback((sessionId: string) => {
    setActiveSessionId(sessionId)
    saveActiveSessionId(sessionId)
  }, [])

  const createSession = useCallback(() => {
    const session = createEmptySession()
    const nextSessions = [session, ...sessionsRef.current]
    persistSessions(nextSessions)
    setActiveSessionId(session.id)
    saveActiveSessionId(session.id)
    void persistSessionRemote(session)
    return session
  }, [persistSessionRemote, persistSessions])

  const deleteSession = useCallback(
    (sessionId: string) => {
      const currentSessions = sessionsRef.current
      const nextSessions = currentSessions.filter(
        (session) => session.id !== sessionId
      )
      const normalized =
        nextSessions.length > 0 ? nextSessions : [createEmptySession()]
      persistSessions(normalized)
      const deletedSession = currentSessions.find(
        (session) => session.id === sessionId
      )
      const remoteId = deletedSession?.remoteId
      if (remoteId || /^\d+$/.test(sessionId)) {
        void deletePlaygroundSessionApi(String(remoteId ?? sessionId))
      }
      if (activeSessionId === sessionId) {
        setActiveSessionId(normalized[0].id)
        saveActiveSessionId(normalized[0].id)
      }
      if (nextSessions.length === 0) {
        void persistSessionRemote(normalized[0])
      }
    },
    [activeSessionId, persistSessionRemote, persistSessions]
  )

  // Clear all messages
  const clearMessages = useCallback(() => {
    updateMessages([])
  }, [updateMessages])

  // Reset config to defaults
  const resetConfig = useCallback(() => {
    setConfig(DEFAULT_CONFIG)
    setParameterEnabled(DEFAULT_PARAMETER_ENABLED)
    saveConfig(DEFAULT_CONFIG)
    saveParameterEnabled(DEFAULT_PARAMETER_ENABLED)
  }, [])

  return {
    // State
    config,
    parameterEnabled,
    sessions,
    activeSessionId,
    activeSession,
    messages,
    selectedImage,
    models,
    groups,

    // Setters
    setModels,
    setGroups,

    // Actions
    updateConfig,
    updateParameterEnabled,
    updateMessages,
    updateSessionMessages,
    updateSelectedImage,
    updateSessionSelectedImage,
    replaceSessions,
    switchSession,
    createSession,
    deleteSession,
    clearMessages,
    resetConfig,
  }
}
