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
  MessageSquareIcon,
  PlusIcon,
  Trash2Icon,
  MoreHorizontalIcon,
  DownloadIcon,
  UploadIcon,
} from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { nanoid } from 'nanoid'
import { motion, AnimatePresence } from 'framer-motion'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  estimatePersistedSessionBytes,
  MAX_PERSISTED_SESSION_BYTES,
} from '../lib'
import type { PlaygroundSession } from '../types'

interface PlaygroundSessionsProps {
  sessions: PlaygroundSession[]
  activeSessionId: string
  onCreateSession: () => void
  onSwitchSession: (sessionId: string) => void
  onDeleteSession: (sessionId: string) => void
  onImportSessions?: (sessions: PlaygroundSession[]) => void
}

function formatSessionTime(timestamp: number) {
  if (!timestamp) return ''
  return new Intl.DateTimeFormat(undefined, {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(new Date(timestamp))
}

function getSessionStoragePercent(session: PlaygroundSession) {
  const bytes = estimatePersistedSessionBytes(session)
  return Math.min(100, Math.max(0, Math.ceil((bytes / MAX_PERSISTED_SESSION_BYTES) * 100)))
}

function groupSessionsByDate(sessions: PlaygroundSession[], t: (key: string) => string) {
  const now = new Date()
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate()).getTime()
  const yesterday = today - 86400000
  const weekAgo = today - 7 * 86400000

  const groups: { label: string; sessions: PlaygroundSession[] }[] = [
    { label: t('Today'), sessions: [] },
    { label: t('Yesterday'), sessions: [] },
    { label: t('Last 7 Days'), sessions: [] },
    { label: t('Earlier'), sessions: [] },
  ]

  for (const s of sessions) {
    const ts = s.updatedAt || s.createdAt
    if (ts >= today) groups[0].sessions.push(s)
    else if (ts >= yesterday) groups[1].sessions.push(s)
    else if (ts >= weekAgo) groups[2].sessions.push(s)
    else groups[3].sessions.push(s)
  }

  return groups.filter((g) => g.sessions.length > 0)
}

export function PlaygroundSessions({
  sessions,
  activeSessionId,
  onCreateSession,
  onSwitchSession,
  onDeleteSession,
  onImportSessions,
}: PlaygroundSessionsProps) {
  const { t } = useTranslation()
  const [searchQuery, setSearchQuery] = useState('')
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleExportAll = () => {
    try {
      const data = {
        version: 1,
        exportedAt: new Date().toISOString(),
        sessions: sessions.map((s) => ({
          title: s.title,
          messages: s.messages.map((m) => ({
            from: m.from,
            content: m.versions[0]?.content || '',
            reasoning: m.reasoning?.content,
            attachments: m.attachments?.map((a) => ({
              id: a.id,
              name: a.name,
              mediaType: a.mediaType,
              textContent: a.textContent,
              url: a.url,
              thumbnailUrl: a.thumbnailUrl,
            })),
            generatedImages: m.generatedImages,
          })),
          createdAt: s.createdAt,
          updatedAt: s.updatedAt,
        })),
      }
      const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `playground-sessions-${new Date().toISOString().slice(0, 10)}.json`
      a.click()
      URL.revokeObjectURL(url)
      toast.success(t('Sessions exported successfully'))
    } catch {
      toast.error(t('Export failed'))
    }
  }

  const handleImport = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return
    try {
      const text = await file.text()
      const data = JSON.parse(text)
      if (data.version !== 1 || !Array.isArray(data.sessions)) {
        toast.error(t('Invalid export file format'))
        return
      }

      const imported = data.sessions.map((s: any) => {
        const now = Date.now()
        return {
          id: nanoid(),
          title: s.title || t('Imported conversation'),
          messages: s.messages.map((m: any) => ({
            key: nanoid(),
            from: m.from,
            versions: [{ id: nanoid(), content: m.content }],
            reasoning: m.reasoning ? { content: m.reasoning, duration: 0 } : undefined,
            attachments: m.attachments,
            generatedImages: m.generatedImages,
            status: 'complete',
          })),
          createdAt: s.createdAt || now,
          updatedAt: s.updatedAt || now,
        }
      })

      onImportSessions?.(imported)
      toast.success(t('Imported successfully'))
    } catch {
      toast.error(t('Import failed'))
    } finally {
      if (fileInputRef.current) fileInputRef.current.value = ''
    }
  }

  const filteredSessions = useMemo(() => {
    if (!searchQuery.trim()) return sessions
    const q = searchQuery.toLowerCase()
    return sessions.filter((s) =>
      s.title.toLowerCase().includes(q) ||
      s.messages.some((m) =>
        m.versions.some((v) => v.content.toLowerCase().includes(q))
      )
    )
  }, [sessions, searchQuery])

  const dateGroups = useMemo(() => {
    return groupSessionsByDate(filteredSessions, t)
  }, [filteredSessions, t])

  return (
    <aside className='border-border/70 bg-background/95 hidden w-64 shrink-0 border-r px-3 py-4 lg:flex lg:flex-col lg:gap-4'>
      <div className='flex gap-2 w-full'>
        <Button
          type='button'
          className='flex-1 justify-start'
          onClick={onCreateSession}
        >
          <PlusIcon className='size-4' />
          {t('New conversation')}
        </Button>
        <DropdownMenu>
          <DropdownMenuTrigger
            render={
              <Button variant='outline' size='icon' className='shrink-0' aria-label={t('More options')} />
            }
          >
            <MoreHorizontalIcon className='size-4' />
          </DropdownMenuTrigger>
          <DropdownMenuContent align='end'>
            <DropdownMenuItem onClick={handleExportAll}>
              <DownloadIcon className='mr-2 size-4' />
              {t('Export All Sessions')}
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => fileInputRef.current?.click()}>
              <UploadIcon className='mr-2 size-4' />
              {t('Import Sessions')}
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
        <input
          ref={fileInputRef}
          type='file'
          accept='.json'
          className='hidden'
          onChange={handleImport}
        />
      </div>

      <div className='relative px-1'>
        <input
          type='text'
          placeholder={t('Search conversation...')}
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className='w-full rounded-md border border-border bg-background px-3 py-1.5 text-xs text-foreground focus:outline-hidden focus:ring-1 focus:ring-primary'
        />
      </div>

      <div className='min-h-0 flex-1 flex flex-col gap-2'>
        <div className='text-muted-foreground px-2 text-xs font-medium flex justify-between items-center'>
          <span>{t('Conversation history')}</span>
          {sessions.length > 0 && (
            <span className='text-[10px] opacity-75'>
              {filteredSessions.length} {t('sessions')}
            </span>
          )}
        </div>
        <div className='flex-1 grid max-h-full gap-2 overflow-y-auto overflow-x-hidden pr-1 pb-4'>
          {dateGroups.map((group) => (
            <div key={group.label} className='space-y-1'>
              <div className='px-2 py-0.5 text-[9px] font-bold text-muted-foreground uppercase tracking-wider bg-muted/20 rounded-sm select-none'>
                {group.label}
              </div>
              <AnimatePresence initial={false}>
                {group.sessions.map((session) => {
                  const isActive = session.id === activeSessionId
                  const storagePercent = getSessionStoragePercent(session)
                  return (
                    <motion.div
                      key={session.id}
                      initial={{ opacity: 0, height: 0 }}
                      animate={{ opacity: 1, height: 'auto' }}
                      exit={{ opacity: 0, height: 0 }}
                      transition={{ duration: 0.15 }}
                      layout
                      className={cn(
                        'group/session relative min-w-0 max-w-full overflow-hidden rounded-lg',
                        isActive && 'bg-muted'
                      )}
                    >
                      <button
                        type='button'
                        className='hover:bg-muted flex w-full min-w-0 max-w-full items-start gap-2 overflow-hidden rounded-lg py-2 pr-9 pl-2 text-left text-sm'
                        onClick={() => onSwitchSession(session.id)}
                      >
                        <MessageSquareIcon className='text-muted-foreground mt-0.5 size-4 shrink-0' />
                        <span className='block min-w-0 flex-1 overflow-hidden'>
                          <span className='block truncate font-medium'>
                            {session.title || t('Untitled conversation')}
                          </span>
                          <span className='text-muted-foreground block truncate text-xs'>
                            {formatSessionTime(session.updatedAt)}
                          </span>
                          <span className='mt-1.5 flex min-w-0 max-w-full items-center gap-2 overflow-hidden'>
                            <span className='bg-muted-foreground/15 h-1.5 min-w-0 flex-1 overflow-hidden rounded-full'>
                              <span
                                className={cn(
                                  'block h-full rounded-full',
                                  storagePercent >= 90
                                    ? 'bg-destructive'
                                    : storagePercent >= 70
                                      ? 'bg-amber-500'
                                      : 'bg-primary'
                                )}
                                style={{ width: `${storagePercent}%` }}
                              />
                            </span>
                            <span className='text-muted-foreground w-8 text-right text-[10px] leading-none'>
                              {storagePercent}%
                            </span>
                          </span>
                        </span>
                      </button>
                      <Button
                        type='button'
                        variant='ghost'
                        size='icon-sm'
                        className='bg-background/90 absolute top-2 right-1 z-10 shrink-0 opacity-0 shadow-sm transition-opacity group-hover/session:opacity-100 focus-visible:opacity-100'
                        onClick={() => onDeleteSession(session.id)}
                        aria-label={t('Delete conversation')}
                      >
                        <Trash2Icon className='size-4' />
                      </Button>
                    </motion.div>
                  )
                })}
              </AnimatePresence>
            </div>
          ))}
          {filteredSessions.length === 0 && (
            <div className='text-center text-xs text-muted-foreground py-8 select-none'>
              {t('No conversations found')}
            </div>
          )}
        </div>
      </div>
    </aside>
  )
}
