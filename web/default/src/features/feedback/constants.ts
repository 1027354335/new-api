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
import type { TFunction } from 'i18next'
import type {
  FeedbackCategory,
  FeedbackPriority,
  FeedbackStatus,
} from './types'

export const FEEDBACK_STATUSES: FeedbackStatus[] = [
  'open',
  'in_progress',
  'resolved',
  'rejected',
  'closed',
]

export const FEEDBACK_CATEGORIES: FeedbackCategory[] = [
  'bug',
  'account',
  'billing',
  'model',
  'other',
]

export const FEEDBACK_PRIORITIES: FeedbackPriority[] = [
  'low',
  'normal',
  'high',
  'urgent',
]

export function getFeedbackStatusLabel(status: FeedbackStatus): string {
  const labels: Record<FeedbackStatus, string> = {
    open: 'Open',
    in_progress: 'In Progress',
    resolved: 'Resolved',
    rejected: 'Rejected',
    closed: 'Closed',
  }
  return labels[status]
}

export function getFeedbackCategoryLabel(category: FeedbackCategory): string {
  const labels: Record<FeedbackCategory, string> = {
    bug: 'Bug Report',
    account: 'Account',
    billing: 'Billing',
    model: 'Model',
    other: 'Other',
  }
  return labels[category]
}

export function getFeedbackPriorityLabel(priority: FeedbackPriority): string {
  const labels: Record<FeedbackPriority, string> = {
    low: 'Low',
    normal: 'Normal',
    high: 'High',
    urgent: 'Urgent',
  }
  return labels[priority]
}

export function getFeedbackStatusVariant(status: FeedbackStatus) {
  const variants: Record<
    FeedbackStatus,
    'info' | 'warning' | 'success' | 'danger' | 'neutral'
  > = {
    open: 'info',
    in_progress: 'warning',
    resolved: 'success',
    rejected: 'danger',
    closed: 'neutral',
  }
  return variants[status]
}

export function getFeedbackPriorityVariant(priority: FeedbackPriority) {
  const variants: Record<
    FeedbackPriority,
    'neutral' | 'info' | 'warning' | 'danger'
  > = {
    low: 'neutral',
    normal: 'info',
    high: 'warning',
    urgent: 'danger',
  }
  return variants[priority]
}

export function getFeedbackStatusOptions(t: TFunction) {
  return FEEDBACK_STATUSES.map((status) => ({
    value: status,
    label: t(getFeedbackStatusLabel(status)),
  }))
}

export function getFeedbackCategoryOptions(t: TFunction) {
  return FEEDBACK_CATEGORIES.map((category) => ({
    value: category,
    label: t(getFeedbackCategoryLabel(category)),
  }))
}

export function getFeedbackPriorityOptions(t: TFunction) {
  return FEEDBACK_PRIORITIES.map((priority) => ({
    value: priority,
    label: t(getFeedbackPriorityLabel(priority)),
  }))
}
