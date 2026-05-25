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
import { useCallback, useEffect, useState } from 'react'
import { Check, ChevronsUpDown, Edit, Plus, Star, Trash2 } from 'lucide-react'
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
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import { ScrollArea } from '@/components/ui/scroll-area'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  createInvoiceTitleCard,
  deleteInvoiceTitleCard,
  getInvoiceTitleCards,
  isApiSuccess,
  updateInvoiceTitleCard,
} from '../../api'
import type { InvoiceTitleCard, InvoiceTitleCardRequest } from '../../types'
import { COUNTRY_OPTIONS } from './invoice-request-dialog'

type InvoiceTitleFormState = InvoiceTitleCardRequest

const emptyForm: InvoiceTitleFormState = {
  name: '',
  billing_type: 'personal',
  title: '',
  tax_id: '',
  email: '',
  street: '',
  address_detail: '',
  city: '',
  zip_code: '',
  country: '',
  is_default: false,
}

interface InvoiceTitleManagementDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onChanged?: () => void
}

export function InvoiceTitleManagementDialog({
  open,
  onOpenChange,
  onChanged,
}: InvoiceTitleManagementDialogProps) {
  const { t } = useTranslation()
  const [cards, setCards] = useState<InvoiceTitleCard[]>([])
  const [loading, setLoading] = useState(false)
  const [saving, setSaving] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [form, setForm] = useState<InvoiceTitleFormState>(emptyForm)
  const [countryOpen, setCountryOpen] = useState(false)

  const loadCards = useCallback(async () => {
    setLoading(true)
    try {
      const response = await getInvoiceTitleCards()
      if (isApiSuccess(response)) {
        setCards(response.data || [])
      } else {
        toast.error(response.message || t('Failed to load invoice title cards'))
      }
    } catch {
      toast.error(t('Failed to load invoice title cards'))
    } finally {
      setLoading(false)
    }
  }, [t])

  useEffect(() => {
    if (open) {
      loadCards()
      setEditingId(null)
      setForm(emptyForm)
      setCountryOpen(false)
    }
  }, [loadCards, open])

  const updateField = <K extends keyof InvoiceTitleFormState>(
    key: K,
    value: InvoiceTitleFormState[K]
  ) => {
    setForm((current) => ({ ...current, [key]: value }))
  }

  const startEdit = (card: InvoiceTitleCard) => {
    setEditingId(card.id)
    setForm({
      name: card.name,
      billing_type: card.billing_type,
      title: card.title,
      tax_id: card.tax_id,
      email: card.email,
      street: card.street,
      address_detail: card.address_detail || '',
      city: card.city,
      zip_code: card.zip_code,
      country: card.country,
      is_default: card.is_default,
    })
  }

  const resetForm = () => {
    setEditingId(null)
    setForm(emptyForm)
  }

  const validateForm = () => {
    if (!form.name.trim()) return t('Card name is required')
    if (!form.title.trim()) return t('Invoice title is required')
    if (!form.email.trim()) return t('Email is required')
    if (form.billing_type === 'enterprise' && !form.tax_id.trim()) {
      return t('Tax ID is required for enterprise billing')
    }
    const country = form.country.trim().toUpperCase()
    if (
      country &&
      !COUNTRY_OPTIONS.some((option) => option.code === country)
    ) {
      return t('Please select a valid country code')
    }
    return ''
  }

  const handleSave = async () => {
    const error = validateForm()
    if (error) {
      toast.error(error)
      return
    }

    setSaving(true)
    try {
      const payload: InvoiceTitleCardRequest = {
        ...form,
        name: form.name.trim(),
        title: form.title.trim(),
        tax_id: form.tax_id.trim(),
        email: form.email.trim(),
        street: form.street.trim(),
        address_detail: form.address_detail.trim(),
        city: form.city.trim(),
        zip_code: form.zip_code.trim(),
        country: form.country.trim().toUpperCase(),
      }
      const response =
        editingId == null
          ? await createInvoiceTitleCard(payload)
          : await updateInvoiceTitleCard(editingId, payload)
      if (isApiSuccess(response)) {
        toast.success(t('Invoice title card saved'))
        resetForm()
        await loadCards()
        onChanged?.()
      } else {
        toast.error(response.message || t('Failed to save invoice title card'))
      }
    } catch {
      toast.error(t('Failed to save invoice title card'))
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (card: InvoiceTitleCard) => {
    setSaving(true)
    try {
      const response = await deleteInvoiceTitleCard(card.id)
      if (isApiSuccess(response)) {
        toast.success(t('Invoice title card deleted'))
        if (editingId === card.id) resetForm()
        await loadCards()
        onChanged?.()
      } else {
        toast.error(response.message || t('Failed to delete invoice title card'))
      }
    } catch {
      toast.error(t('Failed to delete invoice title card'))
    } finally {
      setSaving(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='flex max-h-[calc(100dvh-2rem)] flex-col max-sm:h-dvh max-sm:w-screen max-sm:max-w-none max-sm:rounded-none max-sm:p-4 sm:max-w-3xl'>
        <DialogHeader>
          <DialogTitle>{t('Invoice Title Management')}</DialogTitle>
          <DialogDescription>
            {t('Manage reusable invoice title cards for invoice requests')}
          </DialogDescription>
        </DialogHeader>

        <div className='grid min-h-0 flex-1 gap-3 overflow-hidden md:grid-cols-[minmax(0,1fr)_minmax(280px,0.9fr)] md:gap-4'>
          <ScrollArea className='min-h-[150px] max-h-[30dvh] rounded-lg border p-3 md:max-h-none'>
            {loading ? (
              <div className='text-muted-foreground text-sm'>
                {t('Loading...')}
              </div>
            ) : cards.length === 0 ? (
              <div className='text-muted-foreground flex h-40 items-center justify-center text-sm'>
                {t('No invoice title cards yet')}
              </div>
            ) : (
              <div className='space-y-2'>
                {cards.map((card) => (
                  <div key={card.id} className='rounded-lg border p-3'>
                    <div className='flex items-start justify-between gap-2'>
                      <div className='min-w-0'>
                        <div className='flex min-w-0 items-center gap-1.5'>
                          <span className='truncate text-sm font-semibold'>
                            {card.name}
                          </span>
                          {card.is_default && (
                            <Star className='h-3.5 w-3.5 fill-current text-amber-500' />
                          )}
                        </div>
                        <div className='text-muted-foreground mt-1 truncate text-xs'>
                          {card.title}
                        </div>
                        <div className='text-muted-foreground mt-1 text-xs'>
                          {card.billing_type === 'enterprise'
                            ? t('Enterprise')
                            : t('Personal')}
                          {' · '}
                          {card.email}
                        </div>
                      </div>
                      <div className='flex shrink-0 items-center gap-1'>
                        <Button
                          type='button'
                          variant='ghost'
                          size='icon'
                          className='h-8 w-8'
                          onClick={() => startEdit(card)}
                          disabled={saving}
                        >
                          <Edit className='h-4 w-4' />
                        </Button>
                        <Button
                          type='button'
                          variant='ghost'
                          size='icon'
                          className='h-8 w-8 text-destructive hover:text-destructive'
                          onClick={() => handleDelete(card)}
                          disabled={saving}
                        >
                          <Trash2 className='h-4 w-4' />
                        </Button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </ScrollArea>

          <div className='min-h-0 space-y-3 overflow-y-auto rounded-lg border p-3'>
            <div className='flex items-center justify-between gap-2'>
              <h3 className='text-sm font-semibold'>
                {editingId == null ? t('Add Title Card') : t('Edit Title Card')}
              </h3>
              {editingId != null && (
                <Button
                  type='button'
                  variant='ghost'
                  size='sm'
                  onClick={resetForm}
                  disabled={saving}
                >
                  <Plus className='mr-1 h-3.5 w-3.5 shrink-0' />
                  {t('New')}
                </Button>
              )}
            </div>

            <div className='grid gap-3'>
              <div className='space-y-1.5'>
                <Label>{t('Card Name')}</Label>
                <Input
                  value={form.name}
                  onChange={(e) => updateField('name', e.target.value)}
                  placeholder={t('Enter card name')}
                />
              </div>
              <div className='space-y-1.5'>
                <Label>{t('Billing Type')}</Label>
                <Select
                  items={[
                    { value: 'personal', label: t('Personal') },
                    { value: 'enterprise', label: t('Enterprise') },
                  ]}
                  value={form.billing_type}
                  onValueChange={(value) =>
                    updateField(
                      'billing_type',
                      (value || 'personal') as InvoiceTitleFormState['billing_type']
                    )
                  }
                >
                  <SelectTrigger className='w-full'>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent alignItemWithTrigger={false}>
                    <SelectGroup>
                      <SelectItem value='personal'>{t('Personal')}</SelectItem>
                      <SelectItem value='enterprise'>
                        {t('Enterprise')}
                      </SelectItem>
                    </SelectGroup>
                  </SelectContent>
                </Select>
              </div>
              <div className='space-y-1.5'>
                <Label>{t('Invoice Title')}</Label>
                <Input
                  value={form.title}
                  onChange={(e) => updateField('title', e.target.value)}
                  placeholder={t('Enter invoice title')}
                />
              </div>
              {form.billing_type === 'enterprise' && (
                <div className='space-y-1.5'>
                  <Label>{t('Tax ID')}</Label>
                  <Input
                    value={form.tax_id}
                    onChange={(e) => updateField('tax_id', e.target.value)}
                    placeholder={t('Enter tax identification number')}
                  />
                </div>
              )}
              <div className='space-y-1.5'>
                <Label>{t('Email')}</Label>
                <Input
                  type='email'
                  value={form.email}
                  onChange={(e) => updateField('email', e.target.value)}
                  placeholder={t('Enter email to receive invoice')}
                />
              </div>
              <div className='space-y-2'>
                <Label>{t('Billing Address')}</Label>
                <div className='space-y-1.5'>
                  <Label className='text-muted-foreground text-xs'>
                    {t('Street Address')}
                  </Label>
                  <Input
                    value={form.street}
                    onChange={(e) => updateField('street', e.target.value)}
                    placeholder={t('Enter street address')}
                  />
                </div>
                <div className='space-y-1.5'>
                  <Label className='text-muted-foreground text-xs'>
                    {t('Detailed Address')}
                  </Label>
                  <Input
                    value={form.address_detail}
                    onChange={(e) =>
                      updateField('address_detail', e.target.value)
                    }
                    placeholder={t('Enter detailed address')}
                  />
                </div>
                <div className='grid gap-2 sm:grid-cols-2'>
                  <div className='space-y-1.5'>
                    <Label className='text-muted-foreground text-xs'>
                      {t('City')}
                    </Label>
                    <Input
                      value={form.city}
                      onChange={(e) => updateField('city', e.target.value)}
                      placeholder={t('Enter city')}
                    />
                  </div>
                  <div className='space-y-1.5'>
                    <Label className='text-muted-foreground text-xs'>
                      {t('Zip Code')}
                    </Label>
                    <Input
                      value={form.zip_code}
                      onChange={(e) => updateField('zip_code', e.target.value)}
                      placeholder={t('Zip Code')}
                    />
                  </div>
                  <div className='space-y-1.5 sm:col-span-2'>
                    <Label className='text-muted-foreground text-xs'>
                      {t('Country Code')}
                    </Label>
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
                              !form.country && 'text-muted-foreground'
                            )}
                          />
                        }
                      >
                        {form.country
                          ? `${COUNTRY_OPTIONS.find((country) => country.code === form.country)?.name || form.country} (${form.country})`
                          : t('Country Code')}
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
                                    updateField('country', country.code)
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
                                      form.country === country.code
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
                  </div>
                </div>
              </div>
              <Button
                type='button'
                variant={form.is_default ? 'default' : 'outline'}
                onClick={() => updateField('is_default', !form.is_default)}
              >
                <Check className='mr-1.5 h-4 w-4' />
                {form.is_default ? t('Default Card') : t('Set as Default')}
              </Button>
            </div>
          </div>
        </div>

        <DialogFooter className='gap-2 border-t pt-3 max-sm:grid max-sm:grid-cols-2'>
          <Button
            type='button'
            variant='outline'
            className='max-sm:w-full'
            onClick={() => onOpenChange(false)}
            disabled={saving}
          >
            {t('Close')}
          </Button>
          <Button
            type='button'
            className='max-sm:w-full'
            onClick={handleSave}
            disabled={saving}
          >
            {saving ? t('Saving...') : t('Save')}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
