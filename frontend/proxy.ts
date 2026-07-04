import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

const BLOCKED_HOSTNAMES = new Set(
  (process.env.LEARNWEAVER_BLOCKED_HOSTS ?? '')
    .split(',')
    .map((host) => host.trim().toLowerCase())
    .filter(Boolean),
)
const PROTECTED_PREFIXES = ['/dashboard', '/courses', '/admin', '/super-admin']
const PUBLIC_PATHS = new Set(['/', '/login', '/auth/callback', '/terms', '/privacy', '/agreements', '/super-admin/login'])
const PUBLIC_PREFIXES = ['/auth/', '/api/', '/_next/', '/favicon']
const DEV_HOSTNAMES = new Set(['dev.learnweavr.com'])
const CSP_REPORT_ONLY = [
  "default-src 'self'",
  "script-src 'self' 'unsafe-inline' 'unsafe-eval' https:",
  "style-src 'self' 'unsafe-inline' https:",
  "img-src 'self' data: blob: https:",
  "font-src 'self' data: https:",
  "connect-src 'self' https: wss:",
  "media-src 'self' blob: https:",
  "worker-src 'self' blob:",
  "frame-src 'self' https:",
  "object-src 'none'",
  "base-uri 'self'",
  "form-action 'self'",
  "frame-ancestors 'none'",
].join('; ')
const CSP_ENFORCE = [
  "object-src 'none'",
  "base-uri 'self'",
  "form-action 'self'",
  "frame-ancestors 'none'",
].join('; ')

function withSearchPolicyHeaders(response: NextResponse, hostname: string) {
  response.headers.set('Content-Security-Policy', CSP_ENFORCE)
  response.headers.set('Content-Security-Policy-Report-Only', CSP_REPORT_ONLY)
  if (DEV_HOSTNAMES.has(hostname)) {
    response.headers.set('X-Robots-Tag', 'noindex, nofollow, noarchive, nosnippet')
  }
  return response
}

export function proxy(request: NextRequest) {
  const forwardedHost = request.headers.get('x-forwarded-host')
  const requestHost = request.headers.get('host')
  const hostname = (forwardedHost ?? requestHost ?? request.nextUrl.hostname).split(':')[0].toLowerCase()
  const { pathname } = request.nextUrl

  if (BLOCKED_HOSTNAMES.has(hostname)) {
    return new NextResponse(null, { status: 404 })
  }

  if (PUBLIC_PATHS.has(pathname)) return withSearchPolicyHeaders(NextResponse.next(), hostname)
  if (PUBLIC_PREFIXES.some((p) => pathname.startsWith(p))) return withSearchPolicyHeaders(NextResponse.next(), hostname)

  const isProtected = PROTECTED_PREFIXES.some((p) => pathname.startsWith(p))
  const isSuperAdminPath = pathname.startsWith('/super-admin')
  const isSuperAdminLoginPath = pathname === '/super-admin/login'

  if (isSuperAdminPath && !isSuperAdminLoginPath) {
    const hasSuperAdminToken = request.cookies.has('super_admin_access_token')
    if (!hasSuperAdminToken) {
      const url = request.nextUrl.clone()
      url.pathname = '/super-admin/login'
      url.searchParams.set('redirect_after', pathname)
      return withSearchPolicyHeaders(NextResponse.redirect(url), hostname)
    }
  }

  if (!isProtected || isSuperAdminPath) return withSearchPolicyHeaders(NextResponse.next(), hostname)

  const hasRefreshToken = request.cookies.has('refresh_token')
  if (!hasRefreshToken) {
    const url = request.nextUrl.clone()
    url.pathname = '/login'
    url.searchParams.set('redirect_after', pathname)
    return withSearchPolicyHeaders(NextResponse.redirect(url), hostname)
  }

  return withSearchPolicyHeaders(NextResponse.next(), hostname)
}

export const config = {
  matcher: '/((?!_next/static|_next/image|favicon.ico).*)',
}
