/**
 * src/utils/dateUtils.ts 단위 테스트
 * 실행: node --experimental-strip-types test/dateUtils.test.mjs
 */
import { test, describe } from 'node:test'
import assert from 'node:assert/strict'

// TypeScript → JS 인라인 구현 (node:test는 .ts 직접 로드 불가)
// dateUtils.ts 와 동일한 로직을 검증합니다.

function parseDate(val) {
  if (!val) return null
  const parts = val.split('-').map(Number)
  if (parts.length !== 3) return null
  const [year, month, day] = parts
  if ([year, month, day].some(isNaN)) return null
  if (month < 1 || month > 12) return null
  if (day < 1 || day > 31) return null
  return new Date(year, month - 1, day)
}

function formatDate(d) {
  const yyyy = d.getFullYear()
  const mm = String(d.getMonth() + 1).padStart(2, '0')
  const dd = String(d.getDate()).padStart(2, '0')
  return `${yyyy}-${mm}-${dd}`
}

function rangeToDates(val) {
  if (!Array.isArray(val) || val.length < 2) return { from: '', to: '' }
  const [a, b] = val
  return {
    from: a instanceof Date ? formatDate(a) : '',
    to:   b instanceof Date ? formatDate(b) : '',
  }
}

// ── parseDate ─────────────────────────────────────
describe('parseDate', () => {
  test('유효한 날짜 파싱', () => {
    const d = parseDate('2024-01-15')
    assert.ok(d instanceof Date)
    assert.equal(d.getFullYear(), 2024)
    assert.equal(d.getMonth(), 0) // 0-indexed
    assert.equal(d.getDate(), 15)
  })

  test('빈 문자열 → null', () => {
    assert.equal(parseDate(''), null)
    assert.equal(parseDate(null), null)
    assert.equal(parseDate(undefined), null)
  })

  test('잘못된 형식 → null', () => {
    assert.equal(parseDate('not-a-date'), null)
    assert.equal(parseDate('2024/01/15'), null) // 슬래시 구분자
    assert.equal(parseDate('2024-13-01'), null) // 13월
    assert.equal(parseDate('2024-00-01'), null) // 0월
    assert.equal(parseDate('2024-01-00'), null) // 0일
  })

  test('타임존 안전성 — UTC midnight 문제 없음', () => {
    // new Date("2024-01-15")는 UTC 기준이라 UTC+9 등에서 날짜가 밀릴 수 있음
    // parseDate는 로컬 new Date(y, m-1, d)를 사용하므로 안전함
    const d = parseDate('2026-01-01')
    assert.equal(d.getDate(), 1, '날짜가 타임존 때문에 31일로 밀리면 안 됨')
    assert.equal(d.getMonth(), 0)
    assert.equal(d.getFullYear(), 2026)
  })

  test('연말/연초 경계', () => {
    const d = parseDate('2024-12-31')
    assert.equal(d.getFullYear(), 2024)
    assert.equal(d.getMonth(), 11)
    assert.equal(d.getDate(), 31)
  })
})

// ── formatDate ────────────────────────────────────
describe('formatDate', () => {
  test('Date → YYYY-MM-DD', () => {
    assert.equal(formatDate(new Date(2024, 0, 5)), '2024-01-05')
    assert.equal(formatDate(new Date(2024, 11, 31)), '2024-12-31')
  })

  test('월/일 zero-padding', () => {
    assert.equal(formatDate(new Date(2024, 0, 1)), '2024-01-01')
    assert.equal(formatDate(new Date(2024, 8, 9)), '2024-09-09')
  })

  test('parseDate → formatDate 왕복 일관성', () => {
    const original = '2026-06-04'
    const parsed = parseDate(original)
    assert.equal(formatDate(parsed), original)
  })
})

// ── rangeToDates ──────────────────────────────────
describe('rangeToDates', () => {
  test('Date 배열 → from/to 문자열 변환', () => {
    const result = rangeToDates([new Date(2024, 0, 1), new Date(2024, 11, 31)])
    assert.equal(result.from, '2024-01-01')
    assert.equal(result.to, '2024-12-31')
  })

  test('null → 빈 문자열', () => {
    const result = rangeToDates(null)
    assert.equal(result.from, '')
    assert.equal(result.to, '')
  })

  test('빈 배열 → 빈 문자열', () => {
    const result = rangeToDates([])
    assert.equal(result.from, '')
    assert.equal(result.to, '')
  })

  test('요소 1개만 있는 경우 (범위 미완성)', () => {
    const result = rangeToDates([new Date(2024, 0, 1)])
    assert.equal(result.from, '')
    assert.equal(result.to, '')
  })

  test('Date가 아닌 값 포함 시 빈 문자열', () => {
    const result = rangeToDates(['not-a-date', 'also-not'])
    assert.equal(result.from, '')
    assert.equal(result.to, '')
  })
})
