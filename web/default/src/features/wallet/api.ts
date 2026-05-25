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
import { api } from '@/lib/api'
import type {
  RedemptionRequest,
  PaymentRequest,
  AmountRequest,
  AffiliateTransferRequest,
  ApiResponse,
  TopupInfoResponse,
  RedemptionResponse,
  AmountResponse,
  PaymentResponse,
  StripePaymentResponse,
  PayPalPaymentResponse,
  AlipayPaymentResponse,
  AffiliateCodeResponse,
  AffiliateTransferResponse,
  BillingHistoryResponse,
  CompleteOrderRequest,
  CreemPaymentRequest,
  CreemPaymentResponse,
  WaffoPaymentRequest,
  WaffoPaymentResponse,
  WaffoPancakePaymentRequest,
  WaffoPancakePaymentResponse,
  ApplyInvoiceRequest,
  InvoiceListResponse,
  CompleteInvoiceRequest,
  RejectInvoiceRequest,
  InvoiceTitleCard,
  InvoiceTitleCardRequest,
} from './types'

// ============================================================================
// Wallet API Functions
// ============================================================================

/**
 * Check if API response is successful
 */
export function isApiSuccess(response: ApiResponse): boolean {
  return response.success === true || response.message === 'success'
}

/**
 * Get topup configuration info
 */
export async function getTopupInfo(): Promise<TopupInfoResponse> {
  const res = await api.get('/api/user/topup/info')
  return res.data
}

/**
 * Redeem a topup code
 */
export async function redeemTopupCode(
  request: RedemptionRequest
): Promise<RedemptionResponse> {
  const res = await api.post('/api/user/topup', request)
  return res.data
}

/**
 * Calculate payment amount for regular payment
 */
export async function calculateAmount(
  request: AmountRequest
): Promise<AmountResponse> {
  const res = await api.post('/api/user/amount', request, {
    skipBusinessError: true,
  } as Record<string, unknown>)
  return res.data
}

/**
 * Calculate payment amount for Stripe payment
 */
export async function calculateStripeAmount(
  request: AmountRequest
): Promise<AmountResponse> {
  const res = await api.post('/api/user/stripe/amount', request, {
    skipBusinessError: true,
  } as Record<string, unknown>)
  return res.data
}

/**
 * Calculate payment amount for PayPal payment
 */
export async function calculatePayPalAmount(
  request: AmountRequest
): Promise<AmountResponse> {
  const res = await api.post('/api/user/paypal/amount', request, {
    skipBusinessError: true,
  } as Record<string, unknown>)
  return res.data
}

/**
 * Calculate payment amount for Alipay payment
 */
export async function calculateAlipayAmount(
  request: AmountRequest
): Promise<AmountResponse> {
  const res = await api.post('/api/user/alipay/amount', request, {
    skipBusinessError: true,
  } as Record<string, unknown>)
  return res.data
}

/**
 * Request regular payment
 */
export async function requestPayment(
  request: PaymentRequest
): Promise<PaymentResponse> {
  const res = await api.post('/api/user/pay', request, {
    skipBusinessError: true,
  } as Record<string, unknown>)
  return {
    ...res.data,
    url: res.data.url || (res as unknown as { url?: string }).url,
  }
}

/**
 * Request Stripe payment
 */
export async function requestStripePayment(
  request: PaymentRequest
): Promise<StripePaymentResponse> {
  const res = await api.post('/api/user/stripe/pay', request, {
    skipBusinessError: true,
  } as Record<string, unknown>)
  return res.data
}

/**
 * Request PayPal payment
 */
export async function requestPayPalPayment(
  request: PaymentRequest
): Promise<PayPalPaymentResponse> {
  const res = await api.post('/api/user/paypal/pay', request, {
    skipBusinessError: true,
  } as Record<string, unknown>)
  return res.data
}

/**
 * Request Alipay payment
 */
export async function requestAlipayPayment(
  request: PaymentRequest
): Promise<AlipayPaymentResponse> {
  const res = await api.post('/api/user/alipay/pay', request, {
    skipBusinessError: true,
  } as Record<string, unknown>)
  return res.data
}

/**
 * Request Creem payment
 */
export async function requestCreemPayment(
  request: CreemPaymentRequest
): Promise<CreemPaymentResponse> {
  const res = await api.post('/api/user/creem/pay', request, {
    skipBusinessError: true,
  } as Record<string, unknown>)
  return res.data
}

/**
 * Request Waffo payment
 */
export async function requestWaffoPayment(
  request: WaffoPaymentRequest
): Promise<WaffoPaymentResponse> {
  const res = await api.post('/api/user/waffo/pay', request, {
    skipBusinessError: true,
  } as Record<string, unknown>)
  return res.data
}

/**
 * Calculate payment amount for Waffo Pancake payment
 */
export async function calculateWaffoPancakeAmount(
  request: AmountRequest
): Promise<AmountResponse> {
  const res = await api.post('/api/user/waffo-pancake/amount', request, {
    skipBusinessError: true,
  } as Record<string, unknown>)
  return res.data
}

/**
 * Request Waffo Pancake payment
 */
export async function requestWaffoPancakePayment(
  request: WaffoPancakePaymentRequest
): Promise<WaffoPancakePaymentResponse> {
  const res = await api.post('/api/user/waffo-pancake/pay', request, {
    skipBusinessError: true,
  } as Record<string, unknown>)
  return res.data
}

/**
 * Get affiliate code
 */
export async function getAffiliateCode(): Promise<AffiliateCodeResponse> {
  const res = await api.get('/api/user/aff')
  return res.data
}

/**
 * Transfer affiliate quota to balance
 */
export async function transferAffiliateQuota(
  request: AffiliateTransferRequest
): Promise<AffiliateTransferResponse> {
  const res = await api.post('/api/user/aff_transfer', request)
  return res.data
}

/**
 * Get billing history for current user
 */
export async function getUserBillingHistory(
  page: number,
  pageSize: number,
  keyword?: string
): Promise<ApiResponse<BillingHistoryResponse>> {
  const params = new URLSearchParams({
    p: page.toString(),
    page_size: pageSize.toString(),
  })
  if (keyword) {
    params.append('keyword', keyword)
  }
  const res = await api.get(`/api/user/topup/self?${params.toString()}`)
  return res.data
}

/**
 * Get billing history for all users (admin only)
 */
export async function getAllBillingHistory(
  page: number,
  pageSize: number,
  keyword?: string
): Promise<ApiResponse<BillingHistoryResponse>> {
  const params = new URLSearchParams({
    p: page.toString(),
    page_size: pageSize.toString(),
  })
  if (keyword) {
    params.append('keyword', keyword)
  }
  const res = await api.get(`/api/user/topup?${params.toString()}`)
  return res.data
}

/**
 * Complete a pending order (admin only)
 */
export async function completeOrder(
  request: CompleteOrderRequest
): Promise<ApiResponse> {
  const res = await api.post('/api/user/topup/complete', request)
  return res.data
}

// ============================================================================
// Invoice API
// ============================================================================

/**
 * Apply for an invoice
 */
export async function applyInvoice(
  data: ApplyInvoiceRequest
): Promise<ApiResponse<null>> {
  const res = await api.post('/api/invoice/request', data, {
    skipBusinessError: true,
  } as Record<string, unknown>)
  return res.data
}

/**
 * Get current user's invoices
 */
export async function getMyInvoices(
  page: number,
  pageSize: number
): Promise<ApiResponse<InvoiceListResponse>> {
  const params = new URLSearchParams({
    p: String(page),
    page_size: String(pageSize),
  })
  const res = await api.get(`/api/invoice/my?${params.toString()}`, {
    skipBusinessError: true,
  } as Record<string, unknown>)
  return res.data
}

/**
 * Get all invoices (admin only)
 */
export async function getAdminInvoices(
  page: number,
  pageSize: number,
  status?: string,
  paymentMethod?: string,
  keyword?: string
): Promise<ApiResponse<InvoiceListResponse>> {
  const params = new URLSearchParams({
    p: String(page),
    page_size: String(pageSize),
  })
  if (status) params.append('status', status)
  if (paymentMethod) params.append('payment_method', paymentMethod)
  if (keyword) params.append('keyword', keyword)
  const res = await api.get(`/api/admin/invoice/list?${params.toString()}`, {
    skipBusinessError: true,
  } as Record<string, unknown>)
  return res.data
}

/**
 * Upload invoice file (admin only)
 */
export async function uploadInvoiceFile(
  invoiceId: number,
  file: File
): Promise<ApiResponse<{ file_path?: string; url?: string }>> {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('invoice_id', String(invoiceId))
  const res = await api.post('/api/admin/invoice/upload', formData, {
    skipBusinessError: true,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  } as Record<string, unknown>)
  return res.data
}

/**
 * Download an invoice file.
 */
export async function downloadInvoiceFile(invoiceId: number): Promise<void> {
  const res = await api.get(`/api/invoice/download?id=${invoiceId}`, {
    responseType: 'blob',
    disableDuplicate: true,
  } as Record<string, unknown>)

  const contentDisposition = res.headers['content-disposition']
  const filenameMatch =
    typeof contentDisposition === 'string'
      ? contentDisposition.match(/filename="?([^";]+)"?/)
      : null
  const filename = filenameMatch?.[1] || `invoice_${invoiceId}.pdf`
  const blobUrl = window.URL.createObjectURL(res.data)
  const link = document.createElement('a')
  link.href = blobUrl
  link.download = filename
  document.body.appendChild(link)
  link.click()
  link.remove()
  window.URL.revokeObjectURL(blobUrl)
}

export async function getInvoiceTitleCards(): Promise<
  ApiResponse<InvoiceTitleCard[]>
> {
  const res = await api.get('/api/invoice/titles', {
    skipBusinessError: true,
  } as Record<string, unknown>)
  return res.data
}

export async function createInvoiceTitleCard(
  data: InvoiceTitleCardRequest
): Promise<ApiResponse<InvoiceTitleCard>> {
  const res = await api.post('/api/invoice/titles', data, {
    skipBusinessError: true,
  } as Record<string, unknown>)
  return res.data
}

export async function updateInvoiceTitleCard(
  id: number,
  data: InvoiceTitleCardRequest
): Promise<ApiResponse<InvoiceTitleCard>> {
  const res = await api.put(`/api/invoice/titles/${id}`, data, {
    skipBusinessError: true,
  } as Record<string, unknown>)
  return res.data
}

export async function deleteInvoiceTitleCard(
  id: number
): Promise<ApiResponse<null>> {
  const res = await api.delete(`/api/invoice/titles/${id}`, {
    skipBusinessError: true,
  } as Record<string, unknown>)
  return res.data
}

/**
 * Complete an invoice (admin only)
 */
export async function completeInvoice(
  data: CompleteInvoiceRequest
): Promise<ApiResponse<null>> {
  const res = await api.post('/api/admin/invoice/complete', data, {
    skipBusinessError: true,
  } as Record<string, unknown>)
  return res.data
}

/**
 * Reject an invoice (admin only)
 */
export async function rejectInvoice(
  data: RejectInvoiceRequest
): Promise<ApiResponse<null>> {
  const res = await api.post('/api/admin/invoice/reject', data, {
    skipBusinessError: true,
  } as Record<string, unknown>)
  return res.data
}
