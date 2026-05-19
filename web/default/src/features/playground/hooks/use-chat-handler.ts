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
import { useCallback } from 'react'
import { toast } from 'sonner'
import { sendChatCompletion } from '../api'
import { MESSAGE_STATUS, ERROR_MESSAGES } from '../constants'
import {
  buildChatCompletionPayload,
  updateAssistantMessageWithError,
  updateLastAssistantMessage,
  processStreamingContent,
  finalizeMessage,
} from '../lib'
import type {
  AssistantPostProcessResult,
  Message,
  PlaygroundConfig,
  ParameterEnabled,
} from '../types'
import { useStreamRequest } from './use-stream-request'

interface UseChatHandlerOptions {
  config: PlaygroundConfig
  parameterEnabled: ParameterEnabled
  onMessageUpdate: (updater: (prev: Message[]) => Message[]) => void
  onAssistantComplete?: (
    content: string
  ) => Promise<AssistantPostProcessResult | null> | AssistantPostProcessResult | null
}

/**
 * Hook for handling chat message sending and receiving
 */
export function useChatHandler({
  config,
  parameterEnabled,
  onMessageUpdate,
  onAssistantComplete,
}: UseChatHandlerOptions) {
  const { sendStreamRequest, stopStream, isStreaming } = useStreamRequest()

  const postProcessAssistantMessage = useCallback(
    (messageKey: string, content: string) => {
      if (!onAssistantComplete || !content.trim()) return

      void Promise.resolve(onAssistantComplete(content))
        .then((result) => {
          if (!result) return
          onMessageUpdate((prev) =>
            prev.map((message) => {
              if (message.key !== messageKey) return message
              return {
                ...message,
                versions: result.content
                  ? [
                      {
                        ...message.versions[0],
                        content: result.content,
                      },
                    ]
                  : message.versions,
                generatedFiles: result.generatedFile
                  ? [
                      ...(message.generatedFiles ?? []),
                      result.generatedFile,
                    ]
                  : message.generatedFiles,
              }
            })
          )
        })
        .catch((error: unknown) => {
          const err = error as { message?: string }
          toast.error(err?.message || ERROR_MESSAGES.API_REQUEST_ERROR)
        })
    },
    [onAssistantComplete, onMessageUpdate]
  )

  // Handle stream update
  const handleStreamUpdate = useCallback(
    (type: 'reasoning' | 'content', chunk: string) => {
      onMessageUpdate((prev) =>
        updateLastAssistantMessage(prev, (message) => {
          if (message.status === MESSAGE_STATUS.ERROR) return message

          if (type === 'reasoning') {
            // Direct API reasoning_content
            return {
              ...message,
              reasoning: {
                content: (message.reasoning?.content || '') + chunk,
                duration: 0,
              },
              isReasoningStreaming: true,
              status: MESSAGE_STATUS.STREAMING,
            }
          }

          // Content streaming: handle <think> tags
          return {
            ...processStreamingContent(message, chunk),
            status: MESSAGE_STATUS.STREAMING,
          }
        })
      )
    },
    [onMessageUpdate]
  )

  // Handle stream complete
  const handleStreamComplete = useCallback(() => {
    let completedMessageKey = ''
    let completedContent = ''
    onMessageUpdate((prev) =>
      updateLastAssistantMessage(prev, (message) => {
        if (
          message.status === MESSAGE_STATUS.COMPLETE ||
          message.status === MESSAGE_STATUS.ERROR
        ) {
          return message
        }
        const completedMessage = {
          ...finalizeMessage(message),
          status: MESSAGE_STATUS.COMPLETE,
        }
        completedMessageKey = completedMessage.key
        completedContent = completedMessage.versions[0]?.content || ''
        return completedMessage
      })
    )
    if (completedMessageKey) {
      postProcessAssistantMessage(completedMessageKey, completedContent)
    }
  }, [onMessageUpdate, postProcessAssistantMessage])

  // Handle stream error
  const handleStreamError = useCallback(
    (error: string, errorCode?: string) => {
      toast.error(error)
      onMessageUpdate((prev) =>
        updateAssistantMessageWithError(prev, error, errorCode)
      )
    },
    [onMessageUpdate]
  )

  // Send streaming chat request
  const sendStreamingChat = useCallback(
    (messages: Message[]) => {
      const payload = buildChatCompletionPayload(
        messages,
        config,
        parameterEnabled
      )
      sendStreamRequest(
        payload,
        handleStreamUpdate,
        handleStreamComplete,
        handleStreamError
      )
    },
    [
      config,
      parameterEnabled,
      sendStreamRequest,
      handleStreamUpdate,
      handleStreamComplete,
      handleStreamError,
    ]
  )

  // Send non-streaming chat request
  const sendNonStreamingChat = useCallback(
    async (messages: Message[]) => {
      const payload = buildChatCompletionPayload(
        messages,
        config,
        parameterEnabled
      )

      try {
        const response = await sendChatCompletion(payload)
        const choice = response.choices?.[0]
        if (!choice) return

        let completedMessageKey = ''
        let completedContent = ''
        onMessageUpdate((prev) =>
          updateLastAssistantMessage(prev, (message) => {
            const completedMessage = {
              ...finalizeMessage(
                {
                  ...message,
                  versions: [
                    {
                      ...message.versions[0],
                      content: choice.message?.content || '',
                    },
                  ],
                },
                choice.message?.reasoning_content
              ),
              status: MESSAGE_STATUS.COMPLETE,
            }
            completedMessageKey = completedMessage.key
            completedContent = completedMessage.versions[0]?.content || ''
            return completedMessage
          })
        )
        if (completedMessageKey) {
          postProcessAssistantMessage(completedMessageKey, completedContent)
        }
      } catch (error: unknown) {
        const err = error as {
          response?: {
            data?: { message?: string; error?: { code?: string } }
          }
          message?: string
        }
        handleStreamError(
          err?.response?.data?.message ||
            err?.message ||
            ERROR_MESSAGES.API_REQUEST_ERROR,
          err?.response?.data?.error?.code || undefined
        )
      }
    },
    [
      config,
      parameterEnabled,
      onMessageUpdate,
      handleStreamError,
      postProcessAssistantMessage,
    ]
  )

  // Send chat request (stream or non-stream based on config)
  const sendChat = useCallback(
    (messages: Message[]) => {
      if (config.stream) {
        sendStreamingChat(messages)
      } else {
        sendNonStreamingChat(messages)
      }
    },
    [config.stream, sendStreamingChat, sendNonStreamingChat]
  )

  // Stop generation
  const stopGeneration = useCallback(() => {
    stopStream()
    onMessageUpdate((prev) =>
      updateLastAssistantMessage(prev, (message) =>
        message.status === MESSAGE_STATUS.LOADING ||
        message.status === MESSAGE_STATUS.STREAMING
          ? { ...finalizeMessage(message), status: MESSAGE_STATUS.COMPLETE }
          : message
      )
    )
  }, [stopStream, onMessageUpdate])

  return {
    sendChat,
    stopGeneration,
    isGenerating: isStreaming,
  }
}
