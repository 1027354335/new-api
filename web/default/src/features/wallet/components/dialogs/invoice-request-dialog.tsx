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
import { useEffect, useMemo, useState } from 'react'
import { z } from 'zod'
import { type Resolver, useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Check, ChevronsUpDown, WalletCards } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/components/ui/command'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { applyInvoice, getInvoiceTitleCards, isApiSuccess } from '../../api'
import type { InvoiceTitleCard, TopupRecord } from '../../types'

// ============================================================================
// Invoice Request Form Schema
// ============================================================================

type InvoiceFormValues = {
  billing_type: 'personal' | 'enterprise'
  title: string
  tax_id?: string
  email: string
  street?: string
  address_detail?: string
  city?: string
  zip_code?: string
  country?: string
}

const emptyInvoiceFormValues: InvoiceFormValues = {
  billing_type: 'personal',
  title: '',
  tax_id: '',
  email: '',
  street: '',
  address_detail: '',
  city: '',
  zip_code: '',
  country: '',
}

export const COUNTRY_OPTIONS = [
  { code: 'AD', name: 'Andorra' },
  { code: 'AE', name: 'United Arab Emirates' },
  { code: 'AF', name: 'Afghanistan' },
  { code: 'AG', name: 'Antigua and Barbuda' },
  { code: 'AI', name: 'Anguilla' },
  { code: 'AL', name: 'Albania' },
  { code: 'AM', name: 'Armenia' },
  { code: 'AO', name: 'Angola' },
  { code: 'AR', name: 'Argentina' },
  { code: 'AT', name: 'Austria' },
  { code: 'AU', name: 'Australia' },
  { code: 'AW', name: 'Aruba' },
  { code: 'AZ', name: 'Azerbaijan' },
  { code: 'BA', name: 'Bosnia and Herzegovina' },
  { code: 'BB', name: 'Barbados' },
  { code: 'BD', name: 'Bangladesh' },
  { code: 'BE', name: 'Belgium' },
  { code: 'BF', name: 'Burkina Faso' },
  { code: 'BG', name: 'Bulgaria' },
  { code: 'BH', name: 'Bahrain' },
  { code: 'BI', name: 'Burundi' },
  { code: 'BJ', name: 'Benin' },
  { code: 'BM', name: 'Bermuda' },
  { code: 'BN', name: 'Brunei Darussalam' },
  { code: 'BO', name: 'Bolivia' },
  { code: 'BR', name: 'Brazil' },
  { code: 'BS', name: 'Bahamas' },
  { code: 'BT', name: 'Bhutan' },
  { code: 'BW', name: 'Botswana' },
  { code: 'BY', name: 'Belarus' },
  { code: 'BZ', name: 'Belize' },
  { code: 'CA', name: 'Canada' },
  { code: 'CD', name: 'Congo, The Democratic Republic of the' },
  { code: 'CF', name: 'Central African Republic' },
  { code: 'CG', name: 'Congo' },
  { code: 'CH', name: 'Switzerland' },
  { code: 'CI', name: "Cote d'Ivoire" },
  { code: 'CL', name: 'Chile' },
  { code: 'CM', name: 'Cameroon' },
  { code: 'CN', name: 'China' },
  { code: 'CO', name: 'Colombia' },
  { code: 'CR', name: 'Costa Rica' },
  { code: 'CU', name: 'Cuba' },
  { code: 'CV', name: 'Cape Verde' },
  { code: 'CY', name: 'Cyprus' },
  { code: 'CZ', name: 'Czechia' },
  { code: 'DE', name: 'Germany' },
  { code: 'DJ', name: 'Djibouti' },
  { code: 'DK', name: 'Denmark' },
  { code: 'DM', name: 'Dominica' },
  { code: 'DO', name: 'Dominican Republic' },
  { code: 'DZ', name: 'Algeria' },
  { code: 'EC', name: 'Ecuador' },
  { code: 'EE', name: 'Estonia' },
  { code: 'EG', name: 'Egypt' },
  { code: 'ES', name: 'Spain' },
  { code: 'ET', name: 'Ethiopia' },
  { code: 'FI', name: 'Finland' },
  { code: 'FJ', name: 'Fiji' },
  { code: 'FO', name: 'Faroe Islands' },
  { code: 'FR', name: 'France' },
  { code: 'GA', name: 'Gabon' },
  { code: 'GB', name: 'United Kingdom' },
  { code: 'GD', name: 'Grenada' },
  { code: 'GE', name: 'Georgia' },
  { code: 'GH', name: 'Ghana' },
  { code: 'GI', name: 'Gibraltar' },
  { code: 'GL', name: 'Greenland' },
  { code: 'GM', name: 'Gambia' },
  { code: 'GN', name: 'Guinea' },
  { code: 'GQ', name: 'Equatorial Guinea' },
  { code: 'GR', name: 'Greece' },
  { code: 'GT', name: 'Guatemala' },
  { code: 'GW', name: 'Guinea-Bissau' },
  { code: 'GY', name: 'Guyana' },
  { code: 'HK', name: 'Hong Kong' },
  { code: 'HN', name: 'Honduras' },
  { code: 'HR', name: 'Croatia' },
  { code: 'HT', name: 'Haiti' },
  { code: 'HU', name: 'Hungary' },
  { code: 'ID', name: 'Indonesia' },
  { code: 'IE', name: 'Ireland' },
  { code: 'IL', name: 'Israel' },
  { code: 'IN', name: 'India' },
  { code: 'IQ', name: 'Iraq' },
  { code: 'IR', name: 'Iran' },
  { code: 'IS', name: 'Iceland' },
  { code: 'IT', name: 'Italy' },
  { code: 'JM', name: 'Jamaica' },
  { code: 'JO', name: 'Jordan' },
  { code: 'JP', name: 'Japan' },
  { code: 'KE', name: 'Kenya' },
  { code: 'KG', name: 'Kyrgyzstan' },
  { code: 'KH', name: 'Cambodia' },
  { code: 'KM', name: 'Comoros' },
  { code: 'KN', name: 'Saint Kitts and Nevis' },
  { code: 'KR', name: 'Korea, Republic of' },
  { code: 'KW', name: 'Kuwait' },
  { code: 'KY', name: 'Cayman Islands' },
  { code: 'KZ', name: 'Kazakhstan' },
  { code: 'LA', name: 'Lao People Democratic Republic' },
  { code: 'LB', name: 'Lebanon' },
  { code: 'LC', name: 'Saint Lucia' },
  { code: 'LI', name: 'Liechtenstein' },
  { code: 'LK', name: 'Sri Lanka' },
  { code: 'LR', name: 'Liberia' },
  { code: 'LS', name: 'Lesotho' },
  { code: 'LT', name: 'Lithuania' },
  { code: 'LU', name: 'Luxembourg' },
  { code: 'LV', name: 'Latvia' },
  { code: 'LY', name: 'Libya' },
  { code: 'MA', name: 'Morocco' },
  { code: 'MC', name: 'Monaco' },
  { code: 'MD', name: 'Moldova' },
  { code: 'ME', name: 'Montenegro' },
  { code: 'MG', name: 'Madagascar' },
  { code: 'MK', name: 'North Macedonia' },
  { code: 'ML', name: 'Mali' },
  { code: 'MM', name: 'Myanmar' },
  { code: 'MN', name: 'Mongolia' },
  { code: 'MO', name: 'Macao' },
  { code: 'MR', name: 'Mauritania' },
  { code: 'MT', name: 'Malta' },
  { code: 'MU', name: 'Mauritius' },
  { code: 'MV', name: 'Maldives' },
  { code: 'MW', name: 'Malawi' },
  { code: 'MX', name: 'Mexico' },
  { code: 'MY', name: 'Malaysia' },
  { code: 'MZ', name: 'Mozambique' },
  { code: 'NA', name: 'Namibia' },
  { code: 'NE', name: 'Niger' },
  { code: 'NG', name: 'Nigeria' },
  { code: 'NI', name: 'Nicaragua' },
  { code: 'NL', name: 'Netherlands' },
  { code: 'NO', name: 'Norway' },
  { code: 'NP', name: 'Nepal' },
  { code: 'NZ', name: 'New Zealand' },
  { code: 'OM', name: 'Oman' },
  { code: 'PA', name: 'Panama' },
  { code: 'PE', name: 'Peru' },
  { code: 'PG', name: 'Papua New Guinea' },
  { code: 'PH', name: 'Philippines' },
  { code: 'PK', name: 'Pakistan' },
  { code: 'PL', name: 'Poland' },
  { code: 'PR', name: 'Puerto Rico' },
  { code: 'PT', name: 'Portugal' },
  { code: 'PY', name: 'Paraguay' },
  { code: 'QA', name: 'Qatar' },
  { code: 'RO', name: 'Romania' },
  { code: 'RS', name: 'Serbia' },
  { code: 'RU', name: 'Russian Federation' },
  { code: 'RW', name: 'Rwanda' },
  { code: 'SA', name: 'Saudi Arabia' },
  { code: 'SC', name: 'Seychelles' },
  { code: 'SD', name: 'Sudan' },
  { code: 'SE', name: 'Sweden' },
  { code: 'SG', name: 'Singapore' },
  { code: 'SI', name: 'Slovenia' },
  { code: 'SK', name: 'Slovakia' },
  { code: 'SL', name: 'Sierra Leone' },
  { code: 'SM', name: 'San Marino' },
  { code: 'SN', name: 'Senegal' },
  { code: 'SO', name: 'Somalia' },
  { code: 'SR', name: 'Suriname' },
  { code: 'SV', name: 'El Salvador' },
  { code: 'SY', name: 'Syrian Arab Republic' },
  { code: 'SZ', name: 'Eswatini' },
  { code: 'TD', name: 'Chad' },
  { code: 'TG', name: 'Togo' },
  { code: 'TH', name: 'Thailand' },
  { code: 'TJ', name: 'Tajikistan' },
  { code: 'TN', name: 'Tunisia' },
  { code: 'TR', name: 'Turkey' },
  { code: 'TT', name: 'Trinidad and Tobago' },
  { code: 'TW', name: 'Taiwan' },
  { code: 'TZ', name: 'Tanzania' },
  { code: 'UA', name: 'Ukraine' },
  { code: 'UG', name: 'Uganda' },
  { code: 'US', name: 'United States' },
  { code: 'UY', name: 'Uruguay' },
  { code: 'UZ', name: 'Uzbekistan' },
  { code: 'VA', name: 'Holy See' },
  { code: 'VC', name: 'Saint Vincent and the Grenadines' },
  { code: 'VE', name: 'Venezuela' },
  { code: 'VN', name: 'Viet Nam' },
  { code: 'ZA', name: 'South Africa' },
  { code: 'ZM', name: 'Zambia' },
  { code: 'ZW', name: 'Zimbabwe' },
]

const createInvoiceFormSchema = (
  isPayPal: boolean,
  translate: (key: string) => string
) =>
  z
    .object({
      billing_type: z.enum(['personal', 'enterprise']),
      title: z.string().trim().min(1, translate('Invoice title is required')),
      tax_id: z.string().optional(),
      email: z.string().trim().email(translate('Please enter a valid email')),
      street: z.string().optional(),
      address_detail: z.string().optional(),
      city: z.string().optional(),
      zip_code: z.string().optional(),
      country: z.string().optional(),
    })
    .superRefine((data, ctx) => {
      if (data.billing_type === 'enterprise' && !data.tax_id?.trim()) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: translate('Tax ID is required for enterprise billing'),
          path: ['tax_id'],
        })
      }
      if (!isPayPal) return

      const paypalRequiredFields: Array<keyof InvoiceFormValues> = [
        'street',
        'city',
        'zip_code',
        'country',
      ]
      paypalRequiredFields.forEach((field) => {
        if (!data[field]?.trim()) {
          ctx.addIssue({
            code: z.ZodIssueCode.custom,
            message: translate('This field is required for PayPal invoices'),
            path: [field],
          })
        }
      })

      if (data.country && !/^[A-Z]{2}$/.test(data.country)) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: translate('Please select a valid country code'),
          path: ['country'],
        })
      }
    })

// ============================================================================
// Invoice Request Dialog
// ============================================================================

interface InvoiceRequestDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  topupRecord: TopupRecord
  onSuccess: () => void
}

export function InvoiceRequestDialog({
  open,
  onOpenChange,
  topupRecord,
  onSuccess,
}: InvoiceRequestDialogProps) {
  const { t } = useTranslation()
  const [submitting, setSubmitting] = useState(false)
  const [countryOpen, setCountryOpen] = useState(false)
  const [titleCards, setTitleCards] = useState<InvoiceTitleCard[]>([])
  const [selectedTitleCardId, setSelectedTitleCardId] = useState<string>('none')

  const isPayPal = topupRecord.payment_method === 'paypal'
  const invoiceFormSchema = useMemo(
    () => createInvoiceFormSchema(isPayPal, t),
    [isPayPal, t]
  )

  const form = useForm<InvoiceFormValues>({
    resolver: zodResolver(
      invoiceFormSchema
    ) as unknown as Resolver<InvoiceFormValues>,
    defaultValues: emptyInvoiceFormValues,
  })

  const billingType = form.watch('billing_type')

  const applyTitleCard = (card: InvoiceTitleCard) => {
    form.reset({
      billing_type: card.billing_type,
      title: card.title,
      tax_id: card.tax_id,
      email: card.email,
      street: card.street,
      address_detail: card.address_detail || '',
      city: card.city,
      zip_code: card.zip_code,
      country: card.country,
    })
  }

  useEffect(() => {
    if (open) {
      setCountryOpen(false)
      setSelectedTitleCardId('none')
      form.reset(emptyInvoiceFormValues)
      getInvoiceTitleCards()
        .then((response) => {
          if (!isApiSuccess(response)) return
          const cards = response.data || []
          setTitleCards(cards)
          const defaultCard = cards.find((card) => card.is_default)
          if (defaultCard) {
            setSelectedTitleCardId(String(defaultCard.id))
            applyTitleCard(defaultCard)
          }
        })
        .catch(() => setTitleCards([]))
    }
  }, [open, form])

  const onSubmit = async (values: InvoiceFormValues) => {
    setSubmitting(true)
    try {
      const response = await applyInvoice({
        topup_id: topupRecord.id,
        billing_type: values.billing_type,
        title: values.title,
        tax_id: values.tax_id || undefined,
        email: values.email,
        street: values.street || undefined,
        address_detail: values.address_detail || undefined,
        city: values.city || undefined,
        zip_code: values.zip_code || undefined,
        country: values.country || undefined,
      })
      if (isApiSuccess(response)) {
        toast.success(t('Invoice request submitted successfully'))
        onOpenChange(false)
        onSuccess()
      } else {
        toast.error(response.message || t('Failed to submit invoice request'))
      }
    } catch (error) {
      // eslint-disable-next-line no-console
      console.error('Failed to apply invoice:', error)
      toast.error(t('Failed to submit invoice request'))
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='max-h-[85vh] overflow-y-auto sm:max-w-lg'>
        <DialogHeader>
          <DialogTitle>{t('Request Invoice')}</DialogTitle>
          <DialogDescription>
            {t('Fill in the invoice information for order {{tradeNo}}', {
              tradeNo: topupRecord.trade_no,
            })}
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className='space-y-4'>
            {titleCards.length > 0 && (
              <div className='space-y-2 rounded-lg border p-3'>
                <Label className='flex items-center gap-1.5 text-sm'>
                  <WalletCards className='h-4 w-4' />
                  {t('Invoice Title Card')}
                </Label>
                <Select
                  items={[
                    { value: 'none', label: t('Do not use a title card') },
                    ...titleCards.map((card) => ({
                      value: String(card.id),
                      label: card.name,
                    })),
                  ]}
                  value={selectedTitleCardId}
                  onValueChange={(value) => {
                    const nextValue = value || 'none'
                    setSelectedTitleCardId(nextValue)
                    const card = titleCards.find(
                      (item) => String(item.id) === nextValue
                    )
                    if (card) applyTitleCard(card)
                    else form.reset(emptyInvoiceFormValues)
                  }}
                >
                  <SelectTrigger className='w-full'>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent alignItemWithTrigger={false}>
                    <SelectGroup>
                      <SelectItem value='none'>
                        {t('Do not use a title card')}
                      </SelectItem>
                      {titleCards.map((card) => (
                        <SelectItem key={card.id} value={String(card.id)}>
                          {card.name}
                          {card.is_default ? ` (${t('Default')})` : ''}
                        </SelectItem>
                      ))}
                    </SelectGroup>
                  </SelectContent>
                </Select>
              </div>
            )}

            {/* Billing Type */}
            <FormField
              control={form.control}
              name='billing_type'
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('Billing Type')}</FormLabel>
                  <Select
                    items={[
                      {
                        value: 'personal',
                        label: t('Personal'),
                      },
                      {
                        value: 'enterprise',
                        label: t('Enterprise'),
                      },
                    ]}
                    value={field.value}
                    onValueChange={field.onChange}
                  >
                    <FormControl>
                      <SelectTrigger className='w-full'>
                        <SelectValue />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent alignItemWithTrigger={false}>
                      <SelectGroup>
                        <SelectItem value='personal'>
                          {t('Personal')}
                        </SelectItem>
                        <SelectItem value='enterprise'>
                          {t('Enterprise')}
                        </SelectItem>
                      </SelectGroup>
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />

            {/* Title */}
            <FormField
              control={form.control}
              name='title'
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('Invoice Title')}</FormLabel>
                  <FormControl>
                    <Input placeholder={t('Enter invoice title')} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            {/* Tax ID (required for enterprise) */}
            {billingType === 'enterprise' && (
              <FormField
                control={form.control}
                name='tax_id'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('Tax ID')}</FormLabel>
                    <FormControl>
                      <Input
                        placeholder={t('Enter tax identification number')}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            )}

            {/* Email */}
            <FormField
              control={form.control}
              name='email'
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('Email')}</FormLabel>
                  <FormControl>
                    <Input
                      type='email'
                      placeholder={t('Enter email to receive invoice')}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            {/* PayPal-specific address fields */}
            {isPayPal && (
              <>
                <FormField
                  control={form.control}
                  name='street'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('Street')}</FormLabel>
                      <FormControl>
                        <Input
                          placeholder={t('Enter street address')}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name='address_detail'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('Detailed Address')}</FormLabel>
                      <FormControl>
                        <Input
                          placeholder={t('Enter detailed address')}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <div className='grid grid-cols-2 gap-4'>
                  <FormField
                    control={form.control}
                    name='city'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('City')}</FormLabel>
                        <FormControl>
                          <Input placeholder={t('Enter city')} {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name='zip_code'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('Zip Code')}</FormLabel>
                        <FormControl>
                          <Input placeholder={t('Enter zip code')} {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>

                <FormField
                  control={form.control}
                  name='country'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('Country Code')}</FormLabel>
                      <Popover open={countryOpen} onOpenChange={setCountryOpen}>
                        <PopoverTrigger
                          render={
                            <Button
                              type='button'
                              variant='outline'
                              role='combobox'
                              aria-expanded={countryOpen}
                              className={cn(
                                'w-full justify-between font-normal',
                                !field.value && 'text-muted-foreground'
                              )}
                            />
                          }
                        >
                          {field.value
                            ? `${COUNTRY_OPTIONS.find((country) => country.code === field.value)?.name || field.value} (${field.value})`
                            : t('Select country code')}
                          <ChevronsUpDown className='ml-2 size-4 shrink-0 opacity-50' />
                        </PopoverTrigger>
                        <PopoverContent
                          className='w-[min(28rem,calc(100vw-2rem))] p-0'
                          align='start'
                          collisionPadding={8}
                        >
                          <Command>
                            <CommandInput
                              placeholder={t('Search country code...')}
                            />
                            <CommandList>
                              <CommandEmpty>
                                {t('No country found.')}
                              </CommandEmpty>
                              <CommandGroup>
                                {COUNTRY_OPTIONS.map((country) => (
                                  <CommandItem
                                    key={country.code}
                                    value={`${country.code} ${country.name}`}
                                    onSelect={() => {
                                      field.onChange(country.code)
                                      setCountryOpen(false)
                                    }}
                                  >
                                    <span className='min-w-0 flex-1 truncate'>
                                      {country.name}
                                    </span>
                                    <span className='text-muted-foreground font-mono text-xs'>
                                      {country.code}
                                    </span>
                                    <Check
                                      className={cn(
                                        'size-4',
                                        field.value === country.code
                                          ? 'opacity-100'
                                          : 'opacity-0'
                                      )}
                                    />
                                  </CommandItem>
                                ))}
                              </CommandGroup>
                            </CommandList>
                          </Command>
                        </PopoverContent>
                      </Popover>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </>
            )}

            <DialogFooter>
              <Button
                type='button'
                variant='outline'
                onClick={() => onOpenChange(false)}
                disabled={submitting}
              >
                {t('Cancel')}
              </Button>
              <Button type='submit' disabled={submitting}>
                {submitting ? t('Submitting...') : t('Submit Request')}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
