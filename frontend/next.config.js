/** @type {import('next').NextConfig} */
const internalApiBase = process.env.INTERNAL_API_URL || 'http://127.0.0.1:8081'

const nextConfig = {
  output: 'standalone',
  allowedDevOrigins: [
    'learncosmos.co.kr',
    'www.learncosmos.co.kr',
    'dev.learnweavr.com',
    'http://127.0.0.1:3001',
    'http://127.0.0.1:3000',
    'https://learncosmos.co.kr',
    'https://www.learncosmos.co.kr',
    'https://dev.learnweavr.com',
    'http://dev.learnweavr.com',
    '127.0.0.1',
  ],
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: `${internalApiBase}/api/:path*`,
      },
    ]
  },
}

module.exports = nextConfig
