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
import React from 'react'
import { Loader2 } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { formatLocalCurrencyAmount } from '@/lib/currency'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { Skeleton } from '@/components/ui/skeleton'
import { DEFAULT_DISCOUNT_RATE } from '../../constants'
import { formatPaymentAmount, getPaymentIcon } from '../../lib'
import type { PaymentMethod } from '../../types'
import { AgreementDialog } from './agreement-dialog'

interface PaymentConfirmDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onConfirm: () => void
  topupAmount: number
  paymentAmount: number
  paymentMethod: PaymentMethod | undefined
  calculating: boolean
  processing: boolean
  discountRate?: number
  usdExchangeRate?: number
  paymentCurrency?: string
  exchangeRate?: number
  creditAmountUsd?: number
  username?: string
  email?: string
}

export function PaymentConfirmDialog({
  open,
  onOpenChange,
  onConfirm,
  topupAmount,
  paymentAmount,
  paymentMethod,
  calculating,
  processing,
  discountRate = DEFAULT_DISCOUNT_RATE,
  usdExchangeRate = 1,
  paymentCurrency,
  exchangeRate,
  creditAmountUsd,
  username = '',
  email = '',
}: PaymentConfirmDialogProps) {
  const { t } = useTranslation()
  const hasDiscount = discountRate > 0 && discountRate < 1 && paymentAmount > 0
  const originalAmount = hasDiscount ? paymentAmount / discountRate : 0
  const discountAmount = hasDiscount ? originalAmount - paymentAmount : 0

  const [checked1, setChecked1] = React.useState(false)
  const [checked2, setChecked2] = React.useState(false)
  const [checked3, setChecked3] = React.useState(false)
  const [agreementOpen, setAgreementOpen] = React.useState(false)

  const { i18n } = useTranslation()
  const currentLang = (i18n.language || 'en').toLowerCase().split('-')[0]

  React.useEffect(() => {
    if (open) {
      setChecked1(false)
      setChecked2(false)
      setChecked3(false)
    }
  }, [open])

  // Localized agreement checkbox text mappings
  const termsTextMap: Record<
    string,
    {
      clause1Part1: string
      agreementLink: string
      clause1Part2: string
      clause2: string
      clause3: string
    }
  > = {
    zh: {
      clause1Part1: '我已阅读、理解并同意',
      agreementLink: '《AI Token 购买及使用协议》',
      clause1Part2: '、《隐私政策》及页面展示的价格、Token 数量、有效期、消耗规则和退款规则。',
      clause2:
        '我确认 AI Token 属于可即时交付、即时消耗的数字服务额度，非货币、虚拟货币、金融产品或可流通资产。',
      clause3:
        '我理解点击 “确认付款” 即构成电子签署，平台可按订单充值 Token 并按实际使用扣减额度。',
    },
    en: {
      clause1Part1: 'I have read, understood and agree to the ',
      agreementLink: '"AI Token Purchase and Use Agreement"',
      clause1Part2:
        ', "Privacy Policy", and the pricing, token quantity, validity, consumption and refund rules displayed on the page.',
      clause2:
        'I confirm that AI Tokens are instantly delivered and consumed digital service quotas, not currency, virtual currency, financial products, or tradable assets.',
      clause3:
        'I understand that clicking "Confirm Payment" constitutes electronic signing, and the platform will recharge tokens and deduct quota based on actual usage.',
    },
    fr: {
      clause1Part1: "J'ai lu, compris et j'accepte le ",
      agreementLink: '"Contrat d\'Achat et d\'Utilisation de Jetons d\'IA"',
      clause1Part2:
        ', la "Politique de Confidentialité" et les règles de prix, de quantité de jetons, de validité, de consommation et de remboursement affichées.',
      clause2:
        "Je confirme que les jetons d'IA sont des quotas de services numériques immédiatement livrés et consommés, et non des monnaies ou des actifs financiers.",
      clause3:
        'Je comprends que cliquer sur "Confirmer le paiement" constitue une signature électronique, et la plateforme créditera les jetons et déduira le quota.',
    },
    ja: {
      clause1Part1: '「',
      agreementLink: 'AI トークン購入および利用規約',
      clause1Part2:
        '」と「プライバシーポリシー」、およびページに表示されている価格、トークン数量、有効期限、消費および返金ルールを読み、理解し、同意します。',
      clause2:
        'AI トークンは即時に納品および消費されるデジタルサービス利用枠であり、通貨、暗号資産、金融商品、または流通可能な資産ではないことを確認します。',
      clause3:
        '「支払いの確認」をクリックすることで電子署名が成立し、プラットフォームがトークンをチャージし使用実績に応じて差し引くことを理解します。',
    },
    ru: {
      clause1Part1: 'Я прочитал, понял и согласен с ',
      agreementLink: '"Согласием на покупку и использование AI Токенов"',
      clause1Part2:
        ', "Политикой конфиденциальности" и правилами ценообразования, количества токенов, срока действия, потребления и возврата средств.',
      clause2:
        'Я подтверждаю, что AI токены являются мгновенно доставляемыми цифровыми услугами, а не валютой или финансовыми активами.',
      clause3:
        'Я понимаю, что нажатие кнопки "Подтвердить оплату" означает электронную подпись, и платформа начислит токены.',
    },
    vi: {
      clause1Part1: 'Tôi đã đọc, hiểu và đồng ý với ',
      agreementLink: '"Thỏa thuận Mua và Sử dụng Token AI"',
      clause1Part2:
        ', "Chính sách Bảo mật" và các quy định về giá, số lượng token, thời hạn, tiêu thụ và hoàn tiền hiển thị trên trang.',
      clause2:
        'Tôi xác nhận Token AI là hạn mức dịch vụ số được bàn giao và tiêu thụ ngay lập tức, không phải tiền tệ hay tài sản tài chính.',
      clause3:
        'Tôi hiểu rằng nhấp vào "Xác nhận Thanh toán" sẽ cấu thành ký kết điện tử, và nền tảng sẽ cấp token.',
    },
  }

  const termsText = termsTextMap[currentLang] || termsTextMap.en

  return (
    <>
      <AlertDialog open={open} onOpenChange={onOpenChange}>
        <AlertDialogContent className='max-sm:w-[calc(100vw-1.5rem)] sm:max-w-md'>
          <AlertDialogHeader>
            <AlertDialogTitle className='text-xl font-semibold'>
              {t('Confirm Payment')}
            </AlertDialogTitle>
            <AlertDialogDescription>
              {t('Review your payment details')}
            </AlertDialogDescription>
          </AlertDialogHeader>

          <div className='space-y-3 py-3 sm:space-y-4 sm:py-4'>
            <div className='flex items-center justify-between'>
              <span className='text-muted-foreground text-sm'>
                {t('Topup Amount')}
              </span>
              <span className='text-lg font-semibold'>
                {formatLocalCurrencyAmount(topupAmount * usdExchangeRate, {
                  digitsLarge: 2,
                  digitsSmall: 2,
                  abbreviate: false,
                })}
              </span>
            </div>

            <div className='flex items-center justify-between'>
              <span className='text-muted-foreground text-sm'>
                {t('You Pay')}
              </span>
              {calculating ? (
                <Skeleton className='h-6 w-24' />
              ) : (
                <div className='flex items-baseline gap-2'>
                  <span className='text-2xl font-semibold'>
                    {formatPaymentAmount(paymentAmount, paymentCurrency)}
                  </span>
                  {hasDiscount && (
                    <span className='text-muted-foreground text-sm line-through'>
                      {formatPaymentAmount(originalAmount, paymentCurrency)}
                    </span>
                  )}
                </div>
              )}
            </div>

            {paymentCurrency && exchangeRate && (
              <div className='flex items-center justify-between text-sm'>
                <span className='text-muted-foreground'>
                  {t('Exchange Rate')}
                </span>
                <span className='font-medium'>
                  {t('1 USD = {{rate}} {{currency}}', {
                    rate: exchangeRate,
                    currency: paymentCurrency,
                  })}
                </span>
              </div>
            )}

            {creditAmountUsd && paymentCurrency && (
              <div className='flex items-center justify-between text-sm'>
                <span className='text-muted-foreground'>
                  {t('Credited USD')}
                </span>
                <span className='font-medium'>
                  {formatPaymentAmount(creditAmountUsd, 'USD')}
                </span>
              </div>
            )}

            {hasDiscount && !calculating && (
              <div className='bg-muted/50 rounded-lg p-3'>
                <div className='flex items-center justify-between text-sm'>
                  <span className='text-muted-foreground'>{t('You save')}</span>
                  <span className='font-semibold text-green-600'>
                    {formatPaymentAmount(discountAmount, paymentCurrency)}
                  </span>
                </div>
              </div>
            )}

            <div className='border-t pt-4'>
              <div className='flex items-center justify-between'>
                <span className='text-muted-foreground text-sm'>
                  {t('Payment Method')}
                </span>
                <div className='flex items-center gap-2'>
                  {getPaymentIcon(
                    paymentMethod?.type,
                    'h-4 w-4',
                    paymentMethod?.icon,
                    paymentMethod?.name
                  )}
                  <span className='font-medium'>{paymentMethod?.name}</span>
                </div>
              </div>
            </div>

            {/* Terms and conditions checking list */}
            <div className='flex flex-col gap-3 border-t pt-4 text-xs text-muted-foreground/80 leading-relaxed select-none'>
              <label className='flex items-start gap-2.5 cursor-pointer'>
                <input
                  type='checkbox'
                  checked={checked1}
                  onChange={(e) => setChecked1(e.target.checked)}
                  className='mt-0.5 h-3.5 w-3.5 accent-indigo-600 cursor-pointer rounded'
                />
                <span>
                  {termsText.clause1Part1}
                  <span
                    onClick={(e) => {
                      e.preventDefault()
                      e.stopPropagation()
                      setAgreementOpen(true)
                    }}
                    className='text-indigo-600 hover:text-indigo-500 font-semibold underline underline-offset-2 decoration-1 cursor-pointer'
                  >
                    {termsText.agreementLink}
                  </span>
                  {termsText.clause1Part2}
                </span>
              </label>

              <label className='flex items-start gap-2.5 cursor-pointer'>
                <input
                  type='checkbox'
                  checked={checked2}
                  onChange={(e) => setChecked2(e.target.checked)}
                  className='mt-0.5 h-3.5 w-3.5 accent-indigo-600 cursor-pointer rounded'
                />
                <span>{termsText.clause2}</span>
              </label>

              <label className='flex items-start gap-2.5 cursor-pointer'>
                <input
                  type='checkbox'
                  checked={checked3}
                  onChange={(e) => setChecked3(e.target.checked)}
                  className='mt-0.5 h-3.5 w-3.5 accent-indigo-600 cursor-pointer rounded'
                />
                <span>{termsText.clause3}</span>
              </label>
            </div>
          </div>

          <AlertDialogFooter className='grid grid-cols-2 gap-2 sm:flex border-t pt-4'>
            <AlertDialogCancel disabled={processing}>
              {t('Cancel')}
            </AlertDialogCancel>
            <AlertDialogAction
              onClick={onConfirm}
              disabled={
                processing || !checked1 || !checked2 || !checked3 || calculating
              }
            >
              {processing && <Loader2 className='mr-2 h-4 w-4 animate-spin' />}
              {t('Confirm Payment')}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* Agreement detail popup modal */}
      <AgreementDialog
        open={agreementOpen}
        onOpenChange={setAgreementOpen}
        username={username}
        email={email}
      />
    </>
  )
}
