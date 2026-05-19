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
import { nanoid } from 'nanoid'
import { STORAGE_KEYS } from '../constants'
import type {
  GeneratedFile,
  GeneratedImage,
  PlaygroundAttachment,
  PlaygroundConfig,
  ParameterEnabled,
  Message,
  PlaygroundSession,
} from '../types'
import { sanitizeMessagesOnLoad } from './message-utils'

export const MAX_PERSISTED_SESSION_BYTES = 1_800_000
const MAX_PERSISTED_MESSAGE_TEXT = 120_000
const MAX_PERSISTED_ATTACHMENT_TEXT = 20_000

function getCurrentStorageUserId() {
  try {
    if (typeof window === 'undefined') return 'anonymous'
    return window.localStorage.getItem('uid') || 'anonymous'
  } catch {
    return 'anonymous'
  }
}

function getScopedStorageKey(key: string) {
  return `${key}:${getCurrentStorageUserId()}`
}

function getStorageItem(key: string) {
  return localStorage.getItem(getScopedStorageKey(key))
}

function setStorageItem(key: string, value: string) {
  localStorage.setItem(getScopedStorageKey(key), value)
}

function removeStorageItem(key: string) {
  localStorage.removeItem(getScopedStorageKey(key))
}

/**
 * Load playground config from localStorage
 */
export function loadConfig(): Partial<PlaygroundConfig> {
  try {
    const saved = getStorageItem(STORAGE_KEYS.CONFIG)
    if (saved) {
      return JSON.parse(saved)
    }
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Failed to load config:', error)
  }
  return {}
}

/**
 * Save playground config to localStorage
 */
export function saveConfig(config: Partial<PlaygroundConfig>): void {
  try {
    setStorageItem(STORAGE_KEYS.CONFIG, JSON.stringify(config))
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Failed to save config:', error)
  }
}

/**
 * Load parameter enabled state from localStorage
 */
export function loadParameterEnabled(): Partial<ParameterEnabled> {
  try {
    const saved = getStorageItem(STORAGE_KEYS.PARAMETER_ENABLED)
    if (saved) {
      return JSON.parse(saved)
    }
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Failed to load parameter enabled:', error)
  }
  return {}
}

/**
 * Save parameter enabled state to localStorage
 */
export function saveParameterEnabled(
  parameterEnabled: Partial<ParameterEnabled>
): void {
  try {
    localStorage.setItem(
      getScopedStorageKey(STORAGE_KEYS.PARAMETER_ENABLED),
      JSON.stringify(parameterEnabled)
    )
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Failed to save parameter enabled:', error)
  }
}

/**
 * Load messages from localStorage
 */
export function loadMessages(): Message[] | null {
  try {
    const saved = getStorageItem(STORAGE_KEYS.MESSAGES)
    if (saved) {
      const parsed: unknown = JSON.parse(saved)
      if (!Array.isArray(parsed)) {
        removeStorageItem(STORAGE_KEYS.MESSAGES)
        return null
      }
      const sanitized = sanitizeMessagesOnLoad(parsed as Message[])
      // Persist sanitized result to avoid re-sanitizing on subsequent loads
      if (sanitized !== parsed) {
        saveMessages(sanitized)
      }
      return sanitized
    }
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Failed to load messages:', error)
  }
  return null
}

/**
 * Save messages to localStorage
 */
export function saveMessages(messages: Message[]): void {
  try {
    setStorageItem(STORAGE_KEYS.MESSAGES, JSON.stringify(messages))
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Failed to save messages:', error)
  }
}

export function getSessionTitle(messages: Message[]): string {
  const firstUserMessage = messages.find((message) => message.from === 'user')
  const content = firstUserMessage?.versions[0]?.content?.trim()
  if (!content) return 'Untitled conversation'
  return content.length > 36 ? `${content.slice(0, 36)}...` : content
}

export function createEmptySession(): PlaygroundSession {
  const now = Date.now()
  return {
    id: nanoid(),
    title: 'Untitled conversation',
    messages: [],
    selectedImage: null,
    createdAt: now,
    updatedAt: now,
  }
}

function isInlineImageUrl(url?: string) {
  return typeof url === 'string' && url.startsWith('data:image/')
}

function isEphemeralUrl(url?: string) {
  return typeof url === 'string' && (url.startsWith('blob:') || isInlineImageUrl(url))
}

function stripStoredImage(image?: GeneratedImage | null) {
  if (!image || isInlineImageUrl(image.url)) return null
  return image
}

function stripStoredFile(file?: GeneratedFile | null) {
  if (!file || isEphemeralUrl(file.url)) return null
  return file
}

function stripStoredAttachment(
  attachment: NonNullable<Message['attachments']>[number]
): PlaygroundAttachment | null {
  if (isInlineImageUrl(attachment.url)) return null
  const trimmed: PlaygroundAttachment = {
    ...attachment,
  }
  if (
    trimmed.textContent &&
    trimmed.textContent.length > MAX_PERSISTED_ATTACHMENT_TEXT
  ) {
    trimmed.textContent = `${trimmed.textContent.slice(0, MAX_PERSISTED_ATTACHMENT_TEXT)}\n\n[Content truncated for storage]`
  }
  return trimmed
}

function stripStoredMessage(message: Message): Message {
  const versions = message.versions.map((version) => ({
    ...version,
    content:
      version.content.length > MAX_PERSISTED_MESSAGE_TEXT
        ? `${version.content.slice(0, MAX_PERSISTED_MESSAGE_TEXT)}\n\n[Content truncated for storage]`
        : version.content,
  }))

  return {
    ...message,
    versions,
    attachments: message.attachments
      ?.map(stripStoredAttachment)
      .filter(
        (
          attachment
        ): attachment is NonNullable<Message['attachments']>[number] =>
          Boolean(attachment)
      ),
    generatedImages: message.generatedImages
      ?.map(stripStoredImage)
      .filter((image): image is GeneratedImage => Boolean(image)),
    generatedFiles: message.generatedFiles
      ?.map(stripStoredFile)
      .filter((file): file is GeneratedFile => Boolean(file)),
    imageRequest: message.imageRequest
      ? {
          ...message.imageRequest,
          sourceImages: message.imageRequest.sourceImages?.filter(
            (url) => !isInlineImageUrl(url)
          ),
        }
      : undefined,
  }
}

function estimateSessionBytes(session: PlaygroundSession) {
  return new Blob([JSON.stringify(session)]).size
}

export function estimatePersistedSessionBytes(session: PlaygroundSession) {
  return estimateSessionBytes({
    ...session,
    messages: session.messages.map(stripStoredMessage),
    selectedImage: stripStoredImage(session.selectedImage),
  })
}

function compactMessagesForPersistence(
  session: PlaygroundSession
): Message[] {
  let messages = session.messages
  let candidate = { ...session, messages }

  while (messages.length > 1 && estimateSessionBytes(candidate) > MAX_PERSISTED_SESSION_BYTES) {
    messages = messages.slice(1)
    candidate = { ...candidate, messages }
  }

  return messages
}

export function prepareSessionForPersistence(
  session: PlaygroundSession
): PlaygroundSession {
  const prepared = {
    ...session,
    messages: session.messages.map(stripStoredMessage),
    selectedImage: stripStoredImage(session.selectedImage),
  }
  return {
    ...prepared,
    messages: compactMessagesForPersistence(prepared),
  }
}

function prepareSessionsForStorage(sessions: PlaygroundSession[]) {
  return sessions.map(prepareSessionForPersistence)
}

export function loadSessions(): PlaygroundSession[] {
  try {
    const saved = getStorageItem(STORAGE_KEYS.SESSIONS)
    if (saved) {
      const parsed: unknown = JSON.parse(saved)
      if (Array.isArray(parsed)) {
        const sessions = parsed
          .map((session) => session as Partial<PlaygroundSession>)
          .filter((session) => typeof session.id === 'string')
          .map((session) => ({
            id: session.id!,
            remoteId: session.remoteId,
            title: session.title || getSessionTitle(session.messages || []),
            messages: sanitizeMessagesOnLoad(session.messages || []),
            selectedImage: session.selectedImage ?? null,
            createdAt: session.createdAt || Date.now(),
            updatedAt: session.updatedAt || Date.now(),
          }))
        if (sessions.length > 0) return sessions
      }
    }

    const legacyMessages = loadMessages()
    if (legacyMessages?.length) {
      const now = Date.now()
      return [
        {
          id: nanoid(),
          title: getSessionTitle(legacyMessages),
          messages: legacyMessages,
          selectedImage: null,
          createdAt: now,
          updatedAt: now,
        },
      ]
    }
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Failed to load sessions:', error)
  }
  return [createEmptySession()]
}

export function saveSessions(sessions: PlaygroundSession[]): void {
  try {
    localStorage.setItem(
      getScopedStorageKey(STORAGE_KEYS.SESSIONS),
      JSON.stringify(prepareSessionsForStorage(sessions))
    )
  } catch (error) {
    try {
      const compactSessions = prepareSessionsForStorage(sessions).map(
        (session) => ({
          ...session,
          messages: session.messages.slice(-10),
        })
      )
      localStorage.setItem(
        getScopedStorageKey(STORAGE_KEYS.SESSIONS),
        JSON.stringify(compactSessions)
      )
    } catch (compactError) {
      // eslint-disable-next-line no-console
      console.error('Failed to save sessions:', compactError || error)
    }
  }
}

export function loadActiveSessionId(): string | null {
  try {
    return getStorageItem(STORAGE_KEYS.ACTIVE_SESSION_ID)
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Failed to load active session:', error)
  }
  return null
}

export function saveActiveSessionId(sessionId: string): void {
  try {
    setStorageItem(STORAGE_KEYS.ACTIVE_SESSION_ID, sessionId)
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Failed to save active session:', error)
  }
}

/**
 * Clear all playground data
 */
export function clearPlaygroundData(): void {
  try {
    removeStorageItem(STORAGE_KEYS.CONFIG)
    removeStorageItem(STORAGE_KEYS.PARAMETER_ENABLED)
    removeStorageItem(STORAGE_KEYS.MESSAGES)
    removeStorageItem(STORAGE_KEYS.SESSIONS)
    removeStorageItem(STORAGE_KEYS.ACTIVE_SESSION_ID)
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Failed to clear playground data:', error)
  }
}
