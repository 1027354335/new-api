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
import { useTranslation } from 'react-i18next'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { agreementDataMap } from './agreement-data'

// Helper function to safely convert Unicode strings (like Chinese characters) to Base64 in both browser and Node environments.
function safeBtoa(str: string): string {
  try {
    return window.btoa(unescape(encodeURIComponent(str)))
  } catch (e) {
    return typeof Buffer !== 'undefined'
      ? Buffer.from(str, 'utf-8').toString('base64')
      : ''
  }
}

interface AgreementDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  username: string
  email: string
}

export function AgreementDialog({
  open,
  onOpenChange,
  username,
  email,
}: AgreementDialogProps) {
  const { i18n } = useTranslation()

  // Find preferred template based on current system language, fallback to English
  const currentLang = (i18n.language || 'en').toLowerCase().split('-')[0]
  const data = agreementDataMap[currentLang] || agreementDataMap.en

  // Format today's date
  const today = new Date()
  let formattedDate = ''
  if (currentLang === 'zh' || currentLang === 'ja') {
    formattedDate = `${today.getFullYear()}年${String(today.getMonth() + 1).padStart(2, '0')}月${String(today.getDate()).padStart(2, '0')}日`
  } else {
    formattedDate = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-${String(today.getDate()).padStart(2, '0')}`
  }

  const userText = email ? `${username} (${email})` : username

  // Generate background watermark image
  // IMPORTANT: Watermark text is embedded inside SVG XML, so we must escape XML special characters
  // Using parentheses instead of < > to avoid breaking SVG XML parsing
  const rawWatermark = email ? `${username} (${email})` : username
  // Escape any XML special characters that might appear in username or email
  const watermarkText = rawWatermark
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
  // Dynamic SVG width: generous 11px per char + 100px padding, minimum 350px
  const svgWidth = Math.max(350, rawWatermark.length * 11 + 100)
  // Height must be tall enough for rotated text: text at y=140 with -25° rotation extends upward
  // 140px vertical room supports ~330px of text length before hitting top edge
  const svgHeight = 160
  const textY = svgHeight - 20
  const svg = `<svg xmlns="http://www.w3.org/2000/svg" width="${svgWidth}" height="${svgHeight}">
    <text x="10" y="${textY}" font-family="sans-serif" font-size="11" fill="#808080" fill-opacity="0.18" transform="rotate(-25, 10, ${textY})">${watermarkText}</text>
  </svg>`
  
  // Safe Base64 format that works 100% reliably in Chrome/Firefox/Safari, completely avoiding Unicode rendering or security blocking issues
  const watermarkBg = `url("data:image/svg+xml;base64,${safeBtoa(svg)}")`

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='max-sm:w-[calc(100vw-1rem)] sm:max-w-7xl sm:max-h-[85vh] flex flex-col p-0 overflow-hidden'>
        <DialogHeader className='px-6 pt-6 pb-2 border-b'>
          <DialogTitle className='text-xl font-bold text-center'>
            {data.title}
          </DialogTitle>
        </DialogHeader>

        {/* Agreement body - watermark applied as background directly on scrolling container so it covers ALL scrollable content */}
        <div
          className='flex-1 overflow-y-auto px-6 py-4 space-y-6 text-sm'
          style={{
            backgroundImage: watermarkBg,
            backgroundRepeat: 'repeat',
          }}
        >
          {/* Metadata Block */}
          <div className='p-4 rounded-lg space-y-2 text-muted-foreground border leading-relaxed select-text'>
            <div>{data.metaLabels.version}</div>
            <div>
              {data.metaLabels.effective.replace('{{date}}', formattedDate)}
            </div>
            <div>{data.metaLabels.provider}</div>
            <div>{data.metaLabels.address}</div>
            <div>{data.metaLabels.email}</div>
            <div>{data.metaLabels.user.replace('{{user}}', userText)}</div>
            <div>{data.metaLabels.signMethod}</div>
            <div>{data.metaLabels.applicable}</div>
            <div className='font-medium text-foreground/80 mt-1'>
              {data.metaLabels.declaration}
            </div>
          </div>

          {/* Section Clauses */}
          <div className='space-y-6 select-text leading-relaxed text-foreground/90'>
            {data.sections.map((section, index) => {
              return (
                <div key={index} className='space-y-2'>
                  <h3 className='font-bold text-base text-foreground border-b pb-1'>
                    {section.title}
                  </h3>
                  {section.paragraphs.map((para, pIndex) => {
                    const parts = para.split('：')
                    if (parts.length > 1 && parts[0].length < 15) {
                      const prefix = parts[0]
                      const textVal = parts.slice(1).join('：')
                      return (
                        <p key={pIndex} className='text-justify pl-1'>
                          <strong className='text-foreground'>{prefix}：</strong>
                          {textVal}
                        </p>
                      )
                    }
                    return (
                      <p key={pIndex} className='text-justify pl-1'>
                        {para}
                      </p>
                    )
                  })}
                </div>
              )
            })}
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}
