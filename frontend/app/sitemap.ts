import type { MetadataRoute } from 'next'

const baseUrl = process.env.NEXT_PUBLIC_SITE_URL || 'https://learncosmos.co.kr'

export default function sitemap(): MetadataRoute.Sitemap {
  const now = new Date()
  const publicRoutes = [
    { path: '/', priority: 1 },
    { path: '/en', priority: 0.95 },
    { path: '/terms', priority: 0.55 },
    { path: '/privacy', priority: 0.55 },
    { path: '/open-source', priority: 0.45 },
    { path: '/en/terms', priority: 0.55 },
    { path: '/en/privacy', priority: 0.55 },
    { path: '/en/open-source', priority: 0.45 },
    { path: '/en/youtube-api-use', priority: 0.35 },
  ]

  return publicRoutes.map((route) => ({
    url: `${baseUrl}${route.path}`,
    lastModified: now,
    changeFrequency: route.path === '/' || route.path === '/en' ? 'weekly' : 'monthly',
    priority: route.priority,
  }))
}
