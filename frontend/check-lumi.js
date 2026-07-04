const { chromium } = require('playwright');

(async () => {
  const results = {};
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage({ viewport: { width: 1600, height: 1200 } });
  await page.goto('http://127.0.0.1:3000', { waitUntil: 'networkidle' });
  await page.waitForTimeout(1500);

  results.title = await page.title();
  results.url = page.url();
  results.sectionCount = await page.locator('.lumiScrollGuideContainer, .lumiScrollGuideSection, .lumiScrollGuideContainerWithImage').count();

  const selectors = ['.lumiScrollGuideStep', '.lumiScrollGuideImage', '.lumiStickyGuideImage'];
  for (const selector of selectors) {
    const count = await page.locator(selector).count();
    const items = [];
    for (let i=0; i<count; i++) {
      const item = await page.locator(selector).nth(i).evaluate((node) => {
        const cs = getComputedStyle(node);
        const rect = node.getBoundingClientRect();
        return {
          className: node.className,
          dataset: {...node.dataset},
          transform: cs.transform,
          width: rect.width,
          height: rect.height,
          opacity: cs.opacity,
          display: cs.display,
          zIndex: cs.zIndex,
          position: cs.position,
          top: cs.top,
          left: cs.left,
          src: node.getAttribute ? node.getAttribute('src') : null,
        };
      });
      items.push({ index: i, ...item });
    }
    results[selector] = items;
  }

  const stepCount = await page.locator('.lumiScrollGuideStep').count();
  const scrollSamples = [];
  for (let i=0; i<Math.min(stepCount, 6); i++) {
    await page.locator('.lumiScrollGuideStep').nth(i).scrollIntoViewIfNeeded();
    await page.waitForTimeout(350);

    const activeStates = [];
    for (let j=0; j<stepCount; j++) {
      const stepEl = page.locator('.lumiScrollGuideStep').nth(j);
      const active = await stepEl.getAttribute('data-step-active');
      const visible = await stepEl.isVisible();
      activeStates.push({ index: j, dataStepActive: active, visible });
    }

    const stickyImage = await page.locator('.lumiStickyGuideImage').first().evaluate((node) => {
      const cs = getComputedStyle(node);
      const rect = node.getBoundingClientRect();
      return {
        className: node.className,
        dataset: {...node.dataset},
        transform: cs.transform,
        width: rect.width,
        height: rect.height,
        opacity: cs.opacity,
      };
    }).catch(() => null);

    scrollSamples.push({ stepIndex: i, activeStates, stickyImage });
  }
  results.scrollSamples = scrollSamples;

  const section = await page.locator('.lumiScrollGuideSection').first().boundingBox();
  results.sectionBox = section;

  await page.screenshot({ path: '/tmp/lumi-verify.png', fullPage: true });
  await browser.close();

  console.log(JSON.stringify(results, null, 2));
})();
