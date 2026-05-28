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
import { useRef, useState, useMemo } from 'react'
import {
  PaperclipIcon,
  FileIcon,
  ImageIcon,
  ScreenShareIcon,
  CameraIcon,
  GlobeIcon,
  SendIcon,
  SquareIcon,
  BarChartIcon,
  BoxIcon,
  NotepadTextIcon,
  CodeSquareIcon,
  GraduationCapIcon,
} from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { cn } from '@/lib/utils'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  PromptInput,
  PromptInputAttachment,
  PromptInputAttachments,
  PromptInputButton,
  PromptInputFooter,
  PromptInputTextarea,
  PromptInputTools,
  usePromptInputAttachments,
  type PromptInputMessage,
} from '@/components/ai-elements/prompt-input'
import { Suggestion, Suggestions } from '@/components/ai-elements/suggestion'
import { ModelGroupSelector } from '@/components/model-group-selector'
import { extractAttachmentText } from '../lib'
import type { ModelOption, GroupOption, PlaygroundAttachment } from '../types'

interface PlaygroundInputProps {
  onSubmit: (
    text: string,
    attachments?: PlaygroundAttachment[]
  ) => void | Promise<void>
  onStop?: () => void
  disabled?: boolean
  isGenerating?: boolean
  isImageMode?: boolean
  isEditingImage?: boolean
  models: ModelOption[]
  modelValue: string
  onModelChange: (value: string) => void
  isModelLoading?: boolean
  groups: GroupOption[]
  groupValue: string
  onGroupChange: (value: string) => void
  hasMessages?: boolean
  inputRef?: React.RefObject<HTMLTextAreaElement | null>
}

function compressImageToThumbnail(dataUrl: string, maxDim = 100): Promise<string> {
  return new Promise((resolve) => {
    const img = new Image()
    img.onload = () => {
      const canvas = document.createElement('canvas')
      const ratio = Math.min(maxDim / img.width, maxDim / img.height, 1)
      canvas.width = img.width * ratio
      canvas.height = img.height * ratio
      const ctx = canvas.getContext('2d')
      if (ctx) {
        ctx.drawImage(img, 0, 0, canvas.width, canvas.height)
        resolve(canvas.toDataURL('image/jpeg', 0.6))
      } else {
        resolve(dataUrl)
      }
    }
    img.onerror = () => {
      resolve(dataUrl)
    }
    img.src = dataUrl
  })
}

const suggestions = [
  { icon: BarChartIcon, text: 'Analyze data', color: '#76d0eb' },
  { icon: BoxIcon, text: 'Surprise me', color: '#76d0eb' },
  { icon: NotepadTextIcon, text: 'Summarize text', color: '#ea8444' },
  { icon: CodeSquareIcon, text: 'Code', color: '#6c71ff' },
  { icon: GraduationCapIcon, text: 'Get advice', color: '#76d0eb' },
  { icon: null, text: 'More' },
]

async function filePartsToAttachments(
  files: PromptInputMessage['files'] = []
): Promise<PlaygroundAttachment[]> {
  return Promise.all(
    files
      .filter((file) => !!file.url)
      .map(async (file) => {
        let thumbnailUrl: string | undefined
        if (file.mediaType?.startsWith('image/') && file.url) {
          try {
            thumbnailUrl = await compressImageToThumbnail(file.url)
          } catch {
            // ignore
          }
        }
        return {
          id: 'id' in file && typeof file.id === 'string' ? file.id : file.url!,
          url: file.url!,
          name: file.filename,
          mediaType: file.mediaType,
          textContent: await extractAttachmentText(file),
          thumbnailUrl,
        }
      })
  )
}

function AttachmentMenu({ disabled }: { disabled?: boolean }) {
  const { t } = useTranslation()
  const attachments = usePromptInputAttachments()
  const fileInputRef = useRef<HTMLInputElement>(null)
  const imageInputRef = useRef<HTMLInputElement>(null)
  const cameraInputRef = useRef<HTMLInputElement>(null)

  const addFiles = (files: FileList | null) => {
    if (!files?.length) return
    attachments.add(files)
  }

  return (
    <>
      <input
        ref={fileInputRef}
        type='file'
        className='hidden'
        multiple
        onChange={(event) => addFiles(event.currentTarget.files)}
      />
      <input
        ref={imageInputRef}
        type='file'
        className='hidden'
        accept='image/*'
        multiple
        onChange={(event) => addFiles(event.currentTarget.files)}
      />
      <input
        ref={cameraInputRef}
        type='file'
        className='hidden'
        accept='image/*'
        capture='environment'
        onChange={(event) => addFiles(event.currentTarget.files)}
      />
      <DropdownMenu>
        <DropdownMenuTrigger
          render={
            <PromptInputButton
              className='border font-medium'
              disabled={disabled}
              variant='outline'
            />
          }
        >
          <PaperclipIcon size={16} />
          <span className='hidden sm:inline'>{t('Attach')}</span>
          <span className='sr-only sm:hidden'>{t('Attach')}</span>
        </DropdownMenuTrigger>
        <DropdownMenuContent align='start'>
          <DropdownMenuItem
            onSelect={(event) => {
              event.preventDefault()
              fileInputRef.current?.click()
            }}
          >
            <FileIcon className='mr-2' size={16} />
            {t('Upload file')}
          </DropdownMenuItem>
          <DropdownMenuItem
            onSelect={(event) => {
              event.preventDefault()
              imageInputRef.current?.click()
            }}
          >
            <ImageIcon className='mr-2' size={16} />
            {t('Upload photo')}
          </DropdownMenuItem>
          <DropdownMenuItem
            onSelect={(event) => {
              event.preventDefault()
              imageInputRef.current?.click()
            }}
          >
            <ScreenShareIcon className='mr-2' size={16} />
            {t('Take screenshot')}
          </DropdownMenuItem>
          <DropdownMenuItem
            onSelect={(event) => {
              event.preventDefault()
              cameraInputRef.current?.click()
            }}
          >
            <CameraIcon className='mr-2' size={16} />
            {t('Take photo')}
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </>
  )
}

function SubmitButton({
  disabled,
  text,
  isImageMode,
}: {
  disabled?: boolean
  text: string
  isImageMode: boolean
}) {
  const { t } = useTranslation()
  const attachments = usePromptInputAttachments()
  const canSubmit = text.trim() || attachments.files.length > 0

  return (
    <PromptInputButton
      className='text-foreground font-medium'
      disabled={disabled || !canSubmit}
      type='submit'
      variant='secondary'
    >
      <SendIcon size={16} />
      <span className='hidden sm:inline'>
        {t(isImageMode ? 'Generate' : 'Send')}
      </span>
      <span className='sr-only sm:hidden'>
        {t(isImageMode ? 'Generate' : 'Send')}
      </span>
    </PromptInputButton>
  )
}

export function PlaygroundInput({
  onSubmit,
  onStop,
  disabled,
  isGenerating,
  isImageMode = false,
  isEditingImage = false,
  models,
  modelValue,
  onModelChange,
  isModelLoading = false,
  groups,
  groupValue,
  onGroupChange,
  hasMessages = false,
  inputRef,
}: PlaygroundInputProps) {
  const { t } = useTranslation()
  const [text, setText] = useState('')
  const [isParsing, setIsParsing] = useState(false)
  const historyIndexRef = useRef(-1)

  const isInputDisabled = disabled || isParsing
  const isModelSelectDisabled =
    isInputDisabled || isModelLoading || models.length === 0
  const isGroupSelectDisabled = isInputDisabled || groups.length === 0

  const estimatedTokens = useMemo(() => {
    if (!text.trim()) return 0
    const cjk = (text.match(/[\u4e00-\u9fff\u3040-\u309f\u30a0-\u30ff]/g) || []).length
    return Math.ceil((text.length - cjk) / 4 + cjk / 2)
  }, [text])

  const handleSubmit = async (message: PromptInputMessage) => {
    setIsParsing(true)
    try {
      const attachments = await filePartsToAttachments(message.files)
      if ((!message.text?.trim() && attachments.length === 0) || isInputDisabled) return

      // Save user message to history
      if (message.text?.trim()) {
        const textVal = message.text.trim()
        try {
          const savedHistory = JSON.parse(localStorage.getItem('playground_msg_history') || '[]')
          const updated = [textVal, ...savedHistory.filter((m: string) => m !== textVal)].slice(0, 20)
          localStorage.setItem('playground_msg_history', JSON.stringify(updated))
        } catch {
          // ignore
        }
      }

      const result = await onSubmit(message.text ?? '', attachments)
      setText('')
      return result
    } finally {
      setIsParsing(false)
    }
  }

  const handleSuggestionClick = (suggestion: string) => {
    onSubmit(suggestion)
  }

  return (
    <div className='grid shrink-0 gap-4 px-1 md:pb-4'>
      <PromptInput
        globalDrop
        groupClassName='rounded-xl'
        multiple
        onError={(error) => toast.error(error.message)}
        onSubmit={handleSubmit}
      >
        <div className='flex flex-wrap gap-2 px-3 pt-3'>
          <PromptInputAttachments>
            {(attachment) => <PromptInputAttachment data={attachment} />}
          </PromptInputAttachments>
          {isParsing && (
            <div className='flex items-center space-x-2 p-1.5 bg-muted/40 border border-border/50 rounded-md animate-pulse text-xs text-muted-foreground select-none'>
              <div className='w-2 h-2 rounded-full bg-primary animate-ping' />
              <span>{t('Extracting attachments...')}</span>
            </div>
          )}
        </div>
        <PromptInputTextarea
          ref={inputRef}
          autoComplete='off'
          autoCorrect='off'
          autoCapitalize='off'
          spellCheck={false}
          className='px-5 md:text-base'
          disabled={isInputDisabled}
          onChange={(event) => {
            setText(event.target.value)
            historyIndexRef.current = -1
          }}
          onKeyDown={(e) => {
            if (e.key === 'ArrowUp' && !text.trim()) {
              e.preventDefault()
              try {
                const history = JSON.parse(localStorage.getItem('playground_msg_history') || '[]')
                if (history.length > 0) {
                  historyIndexRef.current = Math.min(historyIndexRef.current + 1, history.length - 1)
                  const histVal = history[historyIndexRef.current]
                  if (histVal) setText(histVal)
                }
              } catch {
                // ignore
              }
            } else if (e.key === 'ArrowDown' && historyIndexRef.current >= 0) {
              e.preventDefault()
              try {
                const history = JSON.parse(localStorage.getItem('playground_msg_history') || '[]')
                historyIndexRef.current -= 1
                if (historyIndexRef.current >= 0) {
                  setText(history[historyIndexRef.current])
                } else {
                  setText('')
                }
              } catch {
                // ignore
              }
            }
          }}
          placeholder={
            isImageMode
              ? isEditingImage
                ? t('Describe how to edit the selected image')
                : t('Describe the image you want to generate')
              : t('Ask anything')
          }
          value={text}
        />

        <PromptInputFooter className='p-2.5'>
          <PromptInputTools>
            <AttachmentMenu disabled={isInputDisabled} />

            <PromptInputButton
              className='border font-medium'
              disabled={isInputDisabled}
              onClick={() => toast.info(t('Search feature in development'))}
              variant='outline'
            >
              <GlobeIcon size={16} />
              <span className='hidden sm:inline'>{t('Search')}</span>
              <span className='sr-only sm:hidden'>{t('Search')}</span>
            </PromptInputButton>
          </PromptInputTools>

          <div className='flex items-center gap-1.5 md:gap-2'>
            {estimatedTokens > 0 && (
              <span
                className={cn(
                  'text-[10px] text-muted-foreground select-none transition-colors mr-2 self-center font-mono',
                  estimatedTokens > 4000 && 'text-amber-500',
                  estimatedTokens > 8000 && 'text-destructive'
                )}
              >
                ~{estimatedTokens} tokens
              </span>
            )}
            <ModelGroupSelector
              selectedModel={modelValue}
              models={models}
              onModelChange={onModelChange}
              selectedGroup={groupValue}
              groups={groups}
              onGroupChange={onGroupChange}
              disabled={isModelSelectDisabled || isGroupSelectDisabled}
            />

            {isGenerating && onStop ? (
              <PromptInputButton
                className='text-foreground font-medium'
                onClick={onStop}
                variant='secondary'
              >
                <SquareIcon className='fill-current' size={16} />
                <span className='hidden sm:inline'>{t('Stop')}</span>
                <span className='sr-only sm:hidden'>{t('Stop')}</span>
              </PromptInputButton>
            ) : (
              <SubmitButton
                disabled={isInputDisabled}
                isImageMode={isImageMode}
                text={text}
              />
            )}
          </div>
        </PromptInputFooter>
      </PromptInput>

      {!hasMessages && (
        <Suggestions>
          {suggestions.map(({ icon: Icon, text, color }) => (
            <Suggestion
              className={`text-xs font-normal sm:text-sm ${
                text === 'More' ? 'hidden sm:flex' : ''
              }`}
              key={text}
              onClick={() => handleSuggestionClick(text)}
              suggestion={text}
            >
              {Icon && <Icon size={16} style={{ color }} />}
              {text}
            </Suggestion>
          ))}
        </Suggestions>
      )}
    </div>
  )
}
