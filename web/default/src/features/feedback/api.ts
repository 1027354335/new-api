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
  AdminReplyFeedbackRequest,
  ApiResponse,
  CreateFeedbackRequest,
  Feedback,
  FeedbackCategory,
  FeedbackListResponse,
  FeedbackPriority,
  FeedbackStatus,
} from './types'

export async function createFeedback(
  data: CreateFeedbackRequest
): Promise<ApiResponse<Feedback>> {
  const res = await api.post('/api/feedback', data, {
    skipBusinessError: true,
  } as Record<string, unknown>)
  return res.data
}

export async function getMyFeedbacks(params: {
  page: number
  pageSize: number
  status?: FeedbackStatus | 'all'
}): Promise<ApiResponse<FeedbackListResponse>> {
  const query = new URLSearchParams({
    p: String(params.page),
    page_size: String(params.pageSize),
  })
  if (params.status && params.status !== 'all') {
    query.set('status', params.status)
  }
  const res = await api.get(`/api/feedback/self?${query.toString()}`, {
    skipBusinessError: true,
  } as Record<string, unknown>)
  return res.data
}

export async function closeMyFeedback(
  id: number
): Promise<ApiResponse<Feedback>> {
  const res = await api.post(`/api/feedback/self/${id}/close`, null, {
    skipBusinessError: true,
  } as Record<string, unknown>)
  return res.data
}

export async function getAdminFeedbacks(params: {
  page: number
  pageSize: number
  status?: FeedbackStatus | 'all'
  category?: FeedbackCategory | 'all'
  priority?: FeedbackPriority | 'all'
  keyword?: string
}): Promise<ApiResponse<FeedbackListResponse>> {
  const query = new URLSearchParams({
    p: String(params.page),
    page_size: String(params.pageSize),
  })
  if (params.status && params.status !== 'all')
    query.set('status', params.status)
  if (params.category && params.category !== 'all') {
    query.set('category', params.category)
  }
  if (params.priority && params.priority !== 'all') {
    query.set('priority', params.priority)
  }
  if (params.keyword?.trim()) query.set('keyword', params.keyword.trim())

  const res = await api.get(`/api/admin/feedback?${query.toString()}`, {
    skipBusinessError: true,
  } as Record<string, unknown>)
  return res.data
}

export async function adminReplyFeedback(
  id: number,
  data: AdminReplyFeedbackRequest
): Promise<ApiResponse<Feedback>> {
  const res = await api.post(`/api/admin/feedback/${id}/reply`, data, {
    skipBusinessError: true,
  } as Record<string, unknown>)
  return res.data
}

export async function adminUpdateFeedbackStatus(
  id: number,
  status: FeedbackStatus
): Promise<ApiResponse<Feedback>> {
  const res = await api.patch(`/api/admin/feedback/${id}/status`, { status }, {
    skipBusinessError: true,
  } as Record<string, unknown>)
  return res.data
}
