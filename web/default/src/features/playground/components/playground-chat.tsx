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
import { useEffect, useMemo, useState } from 'react'
import {
  DownloadIcon,
  FileSpreadsheetIcon,
  FileTextIcon,
  ImageIcon,
  PresentationIcon,
} from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import {
  Branch,
  BranchMessages,
  BranchNext,
  BranchPage,
  BranchPrevious,
  BranchSelector,
} from '@/components/ai-elements/branch'
import {
  Conversation,
  ConversationContent,
  ConversationScrollButton,
} from '@/components/ai-elements/conversation'
import { Loader } from '@/components/ai-elements/loader'
import { Message, MessageContent } from '@/components/ai-elements/message'
import {
  Reasoning,
  ReasoningContent,
  ReasoningTrigger,
} from '@/components/ai-elements/reasoning'
import { Response } from '@/components/ai-elements/response'
import { Shimmer } from '@/components/ai-elements/shimmer'
import {
  Source,
  Sources,
  SourcesContent,
  SourcesTrigger,
} from '@/components/ai-elements/sources'
import { MESSAGE_ROLES } from '../constants'
import { getMessageContentStyles } from '../lib/message-styles'
import { parseThinkTags } from '../lib/message-utils'
import type {
  GeneratedFile,
  GeneratedImage,
  Message as MessageType,
} from '../types'
import { MessageActions } from './message-actions'
import { MessageError } from './message-error'

interface PlaygroundChatProps {
  messages: MessageType[]
  onCopyMessage?: (message: MessageType) => void
  onRegenerateMessage?: (message: MessageType) => void
  onEditMessage?: (message: MessageType) => void
  onDeleteMessage?: (message: MessageType) => void
  isGenerating?: boolean
  editingKey?: string | null
  onSaveEdit?: (newContent: string) => void
  onCancelEdit?: (open: boolean) => void
  onSaveEditAndSubmit?: (newContent: string) => void
  selectedImageUrl?: string | null
  onSelectGeneratedImage?: (image: GeneratedImage) => void
  onSelectAttachmentImage?: (image: GeneratedImage) => void
}

export function PlaygroundChat({
  messages,
  onCopyMessage,
  onRegenerateMessage,
  onEditMessage,
  onDeleteMessage,
  isGenerating = false,
  editingKey,
  onSaveEdit,
  onCancelEdit,
  onSaveEditAndSubmit,
  selectedImageUrl,
  onSelectGeneratedImage,
  onSelectAttachmentImage,
}: PlaygroundChatProps) {
  const { t } = useTranslation()
  const [editText, setEditText] = useState('')
  const [originalText, setOriginalText] = useState('')

  useEffect(() => {
    if (!editingKey) return
    const message = messages.find((m) => m.key === editingKey)
    const content = message?.versions?.[0]?.content || ''
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setEditText(content)

    setOriginalText(content)
  }, [editingKey, messages])

  const isEditing = (key: string) => editingKey === key
  const isEmpty = useMemo(() => !editText.trim(), [editText])
  const isChanged = useMemo(
    () => editText !== originalText,
    [editText, originalText]
  )
  const getFileIcon = (file: GeneratedFile) => {
    if (file.kind === 'excel') return FileSpreadsheetIcon
    if (file.kind === 'powerpoint') return PresentationIcon
    return FileTextIcon
  }

  return (
    <Conversation>
      {/* Remove outer padding; apply padding to inner centered container to align with input */}
      <ConversationContent className='p-0'>
        <div className='mx-auto w-full max-w-4xl px-4 py-4'>
          {messages.map((message, messageIndex) => {
            const { versions = [] } = message
            const isLastAssistantMessage =
              messageIndex === messages.length - 1 &&
              message.from === MESSAGE_ROLES.ASSISTANT
            return (
              <Branch defaultBranch={0} key={message.key}>
                <BranchMessages>
                  {versions.map((version, versionIndex) => (
                    <Message
                      className='group flex-row-reverse'
                      from={message.from}
                      key={`${message.key}-${version.id}-${versionIndex}`}
                    >
                      <div className='w-full min-w-0 flex-1 basis-full py-1'>
                        {isEditing(message.key) ? (
                          <div className='space-y-2'>
                            <Textarea
                              value={editText}
                              onChange={(e) => setEditText(e.target.value)}
                              className='font-mono text-sm'
                              rows={8}
                            />
                            <div className='flex gap-2'>
                              {/* Save & Submit only makes sense for user messages */}
                              {message.from === MESSAGE_ROLES.USER && (
                                <Button
                                  size='sm'
                                  onClick={() =>
                                    onSaveEditAndSubmit?.(editText)
                                  }
                                  disabled={isEmpty || !isChanged}
                                >
                                  {t('Save & Submit')}
                                </Button>
                              )}
                              <Button
                                size='sm'
                                onClick={() => onSaveEdit?.(editText)}
                                disabled={isEmpty || !isChanged}
                              >
                                {t('Save')}
                              </Button>
                              <Button
                                size='sm'
                                variant='outline'
                                onClick={() => onCancelEdit?.(false)}
                              >
                                {t('Cancel')}
                              </Button>
                            </div>
                          </div>
                        ) : (
                          <>
                            {(() => {
                              const isAssistant =
                                message.from === MESSAGE_ROLES.ASSISTANT
                              const hasSources = !!message.sources?.length
                              const showReasoning =
                                isAssistant && !!message.reasoning?.content
                              const showLoader =
                                isAssistant &&
                                !message.isReasoningStreaming &&
                                (message.status === 'loading' ||
                                  (message.status === 'streaming' &&
                                    !version.content))
                              const showMessageContent =
                                (message.from === MESSAGE_ROLES.USER ||
                                  !message.isReasoningStreaming) &&
                                !!version.content
                              const generatedFiles =
                                message.generatedFiles ?? []

                              // Extract visible content (remove <think> tags for assistant messages)
                              const displayContent = isAssistant
                                ? parseThinkTags(version.content).visibleContent
                                : version.content

                              const actions = (
                                <MessageActions
                                  message={message}
                                  onCopy={onCopyMessage}
                                  onRegenerate={onRegenerateMessage}
                                  onEdit={onEditMessage}
                                  onDelete={onDeleteMessage}
                                  isGenerating={isGenerating}
                                  alwaysVisible={isLastAssistantMessage}
                                  className='mt-1'
                                />
                              )
                              const imageAttachments =
                                message.attachments?.filter((attachment) =>
                                  attachment.mediaType?.startsWith('image/')
                                ) ?? []
                              const fileAttachments =
                                message.attachments?.filter(
                                  (attachment) =>
                                    !attachment.mediaType?.startsWith('image/')
                                ) ?? []

                              return (
                                <>
                                  {/* Sources */}
                                  {hasSources && (
                                    <Sources>
                                      <SourcesTrigger
                                        count={message.sources!.length}
                                      />
                                      <SourcesContent>
                                        {message.sources!.map(
                                          (source, sourceIndex) => (
                                            <Source
                                              href={source.href}
                                              key={`${message.key}-source-${sourceIndex}`}
                                              title={source.title}
                                            />
                                          )
                                        )}
                                      </SourcesContent>
                                    </Sources>
                                  )}

                                  {/* Reasoning */}
                                  {showReasoning && (
                                    <Reasoning
                                      defaultOpen={true}
                                      isStreaming={message.isReasoningStreaming}
                                    >
                                      <ReasoningTrigger />
                                      <ReasoningContent>
                                        {message.reasoning!.content}
                                      </ReasoningContent>
                                    </Reasoning>
                                  )}

                                  {/* Loader */}
                                  {showLoader && (
                                    <div className='flex items-center gap-2 py-2'>
                                      <Loader />
                                      <Shimmer className='text-sm' duration={1}>
                                        {t('Responding...')}
                                      </Shimmer>
                                    </div>
                                  )}

                                  {/* Error or Content */}
                                  {message.status === 'error' ? (
                                    <>
                                      <MessageError
                                        message={message}
                                        className='mb-2'
                                      />
                                      {actions}
                                    </>
                                  ) : (
                                    (showMessageContent ||
                                      imageAttachments.length > 0 ||
                                      fileAttachments.length > 0 ||
                                      generatedFiles.length > 0) && (
                                      <>
                                        {imageAttachments.length > 0 && (
                                          <div className='mb-3 grid grid-cols-2 gap-2 sm:grid-cols-3'>
                                            {imageAttachments.map(
                                              (attachment) => (
                                                <button
                                                  type='button'
                                                  key={attachment.id}
                                                  onClick={() =>
                                                    onSelectAttachmentImage?.({
                                                      id: attachment.id,
                                                      url: attachment.url,
                                                      revisedPrompt:
                                                        version.content,
                                                    })
                                                  }
                                                  className={cn(
                                                    'border-border bg-muted/20 hover:border-primary/50 focus-visible:border-ring focus-visible:ring-ring/50 relative overflow-hidden rounded-lg border text-left transition focus-visible:ring-3 focus-visible:outline-hidden',
                                                    selectedImageUrl ===
                                                      attachment.url &&
                                                      'border-primary ring-primary/20 ring-2'
                                                  )}
                                                >
                                                  <img
                                                    src={attachment.url}
                                                    alt={
                                                      attachment.name ||
                                                      t('Attached image')
                                                    }
                                                    className='aspect-square w-full object-cover'
                                                  />
                                                  <span className='bg-background/90 text-foreground absolute right-2 bottom-2 inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs shadow-sm'>
                                                    <ImageIcon className='size-3' />
                                                    {selectedImageUrl ===
                                                    attachment.url
                                                      ? t('Selected')
                                                      : t('Use for edit')}
                                                  </span>
                                                </button>
                                              )
                                            )}
                                          </div>
                                        )}
                                        {fileAttachments.length > 0 && (
                                          <div className='mb-3 flex flex-wrap gap-2'>
                                            {fileAttachments.map(
                                              (attachment) => (
                                                <span
                                                  key={attachment.id}
                                                  className='border-border bg-muted/40 text-muted-foreground inline-flex max-w-56 items-center gap-1.5 rounded-md border px-2 py-1 text-xs'
                                                >
                                                  <span className='truncate'>
                                                    {attachment.name ||
                                                      t('Attachment')}
                                                  </span>
                                                </span>
                                              )
                                            )}
                                          </div>
                                        )}
                                        {generatedFiles.length > 0 && (
                                          <div className='mb-3 flex flex-col gap-2'>
                                            {generatedFiles.map((file) => {
                                              const FileIcon =
                                                getFileIcon(file)
                                              return (
                                                <div
                                                  key={file.id}
                                                  className='border-border bg-muted/30 flex min-w-0 items-center justify-between gap-3 rounded-lg border px-3 py-2'
                                                >
                                                  <div className='flex min-w-0 items-center gap-2'>
                                                    <FileIcon className='text-muted-foreground size-4 shrink-0' />
                                                    <div className='min-w-0'>
                                                      <div className='truncate text-sm font-medium'>
                                                        {file.name}
                                                      </div>
                                                      {!!file.size && (
                                                        <div className='text-muted-foreground text-xs'>
                                                          {Math.ceil(
                                                            file.size / 1024
                                                          )}{' '}
                                                          KB
                                                        </div>
                                                      )}
                                                    </div>
                                                  </div>
                                                  <a
                                                    href={file.url}
                                                    download={file.name}
                                                    className='border-input bg-background hover:bg-accent hover:text-accent-foreground inline-flex h-8 shrink-0 items-center justify-center gap-1.5 rounded-md border px-3 text-sm font-medium transition-colors'
                                                  >
                                                    <DownloadIcon className='size-4' />
                                                    {t('Download')}
                                                  </a>
                                                </div>
                                              )
                                            })}
                                          </div>
                                        )}
                                        {!!message.generatedImages?.length && (
                                          <div className='mb-3 grid grid-cols-1 gap-3 sm:grid-cols-2'>
                                            {message.generatedImages.map(
                                              (image, imageIndex) => (
                                                <button
                                                  type='button'
                                                  key={image.id}
                                                  onClick={() =>
                                                    onSelectGeneratedImage?.(
                                                      image
                                                    )
                                                  }
                                                  className={cn(
                                                    'border-border bg-muted/20 group/image hover:border-primary/50 focus-visible:border-ring focus-visible:ring-ring/50 relative overflow-hidden rounded-lg border text-left transition focus-visible:ring-3 focus-visible:outline-hidden',
                                                    selectedImageUrl ===
                                                      image.url &&
                                                      'border-primary ring-primary/20 ring-2'
                                                  )}
                                                >
                                                  <img
                                                    src={image.url}
                                                    alt={`${t('Generated image')} ${imageIndex + 1}`}
                                                    className='aspect-square w-full object-cover'
                                                  />
                                                  <span className='bg-background/90 text-foreground absolute right-2 bottom-2 inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs shadow-sm'>
                                                    <ImageIcon className='size-3' />
                                                    {selectedImageUrl ===
                                                    image.url
                                                      ? t('Selected')
                                                      : t('Use for edit')}
                                                  </span>
                                                </button>
                                              )
                                            )}
                                          </div>
                                        )}
                                        <MessageContent
                                          variant='flat'
                                          className={cn(
                                            getMessageContentStyles()
                                          )}
                                        >
                                          <Response>{displayContent}</Response>
                                        </MessageContent>
                                        {actions}
                                      </>
                                    )
                                  )}
                                </>
                              )
                            })()}
                          </>
                        )}
                      </div>
                    </Message>
                  ))}
                </BranchMessages>

                {/* Branch selector for multiple versions */}
                {versions.length > 1 && (
                  <BranchSelector className='px-0' from={message.from}>
                    <BranchPrevious />
                    <BranchPage />
                    <BranchNext />
                  </BranchSelector>
                )}
              </Branch>
            )
          })}
        </div>
      </ConversationContent>
      <ConversationScrollButton />
    </Conversation>
  )
}
