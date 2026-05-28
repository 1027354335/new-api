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
import { useState, useEffect } from 'react'
import { useBlocker } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import { ConfirmDialog } from '@/components/confirm-dialog'

type PlaygroundNavigationGuardProps = {
  when: boolean
  onProceed: () => void
}

/**
 * Navigation guard for Playground to protect active AI generation
 *
 * Prevents route transitions while AI is generating a response.
 * Uses the project's native ConfirmDialog and stops the stream on proceed.
 */
export function PlaygroundNavigationGuard({
  when,
  onProceed,
}: PlaygroundNavigationGuardProps) {
  const { t } = useTranslation()
  const blocker = useBlocker({ condition: when })
  const [showDialog, setShowDialog] = useState(false)

  // Listen to blocker status changes
  useEffect(() => {
    if (blocker.status === 'blocked') {
      setShowDialog(true)
    }
  }, [blocker.status])

  const handleConfirm = () => {
    setShowDialog(false)
    onProceed() // Call stopGeneration first to finalize and save the stream
    blocker.proceed?.()
  }

  const handleCancel = () => {
    setShowDialog(false)
    blocker.reset?.()
  }

  // Handle browser navigation (refresh, close tab)
  useEffect(() => {
    if (!when) return

    const handleBeforeUnload = (e: BeforeUnloadEvent) => {
      e.preventDefault()
      e.returnValue = ''
      return ''
    }

    window.addEventListener('beforeunload', handleBeforeUnload)
    return () => window.removeEventListener('beforeunload', handleBeforeUnload)
  }, [when])

  return (
    <ConfirmDialog
      open={showDialog}
      onOpenChange={(open) => {
        if (!open) handleCancel()
      }}
      title={t('Interruption Warning')}
      desc={t('AI is generating a response. Leaving now will interrupt the generation. Are you sure you want to leave?')}
      confirmText={t('Leave')}
      cancelBtnText={t('Stay')}
      destructive
      handleConfirm={handleConfirm}
    />
  )
}
