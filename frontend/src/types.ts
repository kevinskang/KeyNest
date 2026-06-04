export interface APIKeyDTO {
  id: number
  keyName: string
  keyValue: string
  url: string
  expiryDate: string
  registeredDate: string
  memo: string
  createdAt: string
  updatedAt: string
  expiryStatus: number // 0=none 1=normal 2=expiring(≤30d) 3=expired
}

export interface KeyFilter {
  keyName: string
  dateFrom: string
  dateTo: string
}

export interface CreateKeyRequest {
  keyName: string
  keyValue: string
  url: string
  expiryDate: string
  registeredDate: string
  memo: string
}

export interface UpdateKeyRequest {
  id: number
  keyName: string
  keyValue: string
  url: string
  expiryDate: string
  registeredDate: string
  memo: string
}
