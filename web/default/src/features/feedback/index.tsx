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
import { useCallback, useMemo, useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import {
  type ColumnDef,
  getCoreRowModel,
  useReactTable,
} from '@tanstack/react-table'
import {
  Check,
  Eye,
  MessageSquareReply,
  RefreshCw,
  Search,
  X,
} from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'
import { DataTableColumnHeader, DataTablePage } from '@/components/data-table'
import { SectionPageLayout } from '@/components/layout'
import { StatusBadge } from '@/components/status-badge'
import {
  adminReplyFeedback,
  adminUpdateFeedbackStatus,
  closeMyFeedback,
  createFeedback,
  getAdminFeedbacks,
  getMyFeedbacks,
} from './api'
import {
  getFeedbackCategoryLabel,
  getFeedbackCategoryOptions,
  getFeedbackPriorityLabel,
  getFeedbackPriorityOptions,
  getFeedbackPriorityVariant,
  getFeedbackStatusLabel,
  getFeedbackStatusOptions,
  getFeedbackStatusVariant,
} from './constants'
import type {
  Feedback,
  FeedbackCategory,
  FeedbackPriority,
  FeedbackStatus,
} from './types'

const DEFAULT_PAGE_SIZE = 10

function isApiSuccess(response: { success?: boolean; message?: string }) {
  return response.success === true || response.message === 'success'
}

function formatTimestamp(timestamp?: number) {
  if (!timestamp) return '-'
  return new Date(timestamp * 1000).toLocaleString()
}

function FeedbackStatusBadge({ status }: { status: FeedbackStatus }) {
  const { t } = useTranslation()
  return (
    <StatusBadge
      label={t(getFeedbackStatusLabel(status))}
      variant={getFeedbackStatusVariant(status)}
      copyable={false}
      showDot
    />
  )
}

function FeedbackPriorityBadge({ priority }: { priority: FeedbackPriority }) {
  const { t } = useTranslation()
  return (
    <StatusBadge
      label={t(getFeedbackPriorityLabel(priority))}
      variant={getFeedbackPriorityVariant(priority)}
      copyable={false}
      showDot
    />
  )
}

function FeedbackDetailsDialog(props: {
  feedback: Feedback | null
  open: boolean
  onOpenChange: (open: boolean) => void
}) {
  const { t } = useTranslation()
  const feedback = props.feedback

  return (
    <Dialog open={props.open} onOpenChange={props.onOpenChange}>
      <DialogContent className='max-h-[85vh] overflow-y-auto sm:max-w-2xl'>
        <DialogHeader>
          <DialogTitle>{t('Feedback Details')}</DialogTitle>
        </DialogHeader>
        {feedback && (
          <div className='space-y-4 text-sm'>
            <div className='grid gap-3 sm:grid-cols-2'>
              <div>
                <Label className='text-muted-foreground'>{t('Title')}</Label>
                <div className='font-medium break-words'>{feedback.title}</div>
              </div>
              <div>
                <Label className='text-muted-foreground'>{t('Status')}</Label>
                <div className='mt-1'>
                  <FeedbackStatusBadge status={feedback.status} />
                </div>
              </div>
              <div>
                <Label className='text-muted-foreground'>{t('Category')}</Label>
                <div>{t(getFeedbackCategoryLabel(feedback.category))}</div>
              </div>
              <div>
                <Label className='text-muted-foreground'>{t('Priority')}</Label>
                <div className='mt-1'>
                  <FeedbackPriorityBadge priority={feedback.priority} />
                </div>
              </div>
              <div>
                <Label className='text-muted-foreground'>
                  {t('Created At')}
                </Label>
                <div>{formatTimestamp(feedback.create_time)}</div>
              </div>
              <div>
                <Label className='text-muted-foreground'>
                  {t('Updated At')}
                </Label>
                <div>{formatTimestamp(feedback.update_time)}</div>
              </div>
            </div>

            <div className='space-y-1'>
              <Label className='text-muted-foreground'>
                {t('Feedback Content')}
              </Label>
              <div className='bg-muted/30 rounded-md border p-3 break-words whitespace-pre-wrap'>
                {feedback.content}
              </div>
            </div>

            {feedback.reply && (
              <div className='space-y-1'>
                <Label className='text-muted-foreground'>
                  {t('Admin Reply')}
                </Label>
                <div className='bg-muted/30 rounded-md border p-3 break-words whitespace-pre-wrap'>
                  {feedback.reply}
                </div>
                <div className='text-muted-foreground text-xs'>
                  {feedback.admin_name
                    ? `${t('Handled By')}: ${feedback.admin_name}`
                    : null}
                </div>
              </div>
            )}
          </div>
        )}
        <DialogFooter>
          <Button onClick={() => props.onOpenChange(false)}>
            {t('Close')}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

export function FeedbackPage() {
  const { t } = useTranslation()
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE)
  const [status, setStatus] = useState<FeedbackStatus | 'all'>('all')
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const [category, setCategory] = useState<FeedbackCategory>('other')
  const [priority, setPriority] = useState<FeedbackPriority>('normal')
  const [submitting, setSubmitting] = useState(false)
  const [createOpen, setCreateOpen] = useState(false)
  const [selectedFeedback, setSelectedFeedback] = useState<Feedback | null>(
    null
  )
  const [detailsOpen, setDetailsOpen] = useState(false)

  const feedbackQuery = useQuery({
    queryKey: ['my-feedbacks', page, pageSize, status, t],
    queryFn: async () => {
      const response = await getMyFeedbacks({ page, pageSize, status })
      if (!isApiSuccess(response)) {
        toast.error(response.message || t('Failed to load feedback'))
        return { items: [], total: 0 }
      }
      return {
        items: response.data?.items || [],
        total: response.data?.total || 0,
      }
    },
    placeholderData: (previousData) => previousData,
  })
  const {
    data: feedbackData,
    isFetching: isFeedbackFetching,
    isLoading: isFeedbackLoading,
    refetch: refetchFeedbacks,
  } = feedbackQuery

  const refreshFeedbacks = useCallback(() => {
    refetchFeedbacks()
  }, [refetchFeedbacks])

  const handleSubmit = async () => {
    const trimmedTitle = title.trim()
    const trimmedContent = content.trim()
    if (!trimmedTitle) {
      toast.error(t('Please enter feedback title'))
      return
    }
    if (!trimmedContent) {
      toast.error(t('Please enter feedback content'))
      return
    }

    setSubmitting(true)
    try {
      const response = await createFeedback({
        title: trimmedTitle,
        content: trimmedContent,
        category,
        priority,
      })
      if (isApiSuccess(response)) {
        toast.success(t('Feedback submitted successfully'))
        setTitle('')
        setContent('')
        setCategory('other')
        setPriority('normal')
        setCreateOpen(false)
        setPage(1)
        refreshFeedbacks()
      } else {
        toast.error(response.message || t('Failed to submit feedback'))
      }
    } finally {
      setSubmitting(false)
    }
  }

  const handleOpenDetails = (feedback: Feedback) => {
    setSelectedFeedback(feedback)
    setDetailsOpen(true)
  }

  const handleCloseFeedback = useCallback(
    async (feedback: Feedback) => {
      const response = await closeMyFeedback(feedback.id)
      if (isApiSuccess(response)) {
        toast.success(t('Feedback closed'))
        refreshFeedbacks()
      } else {
        toast.error(response.message || t('Failed to close feedback'))
      }
    },
    [refreshFeedbacks, t]
  )

  const columns = useMemo<ColumnDef<Feedback>[]>(
    () => [
      {
        accessorKey: 'title',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('Feedback')} />
        ),
        cell: ({ row }) => (
          <div className='min-w-[220px] space-y-1'>
            <div className='max-w-[320px] truncate font-medium'>
              {row.original.title}
            </div>
            <div className='text-muted-foreground text-xs'>
              {t(getFeedbackCategoryLabel(row.original.category))}
            </div>
          </div>
        ),
        meta: { label: t('Feedback'), mobileTitle: true },
      },
      {
        accessorKey: 'priority',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('Priority')} />
        ),
        cell: ({ row }) => (
          <FeedbackPriorityBadge priority={row.original.priority} />
        ),
        meta: { label: t('Priority'), mobileBadge: true },
      },
      {
        accessorKey: 'status',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('Status')} />
        ),
        cell: ({ row }) => <FeedbackStatusBadge status={row.original.status} />,
        meta: { label: t('Status'), mobileBadge: true },
      },
      {
        accessorKey: 'update_time',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('Updated At')} />
        ),
        cell: ({ row }) => (
          <span className='text-muted-foreground text-xs whitespace-nowrap'>
            {formatTimestamp(row.original.update_time)}
          </span>
        ),
        meta: { label: t('Updated At') },
      },
      {
        id: 'actions',
        header: () => <div className='text-right'>{t('Actions')}</div>,
        cell: ({ row }) => (
          <div className='flex items-center justify-end gap-1'>
            <Button
              variant='ghost'
              size='icon'
              onClick={() => handleOpenDetails(row.original)}
              title={t('View details')}
            >
              <Eye className='h-4 w-4' />
            </Button>
            {row.original.status !== 'closed' && (
              <Button
                variant='outline'
                size='sm'
                onClick={() => handleCloseFeedback(row.original)}
              >
                {t('Close')}
              </Button>
            )}
          </div>
        ),
      },
    ],
    [handleCloseFeedback, t]
  )

  const table = useReactTable({
    data: feedbackData?.items || [],
    columns,
    state: {
      pagination: {
        pageIndex: page - 1,
        pageSize,
      },
    },
    pageCount: Math.max(1, Math.ceil((feedbackData?.total || 0) / pageSize)),
    manualPagination: true,
    getCoreRowModel: getCoreRowModel(),
    onPaginationChange: (updater) => {
      const current = { pageIndex: page - 1, pageSize }
      const next = typeof updater === 'function' ? updater(current) : updater
      setPage(next.pageIndex + 1)
      setPageSize(next.pageSize)
    },
  })

  const categoryOptions = useMemo(() => getFeedbackCategoryOptions(t), [t])
  const priorityOptions = useMemo(() => getFeedbackPriorityOptions(t), [t])

  return (
    <SectionPageLayout>
      <SectionPageLayout.Title>{t('Feedback')}</SectionPageLayout.Title>
      <SectionPageLayout.Actions>
        <Button onClick={() => setCreateOpen(true)}>
          {t('Submit Feedback')}
        </Button>
        <Button
          variant='outline'
          size='icon'
          onClick={refreshFeedbacks}
          disabled={isFeedbackFetching}
        >
          <RefreshCw
            className={`h-4 w-4 ${isFeedbackFetching ? 'animate-spin' : ''}`}
          />
        </Button>
      </SectionPageLayout.Actions>
      <SectionPageLayout.Content>
        <div className='flex flex-col gap-4'>
          <DataTablePage
            table={table}
            columns={columns}
            isLoading={isFeedbackLoading}
            isFetching={isFeedbackFetching}
            emptyTitle={t('No Feedback Found')}
            emptyDescription={t('Submitted feedback will appear here.')}
            skeletonKeyPrefix='feedback-skeleton'
            toolbar={
              <div className='flex justify-end'>
                <Select
                  items={[
                    { value: 'all', label: t('All Statuses') },
                    ...getFeedbackStatusOptions(t),
                  ]}
                  value={status}
                  onValueChange={(value) => {
                    setStatus((value as FeedbackStatus | 'all') || 'all')
                    setPage(1)
                  }}
                >
                  <SelectTrigger className='h-8 w-[150px]'>
                    <SelectValue placeholder={t('Status')} />
                  </SelectTrigger>
                  <SelectContent alignItemWithTrigger={false}>
                    <SelectGroup>
                      <SelectItem value='all'>{t('All Statuses')}</SelectItem>
                      {getFeedbackStatusOptions(t).map((option) => (
                        <SelectItem key={option.value} value={option.value}>
                          {option.label}
                        </SelectItem>
                      ))}
                    </SelectGroup>
                  </SelectContent>
                </Select>
              </div>
            }
          />
        </div>

        <Dialog open={createOpen} onOpenChange={setCreateOpen}>
          <DialogContent className='max-h-[90vh] overflow-y-auto sm:max-w-3xl'>
            <DialogHeader>
              <DialogTitle>{t('Submit Feedback')}</DialogTitle>
            </DialogHeader>
            <div className='space-y-4 py-2'>
              <div className='space-y-2'>
                <Label>{t('Title')}</Label>
                <Input
                  value={title}
                  maxLength={255}
                  onChange={(event) => setTitle(event.target.value)}
                  placeholder={t('Briefly describe your issue')}
                />
              </div>
              <div className='grid gap-3 sm:grid-cols-2'>
                <div className='space-y-2'>
                  <Label>{t('Category')}</Label>
                  <Select
                    items={categoryOptions}
                    value={category}
                    onValueChange={(value) =>
                      setCategory((value as FeedbackCategory) || 'other')
                    }
                  >
                    <SelectTrigger>
                      <SelectValue placeholder={t('Category')} />
                    </SelectTrigger>
                    <SelectContent alignItemWithTrigger={false}>
                      <SelectGroup>
                        {categoryOptions.map((option) => (
                          <SelectItem key={option.value} value={option.value}>
                            {option.label}
                          </SelectItem>
                        ))}
                      </SelectGroup>
                    </SelectContent>
                  </Select>
                </div>
                <div className='space-y-2'>
                  <Label>{t('Priority')}</Label>
                  <Select
                    items={priorityOptions}
                    value={priority}
                    onValueChange={(value) =>
                      setPriority((value as FeedbackPriority) || 'normal')
                    }
                  >
                    <SelectTrigger>
                      <SelectValue placeholder={t('Priority')} />
                    </SelectTrigger>
                    <SelectContent alignItemWithTrigger={false}>
                      <SelectGroup>
                        {priorityOptions.map((option) => (
                          <SelectItem key={option.value} value={option.value}>
                            {option.label}
                          </SelectItem>
                        ))}
                      </SelectGroup>
                    </SelectContent>
                  </Select>
                </div>
              </div>
              <div className='space-y-2'>
                <Label>{t('Content')}</Label>
                <Textarea
                  value={content}
                  maxLength={5000}
                  rows={12}
                  className='min-h-[300px] resize-y'
                  onChange={(event) => setContent(event.target.value)}
                  placeholder={t(
                    'Include steps, expected result, and actual result.'
                  )}
                />
              </div>
            </div>
            <DialogFooter>
              <Button
                variant='outline'
                onClick={() => setCreateOpen(false)}
                disabled={submitting}
              >
                {t('Cancel')}
              </Button>
              <Button onClick={handleSubmit} disabled={submitting}>
                {submitting ? t('Submitting...') : t('Submit Feedback')}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        <FeedbackDetailsDialog
          feedback={selectedFeedback}
          open={detailsOpen}
          onOpenChange={setDetailsOpen}
        />
      </SectionPageLayout.Content>
    </SectionPageLayout>
  )
}

export function FeedbackManagementPage() {
  const { t } = useTranslation()
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE)
  const [status, setStatus] = useState<FeedbackStatus | 'all'>('all')
  const [category, setCategory] = useState<FeedbackCategory | 'all'>('all')
  const [priority, setPriority] = useState<FeedbackPriority | 'all'>('all')
  const [keyword, setKeyword] = useState('')
  const [selectedFeedback, setSelectedFeedback] = useState<Feedback | null>(
    null
  )
  const [detailsOpen, setDetailsOpen] = useState(false)
  const [replyOpen, setReplyOpen] = useState(false)
  const [replyStatus, setReplyStatus] = useState<'resolved' | 'rejected'>(
    'resolved'
  )
  const [reply, setReply] = useState('')
  const [processing, setProcessing] = useState(false)

  const feedbackQuery = useQuery({
    queryKey: [
      'admin-feedbacks',
      page,
      pageSize,
      status,
      category,
      priority,
      keyword,
      t,
    ],
    queryFn: async () => {
      const response = await getAdminFeedbacks({
        page,
        pageSize,
        status,
        category,
        priority,
        keyword,
      })
      if (!isApiSuccess(response)) {
        toast.error(response.message || t('Failed to load feedback'))
        return { items: [], total: 0 }
      }
      return {
        items: response.data?.items || [],
        total: response.data?.total || 0,
      }
    },
    placeholderData: (previousData) => previousData,
  })
  const {
    data: feedbackData,
    isFetching: isFeedbackFetching,
    isLoading: isFeedbackLoading,
    refetch: refetchFeedbacks,
  } = feedbackQuery

  const refreshFeedbacks = useCallback(() => {
    refetchFeedbacks()
  }, [refetchFeedbacks])

  const handleOpenDetails = (feedback: Feedback) => {
    setSelectedFeedback(feedback)
    setDetailsOpen(true)
  }

  const handleOpenReply = (feedback: Feedback) => {
    setSelectedFeedback(feedback)
    setReplyStatus(feedback.status === 'rejected' ? 'rejected' : 'resolved')
    setReply(feedback.reply || '')
    setReplyOpen(true)
  }

  const handleMarkInProgress = useCallback(
    async (feedback: Feedback) => {
      const response = await adminUpdateFeedbackStatus(
        feedback.id,
        'in_progress'
      )
      if (isApiSuccess(response)) {
        toast.success(t('Feedback marked as in progress'))
        refreshFeedbacks()
      } else {
        toast.error(response.message || t('Failed to update feedback status'))
      }
    },
    [refreshFeedbacks, t]
  )

  const handleReply = async () => {
    if (!selectedFeedback) return
    if (!reply.trim()) {
      toast.error(t('Please enter reply content'))
      return
    }
    setProcessing(true)
    try {
      const response = await adminReplyFeedback(selectedFeedback.id, {
        status: replyStatus,
        reply: reply.trim(),
      })
      if (isApiSuccess(response)) {
        toast.success(t('Feedback replied successfully'))
        setReplyOpen(false)
        refreshFeedbacks()
      } else {
        toast.error(response.message || t('Failed to reply feedback'))
      }
    } finally {
      setProcessing(false)
    }
  }

  const columns = useMemo<ColumnDef<Feedback>[]>(
    () => [
      {
        accessorKey: 'username',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('User')} />
        ),
        cell: ({ row }) => (
          <div className='min-w-[150px]'>
            <div className='font-medium'>{row.original.username}</div>
            <div className='text-muted-foreground text-xs'>
              ID: {row.original.user_id}
            </div>
          </div>
        ),
        meta: { label: t('User'), mobileTitle: true },
      },
      {
        accessorKey: 'title',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('Feedback')} />
        ),
        cell: ({ row }) => (
          <div className='min-w-[240px] space-y-1'>
            <div className='max-w-[360px] truncate font-medium'>
              {row.original.title}
            </div>
            <div className='text-muted-foreground flex gap-2 text-xs'>
              <span>{t(getFeedbackCategoryLabel(row.original.category))}</span>
              <span>{row.original.email}</span>
            </div>
          </div>
        ),
        meta: { label: t('Feedback') },
      },
      {
        accessorKey: 'priority',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('Priority')} />
        ),
        cell: ({ row }) => (
          <FeedbackPriorityBadge priority={row.original.priority} />
        ),
        meta: { label: t('Priority'), mobileBadge: true },
      },
      {
        accessorKey: 'status',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('Status')} />
        ),
        cell: ({ row }) => <FeedbackStatusBadge status={row.original.status} />,
        meta: { label: t('Status'), mobileBadge: true },
      },
      {
        accessorKey: 'update_time',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('Updated At')} />
        ),
        cell: ({ row }) => (
          <span className='text-muted-foreground text-xs whitespace-nowrap'>
            {formatTimestamp(row.original.update_time)}
          </span>
        ),
        meta: { label: t('Updated At') },
      },
      {
        id: 'actions',
        header: () => <div className='text-right'>{t('Actions')}</div>,
        cell: ({ row }) => (
          <div className='flex items-center justify-end gap-1'>
            <Button
              variant='ghost'
              size='icon'
              onClick={() => handleOpenDetails(row.original)}
              title={t('View details')}
            >
              <Eye className='h-4 w-4' />
            </Button>
            {row.original.status === 'open' && (
              <Button
                variant='outline'
                size='sm'
                onClick={() => handleMarkInProgress(row.original)}
              >
                {t('Start')}
              </Button>
            )}
            {(row.original.status === 'open' ||
              row.original.status === 'in_progress') && (
              <Button size='sm' onClick={() => handleOpenReply(row.original)}>
                <MessageSquareReply className='mr-1 h-3.5 w-3.5' />
                {t('Reply')}
              </Button>
            )}
          </div>
        ),
      },
    ],
    [handleMarkInProgress, t]
  )

  const table = useReactTable({
    data: feedbackData?.items || [],
    columns,
    state: {
      pagination: {
        pageIndex: page - 1,
        pageSize,
      },
    },
    pageCount: Math.max(1, Math.ceil((feedbackData?.total || 0) / pageSize)),
    manualPagination: true,
    getCoreRowModel: getCoreRowModel(),
    onPaginationChange: (updater) => {
      const current = { pageIndex: page - 1, pageSize }
      const next = typeof updater === 'function' ? updater(current) : updater
      setPage(next.pageIndex + 1)
      setPageSize(next.pageSize)
    },
  })

  const statusOptions = useMemo(() => getFeedbackStatusOptions(t), [t])
  const categoryOptions = useMemo(() => getFeedbackCategoryOptions(t), [t])
  const priorityOptions = useMemo(() => getFeedbackPriorityOptions(t), [t])

  return (
    <SectionPageLayout>
      <SectionPageLayout.Title>
        {t('Feedback Management')}
      </SectionPageLayout.Title>
      <SectionPageLayout.Actions>
        <Button
          variant='outline'
          size='icon'
          onClick={refreshFeedbacks}
          disabled={isFeedbackFetching}
        >
          <RefreshCw
            className={`h-4 w-4 ${isFeedbackFetching ? 'animate-spin' : ''}`}
          />
        </Button>
      </SectionPageLayout.Actions>
      <SectionPageLayout.Content>
        <DataTablePage
          table={table}
          columns={columns}
          isLoading={isFeedbackLoading}
          isFetching={isFeedbackFetching}
          emptyTitle={t('No Feedback Found')}
          emptyDescription={t('No feedback matches the current filters.')}
          skeletonKeyPrefix='admin-feedback-skeleton'
          toolbar={
            <div className='flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between'>
              <div className='relative min-w-0 flex-1'>
                <Search className='text-muted-foreground absolute top-1/2 left-2.5 h-4 w-4 -translate-y-1/2' />
                <Input
                  value={keyword}
                  onChange={(event) => {
                    setKeyword(event.target.value)
                    setPage(1)
                  }}
                  placeholder={t(
                    'Search by title, content, username or email...'
                  )}
                  className='h-8 w-full pl-8 sm:max-w-sm'
                />
              </div>
              <div className='flex flex-wrap gap-2'>
                <Select
                  items={[
                    { value: 'all', label: t('All Statuses') },
                    ...statusOptions,
                  ]}
                  value={status}
                  onValueChange={(value) => {
                    setStatus((value as FeedbackStatus | 'all') || 'all')
                    setPage(1)
                  }}
                >
                  <SelectTrigger className='h-8 w-[145px]'>
                    <SelectValue placeholder={t('Status')} />
                  </SelectTrigger>
                  <SelectContent alignItemWithTrigger={false}>
                    <SelectGroup>
                      <SelectItem value='all'>{t('All Statuses')}</SelectItem>
                      {statusOptions.map((option) => (
                        <SelectItem key={option.value} value={option.value}>
                          {option.label}
                        </SelectItem>
                      ))}
                    </SelectGroup>
                  </SelectContent>
                </Select>
                <Select
                  items={[
                    { value: 'all', label: t('All Categories') },
                    ...categoryOptions,
                  ]}
                  value={category}
                  onValueChange={(value) => {
                    setCategory((value as FeedbackCategory | 'all') || 'all')
                    setPage(1)
                  }}
                >
                  <SelectTrigger className='h-8 w-[150px]'>
                    <SelectValue placeholder={t('Category')} />
                  </SelectTrigger>
                  <SelectContent alignItemWithTrigger={false}>
                    <SelectGroup>
                      <SelectItem value='all'>{t('All Categories')}</SelectItem>
                      {categoryOptions.map((option) => (
                        <SelectItem key={option.value} value={option.value}>
                          {option.label}
                        </SelectItem>
                      ))}
                    </SelectGroup>
                  </SelectContent>
                </Select>
                <Select
                  items={[
                    { value: 'all', label: t('All Priorities') },
                    ...priorityOptions,
                  ]}
                  value={priority}
                  onValueChange={(value) => {
                    setPriority((value as FeedbackPriority | 'all') || 'all')
                    setPage(1)
                  }}
                >
                  <SelectTrigger className='h-8 w-[145px]'>
                    <SelectValue placeholder={t('Priority')} />
                  </SelectTrigger>
                  <SelectContent alignItemWithTrigger={false}>
                    <SelectGroup>
                      <SelectItem value='all'>{t('All Priorities')}</SelectItem>
                      {priorityOptions.map((option) => (
                        <SelectItem key={option.value} value={option.value}>
                          {option.label}
                        </SelectItem>
                      ))}
                    </SelectGroup>
                  </SelectContent>
                </Select>
              </div>
            </div>
          }
        />

        <FeedbackDetailsDialog
          feedback={selectedFeedback}
          open={detailsOpen}
          onOpenChange={setDetailsOpen}
        />

        <Dialog open={replyOpen} onOpenChange={setReplyOpen}>
          <DialogContent className='sm:max-w-lg'>
            <DialogHeader>
              <DialogTitle>{t('Reply Feedback')}</DialogTitle>
            </DialogHeader>
            {selectedFeedback && (
              <div className='space-y-4 py-2'>
                <div className='bg-muted/30 rounded-md border p-3 text-sm'>
                  <div className='font-medium break-words'>
                    {selectedFeedback.title}
                  </div>
                  <div className='text-muted-foreground mt-1 text-xs'>
                    {selectedFeedback.username} ·{' '}
                    {t(getFeedbackCategoryLabel(selectedFeedback.category))}
                  </div>
                </div>
                <div className='space-y-2'>
                  <Label>{t('Result Status')}</Label>
                  <div className='flex gap-3'>
                    <label className='flex cursor-pointer items-center gap-2 text-sm'>
                      <input
                        type='radio'
                        name='reply_status'
                        checked={replyStatus === 'resolved'}
                        onChange={() => setReplyStatus('resolved')}
                      />
                      <Check className='text-success h-4 w-4' />
                      {t('Resolved')}
                    </label>
                    <label className='flex cursor-pointer items-center gap-2 text-sm'>
                      <input
                        type='radio'
                        name='reply_status'
                        checked={replyStatus === 'rejected'}
                        onChange={() => setReplyStatus('rejected')}
                      />
                      <X className='text-destructive h-4 w-4' />
                      {t('Rejected')}
                    </label>
                  </div>
                </div>
                <div className='space-y-2'>
                  <Label>{t('Reply Content')}</Label>
                  <Textarea
                    value={reply}
                    maxLength={5000}
                    rows={6}
                    onChange={(event) => setReply(event.target.value)}
                    placeholder={t(
                      'Enter handling result. It will be emailed to the user.'
                    )}
                  />
                </div>
              </div>
            )}
            <DialogFooter>
              <Button
                variant='outline'
                onClick={() => setReplyOpen(false)}
                disabled={processing}
              >
                {t('Cancel')}
              </Button>
              <Button onClick={handleReply} disabled={processing}>
                {processing ? t('Processing...') : t('Send Reply')}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </SectionPageLayout.Content>
    </SectionPageLayout>
  )
}
