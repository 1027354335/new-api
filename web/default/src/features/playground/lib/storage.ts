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
  PlaygroundConfig,
  ParameterEnabled,
  Message,
  PlaygroundSession,
} from '../types'
import { sanitizeMessagesOnLoad } from './message-utils'

/**
 * Load playground config from localStorage
 */
export function loadConfig(): Partial<PlaygroundConfig> {
  try {
    const saved = localStorage.getItem(STORAGE_KEYS.CONFIG)
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
    localStorage.setItem(STORAGE_KEYS.CONFIG, JSON.stringify(config))
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
    const saved = localStorage.getItem(STORAGE_KEYS.PARAMETER_ENABLED)
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
      STORAGE_KEYS.PARAMETER_ENABLED,
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
    const saved = localStorage.getItem(STORAGE_KEYS.MESSAGES)
    if (saved) {
      const parsed: unknown = JSON.parse(saved)
      if (!Array.isArray(parsed)) {
        localStorage.removeItem(STORAGE_KEYS.MESSAGES)
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
    localStorage.setItem(STORAGE_KEYS.MESSAGES, JSON.stringify(messages))
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

function stripStoredMessage(message: Message): Message {
  return {
    ...message,
    attachments: message.attachments
      ?.map((attachment) =>
        isInlineImageUrl(attachment.url) ? null : attachment
      )
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

export function prepareSessionForPersistence(
  session: PlaygroundSession
): PlaygroundSession {
  return {
    ...session,
    messages: session.messages.map(stripStoredMessage),
    selectedImage: stripStoredImage(session.selectedImage),
  }
}

function prepareSessionsForStorage(sessions: PlaygroundSession[]) {
  return sessions.map(prepareSessionForPersistence)
}

export function loadSessions(): PlaygroundSession[] {
  try {
    const saved = localStorage.getItem(STORAGE_KEYS.SESSIONS)
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
      STORAGE_KEYS.SESSIONS,
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
        STORAGE_KEYS.SESSIONS,
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
    return localStorage.getItem(STORAGE_KEYS.ACTIVE_SESSION_ID)
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Failed to load active session:', error)
  }
  return null
}

export function saveActiveSessionId(sessionId: string): void {
  try {
    localStorage.setItem(STORAGE_KEYS.ACTIVE_SESSION_ID, sessionId)
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
    localStorage.removeItem(STORAGE_KEYS.CONFIG)
    localStorage.removeItem(STORAGE_KEYS.PARAMETER_ENABLED)
    localStorage.removeItem(STORAGE_KEYS.MESSAGES)
    localStorage.removeItem(STORAGE_KEYS.SESSIONS)
    localStorage.removeItem(STORAGE_KEYS.ACTIVE_SESSION_ID)
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Failed to clear playground data:', error)
  }
}
