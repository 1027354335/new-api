/*
Copyright (C) 2023-2026 OSS API

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
*/
import { useState, useMemo } from 'react'
import * as z from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { api } from '@/lib/api'
import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import { SettingsSection } from '../components/settings-section'
import { useResetForm } from '../hooks/use-reset-form'
import { useUpdateOption } from '../hooks/use-update-option'

const createStorageSchema = (t: (key: string) => string) =>
  z.object({
    storage_setting: z.object({
      enabled: z.boolean(),
      endpoint: z.string().refine((value) => {
        const trimmed = value.trim()
        if (!trimmed) return true
        // Remove protocol prefix in checking if user entered http:// or https:// manually, we just want domain:port or domain
        return !/^https?:\/\//.test(trimmed)
      }, t('Provide endpoint as host:port or host, without http:// or https:// prefix')),
      bucket: z.string().min(1, t('Bucket name is required')),
      access_key: z.string().min(1, t('Access Key is required')),
      secret_key: z.string().min(1, t('Secret Key is required')),
      use_ssl: z.boolean(),
      region: z.string().optional(),
      url_prefix: z.string().optional(),
    }),
  })

type StorageFormValues = z.output<ReturnType<typeof createStorageSchema>>
type StorageFormInput = z.input<ReturnType<typeof createStorageSchema>>

type StorageFlatValues = {
  'storage_setting.enabled': boolean
  'storage_setting.endpoint': string
  'storage_setting.bucket': string
  'storage_setting.access_key': string
  'storage_setting.secret_key': string
  'storage_setting.use_ssl': boolean
  'storage_setting.region': string
  'storage_setting.url_prefix': string
}

type StorageSettingsSectionProps = {
  defaultValues: StorageFlatValues
}

const buildFormDefaults = (
  defaults: StorageFlatValues
): StorageFormInput => ({
  storage_setting: {
    enabled: defaults['storage_setting.enabled'] ?? false,
    endpoint: defaults['storage_setting.endpoint'] ?? '',
    bucket: defaults['storage_setting.bucket'] ?? '',
    access_key: defaults['storage_setting.access_key'] ?? '',
    secret_key: defaults['storage_setting.secret_key'] ?? '',
    use_ssl: defaults['storage_setting.use_ssl'] ?? false,
    region: defaults['storage_setting.region'] ?? '',
    url_prefix: defaults['storage_setting.url_prefix'] ?? '',
  },
})

const normalizeFormValues = (
  values: StorageFormValues
): StorageFlatValues => ({
  'storage_setting.enabled': values.storage_setting.enabled,
  'storage_setting.endpoint': values.storage_setting.endpoint.trim(),
  'storage_setting.bucket': values.storage_setting.bucket.trim(),
  'storage_setting.access_key': values.storage_setting.access_key.trim(),
  'storage_setting.secret_key': values.storage_setting.secret_key.trim(),
  'storage_setting.use_ssl': values.storage_setting.use_ssl,
  'storage_setting.region': (values.storage_setting.region || '').trim(),
  'storage_setting.url_prefix': (values.storage_setting.url_prefix || '').trim(),
})

export function StorageSettingsSection({
  defaultValues,
}: StorageSettingsSectionProps) {
  const { t } = useTranslation()
  const updateOption = useUpdateOption()
  const storageSchema = createStorageSchema(t)
  const [testing, setTesting] = useState(false)

  const formDefaults = useMemo(
    () => buildFormDefaults(defaultValues),
    [defaultValues]
  )

  const form = useForm<StorageFormInput, unknown, StorageFormValues>({
    resolver: zodResolver(storageSchema),
    defaultValues: formDefaults,
  })

  useResetForm(form as any, formDefaults)

  const saveOptions = async (values: StorageFormValues) => {
    const flatValues = normalizeFormValues(values)
    const entries = Object.entries(flatValues) as [keyof StorageFlatValues, unknown][]
    const updates = entries.filter(
      ([key, value]) =>
        value !== (defaultValues[key] as unknown)
    )

    if (updates.length > 0) {
      for (const [key, value] of updates) {
        await updateOption.mutateAsync({ key, value: value as string | number | boolean })
      }
    }
  }

  const onSubmit = async (values: StorageFormValues) => {
    try {
      await saveOptions(values)
      toast.success(t('Storage settings saved successfully'))
    } catch (err: any) {
      toast.error(err.message || t('Failed to save storage settings'))
    }
  }

  const testConnection = async () => {
    const values = form.getValues()
    const isValid = await form.trigger()
    if (!isValid) {
      toast.error(t('Please fix validation errors before testing connection'))
      return
    }

    setTesting(true)
    try {
      // Save settings first so backend has current config values
      await saveOptions(values)
      const res = await api.post('/api/storage/test')
      if (res.data.success) {
        toast.success(t('MinIO storage connection test succeeded!'))
      } else {
        toast.error(res.data.message || t('MinIO connection failed'))
      }
    } catch (err: any) {
      toast.error(err.message || t('MinIO connection failed'))
    } finally {
      setTesting(false)
    }
  }

  const isEnabled = form.watch('storage_setting.enabled')

  return (
    <SettingsSection
      title={t('File Storage Settings')}
      description={t('Configure MinIO or S3 compatible file storage for chat playground media')}
    >
      <Form {...form}>
        <form
          onSubmit={form.handleSubmit(onSubmit)}
          autoComplete='off'
          className='space-y-6'
        >
          <FormField
            control={form.control}
            name='storage_setting.enabled'
            render={({ field }) => (
              <FormItem className='flex flex-row items-center justify-between rounded-lg border p-4'>
                <div className='space-y-0.5'>
                  <FormLabel className='text-base'>
                    {t('Enable File Storage')}
                  </FormLabel>
                  <FormDescription>
                    {t('Save generated and user chat images to MinIO storage')}
                  </FormDescription>
                </div>
                <FormControl>
                  <Switch
                    checked={field.value}
                    onCheckedChange={field.onChange}
                  />
                </FormControl>
              </FormItem>
            )}
          />

          {isEnabled && (
            <>
              <div className='grid gap-6 md:grid-cols-2'>
                <FormField
                  control={form.control}
                  name='storage_setting.endpoint'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('Storage Endpoint')}</FormLabel>
                      <FormControl>
                        <Input
                          placeholder='play.minio.io:9000'
                          autoComplete='off'
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>
                        {t('MinIO API server host and port (do not include http:// or https://)')}
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name='storage_setting.bucket'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('Storage Bucket')}</FormLabel>
                      <FormControl>
                        <Input
                          placeholder='playground'
                          autoComplete='off'
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>
                        {t('Bucket name in MinIO storage')}
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>

              <div className='grid gap-6 md:grid-cols-2'>
                <FormField
                  control={form.control}
                  name='storage_setting.access_key'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('Access Key')}</FormLabel>
                      <FormControl>
                        <Input
                          placeholder='S3 Access Key'
                          autoComplete='off'
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name='storage_setting.secret_key'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('Secret Key')}</FormLabel>
                      <FormControl>
                        <Input
                          type='password'
                          placeholder='S3 Secret Key'
                          autoComplete='new-password'
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>

              <div className='grid gap-6 md:grid-cols-2'>
                <FormField
                  control={form.control}
                  name='storage_setting.region'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('Region (Optional)')}</FormLabel>
                      <FormControl>
                        <Input
                          placeholder='us-east-1'
                          autoComplete='off'
                          {...field}
                          value={field.value ?? ''}
                        />
                      </FormControl>
                      <FormDescription>
                        {t('MinIO bucket region, defaults to us-east-1 if left blank')}
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name='storage_setting.use_ssl'
                  render={({ field }) => (
                    <FormItem className='flex flex-row items-center justify-between rounded-lg border p-4'>
                      <div className='space-y-0.5'>
                        <FormLabel className='text-base'>
                          {t('Use SSL (HTTPS)')}
                        </FormLabel>
                        <FormDescription>
                          {t('Connect using SSL/TLS secure connection')}
                        </FormDescription>
                      </div>
                      <FormControl>
                        <Switch
                          checked={field.value}
                          onCheckedChange={field.onChange}
                        />
                      </FormControl>
                    </FormItem>
                  )}
                />
              </div>

              <div className='grid gap-6 md:grid-cols-2'>
                <FormField
                  control={form.control}
                  name='storage_setting.url_prefix'
                  render={({ field }) => (
                    <FormItem className='md:col-span-2'>
                      <FormLabel>{t('File Access URL Prefix')}</FormLabel>
                      <FormControl>
                        <Input
                          placeholder='https://www.oss-global.de/file'
                          autoComplete='off'
                          {...field}
                          value={field.value ?? ''}
                        />
                      </FormControl>
                      <FormDescription>
                        {t('If configured, returned URLs will use this prefix instead of the local proxy endpoint (e.g. for CDN or public bucket access)')}
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
            </>
          )}

          <div className='flex gap-4'>
            <Button type='submit' disabled={updateOption.isPending}>
              {updateOption.isPending ? t('Saving...') : t('Save Storage settings')}
            </Button>
            <Button
              type='button'
              variant='outline'
              disabled={testing}
              onClick={testConnection}
            >
              {testing ? t('Testing Connection...') : t('Test Storage Connection')}
            </Button>
          </div>
        </form>
      </Form>
    </SettingsSection>
  )
}
