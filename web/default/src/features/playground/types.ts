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
// Message types
export type MessageRole = 'user' | 'assistant' | 'system'

export type MessageStatus = 'loading' | 'streaming' | 'complete' | 'error'

export interface MessageVersion {
  id: string
  content: string
}

export interface Message {
  key: string
  from: MessageRole
  versions: MessageVersion[]
  attachments?: PlaygroundAttachment[]
  generatedImages?: GeneratedImage[]
  generatedFiles?: GeneratedFile[]
  imageRequest?: PlaygroundImageRequest
  imageTaskId?: string
  sources?: { href: string; title: string }[]
  reasoning?: {
    content: string
    duration: number
  }
  isReasoningStreaming?: boolean
  isReasoningComplete?: boolean
  isContentComplete?: boolean
  status?: MessageStatus
  errorCode?: string | null
}

// API payload types
export interface ChatCompletionMessage {
  role: MessageRole
  content: string | ContentPart[]
}

export interface ContentPart {
  type: 'text' | 'image_url'
  text?: string
  image_url?: {
    url: string
  }
}

export interface PlaygroundAttachment {
  id: string
  url: string
  name?: string
  mediaType?: string
  textContent?: string
}

export interface ChatCompletionRequest {
  model: string
  group?: string
  messages: ChatCompletionMessage[]
  stream: boolean
  temperature?: number
  top_p?: number
  max_tokens?: number
  frequency_penalty?: number
  presence_penalty?: number
  seed?: number
}

export interface ChatCompletionChunk {
  id: string
  object: string
  created: number
  model: string
  choices: Array<{
    index: number
    delta: {
      role?: MessageRole
      content?: string
      reasoning_content?: string
    }
    finish_reason: string | null
  }>
}

export interface ChatCompletionResponse {
  id: string
  object: string
  created: number
  model: string
  choices: Array<{
    index: number
    message: {
      role: MessageRole
      content: string
      reasoning_content?: string
    }
    finish_reason: string
  }>
  usage?: {
    prompt_tokens: number
    completion_tokens: number
    total_tokens: number
  }
}

export interface GeneratedImage {
  id: string
  url: string
  revisedPrompt?: string
}

export type GeneratedFileKind = 'excel' | 'word' | 'powerpoint'

export interface GeneratedFile {
  id: string
  kind: GeneratedFileKind
  name: string
  url: string
  mimeType: string
  size?: number
}

export interface AssistantPostProcessResult {
  content?: string
  generatedFile?: GeneratedFile
}

export interface PlaygroundSession {
  id: string
  remoteId?: number
  title: string
  messages: Message[]
  selectedImage?: GeneratedImage | null
  createdAt: number
  updatedAt: number
}

export interface PlaygroundSessionRecord {
  id: number
  title: string
  messages: Message[]
  selected_image?: GeneratedImage | null
  created_time: number
  updated_time: number
}

export interface PlaygroundSessionsPage {
  page: number
  page_size: number
  total: number
  items: PlaygroundSessionRecord[]
}

export interface PlaygroundImageRequest {
  prompt: string
  sourceImages?: string[]
  sourceContext?: string
}

export interface ImageGenerationRequest {
  model: string
  group?: string
  prompt: string
  n?: number
  size?: string
  quality?: string
  moderation?: string
  input_fidelity?: string
  images?: string[]
}

export interface ImageGenerationResponse {
  created?: number
  data?: Array<{
    url?: string
    b64_json?: string
    revised_prompt?: string
  }>
  metadata?: unknown
}

export type PlaygroundMode = 'chat' | 'image'

// Configuration types
export interface PlaygroundConfig {
  mode: PlaygroundMode
  model: string
  group: string
  temperature: number
  top_p: number
  max_tokens: number
  frequency_penalty: number
  presence_penalty: number
  seed: number | null
  stream: boolean
  imageSize: string
  imageQuality: string
  imageModeration: string
  imageCount: number
}

export interface ParameterEnabled {
  temperature: boolean
  top_p: boolean
  max_tokens: boolean
  frequency_penalty: boolean
  presence_penalty: boolean
  seed: boolean
}

// Model and group options
export interface ModelOption {
  label: string
  value: string
}

export interface GroupOption {
  label: string
  value: string
  ratio: number
  desc?: string
}
