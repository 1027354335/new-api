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
import JSZip from 'jszip'
import mammoth from 'mammoth'
import * as XLSX from 'xlsx'

const TEXT_ATTACHMENT_LIMIT = 32_000
const SHEET_ROW_LIMIT = 40
const SHEET_COL_LIMIT = 16
const SLIDE_LIMIT = 30

type ReadableAttachment = {
  url?: string
  filename?: string
  mediaType?: string
}

function getExtension(name = '') {
  const index = name.lastIndexOf('.')
  return index === -1 ? '' : name.slice(index + 1).toLowerCase()
}

function decodeXmlText(value: string) {
  const textarea = document.createElement('textarea')
  textarea.innerHTML = value
  return textarea.value
}

function dataUrlToArrayBuffer(url: string) {
  const [, data = ''] = url.split(',', 2)
  const binary = atob(data)
  const bytes = new Uint8Array(binary.length)
  for (let i = 0; i < binary.length; i++) bytes[i] = binary.charCodeAt(i)
  return bytes.buffer
}

function dataUrlToText(url: string) {
  return new TextDecoder('utf-8').decode(dataUrlToArrayBuffer(url))
}

function trimAttachmentText(text: string) {
  const normalized = text.replace(/\r\n/g, '\n').trim()
  if (normalized.length <= TEXT_ATTACHMENT_LIMIT) return normalized
  return `${normalized.slice(0, TEXT_ATTACHMENT_LIMIT)}\n\n[Content truncated]`
}

function extractTextFromXml(xml: string) {
  const matches = [...xml.matchAll(/<a:t[^>]*>([\s\S]*?)<\/a:t>/g)]
  return matches
    .map((match) => decodeXmlText(match[1]))
    .join(' ')
    .replace(/\s+/g, ' ')
    .trim()
}

async function extractSpreadsheet(file: ReadableAttachment) {
  const workbook = XLSX.read(dataUrlToArrayBuffer(file.url!), {
    type: 'array',
    cellDates: true,
  })
  const output = [`Workbook sheets: ${workbook.SheetNames.join(', ')}`]

  for (const sheetName of workbook.SheetNames.slice(0, 8)) {
    const sheet = workbook.Sheets[sheetName]
    const rows = XLSX.utils.sheet_to_json<Array<string | number | boolean>>(
      sheet,
      {
        header: 1,
        blankrows: false,
        defval: '',
      }
    )
    output.push(`\n## ${sheetName}`)
    rows.slice(0, SHEET_ROW_LIMIT).forEach((row) => {
      const values = row
        .slice(0, SHEET_COL_LIMIT)
        .map((cell) => String(cell ?? '').trim())
      if (values.some(Boolean)) output.push(values.join('\t'))
    })
  }

  return trimAttachmentText(output.join('\n'))
}

async function extractWord(file: ReadableAttachment) {
  const result = await mammoth.extractRawText({
    arrayBuffer: dataUrlToArrayBuffer(file.url!),
  })
  const warnings = result.messages?.length
    ? `\n\nWarnings:\n${result.messages
        .map((item: { message: string }) => `- ${item.message}`)
        .join('\n')}`
    : ''
  return trimAttachmentText(`${result.value}${warnings}`)
}

async function extractPowerPoint(file: ReadableAttachment) {
  const zip = await JSZip.loadAsync(dataUrlToArrayBuffer(file.url!))
  const slideFiles = Object.keys(zip.files)
    .filter((name) => /^ppt\/slides\/slide\d+\.xml$/.test(name))
    .sort((a, b) => {
      const ai = Number(a.match(/slide(\d+)\.xml/)?.[1] ?? 0)
      const bi = Number(b.match(/slide(\d+)\.xml/)?.[1] ?? 0)
      return ai - bi
    })

  const output: string[] = [`Slides: ${slideFiles.length}`]
  for (const [index, path] of slideFiles.slice(0, SLIDE_LIMIT).entries()) {
    const xml = await zip.file(path)?.async('string')
    if (!xml) continue
    const text = extractTextFromXml(xml)
    output.push(`\n## Slide ${index + 1}`)
    output.push(text || '[No readable text found]')
  }
  if (slideFiles.length > SLIDE_LIMIT) {
    output.push(`\n[Only first ${SLIDE_LIMIT} slides extracted]`)
  }

  return trimAttachmentText(output.join('\n'))
}

export async function extractAttachmentText(file: ReadableAttachment) {
  if (!file.url || file.mediaType?.startsWith('image/')) return undefined
  const extension = getExtension(file.filename)

  try {
    if (
      ['txt', 'md', 'markdown', 'csv', 'tsv', 'json', 'xml', 'html'].includes(
        extension
      )
    ) {
      return trimAttachmentText(dataUrlToText(file.url))
    }
    if (['xlsx', 'xlsm', 'xlsb', 'ods', 'csv'].includes(extension)) {
      return extractSpreadsheet(file)
    }
    if (extension === 'docx') {
      return extractWord(file)
    }
    if (extension === 'pptx') {
      return extractPowerPoint(file)
    }
    if (extension === 'xls') {
      return 'Legacy .xls binary spreadsheets are not reliably readable in the browser. Please upload .xlsx or export as CSV.'
    }
    if (['doc', 'ppt'].includes(extension)) {
      return `Legacy .${extension} binary documents are not readable in the browser parser. Please upload .${extension}x.`
    }
  } catch (error) {
    return `Failed to extract readable content: ${
      error instanceof Error ? error.message : 'unknown error'
    }`
  }

  return undefined
}
