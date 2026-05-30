import dayjs from 'dayjs'

/**
 * Truncate text to maxLen characters and append "..." if truncated.
 */
export function truncate(text: string, maxLen: number): string {
  if (!text) return ''
  if (text.length <= maxLen) return text
  return text.slice(0, maxLen) + '...'
}

/**
 * Format a date string as relative time (e.g. "刚刚", "5分钟前", "3小时前").
 */
export function timeAgo(date: string): string {
  const d = dayjs(date)
  const now = dayjs()
  const diffMins = now.diff(d, 'minute')
  if (diffMins < 1) return '刚刚'
  if (diffMins < 60) return `${diffMins}分钟前`
  const diffHours = now.diff(d, 'hour')
  if (diffHours < 24) return `${diffHours}小时前`
  const diffDays = now.diff(d, 'day')
  if (diffDays < 7) return `${diffDays}天前`
  return d.format('M月D日')
}

/**
 * Format a date string as "今天", "昨天", or "M月D日 周X".
 */
export function formatDate(date: string): string {
  const d = dayjs(date)
  const today = dayjs()
  const weekDays = ['日', '一', '二', '三', '四', '五', '六']
  if (d.isSame(today, 'day')) return '今天'
  if (d.isSame(today.subtract(1, 'day'), 'day')) return '昨天'
  return `${d.format('M月D日')} 周${weekDays[d.day()]}`
}

/**
 * Return Ant Design color name for a priority level.
 */
export function priorityColor(p: string): string {
  if (p === 'urgent') return 'red'
  if (p === 'important') return 'orange'
  return 'default'
}

/**
 * Return Chinese label for a priority level.
 */
export function priorityLabel(p: string): string {
  if (p === 'urgent') return '紧急'
  if (p === 'important') return '重要'
  return '普通'
}
