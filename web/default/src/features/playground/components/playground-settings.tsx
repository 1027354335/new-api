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
import { useState, useRef, type ReactNode } from 'react'
import {
  ImageIcon,
  MessageSquareTextIcon,
  SlidersHorizontalIcon,
  UploadIcon,
  Loader2Icon,
} from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { uploadPlaygroundImage } from '../api'
import type { PlaygroundConfig, PlaygroundMode } from '../types'

interface PlaygroundSettingsProps {
  config: PlaygroundConfig
  disabled?: boolean
  selectedImageUrl?: string | null
  onConfigChange: <K extends keyof PlaygroundConfig>(
    key: K,
    value: PlaygroundConfig[K]
  ) => void
  onClearSelectedImage: () => void
  onSelectImage?: (image: { id: string; url: string }) => void
}

const modeOptions: Array<{
  value: PlaygroundMode
  label: string
  icon: typeof MessageSquareTextIcon
}> = [
  { value: 'chat', label: 'Chat', icon: MessageSquareTextIcon },
  { value: 'image', label: 'Image', icon: ImageIcon },
]

const sizeOptions = ['1024x1024', '1024x1536', '1536x1024', 'auto']
const qualityOptions = ['auto', 'low', 'medium', 'high', 'standard', 'hd']
const moderationOptions = ['auto', 'low']

function SettingLabel({ children }: { children: ReactNode }) {
  return (
    <div className='text-muted-foreground text-xs font-medium'>{children}</div>
  )
}

export function PlaygroundSettings({
  config,
  disabled,
  selectedImageUrl,
  onConfigChange,
  onClearSelectedImage,
  onSelectImage,
}: PlaygroundSettingsProps) {
  const { t } = useTranslation()
  const isImageMode = config.mode === 'image'

  const [isUploading, setIsUploading] = useState(false)
  const [isDragOver, setIsDragOver] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleUpload = async (file: File) => {
    if (!file.type.startsWith('image/')) {
      toast.error(t('Only image files are supported'))
      return
    }

    setIsUploading(true)
    try {
      const url = await uploadPlaygroundImage(file)
      onSelectImage?.({
        id: Math.random().toString(36).substring(2, 9),
        url,
      })
      toast.success(t('Image uploaded successfully'))
    } catch (error: unknown) {
      const message = error instanceof Error ? error.message : ''
      toast.error(message || t('Upload failed'))
    } finally {
      setIsUploading(false)
      if (fileInputRef.current) {
        fileInputRef.current.value = ''
      }
    }
  }

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      handleUpload(file)
    }
  }

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault()
    if (disabled || isUploading) return
    setIsDragOver(true)
  }

  const handleDragLeave = () => {
    setIsDragOver(false)
  }

  const handleDrop = async (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragOver(false)
    if (disabled || isUploading) return

    const file = e.dataTransfer.files?.[0]
    if (file) {
      await handleUpload(file)
    }
  }

  return (
    <aside className='border-border/70 bg-background/95 hidden w-72 shrink-0 border-l px-4 py-4 lg:flex lg:flex-col lg:gap-5'>
      <div className='flex items-center gap-2'>
        <SlidersHorizontalIcon className='text-muted-foreground size-4' />
        <h2 className='text-sm font-semibold'>{t('Parameters')}</h2>
      </div>

      <div className='grid grid-cols-2 gap-2'>
        {modeOptions.map(({ value, label, icon: Icon }) => (
          <Button
            key={value}
            type='button'
            variant={config.mode === value ? 'secondary' : 'outline'}
            className={cn(
              'justify-start',
              config.mode === value && 'border-primary/30'
            )}
            onClick={() => onConfigChange('mode', value)}
            disabled={disabled}
          >
            <Icon className='size-4' />
            {t(label)}
          </Button>
        ))}
      </div>

      {isImageMode ? (
        <div className='grid gap-4'>
          <div className='grid gap-2'>
            <SettingLabel>{t('Image size')}</SettingLabel>
            <Select
              value={config.imageSize}
              onValueChange={(value) => {
                if (value) onConfigChange('imageSize', value)
              }}
              disabled={disabled}
            >
              <SelectTrigger className='w-full'>
                <SelectValue />
              </SelectTrigger>
              <SelectContent alignItemWithTrigger={false}>
                <SelectGroup>
                  {sizeOptions.map((value) => (
                    <SelectItem key={value} value={value}>
                      {value}
                    </SelectItem>
                  ))}
                </SelectGroup>
              </SelectContent>
            </Select>
          </div>

          <div className='grid gap-2'>
            <SettingLabel>{t('Image quality')}</SettingLabel>
            <Select
              value={config.imageQuality}
              onValueChange={(value) => {
                if (value) onConfigChange('imageQuality', value)
              }}
              disabled={disabled}
            >
              <SelectTrigger className='w-full'>
                <SelectValue />
              </SelectTrigger>
              <SelectContent alignItemWithTrigger={false}>
                <SelectGroup>
                  {qualityOptions.map((value) => (
                    <SelectItem key={value} value={value}>
                      {t(value)}
                    </SelectItem>
                  ))}
                </SelectGroup>
              </SelectContent>
            </Select>
          </div>

          <div className='grid gap-2'>
            <SettingLabel>{t('Sensitivity')}</SettingLabel>
            <Select
              value={config.imageModeration}
              onValueChange={(value) => {
                if (value) onConfigChange('imageModeration', value)
              }}
              disabled={disabled}
            >
              <SelectTrigger className='w-full'>
                <SelectValue />
              </SelectTrigger>
              <SelectContent alignItemWithTrigger={false}>
                <SelectGroup>
                  {moderationOptions.map((value) => (
                    <SelectItem key={value} value={value}>
                      {t(value)}
                    </SelectItem>
                  ))}
                </SelectGroup>
              </SelectContent>
            </Select>
          </div>

          <div className='grid gap-2'>
            <SettingLabel>{t('Image count')}</SettingLabel>
            <Input
              type='number'
              min={1}
              max={10}
              value={config.imageCount}
              disabled={disabled}
              onChange={(event) =>
                onConfigChange(
                  'imageCount',
                  Math.min(10, Math.max(1, Number(event.target.value) || 1))
                )
              }
            />
          </div>

          {selectedImageUrl ? (
            <div className='border-border bg-muted/30 grid gap-3 rounded-lg border p-3'>
              <div className='text-sm font-medium'>{t('Edit source')}</div>
              <img
                src={selectedImageUrl}
                alt={t('Selected image')}
                className='aspect-square w-full rounded-md object-cover'
              />
              <Button
                type='button'
                variant='outline'
                onClick={onClearSelectedImage}
                disabled={disabled}
              >
                {t('Clear image')}
              </Button>
            </div>
          ) : (
            <div
              className={cn(
                'border-border/70 hover:border-primary/30 bg-muted/10 grid cursor-pointer gap-2.5 rounded-lg border border-dashed p-4 text-center transition-colors',
                isDragOver && 'border-primary bg-primary/5'
              )}
              onDragOver={handleDragOver}
              onDragLeave={handleDragLeave}
              onDrop={handleDrop}
              onClick={() => {
                if (!disabled && !isUploading) {
                  fileInputRef.current?.click()
                }
              }}
            >
              <input
                type='file'
                ref={fileInputRef}
                className='hidden'
                accept='image/*'
                onChange={handleFileChange}
                disabled={disabled || isUploading}
              />
              <div className='text-muted-foreground flex flex-col items-center justify-center gap-1.5 py-2'>
                {isUploading ? (
                  <>
                    <Loader2Icon className='text-primary size-5 animate-spin' />
                    <span className='text-xs font-medium'>
                      {t('Uploading...')}
                    </span>
                  </>
                ) : (
                  <>
                    <UploadIcon className='text-muted-foreground/75 size-5' />
                    <span className='text-foreground/90 text-xs font-semibold'>
                      {t('Upload source image')}
                    </span>
                    <span className='text-muted-foreground/70 text-[10px]'>
                      {t('Click or drag image file')}
                    </span>
                  </>
                )}
              </div>
            </div>
          )}
        </div>
      ) : (
        <div className='text-muted-foreground text-sm'>
          {t('Chat parameters use the saved defaults.')}
        </div>
      )}
    </aside>
  )
}
