// These mirror the Go API's JSON wire shapes exactly. The repo package (apps/api/internal/repo)
// has no `json:"..."` struct tags, so encoding/json serializes using the raw (capitalized) Go
// field names — that's not a typo below, it's what the server actually sends.

export interface SignUpRequest {
  username: string
  email: string
  password: string
}

export interface SignUpResponse {
  Email: string
  Username: string
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
  ID: string
  Url: string
  IntervalSeconds: number
}

export interface MonitoredUrl {
  Url: string
  IntervalSeconds: number
}
