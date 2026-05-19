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
        <div className='grid max-h-full gap-1 overflow-y-auto pr-1'>
          {sessions.map((session) => {
            const isActive = session.id === activeSessionId
            return (
              <div
                key={session.id}
                className={cn(
                  'group/session flex items-center gap-1 rounded-lg',
                  isActive && 'bg-muted'
                )}
              >
                <button
                  type='button'
                  className='hover:bg-muted flex min-w-0 flex-1 items-center gap-2 rounded-lg px-2 py-2 text-left text-sm'
                  onClick={() => onSwitchSession(session.id)}
                >
                  <MessageSquareIcon className='text-muted-foreground size-4 shrink-0' />
                  <span className='min-w-0 flex-1'>
                    <span className='block truncate font-medium'>
                      {session.title || t('Untitled conversation')}
                    </span>
                    <span className='text-muted-foreground block truncate text-xs'>
                      {formatSessionTime(session.updatedAt)}
                    </span>
                  </span>
                </button>
                <Button
                  type='button'
                  variant='ghost'
                  size='icon-sm'
                  className='mr-1 opacity-0 group-hover/session:opacity-100'
                  onClick={() => onDeleteSession(session.id)}
                  disabled={sessions.length <= 1}
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
