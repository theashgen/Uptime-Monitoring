// These mirror the Go API's handler-owned response contracts exactly (see
// apps/api/internal/handler/*.go). The API never serializes its sqlc-generated
// repo types directly — every response is a dedicated Go struct the handler maps
// onto — so these types are stable against schema/query changes, not tied to them.

export interface SignUpRequest {
  username: string
  email: string
  password: string
}

export interface SignUpResponse {
  email: string
  username: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  message: string
}

export interface CreateUrlRequest {
  url: string
  interval: number
}

export interface CreateUrlResponse {
  id: string
  url: string
  intervalSeconds: number
}

export interface MonitoredUrl {
  id: string
  url: string
  intervalSeconds: number
  nextCheckAt: string
  isActive: boolean
  createdAt: string
}

export interface UrlCheck {
  id: string
  isUp: boolean
  statusCode: number
  responseTimeMs: number
  error: string | null
  checkedAt: string
}

export interface UrlStatus {
  url: MonitoredUrl
  checks: UrlCheck[]
}
