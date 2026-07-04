import { headers } from 'next/headers'

export async function GET() {
  const headerStore = await headers()
  const host = (headerStore.get('x-forwarded-host') ?? headerStore.get('host') ?? '').split(':')[0].toLowerCase()
  const isDevHost = host === 'dev.learnweavr.com'

  const body = isDevHost
    ? [
        '# LearnCosmos',
        '',
        'This development host is not intended for AI or search indexing.',
        'Use https://learncosmos.co.kr for public information.',
        '',
      ].join('\n')
    : [
        '# LearnCosmos',
        '',
        '> LearnCosmos is a Korean-first self-directed learning platform for turning goals into AI-assisted learning journeys.',
        '',
        'LearnCosmos helps learners clarify goals, generate practical learning paths, collect content, keep learning records, and complete self-directed learning journeys.',
        '',
        '## Public Pages',
        '',
        '- [Korean landing page](https://learncosmos.co.kr/): Main product introduction in Korean.',
        '- [English landing page](https://learncosmos.co.kr/en): English product introduction.',
        '- [Terms of Service](https://learncosmos.co.kr/en/terms): English terms document.',
        '- [Privacy Policy](https://learncosmos.co.kr/en/privacy): English privacy document.',
        '- [Open Source Notices](https://learncosmos.co.kr/en/open-source): Third-party license notices.',
        '',
        '## Product Summary',
        '',
        '- Category: self-directed learning, AI learning journey builder, learning records.',
        '- Core workflow: enter a learning topic, refine a goal, create a journey, collect content, record progress, complete the path.',
        '- Languages: Korean and English public pages; Korean-first product operation.',
        '- Private learner data under /dashboard is not public and should not be indexed.',
        '',
      ].join('\n')

  return new Response(body, {
    headers: {
      'Content-Type': 'text/plain; charset=utf-8',
      ...(isDevHost ? { 'X-Robots-Tag': 'noindex, nofollow, noarchive, nosnippet' } : {}),
    },
  })
}
