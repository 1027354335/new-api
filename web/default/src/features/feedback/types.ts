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
export type FeedbackStatus =
  | 'open'
  | 'in_progress'
  | 'resolved'
  | 'rejected'
  | 'closed'

export type FeedbackCategory = 'bug' | 'account' | 'billing' | 'model' | 'other'

export type FeedbackPriority = 'low' | 'normal' | 'high' | 'urgent'

export interface Feedback {
  id: number
  user_id: number
  username: string
  email: string
  title: string
  content: string
  category: FeedbackCategory
  priority: FeedbackPriority
  status: FeedbackStatus
  admin_id: number
  admin_name: string
  reply: string
  create_time: number
  update_time: number
  resolved_time: number
}

export interface ApiResponse<T = unknown> {
  success?: boolean
  message?: string
  data?: T
}

export interface FeedbackListResponse {
  items: Feedback[]
  total: number
}

export interface CreateFeedbackRequest {
  title: string
  content: string
  category: FeedbackCategory
  priority: FeedbackPriority
}

export interface AdminReplyFeedbackRequest {
  status: 'resolved' | 'rejected'
  reply: string
}
