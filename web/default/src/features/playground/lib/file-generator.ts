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
import { Document, HeadingLevel, Packer, Paragraph, TextRun } from 'docx'
import JSZip from 'jszip'
import { nanoid } from 'nanoid'
import * as XLSX from 'xlsx'
import type {
  AssistantPostProcessResult,
  GeneratedFile,
  GeneratedFileKind,
} from '../types'

type CellValue = string | number | boolean | null

interface FileSheetSpec {
  name?: string
  rows?: CellValue[][]
}

interface FileSlideSpec {
  title?: string
  bullets?: string[]
  notes?: string
}

interface FileSpec {
  type?: GeneratedFileKind
  filename?: string
  title?: string
  sheets?: FileSheetSpec[]
  paragraphs?: string[]
  slides?: FileSlideSpec[]
  summary?: string
}

interface FileSpecEnvelope {
  file?: FileSpec
}

const MIME_TYPES: Record<GeneratedFileKind, string> = {
  excel: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
  word: 'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
  powerpoint:
    'application/vnd.openxmlformats-officedocument.presentationml.presentation',
}

const EXTENSIONS: Record<GeneratedFileKind, string> = {
  excel: 'xlsx',
  word: 'docx',
  powerpoint: 'pptx',
}

export function detectFileGenerationRequest(
  text: string
): GeneratedFileKind | null {
  const normalized = text.toLowerCase()
  if (
    /\b(xlsx|excel|spreadsheet|workbook)\b/i.test(text) ||
    /表格|电子表格|工作簿/.test(text)
  ) {
    return 'excel'
  }
  if (
    /\b(ppt|pptx|powerpoint|slide deck|slides)\b/i.test(text) ||
    /幻灯片|演示文稿|演示稿/.test(text)
  ) {
    return 'powerpoint'
  }
  if (
    /\b(docx|word document|document)\b/i.test(normalized) ||
    /word|文档|报告/.test(text)
  ) {
    return 'word'
  }
  return null
}

export function buildFileInstruction(
  kind: GeneratedFileKind,
  userText: string
): string {
  return `${userText}

Create a downloadable ${kind} file for this request.
Return only a valid JSON object in this exact shape:
{
  "file": {
    "type": "${kind}",
    "filename": "short-file-name",
    "title": "document title",
    "summary": "one short sentence describing the file",
    "sheets": [{"name": "Sheet1", "rows": [["Header 1", "Header 2"], ["Value 1", "Value 2"]]}],
    "paragraphs": ["paragraph for Word documents"],
    "slides": [{"title": "Slide title", "bullets": ["Point 1", "Point 2"], "notes": "optional notes"}]
  }
}
Use sheets for excel files, paragraphs for word files, and slides for powerpoint files. Do not include markdown fences or any extra prose.`
}

export async function createGeneratedFileFromAssistant(
  content: string,
  fallbackKind?: GeneratedFileKind | null
): Promise<AssistantPostProcessResult | null> {
  const spec = parseFileSpec(content)
  const fileSpec = spec?.file
  const kind = normalizeKind(fileSpec?.type || fallbackKind)
  if (!fileSpec || !kind) return null

  const blob = await buildBlob(kind, fileSpec)
  const name = normalizeFilename(fileSpec.filename || fileSpec.title, kind)
  const file: GeneratedFile = {
    id: nanoid(),
    kind,
    name,
    url: URL.createObjectURL(blob),
    mimeType: MIME_TYPES[kind],
    size: blob.size,
  }

  return {
    content:
      fileSpec.summary ||
      fileSpec.title ||
      `Generated ${kind} file: ${name}`,
    generatedFile: file,
  }
}

function parseFileSpec(content: string): FileSpecEnvelope | null {
  const trimmed = content.trim()
  const candidates = [
    trimmed.replace(/^```(?:json)?\s*/i, '').replace(/\s*```$/i, ''),
    ...extractJsonObjects(trimmed),
  ]

  for (const candidate of candidates) {
    try {
      const parsed = JSON.parse(candidate) as FileSpecEnvelope
      if (parsed?.file && typeof parsed.file === 'object') return parsed
    } catch {
      // Try the next candidate.
    }
  }

  return null
}

function extractJsonObjects(text: string): string[] {
  const results: string[] = []
  const starts = [...text.matchAll(/\{/g)].map((match) => match.index ?? 0)

  for (const start of starts) {
    let depth = 0
    let inString = false
    let escaped = false

    for (let index = start; index < text.length; index += 1) {
      const char = text[index]
      if (inString) {
        if (escaped) {
          escaped = false
        } else if (char === '\\') {
          escaped = true
        } else if (char === '"') {
          inString = false
        }
        continue
      }

      if (char === '"') {
        inString = true
      } else if (char === '{') {
        depth += 1
      } else if (char === '}') {
        depth -= 1
        if (depth === 0) {
          results.push(text.slice(start, index + 1))
          break
        }
      }
    }
  }

  return results
}

function normalizeKind(kind?: string | null): GeneratedFileKind | null {
  if (kind === 'excel' || kind === 'word' || kind === 'powerpoint') return kind
  if (kind === 'xlsx') return 'excel'
  if (kind === 'docx') return 'word'
  if (kind === 'ppt' || kind === 'pptx') return 'powerpoint'
  return null
}

async function buildBlob(kind: GeneratedFileKind, spec: FileSpec) {
  if (kind === 'excel') return buildExcelBlob(spec)
  if (kind === 'word') return buildWordBlob(spec)
  return buildPowerPointBlob(spec)
}

function buildExcelBlob(spec: FileSpec): Blob {
  const workbook = XLSX.utils.book_new()
  const sheets =
    spec.sheets?.length && spec.sheets.some((sheet) => sheet.rows?.length)
      ? spec.sheets
      : [
          {
            name: spec.title || 'Sheet1',
            rows: [
              [spec.title || 'Generated file'],
              ...(spec.paragraphs ?? []).map((paragraph) => [paragraph]),
            ],
          },
        ]

  sheets.forEach((sheet, index) => {
    const rows = normalizeRows(sheet.rows)
    const worksheet = XLSX.utils.aoa_to_sheet(rows)
    XLSX.utils.book_append_sheet(
      workbook,
      worksheet,
      normalizeSheetName(sheet.name || `Sheet ${index + 1}`)
    )
  })

  const data = XLSX.write(workbook, { bookType: 'xlsx', type: 'array' })
  return new Blob([data], { type: MIME_TYPES.excel })
}

async function buildWordBlob(spec: FileSpec): Promise<Blob> {
  const children: Paragraph[] = []
  if (spec.title) {
    children.push(
      new Paragraph({
        text: spec.title,
        heading: HeadingLevel.TITLE,
      })
    )
  }

  const paragraphs = spec.paragraphs?.length
    ? spec.paragraphs
    : flattenSpecText(spec)
  paragraphs.forEach((text) => {
    children.push(
      new Paragraph({
        children: [new TextRun(String(text))],
        spacing: { after: 160 },
      })
    )
  })

  const document = new Document({
    sections: [{ children }],
  })
  return Packer.toBlob(document)
}

async function buildPowerPointBlob(spec: FileSpec): Promise<Blob> {
  const slides = spec.slides?.length
    ? spec.slides
    : [{ title: spec.title, bullets: flattenSpecText(spec) }]
  const zip = new JSZip()

  zip.file('[Content_Types].xml', buildPptContentTypes(slides.length))
  zip.file('_rels/.rels', ROOT_RELS_XML)
  zip.file('docProps/core.xml', buildPptCoreXml(spec.title))
  zip.file('docProps/app.xml', buildPptAppXml(slides.length))
  zip.file('ppt/presentation.xml', buildPptPresentationXml(slides.length))
  zip.file(
    'ppt/_rels/presentation.xml.rels',
    buildPptPresentationRelsXml(slides.length)
  )
  zip.file('ppt/theme/theme1.xml', PPT_THEME_XML)
  zip.file('ppt/slideMasters/slideMaster1.xml', PPT_SLIDE_MASTER_XML)
  zip.file(
    'ppt/slideMasters/_rels/slideMaster1.xml.rels',
    PPT_SLIDE_MASTER_RELS_XML
  )
  zip.file('ppt/slideLayouts/slideLayout1.xml', PPT_SLIDE_LAYOUT_XML)
  zip.file(
    'ppt/slideLayouts/_rels/slideLayout1.xml.rels',
    PPT_SLIDE_LAYOUT_RELS_XML
  )

  slides.forEach((slide, index) => {
    zip.file(`ppt/slides/slide${index + 1}.xml`, buildPptSlideXml(slide, spec))
    zip.file(`ppt/slides/_rels/slide${index + 1}.xml.rels`, SLIDE_RELS_XML)
  })

  return zip.generateAsync({
    type: 'blob',
    mimeType: MIME_TYPES.powerpoint,
    compression: 'DEFLATE',
  })
}

function buildPptContentTypes(slideCount: number) {
  const slideOverrides = Array.from({ length: slideCount }, (_, index) =>
    `<Override PartName="/ppt/slides/slide${index + 1}.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>`
  ).join('')

  return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/docProps/core.xml" ContentType="application/vnd.openxmlformats-package.core-properties+xml"/>
  <Override PartName="/docProps/app.xml" ContentType="application/vnd.openxmlformats-officedocument.extended-properties+xml"/>
  <Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>
  <Override PartName="/ppt/slideMasters/slideMaster1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slideMaster+xml"/>
  <Override PartName="/ppt/slideLayouts/slideLayout1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slideLayout+xml"/>
  <Override PartName="/ppt/theme/theme1.xml" ContentType="application/vnd.openxmlformats-officedocument.theme+xml"/>
  ${slideOverrides}
</Types>`
}

function buildPptCoreXml(title?: string) {
  const now = new Date().toISOString()
  return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:dcterms="http://purl.org/dc/terms/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <dc:title>${escapeXml(title || 'Generated presentation')}</dc:title>
  <dc:creator>new-api playground</dc:creator>
  <cp:lastModifiedBy>new-api playground</cp:lastModifiedBy>
  <dcterms:created xsi:type="dcterms:W3CDTF">${now}</dcterms:created>
  <dcterms:modified xsi:type="dcterms:W3CDTF">${now}</dcterms:modified>
</cp:coreProperties>`
}

function buildPptAppXml(slideCount: number) {
  return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Properties xmlns="http://schemas.openxmlformats.org/officeDocument/2006/extended-properties" xmlns:vt="http://schemas.openxmlformats.org/officeDocument/2006/docPropsVTypes">
  <Application>new-api playground</Application>
  <PresentationFormat>On-screen Show (16:9)</PresentationFormat>
  <Slides>${slideCount}</Slides>
</Properties>`
}

function buildPptPresentationXml(slideCount: number) {
  const slideIds = Array.from(
    { length: slideCount },
    (_, index) => `<p:sldId id="${256 + index}" r:id="rId${index + 2}"/>`
  ).join('')
  return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:presentation xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
  <p:sldMasterIdLst><p:sldMasterId id="2147483648" r:id="rId1"/></p:sldMasterIdLst>
  <p:sldIdLst>${slideIds}</p:sldIdLst>
  <p:sldSz cx="12192000" cy="6858000" type="wide"/>
  <p:notesSz cx="6858000" cy="9144000"/>
</p:presentation>`
}

function buildPptPresentationRelsXml(slideCount: number) {
  const slideRels = Array.from(
    { length: slideCount },
    (_, index) =>
      `<Relationship Id="rId${index + 2}" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide${index + 1}.xml"/>`
  ).join('')
  return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster" Target="slideMasters/slideMaster1.xml"/>
  ${slideRels}
  <Relationship Id="rId${slideCount + 2}" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme" Target="theme/theme1.xml"/>
</Relationships>`
}

function buildPptSlideXml(slideSpec: FileSlideSpec, spec: FileSpec) {
  const title = escapeXml(slideSpec.title || spec.title || 'Untitled')
  const bullets = slideSpec.bullets?.length
    ? slideSpec.bullets
    : slideSpec.notes
      ? [slideSpec.notes]
      : ['Generated content']
  const bodyRuns = bullets
    .map(
      (bullet) =>
        `<a:p><a:pPr marL="342900" indent="-171450"><a:buChar char="•"/></a:pPr><a:r><a:rPr lang="en-US" sz="1800"/><a:t>${escapeXml(bullet)}</a:t></a:r></a:p>`
    )
    .join('')

  return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
  <p:cSld>
    <p:spTree>
      <p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr>
      <p:grpSpPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="0" cy="0"/><a:chOff x="0" y="0"/><a:chExt cx="0" cy="0"/></a:xfrm></p:grpSpPr>
      <p:sp>
        <p:nvSpPr><p:cNvPr id="2" name="Title"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr>
        <p:spPr><a:xfrm><a:off x="548640" y="320040"/><a:ext cx="11094720" cy="777240"/></a:xfrm><a:prstGeom prst="rect"><a:avLst/></a:prstGeom></p:spPr>
        <p:txBody><a:bodyPr/><a:lstStyle/><a:p><a:r><a:rPr lang="en-US" sz="2800" b="1"/><a:t>${title}</a:t></a:r></a:p></p:txBody>
      </p:sp>
      <p:sp>
        <p:nvSpPr><p:cNvPr id="3" name="Content"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr>
        <p:spPr><a:xfrm><a:off x="685800" y="1371600"/><a:ext cx="10744200" cy="5029200"/></a:xfrm><a:prstGeom prst="rect"><a:avLst/></a:prstGeom></p:spPr>
        <p:txBody><a:bodyPr/><a:lstStyle/>${bodyRuns}</p:txBody>
      </p:sp>
    </p:spTree>
  </p:cSld>
  <p:clrMapOvr><a:masterClrMapping/></p:clrMapOvr>
</p:sld>`
}

function escapeXml(value: unknown) {
  return String(value ?? '')
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&apos;')
}

const ROOT_RELS_XML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/>
  <Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/>
  <Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/extended-properties" Target="docProps/app.xml"/>
</Relationships>`

const SLIDE_RELS_XML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="../slideLayouts/slideLayout1.xml"/>
</Relationships>`

const PPT_SLIDE_MASTER_RELS_XML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="../slideLayouts/slideLayout1.xml"/>
  <Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme" Target="../theme/theme1.xml"/>
</Relationships>`

const PPT_SLIDE_LAYOUT_RELS_XML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster" Target="../slideMasters/slideMaster1.xml"/>
</Relationships>`

const PPT_THEME_XML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<a:theme xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" name="Office Theme">
  <a:themeElements>
    <a:clrScheme name="Office">
      <a:dk1><a:srgbClr val="000000"/></a:dk1><a:lt1><a:srgbClr val="FFFFFF"/></a:lt1>
      <a:dk2><a:srgbClr val="1F2937"/></a:dk2><a:lt2><a:srgbClr val="F8FAFC"/></a:lt2>
      <a:accent1><a:srgbClr val="2563EB"/></a:accent1><a:accent2><a:srgbClr val="10B981"/></a:accent2>
      <a:accent3><a:srgbClr val="F59E0B"/></a:accent3><a:accent4><a:srgbClr val="EF4444"/></a:accent4>
      <a:accent5><a:srgbClr val="8B5CF6"/></a:accent5><a:accent6><a:srgbClr val="06B6D4"/></a:accent6>
      <a:hlink><a:srgbClr val="2563EB"/></a:hlink><a:folHlink><a:srgbClr val="7C3AED"/></a:folHlink>
    </a:clrScheme>
    <a:fontScheme name="Office"><a:majorFont><a:latin typeface="Arial"/></a:majorFont><a:minorFont><a:latin typeface="Arial"/></a:minorFont></a:fontScheme>
    <a:fmtScheme name="Office"><a:fillStyleLst><a:solidFill><a:schemeClr val="phClr"/></a:solidFill></a:fillStyleLst><a:lnStyleLst><a:ln w="9525"><a:solidFill><a:schemeClr val="phClr"/></a:solidFill></a:ln></a:lnStyleLst><a:effectStyleLst><a:effectStyle><a:effectLst/></a:effectStyle></a:effectStyleLst><a:bgFillStyleLst><a:solidFill><a:schemeClr val="phClr"/></a:solidFill></a:bgFillStyleLst></a:fmtScheme>
  </a:themeElements>
</a:theme>`

const PPT_SLIDE_MASTER_XML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sldMaster xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
  <p:cSld><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld>
  <p:clrMap bg1="lt1" tx1="dk1" bg2="lt2" tx2="dk2" accent1="accent1" accent2="accent2" accent3="accent3" accent4="accent4" accent5="accent5" accent6="accent6" hlink="hlink" folHlink="folHlink"/>
  <p:sldLayoutIdLst><p:sldLayoutId id="2147483649" r:id="rId1"/></p:sldLayoutIdLst>
</p:sldMaster>`

const PPT_SLIDE_LAYOUT_XML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sldLayout xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" type="blank" preserve="1">
  <p:cSld name="Blank"><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld>
  <p:clrMapOvr><a:masterClrMapping/></p:clrMapOvr>
</p:sldLayout>`

function flattenSpecText(spec: FileSpec): string[] {
  const result: string[] = []
  if (spec.summary) result.push(spec.summary)
  spec.sheets?.forEach((sheet) => {
    sheet.rows?.forEach((row) => result.push(row.map(String).join(' | ')))
  })
  spec.slides?.forEach((slide) => {
    if (slide.title) result.push(slide.title)
    slide.bullets?.forEach((bullet) => result.push(bullet))
    if (slide.notes) result.push(slide.notes)
  })
  return result.length ? result : ['Generated content']
}

function normalizeRows(rows?: CellValue[][]): CellValue[][] {
  const normalized = rows
    ?.filter(Array.isArray)
    .map((row) =>
      row.map((cell) =>
        typeof cell === 'string' ||
        typeof cell === 'number' ||
        typeof cell === 'boolean' ||
        cell === null
          ? cell
          : String(cell)
      )
    )
  return normalized?.length ? normalized : [['Generated content']]
}

function normalizeSheetName(name: string): string {
  const cleaned = name.replace(/[\\/?*[\]:]/g, ' ').trim()
  return (cleaned || 'Sheet1').slice(0, 31)
}

function normalizeFilename(name: string | undefined, kind: GeneratedFileKind) {
  const extension = EXTENSIONS[kind]
  const base = (name || `generated-${kind}`)
    .split('')
    .map((char) =>
      char.charCodeAt(0) < 32 || /[<>:"/\\|?*]/.test(char) ? '-' : char
    )
    .join('')
    .replace(/\s+/g, '-')
    .replace(new RegExp(`\\.${extension}$`, 'i'), '')
    .slice(0, 80)
  return `${base || `generated-${kind}`}.${extension}`
}
