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
import { useCallback, useEffect, useState, useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import { nanoid } from 'nanoid'
import { toast } from 'sonner'
import { useTranslation } from 'react-i18next'
import { useAuthStore } from '@/stores/auth-store'
import {
  getUserModels,
  getUserGroups,
  sendImageGeneration,
  getPlaygroundSessions,
} from './api'
import { PlaygroundChat } from './components/playground-chat'
import { PlaygroundInput } from './components/playground-input'
import { PlaygroundSessions } from './components/playground-sessions'
import { PlaygroundSettings } from './components/playground-settings'
import { usePlaygroundState, useChatHandler } from './hooks'
import {
  buildFileInstruction,
  createGeneratedFileFromAssistant,
  createUserMessageWithAttachments,
  createLoadingAssistantMessage,
  createImageAssistantMessage,
  detectFileGenerationRequest,
} from './lib'
import type {
  GeneratedImage,
  ImageGenerationRequest,
  ImageGenerationResponse,
  Message as MessageType,
  PlaygroundAttachment,
  PlaygroundImageRequest,
} from './types'

type ImageTaskState = {
  id: string
  sessionId: string
  edit: boolean
  imageRequest: PlaygroundImageRequest
  status: 'pending' | 'complete' | 'error'
  response?: ImageGenerationResponse
  error?: unknown
}

const imageTaskListeners = new Set<(task: ImageTaskState) => void>()
let activeImageTask: ImageTaskState | null = null

function notifyImageTask(task: ImageTaskState) {
  imageTaskListeners.forEach((listener) => listener(task))
}

function subscribeImageTask(listener: (task: ImageTaskState) => void) {
  imageTaskListeners.add(listener)
  if (activeImageTask) listener(activeImageTask)
  return () => {
    imageTaskListeners.delete(listener)
  }
}

function startImageTask(
  id: string,
  sessionId: string,
  payload: ImageGenerationRequest,
  edit: boolean,
  imageRequest: PlaygroundImageRequest
) {
  const task: ImageTaskState = {
    id,
    sessionId,
    edit,
    imageRequest,
    status: 'pending',
  }
  activeImageTask = task
  notifyImageTask(task)

  sendImageGeneration(payload, edit)
    .then((response) => {
      const completed: ImageTaskState = {
        ...task,
        status: 'complete',
        response,
      }
      activeImageTask = completed
      notifyImageTask(completed)
    })
    .catch((error: unknown) => {
      const failed: ImageTaskState = {
        ...task,
        status: 'error',
        error,
      }
      activeImageTask = failed
      notifyImageTask(failed)
    })
}

export function Playground() {
  const { t } = useTranslation()
  const currentUserId = useAuthStore((state) => state.auth.user?.id)
  const {
    config,
    parameterEnabled,
    sessions,
    activeSessionId,
    messages,
    selectedImage,
    models,
    groups,
    updateMessages,
    updateSessionMessages,
    updateSelectedImage,
    updateSessionSelectedImage,
    replaceSessions,
    switchSession,
    createSession,
    deleteSession,
    setModels,
    setGroups,
    updateConfig,
  } = usePlaygroundState()

  const isImageModel = useCallback((modelName: string): boolean => {
    const name = modelName.toLowerCase()
    const imageKeywords = [
      'dall-e',
      'gpt-image',
      'imagen-',
      'flux-',
      'flux.1-',
      'midjourney',
      'stable-diffusion',
      'sdxl',
      'recraft',
      'playground-',
    ]
    return imageKeywords.some((keyword) => name.includes(keyword))
  }, [])

  // Filter models based on the selected mode
  const filteredModels = useMemo(() => {
    return models.filter((m) => {
      const isImg = isImageModel(m.value)
      return config.mode === 'image' ? isImg : !isImg
    })
  }, [models, config.mode, isImageModel])

  // Clear selected image when switching to chat mode
  useEffect(() => {
    if (config.mode === 'chat') {
      updateSelectedImage(null)
    }
  }, [config.mode, updateSelectedImage])

  // Ensure config.model is valid for the current mode
  useEffect(() => {
    if (filteredModels.length === 0) return
    const isValid = filteredModels.some((m) => m.value === config.model)
    if (!isValid) {
      updateConfig('model', filteredModels[0].value)
    }
  }, [filteredModels, config.model, updateConfig])

  const handleAssistantComplete = useCallback(
    (content: string) => createGeneratedFileFromAssistant(content),
    []
  )

  const prepareMessagesForChatRequest = useCallback(
    (displayMessages: MessageType[]) => {
      const lastUserIndex = displayMessages
        .map((message) => message.from)
        .lastIndexOf('user')
      if (lastUserIndex === -1) return displayMessages

      const userMessage = displayMessages[lastUserIndex]
      const content = userMessage.versions[0]?.content || ''
      const fileKind = detectFileGenerationRequest(content)
      if (!fileKind) return displayMessages

      return displayMessages.map((message, index) =>
        index === lastUserIndex
          ? {
              ...message,
              versions: [
                {
                  ...message.versions[0],
                  content: buildFileInstruction(fileKind, content),
                },
              ],
            }
          : message
      )
    },
    []
  )

  const { sendChat, stopGeneration, isGenerating } = useChatHandler({
    config,
    parameterEnabled,
    onMessageUpdate: updateMessages,
    onAssistantComplete: handleAssistantComplete,
  })

  // Edit dialog state
  const [editingMessageKey, setEditingMessageKey] = useState<string | null>(
    null
  )
  const selectedImageUrl = selectedImage?.url ?? null
  const [isGeneratingImage, setIsGeneratingImage] = useState(false)

  // Load models
  const { data: modelsData, isLoading: isLoadingModels } = useQuery({
    queryKey: ['playground-models', currentUserId],
    queryFn: async () => {
      try {
        return await getUserModels()
      } catch (error) {
        toast.error(
          error instanceof Error
            ? error.message
            : t('Failed to load playground models')
        )
        return []
      }
    },
  })

  // Load groups
  const { data: groupsData } = useQuery({
    queryKey: ['playground-groups', currentUserId],
    queryFn: async () => {
      try {
        return await getUserGroups()
      } catch (error) {
        toast.error(
          error instanceof Error
            ? error.message
            : t('Failed to load playground groups')
        )
        return []
      }
    },
  })

  const { data: sessionsData } = useQuery({
    queryKey: ['playground-sessions', currentUserId],
    queryFn: getPlaygroundSessions,
  })

  useEffect(() => {
    if (sessionsData && sessionsData.length > 0) {
      replaceSessions(sessionsData)
    }
  }, [replaceSessions, sessionsData])

  // Update models when data changes
  useEffect(() => {
    if (!modelsData) return

    setModels(modelsData)

    // Set default model if current model is not available
    const isCurrentModelValid = modelsData.some((m) => m.value === config.model)
    if (modelsData.length > 0 && !isCurrentModelValid) {
      updateConfig('model', modelsData[0].value)
    }
  }, [modelsData, config.model, setModels, updateConfig])

  // Update groups when data changes
  useEffect(() => {
    if (!groupsData) return

    setGroups(groupsData)

    const hasCurrentGroup = groupsData.some((g) => g.value === config.group)
    if (!hasCurrentGroup && groupsData.length > 0) {
      const fallback =
        groupsData.find((g) => g.value === 'default')?.value ??
        groupsData[0].value
      updateConfig('group', fallback)
    }
  }, [groupsData, setGroups, config.group, updateConfig])

  const parseGeneratedImages = useCallback(
    (responseData: ImageGenerationResponse['data']) =>
      (responseData ?? []).reduce<GeneratedImage[]>((result, item) => {
        const url =
          item.url ||
          (item.b64_json ? `data:image/png;base64,${item.b64_json}` : '')
        if (url) {
          result.push({
            id: nanoid(),
            url,
            revisedPrompt: item.revised_prompt,
          })
        }
        return result
      }, []),
    []
  )

  const applyImageTask = useCallback(
    (task: ImageTaskState) => {
      if (task.status === 'pending') {
        setIsGeneratingImage(true)
        return
      }

      setIsGeneratingImage(false)

      if (task.status === 'complete') {
        const images = parseGeneratedImages(task.response?.data)

        if (images.length === 0) {
          toast.error('No image returned')
          updateSessionMessages(task.sessionId, (prev) =>
            prev.map((message) =>
              message.imageTaskId === task.id
                ? {
                    ...message,
                    versions: [
                      {
                        ...message.versions[0],
                        content: 'No image returned',
                      },
                    ],
                    imageRequest: task.imageRequest,
                    status: 'error',
                  }
                : message
            )
          )
          if (activeImageTask?.id === task.id) activeImageTask = null
          return
        }

        updateSessionMessages(task.sessionId, (prev) =>
          prev.map((message) =>
            message.imageTaskId === task.id
              ? createImageAssistantMessage(
                  images,
                  task.edit ? 'Edited image' : 'Generated images',
                  task.imageRequest
                )
              : message
          )
        )
        updateSessionSelectedImage(task.sessionId, images[0] ?? null)
        if (activeImageTask?.id === task.id) activeImageTask = null
        return
      }

      const err = task.error as {
        response?: {
          data?: { message?: string; error?: { message?: string } }
        }
        message?: string
      }
      const message =
        err?.response?.data?.message ||
        err?.response?.data?.error?.message ||
        err?.message ||
        'Image generation failed'
      toast.error(message)
      updateSessionMessages(task.sessionId, (prev) =>
        prev.map((item) =>
          item.imageTaskId === task.id
            ? {
                ...item,
                versions: [{ ...item.versions[0], content: message }],
                imageRequest: task.imageRequest,
                status: 'error',
              }
            : item
        )
      )
      if (activeImageTask?.id === task.id) activeImageTask = null
    },
    [parseGeneratedImages, updateSessionMessages, updateSessionSelectedImage]
  )

  useEffect(() => {
    return subscribeImageTask(applyImageTask)
  }, [applyImageTask])

  const runImageRequest = useCallback(
    (
      prompt: string,
      nextMessages: MessageType[],
      sourceImages: string[] = [],
      sourceContext?: string
    ) => {
      const loadingMessage = nextMessages[nextMessages.length - 1]
      const isEdit = sourceImages.length > 0
      const imageRequest = {
        prompt,
        sourceImages: isEdit ? sourceImages : undefined,
        sourceContext: isEdit ? sourceContext : undefined,
      }
      const finalPrompt =
        isEdit && sourceContext
          ? `Edit the provided image. Preserve the original composition, subject identity, style, lighting, colors, and fine details unless the requested edit explicitly changes them. Original image context: ${sourceContext}. Requested edit: ${prompt}`
          : isEdit
            ? `Edit the provided image. Preserve the original composition, subject identity, style, lighting, colors, and fine details unless the requested edit explicitly changes them. Requested edit: ${prompt}`
            : prompt
      setIsGeneratingImage(true)
      startImageTask(
        loadingMessage.key,
        activeSessionId,
        {
          model: config.model,
          group: config.group,
          prompt: finalPrompt,
          n: config.imageCount,
          size: config.imageSize === 'auto' ? undefined : config.imageSize,
          quality:
            config.imageQuality === 'auto' ? undefined : config.imageQuality,
          moderation:
            config.imageModeration === 'auto'
              ? undefined
              : config.imageModeration,
          input_fidelity: isEdit ? 'high' : undefined,
          images: isEdit ? sourceImages : undefined,
        },
        isEdit,
        imageRequest
      )
    },
    [activeSessionId, config]
  )

  const handleSendMessage = async (
    text: string,
    attachments: PlaygroundAttachment[] = []
  ) => {
    const userMessage = createUserMessageWithAttachments(text, attachments)
    const attachmentImages = attachments.filter((attachment) =>
      attachment.mediaType?.startsWith('image/')
    )

    if (config.mode === 'image') {
      const sourceImages =
        selectedImageUrl || attachmentImages.length > 0
          ? [
              ...(selectedImageUrl ? [selectedImageUrl] : []),
              ...attachmentImages.map((attachment) => attachment.url),
            ]
          : []
      const loadingMessage = {
        ...createLoadingAssistantMessage(),
        imageRequest: {
          prompt: text,
          sourceImages: sourceImages.length > 0 ? sourceImages : undefined,
          sourceContext: selectedImage?.revisedPrompt,
        },
      }
      loadingMessage.imageTaskId = loadingMessage.key
      const newMessages = [...messages, userMessage, loadingMessage]
      updateMessages(newMessages)
      runImageRequest(
        text,
        newMessages,
        sourceImages,
        selectedImage?.revisedPrompt
      )
      return
    }

    const assistantMessage = createLoadingAssistantMessage()

    const newMessages = [...messages, userMessage, assistantMessage]
    updateMessages(newMessages)

    // Send chat request
    sendChat(prepareMessagesForChatRequest(newMessages))
  }

  const handleCopyMessage = (message: MessageType) => {
    // Copy is handled in MessageActions component
    // eslint-disable-next-line no-console
    console.log('Message copied:', message.key)
  }

  const handleRegenerateMessage = async (message: MessageType) => {
    // Find the message index and regenerate from there
    const messageIndex = messages.findIndex((m) => m.key === message.key)
    if (messageIndex === -1) return

    // Remove messages after this one and regenerate
    const messagesUpToHere = messages.slice(0, messageIndex)
    if (message.generatedImages?.length || message.imageRequest) {
      const prompt =
        message.imageRequest?.prompt ||
        messagesUpToHere[messagesUpToHere.length - 1]?.versions[0]?.content ||
        ''
      if (!prompt.trim()) return
      const loadingMessage = {
        ...createLoadingAssistantMessage(),
        imageRequest: message.imageRequest ?? { prompt },
      }
      loadingMessage.imageTaskId = loadingMessage.key
      const newMessages = [...messagesUpToHere, loadingMessage]

      updateMessages(newMessages)
      runImageRequest(
        prompt,
        newMessages,
        message.imageRequest?.sourceImages ?? [],
        message.imageRequest?.sourceContext
      )
      return
    }

    const loadingMessage = createLoadingAssistantMessage()
    const newMessages = [...messagesUpToHere, loadingMessage]

    updateMessages(newMessages)
    sendChat(prepareMessagesForChatRequest(newMessages))
  }

  const handleEditMessage = useCallback((message: MessageType) => {
    setEditingMessageKey(message.key)
  }, [])

  const handleEditOpenChange = useCallback((open: boolean) => {
    if (!open) setEditingMessageKey(null)
  }, [])

  // Apply edit and optionally re-submit from the edited user message
  const applyEdit = useCallback(
    (newContent: string, submit: boolean) => {
      if (!editingMessageKey) return
      const index = messages.findIndex((m) => m.key === editingMessageKey)
      if (index === -1) return

      const updated = messages.map((m) =>
        m.key === editingMessageKey
          ? { ...m, versions: [{ ...m.versions[0], content: newContent }] }
          : m
      )

      setEditingMessageKey(null)

      if (!submit || updated[index].from !== 'user') {
        updateMessages(updated)
        return
      }

      const toSubmit = [
        ...updated.slice(0, index + 1),
        createLoadingAssistantMessage(),
      ]
      updateMessages(toSubmit)
      if (config.mode === 'chat') {
        sendChat(prepareMessagesForChatRequest(toSubmit))
      }
    },
    [
      config.mode,
      editingMessageKey,
      messages,
      updateMessages,
      sendChat,
      prepareMessagesForChatRequest,
    ]
  )

  const handleDeleteMessage = (message: MessageType) => {
    const newMessages = messages.filter((m) => m.key !== message.key)
    updateMessages(newMessages)
  }

  return (
    <div className='relative flex size-full overflow-hidden'>
      <PlaygroundSessions
        sessions={sessions}
        activeSessionId={activeSessionId}
        onCreateSession={createSession}
        onSwitchSession={switchSession}
        onDeleteSession={deleteSession}
      />
      {/* Full-width scroll container: scrolling works even over side whitespace */}
      <div className='flex min-w-0 flex-1 flex-col overflow-hidden'>
        <PlaygroundChat
          messages={messages}
          onCopyMessage={handleCopyMessage}
          onRegenerateMessage={handleRegenerateMessage}
          onEditMessage={handleEditMessage}
          onDeleteMessage={handleDeleteMessage}
          isGenerating={isGenerating}
          editingKey={editingMessageKey}
          onCancelEdit={handleEditOpenChange}
          onSaveEdit={(newContent) => applyEdit(newContent, false)}
          onSaveEditAndSubmit={(newContent) => applyEdit(newContent, true)}
          selectedImageUrl={selectedImageUrl}
          onSelectGeneratedImage={(image) => updateSelectedImage(image)}
          onSelectAttachmentImage={(image) => updateSelectedImage(image)}
        />

        {/* Input area: center content and constrain to the same container width */}
        <div className='mx-auto w-full max-w-4xl'>
          <PlaygroundInput
            disabled={isGenerating || isGeneratingImage}
            groups={groups}
            groupValue={config.group}
            isGenerating={isGenerating || isGeneratingImage}
            isImageMode={config.mode === 'image'}
            isEditingImage={!!selectedImageUrl}
            isModelLoading={isLoadingModels}
            modelValue={config.model}
            models={filteredModels}
            onGroupChange={(value) => updateConfig('group', value)}
            onModelChange={(value) => updateConfig('model', value)}
            onStop={config.mode === 'chat' ? stopGeneration : undefined}
            onSubmit={handleSendMessage}
          />
        </div>
      </div>
      <PlaygroundSettings
        config={config}
        disabled={isGenerating || isGeneratingImage}
        selectedImageUrl={selectedImageUrl}
        onConfigChange={updateConfig}
        onClearSelectedImage={() => updateSelectedImage(null)}
        onSelectImage={updateSelectedImage}
      />
    </div>
  )
}
