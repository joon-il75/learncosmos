import { expect, test } from '@playwright/test';

const publicRoutes = [
  '/',
  '/en',
  '/platform',
  '/platform/news',
  '/login',
  '/en/login',
  '/demo-login',
  '/en/demo-login',
  '/super-admin/login',
];

const guardedRoutes = [
  { path: '/dashboard', redirectPattern: /\/login(?:\?|$)/, redirectAfter: '/dashboard' },
  {
    path: '/dashboard/settings/profile',
    redirectPattern: /\/login(?:\?|$)/,
    redirectAfter: '/dashboard/settings/profile',
  },
  {
    path: '/super-admin',
    redirectPattern: /\/super-admin\/login(?:\?|$)/,
    redirectAfter: '/super-admin',
  },
];

test.describe('public route smoke', () => {
  for (const route of publicRoutes) {
    test(`GET ${route} renders without auth`, async ({ page }) => {
      const response = await page.goto(route, { waitUntil: 'domcontentloaded' });

      expect(response, `missing response for ${route}`).not.toBeNull();
      expect(response!.status(), `HTTP status for ${route}`).toBeLessThan(400);
      await expect(page.locator('body')).toContainText(/\S/);
      await expect(page.locator('body')).not.toContainText(
        /Application error|Unhandled Runtime Error|Internal Server Error/i,
      );
    });
  }

  test('/platform/news remains a read-only public surface', async ({ page }) => {
    const response = await page.goto('/platform/news', { waitUntil: 'domcontentloaded' });

    expect(response, 'missing response for /platform/news').not.toBeNull();
    expect(response!.status()).toBeLessThan(400);
    await expect(page.getByRole('button', { name: /저장|삭제|작성|수정|Save|Delete|Create|Edit/i })).toHaveCount(0);
  });
});

test.describe('public API smoke', () => {
  test('platform notice list returns public list shape', async ({ request }) => {
    const response = await request.get('/api/v1/public/platform-notices?locale=ko&limit=10');

    expect(response.status()).toBe(200);
    const body = await response.json();
    expect(Array.isArray(body.notices)).toBe(true);
    expect(typeof body.total).toBe('number');
    expect(body.locale).toBe('ko');
  });

  test('point policy exposes public cost keys only', async ({ request }) => {
    const response = await request.get('/api/v1/public/point-policy');

    expect(response.status()).toBe(200);
    const body = await response.json();
    expect(typeof body.welcome_points).toBe('number');
    expect(typeof body.course_gen_cost).toBe('number');
    expect(typeof body.lesson_rec_cost).toBe('number');
    expect(body.admin_max_grant).toBeUndefined();
  });
});

test.describe('unauthenticated gate smoke', () => {
  for (const route of guardedRoutes) {
    test(`${route.path} redirects anonymous users`, async ({ page }) => {
      await page.context().clearCookies();
      const response = await page.goto(route.path, { waitUntil: 'domcontentloaded' });

      expect(response, `missing response for ${route.path}`).not.toBeNull();
      expect(response!.status(), `final status for ${route.path}`).toBeLessThan(400);
      await page.waitForURL(route.redirectPattern);
      const redirectedURL = new URL(page.url());
      expect(redirectedURL.searchParams.get('redirect_after')).toBe(route.redirectAfter);
    });
  }
});
