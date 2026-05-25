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
import { useState, useCallback } from 'react'
import i18next from 'i18next'
import { toast } from 'sonner'
import {
  calculateAmount,
  calculatePayPalAmount,
  calculateAlipayAmount,
  calculateStripeAmount,
  calculateWaffoPancakeAmount,
  requestPayment,
  requestPayPalPayment,
  requestAlipayPayment,
  requestStripePayment,
  getTopupInfo,
  isApiSuccess,
} from '../api'
import {
  isStripePayment,
  isPayPalPayment,
  isAlipayPayment,
  isWaffoPancakePayment,
  submitPaymentForm,
} from '../lib'
import type { PaymentAmountQuote } from '../types'

// ============================================================================
// Payment Hook
// ============================================================================

export function usePayment() {
  const [amount, setAmount] = useState<number>(0)
  const [paymentCurrency, setPaymentCurrency] = useState<string>()
  const [exchangeRate, setExchangeRate] = useState<number>()
  const [creditAmountUsd, setCreditAmountUsd] = useState<number>()
  const [calculating, setCalculating] = useState(false)
  const [processing, setProcessing] = useState(false)

  const applyAmountResponse = useCallback(
    (data: string | PaymentAmountQuote) => {
      if (typeof data === 'object' && data !== null) {
        const calculatedAmount = Number.parseFloat(String(data.amount))
        setAmount(calculatedAmount)
        setPaymentCurrency(data.currency || undefined)
        setExchangeRate(data.exchange_rate)
        setCreditAmountUsd(data.credit_amount_usd)
        return calculatedAmount
      }

      const calculatedAmount = Number.parseFloat(data)
      setAmount(calculatedAmount)
      setPaymentCurrency(undefined)
      setExchangeRate(undefined)
      setCreditAmountUsd(undefined)
      return calculatedAmount
    },
    []
  )

  // Calculate payment amount
  const calculatePaymentAmount = useCallback(
    async (topupAmount: number, paymentType: string) => {
      try {
        setCalculating(true)

        const topupInfoRes = await getTopupInfo()
        const enableAlipayTopup = topupInfoRes?.data?.enable_alipay_topup

        const isStripe = isStripePayment(paymentType)
        const isPayPal = isPayPalPayment(paymentType)
        const isAlipay = isAlipayPayment(paymentType, enableAlipayTopup)
        const isPancake = isWaffoPancakePayment(paymentType)
        const response = isStripe
          ? await calculateStripeAmount({ amount: topupAmount })
          : isPayPal
            ? await calculatePayPalAmount({ amount: topupAmount })
            : isAlipay
              ? await calculateAlipayAmount({ amount: topupAmount })
              : isPancake
                ? await calculateWaffoPancakeAmount({ amount: topupAmount })
                : await calculateAmount({ amount: topupAmount })

        if (isApiSuccess(response) && response.data) {
          return applyAmountResponse(response.data)
        }

        // Don't show error for calculation, just set to 0
        setAmount(0)
        setPaymentCurrency(undefined)
        setExchangeRate(undefined)
        setCreditAmountUsd(undefined)
        return 0
      } catch (_error) {
        setAmount(0)
        setPaymentCurrency(undefined)
        setExchangeRate(undefined)
        setCreditAmountUsd(undefined)
        return 0
      } finally {
        setCalculating(false)
      }
    },
    [applyAmountResponse]
  )

  // Process payment
  const processPayment = useCallback(
    async (topupAmount: number, paymentType: string) => {
      try {
        setProcessing(true)

        const topupInfoRes = await getTopupInfo()
        const enableAlipayTopup = topupInfoRes?.data?.enable_alipay_topup

        const isStripe = isStripePayment(paymentType)
        const isPayPal = isPayPalPayment(paymentType)
        const isAlipay = isAlipayPayment(paymentType, enableAlipayTopup)
        const amount = Math.floor(topupAmount)

        const response = isStripe
          ? await requestStripePayment({
              amount,
              payment_method: 'stripe',
            })
          : isPayPal
            ? await requestPayPalPayment({
                amount,
                payment_method: 'paypal',
              })
            : isAlipay
              ? await requestAlipayPayment({
                  amount,
                  payment_method: 'alipay',
                })
              : await requestPayment({
                  amount,
                  payment_method: paymentType,
                })

        if (!isApiSuccess(response)) {
          toast.error(response.message || i18next.t('Payment request failed'))
          return false
        }

        // Handle Stripe payment
        const payLink =
          response.data &&
          typeof response.data === 'object' &&
          'pay_link' in response.data &&
          typeof response.data.pay_link === 'string'
            ? response.data.pay_link
            : ''
        if (isStripe && payLink) {
          window.open(payLink, '_blank')
          toast.success(i18next.t('Redirecting to payment page...'))
          return true
        }

        const approveLink =
          response.data &&
          typeof response.data === 'object' &&
          'approve_link' in response.data &&
          typeof response.data.approve_link === 'string'
            ? response.data.approve_link
            : ''
        if ((isPayPal || isAlipay) && approveLink) {
          window.open(approveLink, '_blank')
          toast.success(i18next.t('Redirecting to payment page...'))
          return true
        }

        // Handle non-Stripe payment
        if (!isStripe && !isPayPal && !isAlipay && response.data) {
          const url = (response as unknown as { url?: string }).url
          if (url) {
            submitPaymentForm(url, response.data)
            toast.success(i18next.t('Redirecting to payment page...'))
            return true
          }
        }

        return false
      } catch (_error) {
        toast.error(i18next.t('Payment request failed'))
        return false
      } finally {
        setProcessing(false)
      }
    },
    []
  )

  return {
    amount,
    paymentCurrency,
    exchangeRate,
    creditAmountUsd,
    calculating,
    processing,
    calculatePaymentAmount,
    processPayment,
    setAmount,
  }
}
