import { useState, useCallback } from 'react'

interface FetchState<T> {
  data: T | null
  loading: boolean
  error: string | null
}

interface FetchOptions extends RequestInit {
  credentials?: RequestCredentials
}

export function useFetchWithError<T = unknown>() {
  const [state, setState] = useState<FetchState<T>>({
    data: null,
    loading: false,
    error: null,
  })

  const execute = useCallback(async (url: string, options?: FetchOptions): Promise<T | null> => {
    setState(prev => ({ ...prev, loading: true, error: null }))
    try {
      const res = await fetch(url, { credentials: 'include', ...options })
      if (!res.ok) {
        const body = await res.json().catch(() => ({}))
        const msg = body?.error ?? `HTTP ${res.status}`
        setState({ data: null, loading: false, error: msg })
        return null
      }
      const data: T = await res.json()
      setState({ data, loading: false, error: null })
      return data
    } catch (err) {
      const msg = err instanceof Error ? err.message : 'network error'
      setState({ data: null, loading: false, error: msg })
      return null
    }
  }, [])

  const reset = useCallback(() => {
    setState({ data: null, loading: false, error: null })
  }, [])

  return { ...state, execute, reset }
}

// Simple one-shot fetch helper (no React state)
export async function fetchJSON<T = unknown>(url: string, options?: FetchOptions): Promise<T> {
  const res = await fetch(url, { credentials: 'include', ...options })
  if (!res.ok) {
    const body = await res.json().catch(() => ({}))
    throw new Error(body?.error ?? `HTTP ${res.status}`)
  }
  return res.json() as Promise<T>
}
