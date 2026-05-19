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
import { api } from '@/lib/api'
import { API_ENDPOINTS } from './constants'
import { prepareSessionForPersistence, sanitizeMessagesOnLoad } from './lib'
import type {
  ChatCompletionRequest,
  ChatCompletionResponse,
  ImageGenerationRequest,
  ImageGenerationResponse,
  ModelOption,
  GroupOption,
  PlaygroundSession,
  PlaygroundSessionRecord,
  PlaygroundSessionsPage,
} from './types'

/**
 * Send chat completion request (non-streaming)
 */
export async function sendChatCompletion(
  payload: ChatCompletionRequest
): Promise<ChatCompletionResponse> {
  const res = await api.post(API_ENDPOINTS.CHAT_COMPLETIONS, payload, {
    skipErrorHandler: true,
  } as Record<string, unknown>)
  return res.data
}

export async function sendImageGeneration(
  payload: ImageGenerationRequest,
  edit = false
): Promise<ImageGenerationResponse> {
  const res = await api.post(
    edit ? API_ENDPOINTS.IMAGE_EDITS : API_ENDPOINTS.IMAGE_GENERATIONS,
    payload,
    {
      skipErrorHandler: true,
    } as Record<string, unknown>
  )
  return res.data
}

function mapSessionRecord(record: PlaygroundSessionRecord): PlaygroundSession {
  return {
    id: String(record.id),
    remoteId: record.id,
    title: record.title,
    messages: sanitizeMessagesOnLoad(record.messages || []),
    selectedImage: record.selected_image ?? null,
    createdAt: record.created_time * 1000,
    updatedAt: record.updated_time * 1000,
  }
}

function unwrapSessionResponse(responseData: {
  success?: boolean
  message?: string
  data?: unknown
}): PlaygroundSession {
  if (!responseData?.success || !responseData.data) {
    throw new Error(responseData?.message || 'Failed to sync session')
  }
  return mapSessionRecord(responseData.data as PlaygroundSessionRecord)
}

function toSessionPayload(session: PlaygroundSession) {
  const persistedSession = prepareSessionForPersistence(session)
  return {
    title: persistedSession.title,
    messages: persistedSession.messages,
    selected_image: persistedSession.selectedImage ?? null,
  }
}

function getRemoteSessionId(session: PlaygroundSession): number | string {
  return session.remoteId ?? session.id
}

export async function getPlaygroundSessions(): Promise<PlaygroundSession[]> {
  try {
    const res = await api.get(API_ENDPOINTS.PLAYGROUND_SESSIONS, {
      params: { p: 1, size: 50 },
      skipErrorHandler: true,
    } as Record<string, unknown>)
    const { data } = res
    if (!data.success || !data.data) return []
    const page = data.data as PlaygroundSessionsPage
    return (page.items || []).map(mapSessionRecord)
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Failed to load playground sessions:', error)
    return []
  }
}

export async function createPlaygroundSession(
  session: PlaygroundSession
): Promise<PlaygroundSession> {
  const res = await api.post(
    API_ENDPOINTS.PLAYGROUND_SESSIONS,
    toSessionPayload(session),
    {
      skipErrorHandler: true,
    } as Record<string, unknown>
  )
  return unwrapSessionResponse(res.data)
}

export async function updatePlaygroundSession(
  session: PlaygroundSession
): Promise<PlaygroundSession> {
  const res = await api.put(
    `${API_ENDPOINTS.PLAYGROUND_SESSIONS}/${getRemoteSessionId(session)}`,
    toSessionPayload(session),
    {
      skipErrorHandler: true,
    } as Record<string, unknown>
  )
  return unwrapSessionResponse(res.data)
}

export async function deletePlaygroundSession(
  sessionId: string
): Promise<void> {
  await api.delete(`${API_ENDPOINTS.PLAYGROUND_SESSIONS}/${sessionId}`, {
    skipErrorHandler: true,
  } as Record<string, unknown>)
}

/**
 * Get user available models
 */
export async function getUserModels(): Promise<ModelOption[]> {
  const res = await api.get(API_ENDPOINTS.USER_MODELS)
  const { data } = res

  if (!data.success || !Array.isArray(data.data)) {
    return []
  }

  return data.data.map((model: string) => ({
    label: model,
    value: model,
  }))
}

/**
 * Get user groups
 */
export async function getUserGroups(): Promise<GroupOption[]> {
  const res = await api.get(API_ENDPOINTS.USER_GROUPS)
  const { data } = res

  if (!data.success || !data.data) {
    return []
  }

  const groupData = data.data as Record<string, { desc: string; ratio: number }>

  // label is for button display (name only); desc is for dropdown content
  return Object.entries(groupData).map(([group, info]) => ({
    label: group,
    value: group,
    ratio: info.ratio,
    desc: info.desc,
  }))
}

/**
 * Upload image to playground storage
 */
export async function uploadPlaygroundImage(file: File): Promise<string> {
  const formData = new FormData()
  formData.append('file', file)
  const res = await api.post('/api/playground/images/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    skipErrorHandler: true,
  } as Record<string, unknown>)
  if (!res.data?.success || !res.data?.url) {
    throw new Error(res.data?.error || 'Upload failed')
  }
  return res.data.url
}
