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
import { useEffect, useState, useCallback, useMemo } from 'react'
import {
  type ColumnDef,
  getCoreRowModel,
  useReactTable,
} from '@tanstack/react-table'
import {
  Search,
  Check,
  X,
  Download,
  AlertCircle,
  Eye,
  RefreshCw,
} from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
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
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import { DataTableColumnHeader, DataTablePage } from '@/components/data-table'
import { StatusBadge } from '@/components/status-badge'
import {
  getAdminInvoices,
  downloadInvoiceFile,
  uploadInvoiceFile,
  completeInvoice,
  rejectInvoice,
  isApiSuccess,
} from '../api'
import { formatPaymentAmount } from '../lib'
import {
  getInvoiceStatusConfig,
  getPaymentMethodName,
  formatTimestamp,
} from '../lib/billing'
import type { InvoiceRecord } from '../types'

function getInvoicePaidAmount(invoice: InvoiceRecord): number {
  return invoice.paid_amount && invoice.paid_amount > 0
    ? invoice.paid_amount
    : invoice.money
}

function getInvoicePaidCurrency(invoice: InvoiceRecord): string {
  if (invoice.paid_currency) return invoice.paid_currency
  if (invoice.payment_method === 'paypal') return 'EUR'
  if (invoice.payment_method === 'alipay') return 'CNY'
  return 'USD'
}

function formatInvoicePaidAmount(invoice: InvoiceRecord): string {
  return formatPaymentAmount(
    getInvoicePaidAmount(invoice),
    getInvoicePaidCurrency(invoice)
  )
}

export function InvoiceManagement() {
  const { t } = useTranslation()
  const [invoices, setInvoices] = useState<InvoiceRecord[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [status, setStatus] = useState<string>('all')
  const [paymentMethod, setPaymentMethod] = useState<string>('all')
  const [keyword, setKeyword] = useState<string>('')
  const [loading, setLoading] = useState(false)

  // Dialog states
  const [selectedInvoice, setSelectedInvoice] = useState<InvoiceRecord | null>(
    null
  )
  const [processDialogOpen, setProcessDialogOpen] = useState(false)
  const [viewDialogOpen, setViewDialogOpen] = useState(false)

  // Process Form states
  const [rejectReason, setRejectReason] = useState('')
  const [showRejectForm, setShowRejectForm] = useState(false)
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const [processing, setProcessing] = useState(false)
  const [paypalMethod, setPaypalMethod] = useState<'lexware' | 'manual'>(
    'lexware'
  )

  const fetchInvoices = useCallback(async () => {
    setLoading(true)
    try {
      const response = await getAdminInvoices(
        page,
        pageSize,
        status === 'all' ? undefined : status,
        paymentMethod === 'all' ? undefined : paymentMethod,
        keyword.trim() || undefined
      )
      if (isApiSuccess(response) && response.data) {
        setInvoices(response.data.items || [])
        setTotal(response.data.total || 0)
      } else {
        toast.error(response.message || t('Failed to load invoices'))
      }
    } catch (err) {
      // eslint-disable-next-line no-console
      console.error('Fetch invoices error:', err)
      toast.error(t('Failed to load invoices'))
    } finally {
      setLoading(false)
    }
  }, [page, pageSize, status, paymentMethod, keyword, t])

  useEffect(() => {
    fetchInvoices()
  }, [fetchInvoices])

  const handleSearchChange = (val: string) => {
    setKeyword(val)
    setPage(1)
  }

  const handleStatusChange = (val: string) => {
    setStatus(val)
    setPage(1)
  }

  const handleMethodChange = (val: string) => {
    setPaymentMethod(val)
    setPage(1)
  }

  const handleOpenProcess = (invoice: InvoiceRecord) => {
    setSelectedInvoice(invoice)
    setShowRejectForm(false)
    setRejectReason('')
    setSelectedFile(null)
    setPaypalMethod('lexware')
    setProcessDialogOpen(true)
  }

  const handleOpenView = (invoice: InvoiceRecord) => {
    setSelectedInvoice(invoice)
    setViewDialogOpen(true)
  }

  const handleDownloadInvoice = useCallback((invoice: InvoiceRecord) => {
    if (!invoice.download_url) return
    if (invoice.download_url.startsWith('http')) {
      window.open(invoice.download_url, '_blank')
      return
    }
    downloadInvoiceFile(invoice.id)
  }, [])

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0]
      if (file.size > 10 * 1024 * 1024) {
        toast.error(t('File size must be less than 10MB'))
        return
      }
      setSelectedFile(file)
    }
  }

  const handleReject = async () => {
    if (!selectedInvoice) return
    if (!rejectReason.trim()) {
      toast.error(t('Please enter a rejection reason'))
      return
    }

    setProcessing(true)
    try {
      const response = await rejectInvoice({
        invoice_id: selectedInvoice.id,
        message: rejectReason,
      })
      if (isApiSuccess(response)) {
        toast.success(t('Invoice request rejected successfully'))
        setProcessDialogOpen(false)
        fetchInvoices()
      } else {
        toast.error(response.message || t('Failed to reject invoice'))
      }
    } catch (err) {
      // eslint-disable-next-line no-console
      console.error('Reject invoice error:', err)
      toast.error(t('Failed to reject invoice'))
    } finally {
      setProcessing(false)
    }
  }

  const handleComplete = async () => {
    if (!selectedInvoice) return

    setProcessing(true)
    try {
      let downloadUrl = ''

      // If alipay or paypal-manual, we need file upload
      const isManual =
        selectedInvoice.payment_method === 'alipay' || paypalMethod === 'manual'

      if (isManual) {
        if (!selectedFile) {
          toast.error(t('Please select an invoice file to upload'))
          setProcessing(false)
          return
        }

        const uploadRes = await uploadInvoiceFile(
          selectedInvoice.id,
          selectedFile
        )
        const uploadedPath = uploadRes.data?.file_path || uploadRes.data?.url
        if (!isApiSuccess(uploadRes) || !uploadedPath) {
          toast.error(uploadRes.message || t('Failed to upload invoice file'))
          setProcessing(false)
          return
        }
        downloadUrl = uploadedPath
      }

      const response = await completeInvoice({
        invoice_id: selectedInvoice.id,
        download_url: downloadUrl,
      })

      if (isApiSuccess(response)) {
        toast.success(t('Invoice processed and completed successfully'))
        setProcessDialogOpen(false)
        fetchInvoices()
      } else {
        toast.error(response.message || t('Failed to complete invoice'))
      }
    } catch (err) {
      // eslint-disable-next-line no-console
      console.error('Complete invoice error:', err)
      toast.error(t('Failed to complete invoice'))
    } finally {
      setProcessing(false)
    }
  }

  const columns = useMemo<ColumnDef<InvoiceRecord>[]>(
    () => [
      {
        accessorKey: 'username',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('User')} />
        ),
        cell: ({ row }) => {
          const invoice = row.original
          return (
            <div className='min-w-[150px]'>
              <div className='font-medium'>{invoice.username}</div>
              <div className='text-muted-foreground text-xs'>
                ID: {invoice.user_id}
              </div>
            </div>
          )
        },
        meta: {
          label: t('User'),
          mobileTitle: true,
        },
      },
      {
        id: 'order',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('Order Details')} />
        ),
        cell: ({ row }) => {
          const invoice = row.original
          return (
            <div className='min-w-[220px] space-y-0.5'>
              <div className='font-semibold text-red-600'>
                {formatInvoicePaidAmount(invoice)}
              </div>
              {invoice.credit_amount_usd && invoice.credit_amount_usd > 0 && (
                <div className='text-muted-foreground text-xs'>
                  {t('Credited USD')}:{' '}
                  {formatPaymentAmount(invoice.credit_amount_usd, 'USD')}
                </div>
              )}
              <div
                className='text-muted-foreground max-w-[220px] truncate font-mono text-xs'
                title={invoice.trade_no}
              >
                {invoice.trade_no}
              </div>
              <div className='text-xs'>
                {getPaymentMethodName(invoice.payment_method, t)}
              </div>
            </div>
          )
        },
        meta: {
          label: t('Order Details'),
        },
      },
      {
        id: 'invoice',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('Invoice Info')} />
        ),
        cell: ({ row }) => {
          const invoice = row.original
          return (
            <div className='min-w-[220px] space-y-0.5'>
              <div
                className='max-w-[220px] truncate font-medium'
                title={invoice.title}
              >
                {invoice.title}
              </div>
              <div className='text-muted-foreground max-w-[220px] truncate text-xs'>
                {invoice.email}
              </div>
              <div className='text-muted-foreground text-xs font-semibold uppercase'>
                {invoice.billing_type === 'enterprise'
                  ? t('Enterprise')
                  : t('Personal')}
              </div>
            </div>
          )
        },
        meta: {
          label: t('Invoice Info'),
        },
      },
      {
        accessorKey: 'status',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('Status')} />
        ),
        cell: ({ row }) => {
          const invoice = row.original
          const badgeConfig = getInvoiceStatusConfig(invoice.status)
          return (
            <div className='flex items-center gap-1.5'>
              <StatusBadge
                label={t(badgeConfig.label)}
                variant={
                  badgeConfig.variant === 'destructive'
                    ? 'danger'
                    : badgeConfig.variant === 'warning'
                      ? 'warning'
                      : 'success'
                }
                showDot
                copyable={false}
              />
              {invoice.status === 'rejected' && invoice.message && (
                <Tooltip>
                  <TooltipTrigger
                    render={
                      <div className='text-muted-foreground hover:text-foreground cursor-pointer'>
                        <AlertCircle className='h-4 w-4' />
                      </div>
                    }
                  />
                  <TooltipContent>
                    <p className='max-w-xs break-all'>{invoice.message}</p>
                  </TooltipContent>
                </Tooltip>
              )}
            </div>
          )
        },
        meta: {
          label: t('Status'),
          mobileBadge: true,
        },
      },
      {
        accessorKey: 'create_time',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('Request Time')} />
        ),
        cell: ({ row }) => (
          <span className='text-muted-foreground text-xs whitespace-nowrap'>
            {formatTimestamp(row.original.create_time)}
          </span>
        ),
        meta: {
          label: t('Request Time'),
        },
      },
      {
        id: 'actions',
        header: () => <div className='text-right'>{t('Actions')}</div>,
        cell: ({ row }) => {
          const invoice = row.original
          return (
            <div className='flex items-center justify-end gap-1'>
              <Button
                variant='ghost'
                size='icon'
                onClick={() => handleOpenView(invoice)}
                title={t('View details')}
              >
                <Eye className='h-4 w-4' />
              </Button>
              {invoice.status === 'pending' ? (
                <Button size='sm' onClick={() => handleOpenProcess(invoice)}>
                  {t('Process')}
                </Button>
              ) : invoice.status === 'completed' && invoice.download_url ? (
                <Button
                  variant='outline'
                  size='sm'
                  onClick={() => handleDownloadInvoice(invoice)}
                >
                  <Download className='mr-1 h-3.5 w-3.5' />
                  {t('Download')}
                </Button>
              ) : null}
            </div>
          )
        },
      },
    ],
    [handleDownloadInvoice, t]
  )

  const table = useReactTable({
    data: invoices,
    columns,
    state: {
      pagination: {
        pageIndex: page - 1,
        pageSize,
      },
    },
    pageCount: Math.max(1, Math.ceil(total / pageSize)),
    manualPagination: true,
    getCoreRowModel: getCoreRowModel(),
    onPaginationChange: (updater) => {
      const current = { pageIndex: page - 1, pageSize }
      const next = typeof updater === 'function' ? updater(current) : updater
      setPage(next.pageIndex + 1)
      setPageSize(next.pageSize)
    },
  })

  return (
    <div className='flex flex-col gap-4 p-4 sm:p-6'>
      <div>
        <div className='flex items-center justify-between gap-3'>
          <h1 className='truncate text-base font-bold tracking-tight sm:text-lg'>
            {t('Invoice Management')}
          </h1>
          <Button
            variant='outline'
            size='icon'
            className='shrink-0'
            onClick={fetchInvoices}
            disabled={loading}
          >
            <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
          </Button>
        </div>
      </div>

      <DataTablePage
        table={table}
        columns={columns}
        isLoading={loading}
        emptyTitle={t('No invoice requests found')}
        emptyDescription={t('No invoice requests match the current filters.')}
        skeletonKeyPrefix='invoice-skeleton'
        toolbar={
          <div className='flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between'>
            <div className='relative min-w-0 flex-1'>
              <Search className='text-muted-foreground absolute top-1/2 left-2.5 h-4 w-4 -translate-y-1/2' />
              <Input
                placeholder={t('Search by username or trade number...')}
                value={keyword}
                onChange={(event) => handleSearchChange(event.target.value)}
                className='h-8 w-full pl-8 sm:max-w-sm'
              />
            </div>
            <div className='flex flex-wrap items-center gap-2'>
              <Select
                items={[
                  { value: 'all', label: t('All') },
                  { value: 'pending', label: t('Pending') },
                  { value: 'completed', label: t('Completed') },
                  { value: 'rejected', label: t('Rejected') },
                ]}
                value={status}
                onValueChange={(val) => handleStatusChange(val ?? 'all')}
              >
                <SelectTrigger className='h-8 w-[130px]'>
                  <SelectValue placeholder={t('Status')} />
                </SelectTrigger>
                <SelectContent alignItemWithTrigger={false}>
                  <SelectGroup>
                    <SelectItem value='all'>{t('All Statuses')}</SelectItem>
                    <SelectItem value='pending'>{t('Pending')}</SelectItem>
                    <SelectItem value='completed'>{t('Completed')}</SelectItem>
                    <SelectItem value='rejected'>{t('Rejected')}</SelectItem>
                  </SelectGroup>
                </SelectContent>
              </Select>

              <Select
                items={[
                  { value: 'all', label: t('All') },
                  { value: 'alipay', label: t('Alipay') },
                  { value: 'paypal', label: t('PayPal') },
                ]}
                value={paymentMethod}
                onValueChange={(val) => handleMethodChange(val ?? 'all')}
              >
                <SelectTrigger className='h-8 w-[140px]'>
                  <SelectValue placeholder={t('Payment Method')} />
                </SelectTrigger>
                <SelectContent alignItemWithTrigger={false}>
                  <SelectGroup>
                    <SelectItem value='all'>{t('All Methods')}</SelectItem>
                    <SelectItem value='alipay'>{t('Alipay')}</SelectItem>
                    <SelectItem value='paypal'>{t('PayPal')}</SelectItem>
                  </SelectGroup>
                </SelectContent>
              </Select>
            </div>
          </div>
        }
      />

      {/* Process Invoice Dialog */}
      <Dialog open={processDialogOpen} onOpenChange={setProcessDialogOpen}>
        <DialogContent className='max-h-[90vh] overflow-y-auto sm:max-w-lg'>
          <DialogHeader>
            <DialogTitle>{t('Process Invoice Request')}</DialogTitle>
            <DialogDescription>
              {t('Review details and select complete or reject option')}
            </DialogDescription>
          </DialogHeader>

          {selectedInvoice && (
            <div className='space-y-4 py-2'>
              {/* Detail Info Card */}
              <div className='bg-muted/30 grid grid-cols-2 gap-3 rounded-lg border p-3 text-xs sm:grid-cols-3'>
                <div>
                  <Label className='text-muted-foreground'>{t('User')}</Label>
                  <div className='font-medium'>{selectedInvoice.username}</div>
                </div>
                <div>
                  <Label className='text-muted-foreground'>
                    {t('Payment Method')}
                  </Label>
                  <div className='font-medium uppercase'>
                    {selectedInvoice.payment_method}
                  </div>
                </div>
                <div>
                  <Label className='text-muted-foreground'>{t('Amount')}</Label>
                  <div className='font-medium font-semibold text-red-600'>
                    {formatInvoicePaidAmount(selectedInvoice)}
                  </div>
                  {selectedInvoice.credit_amount_usd &&
                    selectedInvoice.credit_amount_usd > 0 && (
                      <div className='text-muted-foreground text-[11px]'>
                        {t('Credited USD')}:{' '}
                        {formatPaymentAmount(
                          selectedInvoice.credit_amount_usd,
                          'USD'
                        )}
                      </div>
                    )}
                </div>
                <div>
                  <Label className='text-muted-foreground'>
                    {t('Billing Type')}
                  </Label>
                  <div className='font-medium capitalize'>
                    {t(selectedInvoice.billing_type)}
                  </div>
                </div>
                <div className='col-span-2'>
                  <Label className='text-muted-foreground'>
                    {t('Invoice Title')}
                  </Label>
                  <div className='truncate font-medium'>
                    {selectedInvoice.title}
                  </div>
                </div>
                {selectedInvoice.billing_type === 'enterprise' && (
                  <div>
                    <Label className='text-muted-foreground'>
                      {t('Tax ID')}
                    </Label>
                    <div className='font-medium'>
                      {selectedInvoice.tax_id || '-'}
                    </div>
                  </div>
                )}
                <div className='col-span-2'>
                  <Label className='text-muted-foreground'>{t('Email')}</Label>
                  <div className='font-medium'>{selectedInvoice.email}</div>
                </div>
              </div>

              {selectedInvoice.payment_method === 'paypal' && (
                <div className='space-y-3 border-t pt-3'>
                  <Label className='text-sm font-semibold'>
                    {t('PayPal Address Details')}
                  </Label>
                  <div className='grid grid-cols-2 gap-2 text-xs'>
                    <div>
                      <Label className='text-muted-foreground'>
                        {t('Street')}
                      </Label>
                      <div>{selectedInvoice.street || '-'}</div>
                    </div>
                    <div>
                      <Label className='text-muted-foreground'>
                        {t('Detailed Address')}
                      </Label>
                      <div>{selectedInvoice.address_detail || '-'}</div>
                    </div>
                    <div>
                      <Label className='text-muted-foreground'>
                        {t('City')}
                      </Label>
                      <div>{selectedInvoice.city || '-'}</div>
                    </div>
                    <div>
                      <Label className='text-muted-foreground'>
                        {t('Zip Code')}
                      </Label>
                      <div>{selectedInvoice.zip_code || '-'}</div>
                    </div>
                    <div>
                      <Label className='text-muted-foreground'>
                        {t('Country')}
                      </Label>
                      <div>{selectedInvoice.country || '-'}</div>
                    </div>
                  </div>
                </div>
              )}

              {/* Action selection */}
              {!showRejectForm ? (
                <div className='space-y-4 border-t pt-3'>
                  {selectedInvoice.payment_method === 'paypal' && (
                    <div className='space-y-2'>
                      <Label className='text-sm font-semibold'>
                        {t('PayPal Invoicing Mode')}
                      </Label>
                      <div className='flex gap-4'>
                        <label className='flex cursor-pointer items-center gap-2 text-sm'>
                          <input
                            type='radio'
                            name='paypal_mode'
                            checked={paypalMethod === 'lexware'}
                            onChange={() => setPaypalMethod('lexware')}
                          />
                          {t('Automatic (Lexware)')}
                        </label>
                        <label className='flex cursor-pointer items-center gap-2 text-sm'>
                          <input
                            type='radio'
                            name='paypal_mode'
                            checked={paypalMethod === 'manual'}
                            onChange={() => setPaypalMethod('manual')}
                          />
                          {t('Manual File Upload')}
                        </label>
                      </div>
                    </div>
                  )}

                  {(selectedInvoice.payment_method === 'alipay' ||
                    paypalMethod === 'manual') && (
                    <div className='space-y-2'>
                      <Label className='text-sm font-semibold'>
                        {t('Upload Invoice File')}
                      </Label>
                      <div className='flex items-center gap-3'>
                        <Input
                          type='file'
                          accept='.pdf,.png,.jpg,.jpeg'
                          onChange={handleFileChange}
                          className='cursor-pointer text-xs'
                        />
                        {selectedFile && (
                          <Check className='text-success h-5 w-5' />
                        )}
                      </div>
                      <p className='text-muted-foreground text-[10px]'>
                        {t(
                          'Only PDF, PNG, JPG, JPEG files under 10MB are allowed'
                        )}
                      </p>
                    </div>
                  )}

                  {paypalMethod === 'lexware' &&
                    selectedInvoice.payment_method === 'paypal' && (
                      <div className='rounded border border-blue-100 bg-blue-50 p-2.5 text-xs text-blue-800 dark:border-blue-950/40 dark:bg-blue-950/20 dark:text-blue-200'>
                        <p>
                          {t(
                            'This will process the invoice using Lexware Office. Ensure your LexwareApiKey is configured in Settings.'
                          )}
                        </p>
                      </div>
                    )}
                </div>
              ) : (
                <div className='space-y-3 border-t pt-3'>
                  <Label className='text-destructive text-sm font-semibold'>
                    {t('Rejection Reason')}
                  </Label>
                  <Textarea
                    placeholder={t(
                      'Enter rejection reason here... (Will be emailed to the user)'
                    )}
                    value={rejectReason}
                    onChange={(e) => setRejectReason(e.target.value)}
                    rows={3}
                  />
                </div>
              )}
            </div>
          )}

          <DialogFooter className='gap-2 sm:gap-0'>
            {showRejectForm ? (
              <>
                <Button
                  variant='outline'
                  onClick={() => setShowRejectForm(false)}
                  disabled={processing}
                >
                  {t('Back')}
                </Button>
                <Button
                  variant='destructive'
                  onClick={handleReject}
                  disabled={processing}
                >
                  {processing ? t('Processing...') : t('Confirm Rejection')}
                </Button>
              </>
            ) : (
              <>
                <Button
                  variant='outline'
                  className='text-destructive hover:bg-destructive/10'
                  onClick={() => setShowRejectForm(true)}
                  disabled={processing}
                >
                  <X className='mr-1 h-4 w-4' />
                  {t('Reject')}
                </Button>
                <Button onClick={handleComplete} disabled={processing}>
                  <Check className='mr-1 h-4 w-4' />
                  {processing ? t('Processing...') : t('Complete Invoice')}
                </Button>
              </>
            )}
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* View Details Dialog */}
      <Dialog open={viewDialogOpen} onOpenChange={setViewDialogOpen}>
        <DialogContent className='max-h-[85vh] overflow-y-auto sm:max-w-md'>
          <DialogHeader>
            <DialogTitle>{t('Invoice Request Details')}</DialogTitle>
          </DialogHeader>

          {selectedInvoice && (
            <div className='space-y-4 py-2 text-sm'>
              <div className='grid grid-cols-2 gap-x-2 gap-y-3 border-b pb-3'>
                <div>
                  <span className='text-muted-foreground block text-xs'>
                    {t('Request ID')}
                  </span>
                  <span className='font-medium'>{selectedInvoice.id}</span>
                </div>
                <div>
                  <span className='text-muted-foreground block text-xs'>
                    {t('Status')}
                  </span>
                  <span className='font-medium capitalize'>
                    {t(selectedInvoice.status)}
                  </span>
                </div>
                <div>
                  <span className='text-muted-foreground block text-xs'>
                    {t('User')}
                  </span>
                  <span className='font-medium'>
                    {selectedInvoice.username} (ID: {selectedInvoice.user_id})
                  </span>
                </div>
                <div>
                  <span className='text-muted-foreground block text-xs'>
                    {t('Request Time')}
                  </span>
                  <span className='text-xs font-medium'>
                    {formatTimestamp(selectedInvoice.create_time)}
                  </span>
                </div>
                {selectedInvoice.complete_time > 0 && (
                  <div>
                    <span className='text-muted-foreground block text-xs'>
                      {t('Completed Time')}
                    </span>
                    <span className='text-xs font-medium'>
                      {formatTimestamp(selectedInvoice.complete_time)}
                    </span>
                  </div>
                )}
              </div>

              <div className='grid grid-cols-2 gap-x-2 gap-y-3 border-b pb-3'>
                <div>
                  <span className='text-muted-foreground block text-xs'>
                    {t('Amount Paid')}
                  </span>
                  <span className='font-semibold text-red-600'>
                    {formatInvoicePaidAmount(selectedInvoice)}
                  </span>
                  {selectedInvoice.credit_amount_usd &&
                    selectedInvoice.credit_amount_usd > 0 && (
                      <span className='text-muted-foreground block text-xs'>
                        {t('Credited USD')}:{' '}
                        {formatPaymentAmount(
                          selectedInvoice.credit_amount_usd,
                          'USD'
                        )}
                      </span>
                    )}
                </div>
                <div>
                  <span className='text-muted-foreground block text-xs'>
                    {t('Payment Method')}
                  </span>
                  <span className='font-medium uppercase'>
                    {selectedInvoice.payment_method}
                  </span>
                </div>
                <div className='col-span-2'>
                  <span className='text-muted-foreground block text-xs'>
                    {t('Order Trade Number')}
                  </span>
                  <code className='font-mono text-xs break-all'>
                    {selectedInvoice.trade_no}
                  </code>
                </div>
              </div>

              <div className='grid grid-cols-2 gap-x-2 gap-y-3 border-b pb-3'>
                <div>
                  <span className='text-muted-foreground block text-xs'>
                    {t('Billing Type')}
                  </span>
                  <span className='font-medium capitalize'>
                    {t(selectedInvoice.billing_type)}
                  </span>
                </div>
                <div>
                  <span className='text-muted-foreground block text-xs'>
                    {t('Invoice Title')}
                  </span>
                  <span className='font-medium'>{selectedInvoice.title}</span>
                </div>
                {selectedInvoice.billing_type === 'enterprise' && (
                  <div>
                    <span className='text-muted-foreground block text-xs'>
                      {t('Tax ID')}
                    </span>
                    <span className='font-medium'>
                      {selectedInvoice.tax_id || '-'}
                    </span>
                  </div>
                )}
                <div>
                  <span className='text-muted-foreground block text-xs'>
                    {t('Invoice Email')}
                  </span>
                  <span className='font-medium'>{selectedInvoice.email}</span>
                </div>
              </div>

              {selectedInvoice.payment_method === 'paypal' && (
                <div className='space-y-2 border-b pb-3'>
                  <span className='text-muted-foreground block text-xs font-semibold'>
                    {t('PayPal Billing Address')}
                  </span>
                  <div className='grid grid-cols-2 gap-x-2 gap-y-2 text-xs'>
                    <div>
                      <span className='text-muted-foreground block text-[10px]'>
                        {t('Street')}
                      </span>
                      <span>{selectedInvoice.street || '-'}</span>
                    </div>
                    <div>
                      <span className='text-muted-foreground block text-[10px]'>
                        {t('Detailed Address')}
                      </span>
                      <span>{selectedInvoice.address_detail || '-'}</span>
                    </div>
                    <div>
                      <span className='text-muted-foreground block text-[10px]'>
                        {t('City')}
                      </span>
                      <span>{selectedInvoice.city || '-'}</span>
                    </div>
                    <div>
                      <span className='text-muted-foreground block text-[10px]'>
                        {t('Zip Code')}
                      </span>
                      <span>{selectedInvoice.zip_code || '-'}</span>
                    </div>
                    <div>
                      <span className='text-muted-foreground block text-[10px]'>
                        {t('Country')}
                      </span>
                      <span>{selectedInvoice.country || '-'}</span>
                    </div>
                  </div>
                </div>
              )}

              {selectedInvoice.status === 'completed' &&
                selectedInvoice.download_url && (
                  <div>
                    <span className='text-muted-foreground block text-xs'>
                      {t('Download Link')}
                    </span>
                    <Button
                      variant='outline'
                      size='sm'
                      className='mt-1'
                      onClick={() => handleDownloadInvoice(selectedInvoice)}
                    >
                      <Download className='mr-1 h-3.5 w-3.5' />
                      {t('Download Generated Invoice')}
                    </Button>
                  </div>
                )}

              {selectedInvoice.status === 'rejected' &&
                selectedInvoice.message && (
                  <div className='bg-destructive/10 text-destructive border-destructive/20 rounded-lg border p-3'>
                    <span className='block text-xs font-semibold'>
                      {t('Rejection Reason')}
                    </span>
                    <p className='mt-1 text-xs break-all'>
                      {selectedInvoice.message}
                    </p>
                  </div>
                )}
            </div>
          )}

          <DialogFooter>
            <Button onClick={() => setViewDialogOpen(false)}>
              {t('Close')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
