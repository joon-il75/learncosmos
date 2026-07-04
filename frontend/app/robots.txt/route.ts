import { headers } from 'next/headers'

const PUBLIC_HOSTS = new Set(['learncosmos.co.kr', 'www.learncosmos.co.kr'])
const DEV_HOSTS = new Set(['dev.learnweavr.com'])

export async function GET() {
  const headerStore = await headers()
  const host = (headerStore.get('x-forwarded-host') ?? headerStore.get('host') ?? '').split(':')[0].toLowerCase()
  const isPublicHost = PUBLIC_HOSTS.has(host)
  const isDevHost = DEV_HOSTS.has(host)

  const body = isPublicHost
    ? [
        'User-agent: *',
        'Allow: /',
        'Disallow: /dashboard/',
        'Disallow: /admin/',
        'Disallow: /super-admin/',
        'Disallow: /api/',
        '',
        'Sitemap: https://learncosmos.co.kr/sitemap.xml',
        '',
      ].join('\n')
    : [
        'User-agent: *',
        'Disallow: /',
        '',
        'User-agent: GPTBot',
        'Disallow: /',
        '',
        'User-agent: ChatGPT-User',
        'Disallow: /',
        '',
        'User-agent: ClaudeBot',
        'Disallow: /',
        '',
        'User-agent: anthropic-ai',
        'Disallow: /',
        '',
        'User-agent: PerplexityBot',
        'Disallow: /',
        '',
      ].join('\n')

  return new Response(body, {
    headers: {
      'Content-Type': 'text/plain; charset=utf-8',
      ...(isDevHost ? { 'X-Robots-Tag': 'noindex, nofollow, noarchive, nosnippet' } : {}),
    },
  })
}
