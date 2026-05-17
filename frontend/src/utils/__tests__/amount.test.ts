import { describe, it, expect } from 'vitest'
import { formatAmount, formatLedgerAmount, yuanToCents } from '../amount'

describe('formatAmount', () => {
  it('should convert cents to yuan format', () => {
    expect(formatAmount(3550)).toBe('¥35.50')
  })
  it('should handle zero', () => {
    expect(formatAmount(0)).toBe('¥0.00')
  })
  it('should handle whole yuan', () => {
    expect(formatAmount(10000)).toBe('¥100.00')
  })
  it('should handle large amount', () => {
    expect(formatAmount(99999999)).toBe('¥999999.99')
  })
})

describe('formatLedgerAmount', () => {
  it('income should have + prefix', () => {
    expect(formatLedgerAmount(10000, 'income')).toBe('+¥100.00')
  })
  it('expense should have - prefix', () => {
    expect(formatLedgerAmount(5000, 'expense')).toBe('-¥50.00')
  })
})

describe('yuanToCents', () => {
  it('should convert yuan to cents', () => {
    expect(yuanToCents(35.5)).toBe(3550)
  })
  it('should handle integers', () => {
    expect(yuanToCents(100)).toBe(10000)
  })
  it('should handle floating point precision', () => {
    expect(yuanToCents(0.01)).toBe(1)
    expect(yuanToCents(0.99)).toBe(99)
  })
})
