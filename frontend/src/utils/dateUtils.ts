/**
 * YYYY-MM-DD 문자열을 로컬 타임존 Date로 변환합니다.
 * new Date("2024-01-15")는 UTC 기준이라 타임존에 따라 날짜가 밀릴 수 있으므로
 * 연/월/일을 직접 파싱합니다.
 */
export function parseDate(val: string): Date | null {
  if (!val) return null
  const parts = val.split('-').map(Number)
  if (parts.length !== 3) return null
  const [year, month, day] = parts
  if (isNaN(year) || isNaN(month) || isNaN(day)) return null
  if (month < 1 || month > 12) return null
  if (day < 1 || day > 31) return null
  return new Date(year, month - 1, day)
}

/**
 * Date를 YYYY-MM-DD 문자열로 변환합니다. (로컬 타임존 기준)
 */
export function formatDate(d: Date): string {
  const yyyy = d.getFullYear()
  const mm = String(d.getMonth() + 1).padStart(2, '0')
  const dd = String(d.getDate()).padStart(2, '0')
  return `${yyyy}-${mm}-${dd}`
}

/**
 * VueDatePicker range 모드의 emit 값(Date[] | null)을
 * { from: string, to: string } 으로 변환합니다.
 */
export function rangeToDates(val: unknown): { from: string; to: string } {
  if (!Array.isArray(val) || val.length < 2) {
    return { from: '', to: '' }
  }
  const [a, b] = val as [Date, Date]
  return {
    from: a instanceof Date ? formatDate(a) : '',
    to:   b instanceof Date ? formatDate(b) : '',
  }
}
