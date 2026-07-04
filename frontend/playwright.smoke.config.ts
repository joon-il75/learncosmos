import { defineConfig } from '@playwright/test';

function readConfigValue(name: string, fallback: string, pattern: RegExp): string {
  const value = process.env[name] || fallback;
  if (!pattern.test(value)) {
    throw new Error(`${name} contains unsupported characters`);
  }
  return value;
}

const port = readConfigValue('PLAYWRIGHT_PORT', '3000', /^[0-9]+$/);
const hostPattern = /^[A-Za-z0-9._:-]+$/;
const bindHost = readConfigValue('PLAYWRIGHT_BIND_HOST', '127.0.0.1', hostPattern);
const baseHost = readConfigValue('PLAYWRIGHT_BASE_HOST', '127.0.0.1', hostPattern);
const baseURL = process.env.PLAYWRIGHT_BASE_URL || `http://${baseHost}:${port}`;
const shouldStartWebServer = !process.env.PLAYWRIGHT_BASE_URL;

export default defineConfig({
  testDir: './e2e',
  testMatch: /public-route-smoke\.spec\.ts/,
  timeout: 30_000,
  expect: {
    timeout: 10_000,
  },
  fullyParallel: true,
  reporter: [['list']],
  use: {
    baseURL,
    trace: 'retain-on-failure',
    screenshot: 'only-on-failure',
  },
  webServer: shouldStartWebServer
    ? {
        command: `npm run dev -- --hostname ${bindHost} --port ${port}`,
        url: baseURL,
        reuseExistingServer: true,
        timeout: 120_000,
      }
    : undefined,
});
