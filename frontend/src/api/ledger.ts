import request from '@/utils/request'

export function getLedgers(params: {
  start_date?: string
  end_date?: string
  category_id?: number
  creator_id?: number
  limit?: number
  cursor?: string
}, signal?: AbortSignal) {
  return request.get('/ledgers', { params, signal })
}

export function getLedgerById(id: number) {
  return request.get(`/ledgers/${id}`)
}

export function createLedger(data: {
  amount: number
  note?: string
  category_id: number
  occurred_at?: string
}) {
  return request.post('/ledgers', data)
}

export function updateLedger(id: number, data: {
  amount?: number
  note?: string
  category_id?: number
  occurred_at?: string
}) {
  return request.put(`/ledgers/${id}`, data)
}

export function deleteLedger(id: number) {
  return request.delete(`/ledgers/${id}`)
}
