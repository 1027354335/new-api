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
import { MessageSquareIcon, PlusIcon, Trash2Icon } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
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

export function PlaygroundSessions({
  sessions,
  activeSessionId,
  onCreateSession,
  onSwitchSession,
  onDeleteSession,
}: PlaygroundSessionsProps) {
  const { t } = useTranslation()

  return (
    <aside className='border-border/70 bg-background/95 hidden w-64 shrink-0 border-r px-3 py-4 lg:flex lg:flex-col lg:gap-4'>
      <Button
        type='button'
        className='w-full justify-start'
        onClick={onCreateSession}
      >
        <PlusIcon className='size-4' />
        {t('New conversation')}
      </Button>

      <div className='min-h-0 flex-1'>
        <div className='text-muted-foreground mb-2 px-2 text-xs font-medium'>
          {t('Conversation history')}
        </div>
        <div className='grid max-h-full gap-1 overflow-y-auto overflow-x-hidden pr-1'>
          {sessions.map((session) => {
            const isActive = session.id === activeSessionId
            const storagePercent = getSessionStoragePercent(session)
            return (
              <div
                key={session.id}
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
              </div>
            )
          })}
        </div>
      </div>
    </aside>
  )
}
