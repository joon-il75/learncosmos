const API_BASE = ''

interface AuthStore {
  token: null
  setToken: (_token: string) => void
  clearToken: () => void
}

export const useAuthStore = {
  getState: (): AuthStore => ({
    token: null,
    setToken: () => {},
    clearToken: () => {},
  }),
}

export async function authFetch(input: RequestInfo, init?: RequestInit): Promise<Response> {
  const headers = new Headers(init?.headers)
  const res = await fetch(input, { ...init, headers, credentials: 'include' })

  if (res.status === 401) {
    const refreshRes = await fetch(`${API_BASE}/api/v1/auth/refresh`, {
      method: 'POST',
      credentials: 'include',
    })
    if (refreshRes.ok) {
      return fetch(input, { ...init, headers, credentials: 'include' })
    }
    window.location.href = '/login'
  }

  return res
}

export async function logout(): Promise<void> {
  await fetch(`${API_BASE}/api/v1/auth/logout`, {
    method: 'POST',
    credentials: 'include',
  })
  window.location.href = '/'
}
