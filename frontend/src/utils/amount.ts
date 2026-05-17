export function formatAmount(cents: number): string {
  const yuan = cents / 100
  return `¥${yuan.toFixed(2)}`
}

export function formatLedgerAmount(cents: number, type: 'income' | 'expense'): string {
  const yuan = cents / 100
  return type === 'income' ? `+¥${yuan.toFixed(2)}` : `-¥${yuan.toFixed(2)}`
}

export function yuanToCents(yuan: number): number {
  return Math.round(yuan * 100)
}
