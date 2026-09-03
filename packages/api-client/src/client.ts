import type {
  CreateUrlRequest,
  CreateUrlResponse,
  LoginRequest,
  LoginResponse,
  MonitoredUrl,
  SignUpRequest,
  SignUpResponse,
  UrlStatus,
} from "./types"

export class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
  ) {
    super(message)
    this.name = "ApiError"
  }
}

export interface ApiClientOptions {
  /** Base URL of the Go API, e.g. "http://localhost:3000". No trailing slash. */
  baseUrl: string
}

// Handlers return errors via http.Error, which writes a plain-text body (not JSON) —
// see internal/handler/url.go and user.go in apps/api.
async function request<T>(
  baseUrl: string,
  path: string,
  init?: RequestInit,
): Promise<T> {
  const res = await fetch(`${baseUrl}${path}`, {
    ...init,
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...init?.headers,
    },
  })

  if (!res.ok) {
    const message = (await res.text()).trim()
    throw new ApiError(res.status, message || res.statusText)
  }

  return res.json() as Promise<T>
}

export function createApiClient({ baseUrl }: ApiClientOptions) {
  return {
    signUp: (body: SignUpRequest) =>
      request<SignUpResponse>(baseUrl, "/api/v1/signup", {
        method: "POST",
        body: JSON.stringify(body),
      }),

    login: (body: LoginRequest) =>
      request<LoginResponse>(baseUrl, "/api/v1/login", {
        method: "POST",
        body: JSON.stringify(body),
      }),

    getUrls: () => request<MonitoredUrl[]>(baseUrl, "/api/v1/urls"),

    getUrlStatus: (id: string) =>
      request<UrlStatus>(baseUrl, `/api/v1/urls/${id}`),

    createUrl: (body: CreateUrlRequest) =>
      request<CreateUrlResponse>(baseUrl, "/api/v1/urls", {
        method: "POST",
        body: JSON.stringify(body),
      }),
  }
}

export type ApiClient = ReturnType<typeof createApiClient>
